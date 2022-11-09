package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/go-resty/resty/v2"
)

const (
	GetScore       = "/models/inference"
	DefaultScore   = float64(-1)
	HTTPStatusCode = "http.status_code"
)

type OmniscoreClient interface {
	GetScore(ctx context.Context, charge entities.ChargeRequest) (float64, error)
}

type omniscoreClient struct {
	config config.Config
	logs   logs.Logger
}

func NewOmniscoreClient(cfg config.Config, logger logs.Logger) OmniscoreClient {
	return &omniscoreClient{
		config: cfg,
		logs:   logger,
	}
}

func (o *omniscoreClient) GetScore(ctx context.Context, charge entities.ChargeRequest) (float64, error) {
	var err error
	var resp *resty.Response
	if !o.config.Omniscore.IsEnabled {
		return DefaultScore, nil
	}

	span, _ := tracer.StartSpanFromContext(ctx, "omniscore")
	span.SetTag("http.host", o.config.Omniscore.Host)
	span.SetTag("http.url", GetScore)
	span.SetTag("charge_id", charge.ID)
	span.SetTag("company_id", charge.CompanyID)
	defer span.Finish(tracer.WithError(err))

	scoreEndpointURL := fmt.Sprintf("%s%s", o.config.Omniscore.Host, GetScore)
	var score *float64

	timeout := time.Duration(o.config.Omniscore.TimeoutMilliseconds) * time.Millisecond
	req := resty.New().SetTimeout(timeout).R().
		SetContext(ctx).
		SetBody(charge).
		SetHeader("X-Application-ID", o.config.ProjectName).
		SetResult(&score)

	err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header))
	if err != nil {
		o.logs.Error(ctx, err.Error(), "charge_id", charge.ID,
			"company_id", charge.CompanyID)
		span.SetTag(HTTPStatusCode, http.StatusInternalServerError)
		span.SetTag("error", err)
		return DefaultScore, err
	}

	resp, err = req.Post(scoreEndpointURL)
	if err != nil {
		o.logs.Error(ctx, err.Error(), "charge_id", charge.ID,
			"company_id", charge.CompanyID, text.LogTagMethod, "GetScore")
		span.SetTag(HTTPStatusCode, http.StatusInternalServerError)
		span.SetTag("error", err)
		return DefaultScore, err
	}

	if !resp.IsSuccess() {
		err = fmt.Errorf("error getting omniscore with https status code: %d, body %s", resp.StatusCode(),
			string(resp.Body()))
		o.logs.Error(ctx, err.Error(), "charge_id", charge.ID,
			"company_id", charge.CompanyID, text.LogTagMethod, "GetScore", "http_status", resp.StatusCode())
		span.SetTag(HTTPStatusCode, resp.StatusCode())
		span.SetTag("error", err)
		return DefaultScore, err
	}

	if score == nil {
		err = fmt.Errorf("score was nil with https status code: %d, body %s", resp.StatusCode(),
			string(resp.Body()))
		o.logs.Error(ctx, err.Error(), "charge_id", charge.ID,
			"company_id", charge.CompanyID, text.LogTagMethod, "GetScore", "http_status", resp.StatusCode())
		span.SetTag(HTTPStatusCode, resp.StatusCode())
		span.SetTag("error", err)
		return DefaultScore, err
	}

	span.SetTag(HTTPStatusCode, resp.StatusCode())
	return *score, nil
}
