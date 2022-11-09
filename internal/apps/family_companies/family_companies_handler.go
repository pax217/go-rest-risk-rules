package familycom

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

const handlerName = "family_companies.handler.%s"

type FamilyCompaniesHandler interface {
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
	Get(ctx echo.Context) error
}

type familyCompaniesHandler struct {
	logs    logs.Logger
	service FamilyCompaniesService
}

func NewFamilyCompaniesHandler(service FamilyCompaniesService, logger logs.Logger) FamilyCompaniesHandler {
	return &familyCompaniesHandler{
		service: service,
		logs:    logger,
	}
}

func (handler *familyCompaniesHandler) Create(ctx echo.Context) error {
	familyCompaniesRequest := new(entities.FamilyCompaniesRequest)

	if err := ctx.Bind(familyCompaniesRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(familyCompaniesRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	if err := familyCompaniesRequest.Validate(); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Create(ctx.Request().Context(), familyCompaniesRequest.NewFamilyCompaniesFromPostRequest())

	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusCreated)
}

func (handler *familyCompaniesHandler) Update(ctx echo.Context) error {
	familyCompaniesID := ctx.Param("id")
	if str.IsEmpty(familyCompaniesID) {
		err := customHttp.NewBadRequestError("id cannot be empty to update a family companies")
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	req := new(entities.FamilyCompaniesRequest)
	if err := ctx.Bind(req); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(req); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	if err := req.Validate(); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Update(ctx.Request().Context(), familyCompaniesID, req.NewFamilyCompaniesFromPutRequest())

	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler *familyCompaniesHandler) Get(ctx echo.Context) error {
	var filter entities.FamilyCompaniesFilter
	pagination := entities.NewDefaultPagination()
	ctx.Bind(&pagination)
	ctx.Bind(&filter)

	pagedFamilyCompanies, err := handler.service.Get(ctx.Request().Context(), pagination, filter)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, pagedFamilyCompanies)
}

func (handler *familyCompaniesHandler) Delete(ctx echo.Context) error {
	familyCompaniesID := ctx.Param("id")

	if str.IsEmpty(familyCompaniesID) {
		err := customHttp.NewBadRequestError("error: id cannot be empty to delete a family companies")
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Delete(ctx.Request().Context(), familyCompaniesID)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}
