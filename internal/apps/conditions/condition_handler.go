package conditions

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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const handlerName = "condition.handler.%s"

type ConditionHandler interface {
	Add(ctx echo.Context) error
	GetPaged(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}

type conditionHandler struct {
	logs    logs.Logger
	service ConditionService
}

func NewConditionsHandler(service ConditionService, logger logs.Logger) ConditionHandler {
	return &conditionHandler{
		logs:    logger,
		service: service,
	}
}

func (handler *conditionHandler) Add(ctx echo.Context) error {
	var conditionRequest entities.ConditionRequest
	err := ctx.Bind(&conditionRequest)
	if err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Add"))
		ctx.Error(err)
		return nil
	}

	err = ctx.Validate(&conditionRequest)
	if err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Add"))
		ctx.Error(err)
		return nil
	}

	err = handler.service.Add(ctx.Request().Context(), conditionRequest.NewConditionFromPostRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusCreated)
}

func (handler *conditionHandler) GetPaged(ctx echo.Context) error {
	var filter entities.ConditionsFilter
	pagination := entities.NewDefaultPagination()
	ctx.Bind(&pagination)
	ctx.Bind(&filter)

	pagedConditions, err := handler.service.GetAll(ctx.Request().Context(), filter, pagination)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, pagedConditions)
}

func (handler *conditionHandler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	if str.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	conditionRequest := new(entities.ConditionRequest)
	if err := ctx.Bind(conditionRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(conditionRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Update(ctx.Request().Context(), id, conditionRequest.NewConditionFromPutRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler *conditionHandler) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if str.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err := customHttp.NewBadRequestError(errors.New("invalid id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	err = handler.service.Delete(ctx.Request().Context(), id)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusOK)
}
