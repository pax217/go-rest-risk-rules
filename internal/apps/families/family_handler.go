package families

import (
	"fmt"
	"net/http"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/entities"
	str "github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/labstack/echo/v4"
)

const handlerName = "family.handler.%s"

type FamilyHandler interface {
	Create(ctx echo.Context) error
	Delete(ctx echo.Context) error
	Update(ctx echo.Context) error
	Get(ctx echo.Context) error
}

type familyHandler struct {
	logs    logs.Logger
	service FamilyService
}

func NewFamilyHandler(service FamilyService, logger logs.Logger) FamilyHandler {
	return &familyHandler{
		logs:    logger,
		service: service,
	}
}

func (handler *familyHandler) Create(ctx echo.Context) error {
	familyRequest := new(entities.FamilyRequest)
	if err := ctx.Bind(familyRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(familyRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	if err := familyRequest.Validate(); err != nil {
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Create(ctx.Request().Context(), familyRequest.NewFamilyFromPostRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusCreated)
}

func (handler *familyHandler) Get(ctx echo.Context) error {
	var filter entities.FamilyFilter
	pagination := entities.NewDefaultPagination()
	ctx.Bind(&pagination)
	ctx.Bind(&filter)

	pagedFamilies, err := handler.service.Get(ctx.Request().Context(), pagination, filter)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, pagedFamilies)
}

func (handler *familyHandler) Update(ctx echo.Context) error {
	familyID := ctx.Param("id")
	if str.IsEmpty(familyID) {
		err := customHttp.NewBadRequestError("id cannot be empty to update a family")
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	request := new(entities.FamilyRequest)
	if err := ctx.Bind(request); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(request); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	if err := request.Validate(); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Update(ctx.Request().Context(), familyID, request.NewFamilyFromPutRequest())

	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler *familyHandler) Delete(ctx echo.Context) error {
	familyID := ctx.Param("id")
	if str.IsEmpty(familyID) {
		err := customHttp.NewBadRequestError("error: id cannot be empty to delete a family")
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Delete(ctx.Request().Context(), familyID)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}
