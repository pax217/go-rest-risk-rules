package fields

import (
	"errors"
	"fmt"
	"net/http"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	str "github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const handlerName = "field.handler.%s"

type FieldHandler interface {
	Create(c echo.Context) error
	GetPaged(c echo.Context) error
	Delete(c echo.Context) error
	Update(c echo.Context) error
}

type fieldsHandler struct {
	config  config.Config
	service FieldService
	logs    logs.Logger
}

func NewFieldsHandler(cfg config.Config, service FieldService, logger logs.Logger) FieldHandler {
	return &fieldsHandler{
		config:  cfg,
		service: service,
		logs:    logger,
	}
}

func (h *fieldsHandler) Create(ctx echo.Context) error {
	fieldReq := new(entities.FieldRequest)

	if err := ctx.Bind(fieldReq); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		h.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(fieldReq); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		h.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	err := h.service.AddField(ctx.Request().Context(), fieldReq.NewFieldFromPostRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}
	return ctx.NoContent(http.StatusCreated)
}

func (h *fieldsHandler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	if str.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		h.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	fieldRequest := new(entities.FieldRequest)
	if err := ctx.Bind(fieldRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		h.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(fieldRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		h.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	err := h.service.Update(ctx.Request().Context(), id, fieldRequest.NewFieldFromPutRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *fieldsHandler) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if str.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		h.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err := customHttp.NewBadRequestError(errors.New("invalid id").Error())
		h.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	err = h.service.Delete(ctx.Request().Context(), id)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *fieldsHandler) GetPaged(ctx echo.Context) error {
	var filter entities.FieldsFilter
	pagination := entities.NewDefaultPagination()
	ctx.Bind(&pagination)
	ctx.Bind(&filter)

	pagedFields, err := h.service.GetFields(ctx.Request().Context(), filter, pagination)
	if err != nil {
		ctx.Error(err)
		return nil
	}
	return ctx.JSON(http.StatusOK, pagedFields)
}
