package modules

import (
	"errors"
	"fmt"
	"net/http"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/entities"
	str "github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/labstack/echo/v4"
)

type ModuleHandler interface {
	Add(ctx echo.Context) error
	GetAll(ctx echo.Context) error
	Delete(ctx echo.Context) error
	Update(ctx echo.Context) error
}

const handlerName = "module.handler.%s"

type moduleHandler struct {
	service ModuleService
	logs    logs.Logger
}

func NewModuleHandler(modulesService ModuleService, logger logs.Logger) ModuleHandler {
	return &moduleHandler{
		service: modulesService,
		logs:    logger,
	}
}

func (handler *moduleHandler) Add(ctx echo.Context) error {
	var moduleRequest entities.ModuleRequest
	if err := ctx.Bind(&moduleRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Add"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(&moduleRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Add"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Add(ctx.Request().Context(), moduleRequest.NewModuleFromPostRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusCreated)
}

func (handler *moduleHandler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	if str.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("error: empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	moduleRequest := new(entities.ModuleRequest)
	if err := ctx.Bind(moduleRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(moduleRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Update(ctx.Request().Context(), id, moduleRequest.NewModuleFromPutRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler *moduleHandler) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if str.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("error: empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Delete(ctx.Request().Context(), id)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusOK)
}

func (handler *moduleHandler) GetAll(ctx echo.Context) error {
	var filter entities.ModuleFilter
	pagination := entities.NewDefaultPagination()
	ctx.Bind(&pagination)
	ctx.Bind(&filter)

	pagedModules, err := handler.service.GetAll(ctx.Request().Context(), pagination, filter)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, pagedModules)
}
