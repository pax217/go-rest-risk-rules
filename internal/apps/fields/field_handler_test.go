package fields_test

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
	"github.com/conekta/risk-rules/internal/apps/fields"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const fieldsUri = "/risk-rules/fields/"

func TestFieldsHandler_Create(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("create new field successful", func(t *testing.T) {
		field := testdata.GetFieldRequest()
		body, _ := json.Marshal(field)
		serviceMock := mocks.NewFieldsServiceMock()
		handler := fields.NewFieldsHandler(configs, &serviceMock, logger)

		ctx, rec := echo.SetupAsRecorder(http.MethodPost, fieldsUri, "", string(body))
		serviceMock.On("AddField",
			ctx.Request().Context(),
			mock.AnythingOfType("entities.Field")).
			Return(nil)

		err := handler.Create(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("create new field fail", func(t *testing.T) {
		field := testdata.GetFieldRequestNotValid()
		body, _ := json.Marshal(field)
		serviceMock := mocks.NewFieldsServiceMock()
		handler := fields.NewFieldsHandler(configs, &serviceMock, logger)

		c, _ := echo.SetupAsRecorder(http.MethodPost, fieldsUri, "", string(body))
		serviceMock.On("AddField", c.Request().Context(), field).Return(nil)

		err := handler.Create(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, c.Response().Status)
	})

	t.Run("create new field with malformed json then return BadRequestError", func(t *testing.T) {
		expectedError := "Syntax error"
		fieldBadRequest := `{
						"name": "email",
						"description": "Representa el campo email perteneciente al charge",
						"type": "string",
						}`

		context, recorder := echo.SetupAsRecorder(http.MethodPost, fieldsUri, "", fieldBadRequest)
		handler := fields.NewFieldsHandler(configs, nil, logger)

		handler.Create(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("create new field service error return nil", func(t *testing.T) {
		field := testdata.GetFieldRequest()
		body, _ := json.Marshal(field)
		serviceMock := mocks.NewFieldsServiceMock()
		handler := fields.NewFieldsHandler(configs, &serviceMock, logger)
		expectedError := errors.New("service error")

		ctx, rec := echo.SetupAsRecorder(http.MethodPost, fieldsUri, "", string(body))
		serviceMock.On("AddField",
			ctx.Request().Context(),
			mock.AnythingOfType("entities.Field")).
			Return(expectedError)

		err := handler.Create(ctx)

		httpError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, httpError.Status())
		assert.Equal(t, expectedError.Error(), httpError.Message())
	})
}

func TestFieldHandler_GetPaged(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("get all fields with default pagination", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetFields(),
		}

		q := make(url.Values)
		q.Set("page", "0")
		q.Set("size", "5")
		filterFields := entities.FieldsFilter{ID: ""}
		serviceMock := mocks.NewFieldsServiceMock()

		handler := fields.NewFieldsHandler(configs, &serviceMock, logger)

		c, rec := echo.SetupAsRecorder(http.MethodPost, fieldsUri, "?"+q.Encode(), "")
		serviceMock.On("GetFields",
			c.Request().Context(),
			filterFields,
			entities.NewDefaultPagination()).
			Return(serviceResponse, nil)

		handler.GetPaged(c)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, serviceResponse.Data, interfaceToFields(pagedResponse.Data))
	})

	t.Run("when id param is not valid", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    []entities.Field{},
		}

		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		q.Set("id", "613688c2a95f286d57047d7d")

		fieldsServiceMock := new(mocks.FieldsServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, fieldsUri, "?"+q.Encode(), "")
		handler := fields.NewFieldsHandler(configs, fieldsServiceMock, logger)
		fieldsServiceMock.On("GetFields",
			context.Request().Context(),
			entities.FieldsFilter{
				ID: "613688c2a95f286d57047d7d",
			},
			entities.NewDefaultPagination()).
			Return(entities.PagedResponse{}, nil)

		handler.GetPaged(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, serviceResponse.Data, []entities.Field{})
	})

	t.Run("when repository return an error", func(t *testing.T) {
		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		q.Set("id", "613688c2a95f286d57047d7d")

		expectedError := errors.New("connection lost")

		fieldsServiceMock := new(mocks.FieldsServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, fieldsUri, "?"+q.Encode(), "")
		handler := fields.NewFieldsHandler(configs, fieldsServiceMock, logger)
		fieldsServiceMock.On("GetFields",
			context.Request().Context(),
			entities.FieldsFilter{
				ID: "613688c2a95f286d57047d7d",
			},
			entities.NewDefaultPagination()).
			Return(entities.PagedResponse{}, expectedError)

		err := handler.GetPaged(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestFieldHandler_Update(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when id request is required then return BadRequest", func(t *testing.T) {
		expectedError := "empty id"
		request := `
			{
				"author": "jesus.vega@conekta.com",
				"name": "email",
				"type": "string",
				"description": "Representa el campo email perteneciente al charge",
			}`

		context, recorder := echo.SetupAsRecorder(http.MethodPut, fieldsUri, "", request)
		handler := fields.NewFieldsHandler(configs, nil, logger)

		handler.Update(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedError)
	})

	t.Run("when request cant be decoded then return BadRequest", func(t *testing.T) {
		fieldID := "611c20f61f94f71ad2e181a2"
		q := make(url.Values)
		q.Set("id", fieldID)

		expectedError := "Syntax error"
		badRequest :=
			`{
				"author": "jesus.vega@conekta.com",
				"name": "email",
				"type": "string",
				"description": "Representa el campo email perteneciente al charge",
			}`

		context, recorder := echo.SetupAsRecorder(http.MethodPut, fieldsUri, fieldID, badRequest)
		handler := fields.NewFieldsHandler(configs, nil, logger)

		handler.Update(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedError)
	})

	t.Run("when request is invalid then return BadRequest", func(t *testing.T) {
		fieldID := "611c20f61f94f71ad2e181a2"
		expectedError := "Key: 'FieldRequest.Author' Error:Field validation for 'Author' failed on the 'required' tag"
		fieldRequest := entities.FieldRequest{
			Name:        "email",
			Type:        "string",
			Description: "Representa el campo email perteneciente al charge",
		}
		bodyBytes, _ := json.Marshal(fieldRequest)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, fieldsUri, fieldID, string(bodyBytes))
		handler := fields.NewFieldsHandler(configs, nil, logger)

		handler.Update(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedError)
	})

	t.Run("when service fails then return error", func(t *testing.T) {
		fieldsServiceMock := new(mocks.FieldsServiceMock)
		id := "60f6f32ba0f965ae8ae2c87e"
		expectedError := errors.New("Internal Server Error")
		fieldRequest := entities.FieldRequest{
			Name:        "email",
			Type:        "string",
			Description: "Representa el campo email perteneciente al charge",
			Author:      "carlos.maldonado@conekta.com",
		}
		bodyBytes, _ := json.Marshal(fieldRequest)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, fieldsUri, id, string(bodyBytes))

		fieldsServiceMock.Mock.On("Update",
			context.Request().Context(),
			id,
			mock.AnythingOfType("entities.Field")).Return(expectedError).Once()
		handler := fields.NewFieldsHandler(configs, fieldsServiceMock, logger)

		handler.Update(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedError.Error())
	})

	t.Run("when service fails then return error", func(t *testing.T) {
		id := "60f6f32ba0f965ae8ae2c87e"
		fieldsServiceMock := new(mocks.FieldsServiceMock)
		fieldRequest := entities.FieldRequest{
			Name:        "email",
			Type:        "string",
			Description: "Representa el campo email perteneciente al charge",
			Author:      "carlos.maldonado@conekta.com",
		}

		bodyBytes, _ := json.Marshal(fieldRequest)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, fieldsUri, id, string(bodyBytes))
		fieldsServiceMock.Mock.On("Update",
			context.Request().Context(),
			id,
			mock.AnythingOfType("entities.Field"),
		).Return(nil).Once()
		handler := fields.NewFieldsHandler(configs, fieldsServiceMock, logger)

		err := handler.Update(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
		fieldsServiceMock.AssertExpectations(t)
	})
}

func TestFieldsHandler_Delete(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when process ok then return status ok", func(t *testing.T) {
		fieldsServiceMock := new(mocks.FieldsServiceMock)
		id := "60f6f32ba0f965ae8ae2c87e"
		request := `{}`
		context, recorder := echo.SetupAsRecorder(http.MethodDelete, fieldsUri, id, request)
		fieldsServiceMock.Mock.On("Delete", context.Request().Context(), id).
			Return(nil).Once()
		handler := fields.NewFieldsHandler(configs, fieldsServiceMock, logger)

		err := handler.Delete(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, recorder.Code)
		fieldsServiceMock.AssertExpectations(t)
	})

	t.Run("when id empty then return BadRequest", func(t *testing.T) {
		expectedError := "empty id"
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, fieldsUri, "", request)
		handler := fields.NewFieldsHandler(configs, nil, logger)

		handler.Delete(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("when id is a not valid ObjectID then return BadRequest", func(t *testing.T) {
		id := "60f6f32ba0f965ae8ae2c87e-x"
		expectedError := "invalid id"
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, fieldsUri, id, request)
		handler := fields.NewFieldsHandler(configs, nil, logger)

		handler.Delete(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.True(t, strings.Contains(httpError.Message(), expectedError))
	})

	t.Run("when service fails then return error", func(t *testing.T) {
		fieldID := "60f6f32ba0f965ae8ae2c87e"
		fieldsServiceMock := new(mocks.FieldsServiceMock)
		expectedError := errors.New("Internal Server Error")
		request := `{}`

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, fieldsUri, fieldID, request)
		fieldsServiceMock.Mock.On("Delete",
			context.Request().Context(),
			fieldID).Return(expectedError).Once()

		handler := fields.NewFieldsHandler(configs, fieldsServiceMock, logger)

		handler.Delete(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedError.Error())
		fieldsServiceMock.AssertExpectations(t)
	})
}

func interfaceToFields(in interface{}) []entities.Field {
	interfaceArray := in.([]interface{})
	var fieldItem entities.Field
	var fieldJson []byte
	fields := make([]entities.Field, 0)

	for _, operatorMap := range interfaceArray {
		fieldJson, _ = json.Marshal(operatorMap)
		json.Unmarshal(fieldJson, &fieldItem)
		fields = append(fields, fieldItem)
	}
	return fields
}
