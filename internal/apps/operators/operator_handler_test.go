package operators_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/conekta/risk-rules/pkg/echo"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/operators"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const operatorUri = "/risk-rules/v1/operators"

func TestOperatorHandler_AddOperator(t *testing.T) {
	logger, _ := logs.New()

	t.Run("add operator is successful", func(t *testing.T) {
		operatorRequest := entities.OperatorRequest{
			Author:      "carlos.maldonado@conekta.com",
			Name:        ">",
			Title:       "Mayor a",
			Description: "indica la desigualdad matemática de 2 numeros",
			Type:        "string",
		}
		body, _ := json.Marshal(operatorRequest)
		context, _ := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/operators", "", string(body))
		operatorServiceMock := new(mocks.OperatorServiceMock)

		operator := operatorRequest.NewOperatorFromPostRequest()

		operatorServiceMock.On("AddOperator", context.Request().Context(),
			operator).
			Return(nil).Once()

		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		err := handler.AddOperator(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, context.Response().Status)
		operatorServiceMock.AssertExpectations(t)
	})

	t.Run("add operator return bad request", func(t *testing.T) {
		operatorRequest := entities.OperatorRequest{}
		request := `{
			"author": "carlos.maldonado@conekta.com",
			"name": ">",
			"description": "indica la desigualdad matemática de 2 numeros",
			"type": "string",
			"title": "Mayor a",
		}`

		json.Unmarshal([]byte(request), &operatorRequest)
		context, _ := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/operators", "", string(request))
		operatorServiceMock := new(mocks.OperatorServiceMock)

		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		err := handler.AddOperator(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, context.Response().Status)
		operatorServiceMock.AssertExpectations(t)
	})

	t.Run("add operator fails in validation Name required", func(t *testing.T) {
		operatorRequest := entities.OperatorRequest{
			Author:      "carlos.maldonado@conekta.com",
			Title:       "Mayor a",
			Description: "indica la desigualdad matemática de 2 numeros",
			Type:        "string",
		}
		body, _ := json.Marshal(operatorRequest)
		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/operators", "", string(body))
		operatorServiceMock := new(mocks.OperatorServiceMock)

		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		_ = handler.AddOperator(context)

		assert.Equal(t, http.StatusBadRequest, context.Response().Status)
		restError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Contains(t, restError.Message(), "OperatorRequest.Name")
		operatorServiceMock.AssertExpectations(t)
	})

	t.Run("add operator service return error", func(t *testing.T) {
		expectedError := errors.New("connection lost")

		operatorRequest := entities.OperatorRequest{
			Author:      "carlos.maldonado@conekta.com",
			Name:        ">",
			Title:       "Mayor a",
			Description: "indica la desigualdad matemática de 2 numeros",
			Type:        "string",
		}
		body, _ := json.Marshal(operatorRequest)
		context, _ := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/operators", "", string(body))
		operatorServiceMock := new(mocks.OperatorServiceMock)

		operator := operatorRequest.NewOperatorFromPostRequest()

		operatorServiceMock.On("AddOperator", context.Request().Context(),
			operator).
			Return(expectedError).Once()

		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		err := handler.AddOperator(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, context.Response().Status)
		operatorServiceMock.AssertExpectations(t)
	})
}

func TestOperatorHandler_GetAll(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when list service fails", func(t *testing.T) {
		expectedErr := errors.New("error: database connection lost")

		operatorServiceMock := new(mocks.OperatorServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, operatorUri, "", "")
		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		operatorServiceMock.On("Get", context.Request().Context(),
			entities.OperatorFilter{Type: ""},
			entities.NewDefaultPagination()).
			Return(entities.PagedResponse{}, expectedErr)

		handler.GetAll(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedErr.Error())
		operatorServiceMock.AssertExpectations(t)
	})

	t.Run("when list service is ok", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetOperators(),
		}

		operatorServiceMock := new(mocks.OperatorServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, operatorUri, "", "")
		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		operatorServiceMock.On("Get",
			context.Request().Context(),
			entities.OperatorFilter{Type: ""},
			entities.NewDefaultPagination()).
			Return(serviceResponse, nil)

		handler.GetAll(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.EqualValues(t, serviceResponse.Data, interfaceToOperators(pagedResponse.Data))
		operatorServiceMock.AssertExpectations(t)
	})
}

func TestOperatorHandler_Delete(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when process ok then return status ok", func(t *testing.T) {
		operatorServiceMock := new(mocks.OperatorServiceMock)
		id := "60f6f32ba0f965ae8ae2c87e"
		request := `{}`
		context, recorder := echo.SetupAsRecorder(http.MethodDelete, "/risk-rules/v1/operators/", id, request)
		operatorServiceMock.Mock.On("Delete", context.Request().Context(), id).
			Return(nil).Once()
		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		err := handler.Delete(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, recorder.Code)
		operatorServiceMock.AssertExpectations(t)
	})

	t.Run("when id empty then return BadRequest", func(t *testing.T) {
		expectedError := "empty id"
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, "/risk-rules/v1/operators/", "", request)
		handler := operators.NewOperatorHandler(logger, nil)

		handler.Delete(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("when service fails then return error", func(t *testing.T) {
		operatorServiceMock := new(mocks.OperatorServiceMock)
		expectedError := errors.New("Internal Server Error")
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, "/risk-rules/v1/operators/", "60f6f32ba0f965ae8ae2c87e", request)
		operatorServiceMock.Mock.On("Delete",
			context.Request().Context(),
			"60f6f32ba0f965ae8ae2c87e").Return(expectedError).Once()

		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		handler.Delete(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedError.Error())
		operatorServiceMock.AssertExpectations(t)
	})
}

func TestOperatorHandler_Update(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when id is not send then return BadRequest", func(t *testing.T) {
		expectedError := "empty id"
		request := `{
						"author": "carlos.maldonado@conekta.com",
						"name": "+",
						"description": "indicates the sum"
					}`

		context, recorder := echo.SetupAsRecorder(http.MethodPut, "/risk-rules/v1/operators", "", request)
		handler := operators.NewOperatorHandler(logger, nil)

		handler.Update(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("when id empty then return BadRequest", func(t *testing.T) {
		expectedError := "empty id"
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodPut, "/risk-rules/v1/operators/", "", request)
		handler := operators.NewOperatorHandler(logger, nil)

		handler.Delete(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("when request cant be decoded then return BadRequest", func(t *testing.T) {
		id := "61085ed7cbb88ca2154462fc"
		q := make(url.Values)
		q.Set("id", id)

		expectedError := "Syntax error"
		badRequest := `{
						"author": "carlos.maldonado@conekta.com",
						"name": "+",
						"description": "indicates the sum",
						}`

		context, recorder := echo.SetupAsRecorder(http.MethodPut, "/risk-rules/v1/operators", id, badRequest)
		handler := operators.NewOperatorHandler(logger, nil)

		handler.Update(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("when request is invalid then return BadRequest", func(t *testing.T) {
		expectedError := "Key: 'OperatorRequest.Author' Error:Field validation for 'Author' failed on the 'required' tag"
		conditionRequest := entities.OperatorRequest{
			Name:        "+",
			Title:       "Más",
			Description: "indicates the sum",
		}
		bodyBytes, _ := json.Marshal(conditionRequest)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, "/risk-rules/v1/operators", "61085ed7cbb88ca2154462fc", string(bodyBytes))
		handler := operators.NewOperatorHandler(logger, nil)

		handler.Update(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("when service fails then return Error", func(t *testing.T) {
		operatorServiceMock := new(mocks.OperatorServiceMock)
		expectedError := errors.New("Internal Server Error")
		id := "61085ed7cbb88ca2154462fc"
		operatorRequest := entities.OperatorRequest{
			Author:      "carlos.maldonado@conekta.com",
			Name:        "+",
			Title:       "Más",
			Description: "indicates the sum",
			Type:        "string",
		}
		bodyBytes, _ := json.Marshal(operatorRequest)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, "/risk-rules/v1/operators", id, string(bodyBytes))
		operatorServiceMock.Mock.On("Update",
			context.Request().Context(),
			id,
			mock.AnythingOfType("entities.Operator")).Return(expectedError).Once()
		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		handler.Update(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedError.Error())
		operatorServiceMock.AssertExpectations(t)
	})

	t.Run("when processed ok then return updated", func(t *testing.T) {
		operatorServiceMock := new(mocks.OperatorServiceMock)
		id := "61085ed7cbb88ca2154462fc"
		operator := entities.OperatorRequest{
			Author:      "carlos.maldonado@conekta.com",
			Name:        "+",
			Title:       "Más",
			Description: "indicates the sum",
			Type:        "string",
		}
		bodyBytes, _ := json.Marshal(operator)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, "/risk-rules/v1/operators", id, string(bodyBytes))
		operatorServiceMock.Mock.On("Update",
			context.Request().Context(),
			id,
			mock.AnythingOfType("entities.Operator"),
		).Return(nil).Once()
		handler := operators.NewOperatorHandler(logger, operatorServiceMock)

		err := handler.Update(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
		operatorServiceMock.AssertExpectations(t)
	})
}

func interfaceToOperators(in interface{}) []entities.Operator {
	interfaceArray := in.([]interface{})
	var operatorItem entities.Operator
	var operatorJson []byte
	operators := make([]entities.Operator, 0)

	for _, operatorMap := range interfaceArray {
		operatorJson, _ = json.Marshal(operatorMap)
		json.Unmarshal(operatorJson, &operatorItem)
		operators = append(operators, operatorItem)
	}
	return operators
}
