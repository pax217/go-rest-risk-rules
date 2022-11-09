package operators

import (
	"errors"
	"fmt"
	"net/http"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/go_common/strings"
	"github.com/conekta/risk-rules/internal/entities"
	str "github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/labstack/echo/v4"
)

const handlerName = "operator.handler.%s"

type OperatorHandler interface {
	AddOperator(ctx echo.Context) error
	GetAll(ctx echo.Context) error
	Delete(ctx echo.Context) error
	Update(ctx echo.Context) error
}

type operatorHandler struct {
	logs            logs.Logger
	operatorService OperatorService
}

func NewOperatorHandler(logger logs.Logger, service OperatorService) OperatorHandler {
	return &operatorHandler{logs: logger, operatorService: service}
}

func (handler *operatorHandler) AddOperator(ctx echo.Context) error {
	ctxReq := ctx.Request().Context()
	operatorRequest := new(entities.OperatorRequest)
	if err := ctx.Bind(operatorRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "AddOperator"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(operatorRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "AddOperator"))
		ctx.Error(err)
		return nil
	}
	err := handler.operatorService.AddOperator(ctxReq, operatorRequest.NewOperatorFromPostRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}
	return ctx.NoContent(http.StatusCreated)
}

func (handler *operatorHandler) GetAll(ctx echo.Context) error {
	var operatorFilter entities.OperatorFilter
	pagination := entities.NewDefaultPagination()
	ctx.Bind(&pagination)
	ctx.Bind(&operatorFilter)

	if !strings.IsEmpty(operatorFilter.Type) {
		if err := operatorFilter.Validate(); err != nil {
			err = customHttp.NewBadRequestError(err.Error())
			handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "GetAll"))
			ctx.Error(err)
			return nil
		}
	}

	operators, err := handler.operatorService.Get(ctx.Request().Context(), operatorFilter, pagination)
	if err != nil {
		ctx.Error(err)
		return nil
	}
	return ctx.JSON(http.StatusOK, operators)
}

func (handler *operatorHandler) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	if str.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	operatorRequest := new(entities.OperatorRequest)
	if err := ctx.Bind(operatorRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	if err := ctx.Validate(operatorRequest); err != nil {
		err = customHttp.NewBadRequestError(err.Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		ctx.Error(err)
		return nil
	}

	err := handler.operatorService.Update(ctx.Request().Context(), id, operatorRequest.NewModuleFromPutRequest())
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler *operatorHandler) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	if str.IsEmpty(id) {
		err := customHttp.NewBadRequestError(errors.New("empty id").Error())
		handler.logs.Error(ctx.Request().Context(), err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	err := handler.operatorService.Delete(ctx.Request().Context(), id)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusOK)
}
