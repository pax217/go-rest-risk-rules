package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/go-resty/resty/v2"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	errorGetLists = "error getting lists with %s"
	urlLists      = "%s/risk/lists/v1/lists"
)

type RkListsClient interface {
	ListsSearch(ctx context.Context, listsSearch entities.ListsSearch) ([]entities.List, error)
}

type rkListsRestClient struct {
	config config.Config
	logs   logs.Logger
}

func NewRkListsRestClient(cfg config.Config, logger logs.Logger) RkListsClient {
	return &rkListsRestClient{
		config: cfg,
		logs:   logger,
	}
}

func (r *rkListsRestClient) ListsSearch(ctx context.Context, listsSearch entities.ListsSearch) ([]entities.List, error) {
	var err error
	var resp *resty.Response
	jsonListsSearch, err := json.Marshal(listsSearch)
	if err != nil {
		r.logs.Error(ctx, err.Error(), text.LogTagMethod, "ListsSearch")
		return nil, err
	}

	span, _ := tracer.StartSpanFromContext(ctx, "rklists")
	span.SetTag("http.host", config.Configs.InternalService.Host)
	span.SetTag("http.url", urlLists)
	span.SetTag("listsSearch", string(jsonListsSearch))
	defer span.Finish()

	host := fmt.Sprintf(urlLists, config.Configs.InternalService.Host)

	lists := make([]entities.List, 0)

	timeout := time.Duration(r.config.InternalService.TimeoutMilliseconds) * time.Millisecond
	req := resty.New().SetTimeout(timeout).R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"email":      listsSearch.Email,
			"card_hash":  listsSearch.CardHash,
			"phone":      listsSearch.Phone,
			"company_id": listsSearch.CompanyID,
		}).
		SetHeader("X-Application-ID", r.config.ProjectName).
		SetResult(&lists)

	err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header))
	if err != nil {
		r.logs.Error(ctx, err.Error(), "listsSearch", string(jsonListsSearch), text.LogTagMethod, "ListsSearch")
		span.SetTag(HTTPStatusCode, http.StatusInternalServerError)
		span.SetTag("error", err)
		return nil, err
	}

	resp, err = req.Get(host)
	if err != nil {
		r.logs.Error(ctx, err.Error(), "listsSearch", string(jsonListsSearch), text.LogTagMethod, "ListsSearch")
		return nil, err
	}

	if !resp.IsSuccess() {
		err = fmt.Errorf(errorGetLists, listsSearch)
		r.logs.Error(ctx, err.Error(), "listsSearch", string(jsonListsSearch), text.LogTagMethod,
			"ListsSearch", "http_status", resp.StatusCode())
		return nil, err
	}

	span.SetTag(HTTPStatusCode, resp.StatusCode())
	return lists, nil
}
