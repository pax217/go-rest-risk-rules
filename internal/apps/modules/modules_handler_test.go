package modules_test

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
	"github.com/conekta/risk-rules/internal/apps/modules"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const modulesUri = "/risk-rules/modules/"

func TestAdd_WhenRequestCantBeDecodedThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()
	expectedError := "Syntax error"
	badRequest := `{
    	"module": "policy_compliance",
    	"is_global": false,
	}`

	context, rec := echo.SetupAsRecorder(http.MethodPost, modulesUri, "", badRequest)
	handler := modules.NewModuleHandler(nil, logger)

	handler.Add(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.True(t, strings.Contains(restError.Message(), expectedError))
}

func TestAdd_WhenRequestIsInvalidThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()
	expectedError := "Key: 'ModuleRequest.Description' Error:Field validation for 'Description' failed on the 'required' tag"
	module := entities.ModuleRequest{
		Author: "carlos.maldonado@conekta.com",
		Name:   "policy_compliance",
	}
	bodyBytes, _ := json.Marshal(module)

	context, rec := echo.SetupAsRecorder(http.MethodPost, modulesUri, "", string(bodyBytes))
	handler := modules.NewModuleHandler(nil, logger)

	handler.Add(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedError, restError.Message())
}

func TestAdd_WhenServiceFailsThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	modulesServiceMock := new(mocks.ModuleServiceMock)
	expectedError := errors.New("connection lost")
	module := entities.ModuleRequest{
		Author:      "santiago.ceron@conekta.com",
		Name:        "policy_compliance",
		Description: "Regla para validar contratos con OXXO",
	}
	bodyBytes, _ := json.Marshal(module)

	context, rec := echo.SetupAsRecorder(http.MethodPost, modulesUri, "", string(bodyBytes))
	modulesServiceMock.Mock.On("Add",
		context.Request().Context(),
		mock.AnythingOfType("entities.Module")).Return(expectedError).Once()
	handler := modules.NewModuleHandler(modulesServiceMock, logger)

	handler.Add(context)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	modulesServiceMock.AssertExpectations(t)
}

func TestAdd_WhenProcessedOkThenReturnCreated(t *testing.T) {
	logger, _ := logs.New()
	modulesServiceMock := new(mocks.ModuleServiceMock)
	module := entities.ModuleRequest{
		Author:      "santiago.ceron@conekta.com",
		Name:        "policy_compliance",
		Description: "Regla para validar contratos con OXXO",
	}
	bodyBytes, _ := json.Marshal(module)

	context, recorder := echo.SetupAsRecorder(http.MethodPost, modulesUri, "", string(bodyBytes))
	modulesServiceMock.Mock.On("Add",
		context.Request().Context(),
		mock.AnythingOfType("entities.Module")).Return(nil).Once()
	handler := modules.NewModuleHandler(modulesServiceMock, logger)

	err := handler.Add(context)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	modulesServiceMock.AssertExpectations(t)
}

func TestUpdate_WhenIdRequestIsRequiredThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()
	expectedError := "error: empty id"
	badRequest := `{
    	"name": "Name Module",
		"description": "Description Example",
	}`

	context, rec := echo.SetupAsRecorder(http.MethodPut, modulesUri, "", badRequest)
	handler := modules.NewModuleHandler(nil, logger)

	handler.Update(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.True(t, strings.Contains(restError.Message(), expectedError))
}

func TestUpdate_WhenRequestCantBeDecodedThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()

	q := make(url.Values)
	q.Set("id", "60f6f32ba0f965ae8ae2c87e")

	expectedError := "Syntax error"
	badRequest := `{
    	"name": "Name Module",
		"description": "Description Example",
	}`

	context, rec := echo.SetupAsRecorder(http.MethodPut, modulesUri, "60f6f32ba0f965ae8ae2c87e", badRequest)
	handler := modules.NewModuleHandler(nil, logger)

	handler.Update(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.True(t, strings.Contains(restError.Message(), expectedError))
}

func TestUpdate_WhenRequestIsInvalidThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()
	expectedError := "Key: 'ModuleRequest.Description' Error:Field validation for 'Description' failed on the 'required' tag"
	module := entities.ModuleRequest{
		Name:   "policy_compliance",
		Author: "carlos.maldonado@conekta.com",
	}
	bodyBytes, _ := json.Marshal(module)

	context, rec := echo.SetupAsRecorder(http.MethodPut, modulesUri, "60f6f32ba0f965ae8ae2c87e", string(bodyBytes))
	handler := modules.NewModuleHandler(nil, logger)

	handler.Update(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, expectedError, restError.Message())
}

func TestUpdate_WhenServiceFailsThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	modulesServiceMock := new(mocks.ModuleServiceMock)
	id := "60f6f32ba0f965ae8ae2c87e"
	expectedError := errors.New("Internal Server Error")
	moduleRequest := entities.ModuleRequest{
		Author:      "santiago.ceron@conekta.com",
		Name:        "policy_compliance",
		Description: "Regla para validar contratos con OXXO",
	}
	bodyBytes, _ := json.Marshal(moduleRequest)

	context, recorder := echo.SetupAsRecorder(http.MethodPut, "/risk-rules/modules", id, string(bodyBytes))

	modulesServiceMock.Mock.On("Update",
		context.Request().Context(),
		"60f6f32ba0f965ae8ae2c87e",
		mock.AnythingOfType("entities.Module")).Return(expectedError).Once()
	handler := modules.NewModuleHandler(modulesServiceMock, logger)

	handler.Update(context)

	httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
	assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
	assert.Contains(t, httpErrorResponse.Message(), expectedError.Error())
}

func TestUpdate_WhenProcessedOkThenReturnCreated(t *testing.T) {
	logger, _ := logs.New()
	modulesServiceMock := new(mocks.ModuleServiceMock)
	moduleRequest := entities.ModuleRequest{
		Name:        "policy_compliance",
		Description: "Regla para validar contratos con OXXO",
		Author:      "carlos.maldonado@conekta.com",
	}
	bodyBytes, _ := json.Marshal(moduleRequest)

	context, recorder := echo.SetupAsRecorder(http.MethodPut, modulesUri, "60f6f32ba0f965ae8ae2c87e", string(bodyBytes))
	modulesServiceMock.Mock.On("Update",
		context.Request().Context(),
		"60f6f32ba0f965ae8ae2c87e",
		mock.AnythingOfType("entities.Module"),
	).Return(nil).Once()
	handler := modules.NewModuleHandler(modulesServiceMock, logger)

	err := handler.Update(context)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
	modulesServiceMock.AssertExpectations(t)
}

func TestDelete_WhenProcessedOkThenReturnStatusOK(t *testing.T) {
	logger, _ := logs.New()
	modulesServiceMock := new(mocks.ModuleServiceMock)
	id := "60f6f32ba0f965ae8ae2c87e"
	request := `{}`
	context, recorder := echo.SetupAsRecorder(http.MethodGet, modulesUri, id, request)
	modulesServiceMock.Mock.On("Delete", context.Request().Context(), id).
		Return(nil).Once()
	handler := modules.NewModuleHandler(modulesServiceMock, logger)

	err := handler.Delete(context)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	modulesServiceMock.AssertExpectations(t)
}

func TestDelete_WhenServiceFailsThenReturnError(t *testing.T) {
	logger, _ := logs.New()
	modulesServiceMock := new(mocks.ModuleServiceMock)
	expectedError := errors.New("Internal Server Error")
	request := `{}`

	context, recorder := echo.SetupAsRecorder(http.MethodDelete, modulesUri, "60f6f32ba0f965ae8ae2c87e", request)
	modulesServiceMock.Mock.On("Delete",
		context.Request().Context(),
		"60f6f32ba0f965ae8ae2c87e").Return(expectedError).Once()
	handler := modules.NewModuleHandler(modulesServiceMock, logger)

	handler.Delete(context)

	httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
	assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
	assert.Contains(t, httpErrorResponse.Message(), expectedError.Error())
}

func TestDelete_WhenIdIsEmptyThenReturnBadRequest(t *testing.T) {
	logger, _ := logs.New()
	expectedError := "error: empty id"
	request := `{}`

	context, rec := echo.SetupAsRecorder(http.MethodDelete, modulesUri, "", request)
	handler := modules.NewModuleHandler(nil, logger)

	handler.Delete(context)

	restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.True(t, strings.Contains(restError.Message(), expectedError))
}

func TestModuleHandler_GetAll(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when list service fails", func(t *testing.T) {
		expectedErr := errors.New("error: database connection lost")

		moduleServiceMock := new(mocks.ModuleServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, modulesUri, "", "")
		handler := modules.NewModuleHandler(moduleServiceMock, logger)

		moduleServiceMock.On("GetAll",
			context.Request().Context(),
			entities.NewDefaultPagination(), entities.ModuleFilter{}).
			Return(entities.PagedResponse{}, expectedErr)

		handler.GetAll(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedErr.Error())
	})

	t.Run("when list service is ok", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetModules(),
		}

		moduleServiceMock := new(mocks.ModuleServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, modulesUri, "", "")
		handler := modules.NewModuleHandler(moduleServiceMock, logger)
		moduleServiceMock.On("GetAll",
			context.Request().Context(),
			entities.NewDefaultPagination(), entities.ModuleFilter{}).
			Return(serviceResponse, nil)

		handler.GetAll(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, serviceResponse.Data, interfaceToModules(pagedResponse.Data))
	})
}

func interfaceToModules(in interface{}) []entities.Module {
	interfaceArray := in.([]interface{})
	var moduleItem entities.Module
	var moduleJson []byte
	modules := make([]entities.Module, 0)

	for _, operatorMap := range interfaceArray {
		moduleJson, _ = json.Marshal(operatorMap)
		json.Unmarshal(moduleJson, &moduleItem)
		modules = append(modules, moduleItem)
	}
	return modules
}
