package merchantsscore

import (
	"net/http"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/labstack/echo/v4"
)

type MerchantsScoreHandler interface {
	MerchantScoreProcessing(ctx echo.Context) error
}

type merchantsScoreHandler struct {
	logs    logs.Logger
	config  config.Config
	service MerchantsScoreService
}

func NewMerchantsScoreHandler(cfg config.Config, logger logs.Logger, service MerchantsScoreService) MerchantsScoreHandler {
	return &merchantsScoreHandler{
		logs:    logger,
		config:  cfg,
		service: service,
	}
}

func (handler *merchantsScoreHandler) MerchantScoreProcessing(ctx echo.Context) error {
	err := handler.service.MerchantScoreProcessing(ctx.Request().Context())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusOK)
}
