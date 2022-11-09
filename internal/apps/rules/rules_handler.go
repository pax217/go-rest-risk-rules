package rules

import (
	"errors"
	"fmt"
	"net/http"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/go_common/strings"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/labstack/echo/v4"
)

const handlerName = "rule.handler.%s"

type RuleHandler interface {
	AddRule(c echo.Context) error
	UpdateRule(c echo.Context) error
	RemoveRule(c echo.Context) error
	GetPaged(c echo.Context) error
}

type ruleHandler struct {
	config  config.Config
	service RuleService
	logs    logs.Logger
}

func NewRulesHandler(cfg config.Config, service RuleService, logger logs.Logger) RuleHandler {
	return &ruleHandler{
		config:  cfg,
		service: service,
		logs:    logger,
	}
}

func (handler *ruleHandler) AddRule(ctx echo.Context) error {
	ruleReq := new(entities.RuleRequest)
	if err := ctx.Bind(ruleReq); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "AddRule"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(ruleReq); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "AddRule"))
		ctx.Error(err)
		return nil
	}

	if err := ruleReq.Validate(); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "AddRule"))
		ctx.Error(err)
		return nil
	}

	rule, err := handler.service.AddRule(ctx.Request().Context(), ruleReq.NewRuleFromPostRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, rule)
}

func (handler *ruleHandler) UpdateRule(ctx echo.Context) error {
	ruleID := ctx.Param("id")
	if strings.IsEmpty(ruleID) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "UpdateRule"))
		ctx.Error(err)
		return nil
	}

	ruleReq := new(entities.RuleRequest)
	if err := ctx.Bind(ruleReq); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "UpdateRule"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(ruleReq); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "UpdateRule"))
		ctx.Error(err)
		return nil
	}

	err := ruleReq.Validate()
	if err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "UpdateRule"))
		ctx.Error(err)
		return nil
	}

	err = handler.service.UpdateRule(ctx.Request().Context(), ruleID, ruleReq.NewRuleFromPutRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler *ruleHandler) RemoveRule(ctx echo.Context) error {
	ruleID := ctx.Param("id")
	if strings.IsEmpty(ruleID) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "RemoveRule"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.RemoveRule(ctx.Request().Context(), ruleID)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler *ruleHandler) GetPaged(ctx echo.Context) error {
	var ruleFilter entities.RuleFilter
	pagination := entities.NewDefaultPagination()
	ctx.Bind(&pagination)
	ctx.Bind(&ruleFilter)

	if !ruleFilter.IsIDValid() {
		err := customHttp.NewBadRequestError("invalid id")
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "GetPaged"))
		ctx.Error(err)
		return nil
	}

	pagedRules, err := handler.service.ListRules(ctx.Request().Context(), ruleFilter, pagination)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, pagedRules)
}
