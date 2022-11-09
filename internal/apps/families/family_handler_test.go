package families_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/families"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/echo"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const familiesUri = "/risk-rules/v1/families/"

func TestFamilyHandler_Create(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when request is invalid, then return BadRequest", func(t *testing.T) {
		expectedError := "Syntax error"
		badRequest := `{
    		"name": "Family Name",
    		"mccs": [
				"1111",
				"2222",
				"3333"
			],
			"author": "carlos.maldonado@conekta.com",
		}`
		context, rec := echo.SetupAsRecorder(http.MethodPost, familiesUri, "", badRequest)
		handler := families.NewFamilyHandler(nil, logger)

		handler.Create(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when request author is not send, then return BadRequest", func(t *testing.T) {
		expectedError := "Key: 'FamilyRequest.Author' Error:Field validation for 'Author' failed on the 'required' tag"
		badRequest := `{
    		"name": "Family Name",
    		"mccs": [
				"1111",
				"2222",
				"3333"
			]
		}`

		context, rec := echo.SetupAsRecorder(http.MethodPost, familiesUri, "", badRequest)
		handler := families.NewFamilyHandler(nil, logger)

		handler.Create(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when request mcc is not a number, then return BadRequest", func(t *testing.T) {
		familyRequest := testdata.GetFamilyRequestWithMccCharacters()
		expectedError := fmt.Sprintf("family mcc [%s] is not number value", familyRequest.Mccs[0])
		bodyBytes, _ := json.Marshal(familyRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familiesUri, "", string(bodyBytes))
		handler := families.NewFamilyHandler(nil, logger)

		handler.Create(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when the length of the mcc is different from 4 positions, then return BadRequest", func(t *testing.T) {
		familyRequest := testdata.GetFamilyRequestWithMccLengthIsDifferentForm4()
		expectedError := fmt.Sprintf("family mcc [%s] length must be 4 positions", familyRequest.Mccs[0])
		bodyBytes, _ := json.Marshal(familyRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familiesUri, "", string(bodyBytes))
		handler := families.NewFamilyHandler(nil, logger)

		handler.Create(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when service fails, then return error", func(t *testing.T) {
		familyServiceMock := new(mocks.FamilyServiceMock)
		expectedError := errors.New("connection lost")
		familyRequest := entities.FamilyRequest{
			Name:   "Family Name",
			Mccs:   []string{"1111", "2222", "3333"},
			Author: "santiago.ceron@conekta.com",
		}
		bodyBytes, _ := json.Marshal(familyRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familiesUri, "", string(bodyBytes))
		familyServiceMock.Mock.On("Create",
			context.Request().Context(),
			mock.AnythingOfType("entities.Family")).Return(expectedError).Once()
		handler := families.NewFamilyHandler(familyServiceMock, logger)

		handler.Create(context)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		familyServiceMock.AssertExpectations(t)
	})

	t.Run("when family create proccesed ok, then return created", func(t *testing.T) {
		familyServiceMock := new(mocks.FamilyServiceMock)
		familyRequest := entities.FamilyRequest{
			Name:   "Family Name",
			Mccs:   []string{"1111", "2222", "3333"},
			Author: "santiago.ceron@conekta.com",
		}
		bodyBytes, _ := json.Marshal(familyRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familiesUri, "", string(bodyBytes))
		familyServiceMock.Mock.On("Create",
			context.Request().Context(),
			mock.AnythingOfType("entities.Family")).
			Return(nil).Once()
		handler := families.NewFamilyHandler(familyServiceMock, logger)

		err := handler.Create(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		familyServiceMock.AssertExpectations(t)
	})
}

func TestFamilyHandler_Get(t *testing.T) {
	logger, _ := logs.New()

	t.Run("get all families with default pagination", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetFamilies(),
		}

		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		filter := entities.FamilyFilter{ID: ""}
		serviceMock := new(mocks.FamilyServiceMock)

		handler := families.NewFamilyHandler(serviceMock, logger)

		c, rec := echo.SetupAsRecorder(http.MethodGet, familiesUri, "?"+q.Encode(), "")
		serviceMock.On("Get", c.Request().Context(), entities.NewDefaultPagination(), filter).
			Return(serviceResponse, nil)

		handler.Get(c)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when id param is not valid", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    []entities.Family{},
		}

		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		q.Set("id", "613688c2a95f286d57047d7dx")

		serviceMock := new(mocks.FamilyServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, familiesUri, "?"+q.Encode(), "")
		handler := families.NewFamilyHandler(serviceMock, logger)
		serviceMock.On("Get", context.Request().Context(), entities.NewDefaultPagination(),
			entities.FamilyFilter{
				ID: "613688c2a95f286d57047d7dx",
			}).
			Return(entities.PagedResponse{}, nil)

		handler.Get(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, serviceResponse.Data, []entities.Family{})
	})

	t.Run("when an error occurs in the service return nil", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetFamilies(),
		}

		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		filter := entities.FamilyFilter{ID: ""}
		serviceMock := new(mocks.FamilyServiceMock)

		expectedError := errors.New("service error")

		handler := families.NewFamilyHandler(serviceMock, logger)

		c, rec := echo.SetupAsRecorder(http.MethodGet, familiesUri, "?"+q.Encode(), "")
		serviceMock.On("Get", c.Request().Context(), entities.NewDefaultPagination(), filter).
			Return(serviceResponse, expectedError)

		err := handler.Get(c)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestFamilyHandler_Update(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when id of family is empty, then return BadRequest", func(t *testing.T) {
		expectedError := errors.New("id cannot be empty to update a family")
		request := testdata.GetFamilyUpdateRequest()
		body, _ := json.Marshal(request)

		familyServiceMock := mocks.NewFamilyServiceMock()
		handler := families.NewFamilyHandler(&familyServiceMock, logger)
		ctx, rec := echo.SetupAsRecorder(http.MethodPut, familiesUri, "", string(body))

		handler.Update(ctx)
		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, restError.Error(), expectedError.Error())
	})

	t.Run("when family json request is malformed, then return BadRequest", func(t *testing.T) {
		familyServiceMock := new(mocks.FamilyServiceMock)
		uriWithID := familiesUri

		req := testdata.GetJsonRequestIsMalformed()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, "611709bb70cbe3606baa3f8d", string(request))

		handler := families.NewFamilyHandler(familyServiceMock, logger)
		handler.Update(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		familyServiceMock.AssertExpectations(t)
	})

	t.Run("when family request author not send, then return BadRequest", func(t *testing.T) {
		familyServiceMock := new(mocks.FamilyServiceMock)
		uriWithID := familiesUri
		id := "611709bb70cbe3606baa3f8d"
		famiyRequest := `{
    		"name": "Family Name",
    		"mccs": [
				"1111",
				"2222",
				"3333"
			]
		}`

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, id, famiyRequest)

		handler := families.NewFamilyHandler(familyServiceMock, logger)
		handler.Update(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		familyServiceMock.AssertExpectations(t)
	})

	t.Run("when request mcc is not a number, then return BadRequest", func(t *testing.T) {
		familyServiceMock := new(mocks.FamilyServiceMock)
		uriWithID := familiesUri

		familyRequest := testdata.GetFamilyRequestWithMccCharacters()
		expectedError := fmt.Sprintf("family mcc [%s] is not number value", familyRequest.Mccs[0])
		request, _ := json.Marshal(familyRequest)
		id := "611709bb70cbe3606baa3f8d"

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, id, string(request))
		familyServiceMock.On("Update", context.Request().Context(), id, familyRequest).
			Return(nil)

		handler := families.NewFamilyHandler(familyServiceMock, logger)
		handler.Update(context)

		restError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when the length of the mcc is different from 4 positions, then return BadRequest", func(t *testing.T) {
		familyServiceMock := new(mocks.FamilyServiceMock)
		uriWithID := familiesUri

		familyRequest := testdata.GetFamilyRequestWithMccLengthIsDifferentForm4()
		expectedError := fmt.Sprintf("family mcc [%s] length must be 4 positions", familyRequest.Mccs[0])
		request, _ := json.Marshal(familyRequest)
		id := "611709bb70cbe3606baa3f8d"

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, id, string(request))
		familyServiceMock.On("Update", context.Request().Context(), id, familyRequest).
			Return(nil)

		handler := families.NewFamilyHandler(familyServiceMock, logger)
		handler.Update(context)

		restError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when service fails, then return error", func(t *testing.T) {
		familyServiceMock := new(mocks.FamilyServiceMock)
		expectedError := errors.New("Service error")
		family := testdata.GetFamilyRequest()
		request, _ := json.Marshal(family)

		context, _ := echo.SetupAsRecorder(http.MethodPost, familiesUri, "611709bb70cbe3606baa3f8d", string(request))
		familyServiceMock.On("Update", context.Request().Context(), "611709bb70cbe3606baa3f8d",
			mock.AnythingOfType("entities.Family")).
			Return(expectedError)

		handler := families.NewFamilyHandler(familyServiceMock, logger)
		err := handler.Update(context)

		assert.NotNil(t, expectedError, err)
	})

	t.Run("when family update is success", func(t *testing.T) {
		familyID := "611709bb70cbe3606baa3f8d"

		familyServiceMock := new(mocks.FamilyServiceMock)
		uriWithID := familiesUri + "/"

		req := testdata.GetFamilyRequest()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, familyID, string(request))
		familyServiceMock.On("Update", context.Request().Context(), familyID,
			mock.AnythingOfType("entities.Family")).
			Return(nil)

		handler := families.NewFamilyHandler(familyServiceMock, logger)
		err := handler.Update(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})
}

func Test_FamilyHandler_Delete(t *testing.T) {
	logger, _ := logs.New()
	familyID := "611709bb70cbe3606baa3f8d"

	t.Run("delete family when id is empty", func(t *testing.T) {
		service := new(mocks.FamilyServiceMock)
		expectedError := errors.New("empty id")

		context, recorder := echo.SetupAsRecorder(http.MethodPut, familiesUri, "", "")
		service.On("Delete", context.Request().Context(), "").Return(expectedError)

		handler := families.NewFamilyHandler(service, logger)
		handler.Delete(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
	})

	t.Run("delete family service return error", func(t *testing.T) {
		service := new(mocks.FamilyServiceMock)
		expectedError := errors.New("Service error")

		context, rec := echo.SetupAsRecorder(http.MethodDelete, familiesUri, familyID, "")
		service.On("Delete", context.Request().Context(), familyID).Return(expectedError).Once()

		handler := families.NewFamilyHandler(service, logger)
		handler.Delete(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, restError.Error(), expectedError.Error())
	})

	t.Run("delete family successful", func(t *testing.T) {
		service := new(mocks.FamilyServiceMock)

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, familiesUri, familyID, "")
		service.On("Delete", context.Request().Context(), familyID).Return(nil)

		handler := families.NewFamilyHandler(service, logger)
		err := handler.Delete(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})
}
