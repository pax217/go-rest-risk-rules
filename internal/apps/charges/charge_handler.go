package charges

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/conekta/risk-rules/pkg/metrics"

	"github.com/conekta/go_common/datadog"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/labstack/echo/v4"
)

const handlerName = "charge.handler.%s"

type ChargeHandler interface {
	Evaluate(c echo.Context) error
	GetEvaluation(c echo.Context) error
	EvaluateOnlyRules(c echo.Context) error
	GetEvaluationOnlyRules(c echo.Context) error
}

type chargeHandler struct {
	config  config.Config
	service ChargeService
	logs    logs.Logger
	metrics datadog.Metricer
}

func NewChargeHandler(cfg config.Config, service ChargeService, logger logs.Logger, metricer datadog.Metricer) ChargeHandler {
	return &chargeHandler{
		config:  cfg,
		service: service,
		logs:    logger,
		metrics: metricer,
	}
}

func (handler *chargeHandler) Evaluate(ctx echo.Context) error {
	ctx, request, err := handler.bindChargeRequest(ctx)
	if err != nil {
		return nil
	}

	request.ValidateConsole()

	resp, err := handler.service.EvaluateCharge(ctx.Request().Context(), *request)
	if err != nil {
		ctx.Error(err)
		handler.sendMetricsFail(ctx.Request().Context())
		return nil
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (handler *chargeHandler) EvaluateOnlyRules(ctx echo.Context) error {
	ctx, request, err := handler.bindChargeRequest(ctx)
	if err != nil {
		return nil
	}

	request.ValidateConsoleOnlyRules()

	resp, err := handler.service.EvaluateChargeOnlyRules(ctx.Request().Context(), *request)
	if err != nil {
		ctx.Error(err)
		handler.sendMetricsFail(ctx.Request().Context())
		return nil
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (handler *chargeHandler) GetEvaluation(ctx echo.Context) error {
	id := ctx.Param("id")
	if strings.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "GetEvaluation"))
		ctx.Error(err)
		return nil
	}

	evaluation, err := handler.service.Get(ctx.Request().Context(), id)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, evaluation)
}

func (handler *chargeHandler) GetEvaluationOnlyRules(ctx echo.Context) error {
	id := ctx.Param("id")
	if strings.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "GetEvaluationOnlyRules"))
		ctx.Error(err)
		return nil
	}

	evaluation, err := handler.service.GetOnlyRules(ctx.Request().Context(), id)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, evaluation)
}

func (handler *chargeHandler) bindChargeRequest(ctx echo.Context) (echo.Context, *entities.ChargeRequest, error) {
	request := new(entities.ChargeRequest)
	if err := ctx.Bind(request); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(),
			text.Functionality, "Bind",
			text.LogTagMethod, fmt.Sprintf(handlerName, "Evaluate"))
		handler.sendMetricsFail(ctx.Request().Context())
		ctx.Error(err)
		return ctx, request, err
	}

	if err := ctx.Validate(request); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(),
			text.Functionality, "Validate",
			text.LogTagMethod, fmt.Sprintf(handlerName, "Evaluate"))
		ctx.Error(err)
		handler.sendMetricsFail(ctx.Request().Context())
		return ctx, request, err
	}

	return ctx, request, nil
}

func (handler *chargeHandler) sendMetricsFail(ctx context.Context) {
	metricData := metrics.NewMetricData(ctx, "EvaluateCharge", handlerName, handler.config.Env)
	metricData.SetResult(false)
	metrics.SendAsyncMetrics(handler.metrics, handler.logs, metricData, text.EvaluateChargeMetricName)
}
