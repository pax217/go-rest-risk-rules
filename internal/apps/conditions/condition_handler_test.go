package conditions_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/conekta/risk-rules/pkg/echo"

	"github.com/conekta/risk-rules/test/testdata"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/conditions"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const conditionsUri = "/risk-rules/conditions/"

func TestConditionHandler_Add_fails(t *testing.T) {
	log, _ := logs.New()

	t.Run("json has no sense", func(t *testing.T) {
		expectedError := "Syntax error"
		malformedBody := `{"created_by": "santiago.ceron@conekta.com",}`

		mockCtx, rec := echo.SetupAsRecorder(http.MethodPost, conditionsUri, "", malformedBody)
		handler := conditions.NewConditionsHandler(nil, log)

		handler.Add(mockCtx)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("condition is not valid", func(t *testing.T) {
		expectedError := "Key: 'ConditionRequest.Description' Error:Field validation for 'Description' failed on the 'required' tag"
		request := entities.ConditionRequest{
			Author: "santiago.ceron@conekta.com",
			Name:   "and",
		}
		bodyReq, _ := json.Marshal(request)

		mockCtx, rec := echo.SetupAsRecorder(http.MethodPost, conditionsUri, "", string(bodyReq))
		handler := conditions.NewConditionsHandler(nil, log)

		handler.Add(mockCtx)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectedError, restError.Message())
	})

	t.Run("service fails", func(t *testing.T) {
		log, _ := logs.New()
		expectedErr := errors.New("database connection lost")
		request := entities.ConditionRequest{
			Author:      "santiago.ceron@conekta.com",
			Name:        "and",
			Description: "checks if two or more conditions are true",
		}
		bodyReq, _ := json.Marshal(request)
		mockedService := new(mocks.ConditionServiceMock)

		mockCtx, rec := echo.SetupAsRecorder(http.MethodPost, conditionsUri, "", string(bodyReq))
		mockedService.On("Add", mockCtx.Request().Context(), mock.AnythingOfType("entities.Condition")).
			Return(expectedErr).Once()
		handler := conditions.NewConditionsHandler(mockedService, log)

		handler.Add(mockCtx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.True(t, mockedService.AssertExpectations(t))
	})
}

func TestConditionHandler_Add_success(t *testing.T) {
	log, _ := logs.New()
	mockedService := new(mocks.ConditionServiceMock)
	request := entities.ConditionRequest{
		Author:      "santiago.ceron@conekta.com",
		Name:        "and",
		Description: "checks if two or more conditions are true",
	}
	bodyReq, _ := json.Marshal(request)

	mockCtx, recorder := echo.SetupAsRecorder(http.MethodPost, conditionsUri, "", string(bodyReq))
	mockedService.On("Add",
		mockCtx.Request().Context(),
		mock.AnythingOfType("entities.Condition")).
		Return(nil).Once()
	handler := conditions.NewConditionsHandler(mockedService, log)

	err := handler.Add(mockCtx)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.True(t, mockedService.AssertExpectations(t))
}

func TestUpdate_WhenIdRequestIsRequiredThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()
	expectedError := "empty id"
	request := `{
    	"updated_by": "carlos.maldonado@conekta.com",
    	"name": and,
    	"description": "checks if two or more conditions are true"
	}`

	context, rec := echo.SetupAsRecorder(http.MethodPut, conditionsUri, "", request)
	handler := conditions.NewConditionsHandler(nil, logger)

	handler.Update(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.True(t, strings.Contains(restError.Message(), expectedError))
}

func TestUpdate_WhenRequestCantBeDecodedThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()
	id := "61085ed7cbb88ca2154462fc"
	q := make(url.Values)
	q.Set("id", id)

	expectedError := "Syntax error"
	badRequest := `{
    	"updated_by": "carlos.maldonado@conekta.com",
    	"name": and,
    	"description": "checks if two or more conditions are true",
	}`

	context, rec := echo.SetupAsRecorder(http.MethodPut, conditionsUri, id, badRequest)
	handler := conditions.NewConditionsHandler(nil, logger)

	handler.Update(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.True(t, strings.Contains(restError.Message(), expectedError))
}

func TestUpdate_WhenRequestIsInvalidThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()
	expectedError := "Key: 'ConditionRequest.Author' Error:Field validation for 'Author' failed on the 'required' tag"
	conditionRequest := entities.ConditionRequest{
		Name:        "and",
		Description: "checks if two or more conditions are true",
	}
	bodyBytes, _ := json.Marshal(conditionRequest)

	context, rec := echo.SetupAsRecorder(http.MethodPut, conditionsUri, "61085ed7cbb88ca2154462fc", string(bodyBytes))
	handler := conditions.NewConditionsHandler(nil, logger)

	handler.Update(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedError, restError.Message())
}

func TestUpdate_WhenServiceFailsThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	conditionsServiceMock := new(mocks.ConditionServiceMock)
	expectedError := errors.New("Internal Server Error")
	id := "61085ed7cbb88ca2154462fc"
	conditionRequest := entities.ConditionRequest{
		Author:      "carlos.maldonado@conekta.com",
		Name:        "and",
		Description: "checks if two or more conditions are true",
	}
	bodyBytes, _ := json.Marshal(conditionRequest)

	context, recorder := echo.SetupAsRecorder(http.MethodPut, conditionsUri, id, string(bodyBytes))
	conditionsServiceMock.Mock.On("Update",
		context.Request().Context(),
		id,
		mock.AnythingOfType("entities.Condition")).Return(expectedError).Once()
	handler := conditions.NewConditionsHandler(conditionsServiceMock, logger)

	handler.Update(context)

	httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
	assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
	assert.Contains(t, httpErrorResponse.Message(), expectedError.Error())
}

func TestUpdate_WhenProcessedOkThenReturnCreated(t *testing.T) {
	logger, _ := logs.New()
	conditionsServiceMock := new(mocks.ConditionServiceMock)
	id := "61085ed7cbb88ca2154462fc"
	module := entities.ConditionRequest{
		Author:      "carlos.maldonado@conekta.com",
		Name:        "and",
		Description: "checks if two or more conditions are true",
	}
	bodyBytes, _ := json.Marshal(module)

	context, recorder := echo.SetupAsRecorder(http.MethodPut, conditionsUri, id, string(bodyBytes))
	conditionsServiceMock.Mock.On("Update",
		context.Request().Context(),
		id,
		mock.AnythingOfType("entities.Condition"),
	).Return(nil).Once()
	handler := conditions.NewConditionsHandler(conditionsServiceMock, logger)

	err := handler.Update(context)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
	conditionsServiceMock.AssertExpectations(t)
}

func Test_Delete_Condition_Handler(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when process ok then return status ok", func(t *testing.T) {
		conditionServiceMock := new(mocks.ConditionServiceMock)
		id := "60f6f32ba0f965ae8ae2c87e"
		request := `{}`
		context, recorder := echo.SetupAsRecorder(http.MethodDelete, conditionsUri, id, request)
		conditionServiceMock.Mock.On("Delete", context.Request().Context(), id).
			Return(nil).Once()
		handler := conditions.NewConditionsHandler(conditionServiceMock, logger)

		err := handler.Delete(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, recorder.Code)
		conditionServiceMock.AssertExpectations(t)
	})

	t.Run("when service fails then return error", func(t *testing.T) {
		conditionServiceMock := new(mocks.ConditionServiceMock)
		expectedError := errors.New("Internal Server Error")
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, conditionsUri, "60f6f32ba0f965ae8ae2c87e", request)
		conditionServiceMock.Mock.On("Delete",
			context.Request().Context(),
			"60f6f32ba0f965ae8ae2c87e").Return(expectedError).Once()

		handler := conditions.NewConditionsHandler(conditionServiceMock, logger)

		handler.Delete(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedError.Error())
		conditionServiceMock.AssertExpectations(t)
	})

	t.Run("when id request invalid then return BadRequest", func(t *testing.T) {
		expectedError := "invalid id"
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, conditionsUri, "908b9470-3985-48ad-8eb2-d6880329f234", request)
		handler := conditions.NewConditionsHandler(nil, logger)

		handler.Delete(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("when id empty then return BadRequest", func(t *testing.T) {
		expectedError := "empty id"
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, conditionsUri, "", request)
		handler := conditions.NewConditionsHandler(nil, logger)

		handler.Delete(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})
}

func TestConditionHandler_GetAll(t *testing.T) {
	logger, _ := logs.New()

	t.Run("get all conditions with default pagination", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetConditions(),
		}

		q := make(url.Values)
		q.Set("page", "0")
		q.Set("size", "5")
		filterConditions := entities.ConditionsFilter{ID: ""}
		serviceMock := new(mocks.ConditionServiceMock)

		handler := conditions.NewConditionsHandler(serviceMock, logger)

		c, rec := echo.SetupAsRecorder(http.MethodPost, conditionsUri, "?"+q.Encode(), "")
		serviceMock.On("GetAll",
			c.Request().Context(),
			filterConditions,
			entities.NewDefaultPagination()).
			Return(serviceResponse, nil)

		handler.GetPaged(c)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, serviceResponse.Data, interfaceToConditions(pagedResponse.Data))
	})

	t.Run("when id param is not valid", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    []entities.Condition{},
		}

		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		q.Set("id", "613688c2a95f286d57047d7dx")

		conditionsServiceMock := new(mocks.ConditionServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, conditionsUri, "?"+q.Encode(), "")
		handler := conditions.NewConditionsHandler(conditionsServiceMock, logger)
		conditionsServiceMock.On("GetAll",
			context.Request().Context(),
			entities.ConditionsFilter{
				ID: "613688c2a95f286d57047d7dx",
			},
			entities.NewDefaultPagination()).
			Return(entities.PagedResponse{}, nil)

		handler.GetPaged(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, serviceResponse.Data, []entities.Condition{})
	})

	t.Run("when an error occurs in the service return nil", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetConditions(),
		}

		q := make(url.Values)
		q.Set("page", "0")
		q.Set("size", "5")
		filterConditions := entities.ConditionsFilter{ID: ""}
		serviceMock := new(mocks.ConditionServiceMock)

		expectedError := errors.New("service error")

		handler := conditions.NewConditionsHandler(serviceMock, logger)

		c, rec := echo.SetupAsRecorder(http.MethodPost, conditionsUri, "?"+q.Encode(), "")
		serviceMock.On("GetAll",
			c.Request().Context(),
			filterConditions,
			entities.NewDefaultPagination()).
			Return(serviceResponse, expectedError)

		err := handler.GetPaged(c)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func interfaceToConditions(in interface{}) []entities.Condition {
	interfaceArray := in.([]interface{})
	var conditionItem entities.Condition
	var conditionJson []byte
	conditions := make([]entities.Condition, 0)

	for _, operatorMap := range interfaceArray {
		conditionJson, _ = json.Marshal(operatorMap)
		json.Unmarshal(conditionJson, &conditionItem)
		conditions = append(conditions, conditionItem)
	}
	return conditions
}
