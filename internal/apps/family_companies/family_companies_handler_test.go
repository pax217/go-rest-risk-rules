package familycom_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/family_companies"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/echo"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const familyCompaniesUri = "/risk-rules/v1/family_companies"

func TestFamilyCompaniesHandler_Create(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when request is invalid, then return BadRequest", func(t *testing.T) {
		expectedError := "Syntax error"
		badRequest := `{
    		"name": "Family Name",
    		"company_ids": [
				"61e4dd6da5997ad4d9e76945",
				"61e4dd7320fbfc5f0849fba5",
			],
			"author": "carlos.maldonado@conekta.com",
		}`
		context, rec := echo.SetupAsRecorder(http.MethodPost, familyCompaniesUri, "", badRequest)
		handler := familycom.NewFamilyCompaniesHandler(nil, logger)

		handler.Create(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when request author isn not send, then return BadRequest", func(t *testing.T) {
		expectedError := "Key: 'FamilyCompaniesRequest.Author' Error:Field validation for 'Author' failed on the 'required' tag"
		familyCompaniesRequest := testdata.GetFamilyCompaniesRequestWithOutAuthor()
		bodyBytes, _ := json.Marshal(familyCompaniesRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familyCompaniesUri, "", string(bodyBytes))
		handler := familycom.NewFamilyCompaniesHandler(nil, logger)

		handler.Create(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when request company_ids is not id valid, then return BadRequest", func(t *testing.T) {
		familyCompaniesRequest := testdata.GetFamilyCompaniesRequestWithBadCompanyId()
		expectedError := fmt.Sprintf(
			"company id [%s] is not a valid format of type mongo id",
			familyCompaniesRequest.CompanyIDs[0])
		bodyBytes, _ := json.Marshal(familyCompaniesRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familyCompaniesUri, "", string(bodyBytes))
		handler := familycom.NewFamilyCompaniesHandler(nil, logger)

		handler.Create(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when request company_ids valid mongo id, then return created", func(t *testing.T) {
		familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
		familyCompaniesRequest := entities.FamilyCompaniesRequest{
			Name:       "Family Companies Name",
			CompanyIDs: []string{"62046bf28e99e83b0554fbae", "62046bf7ef5fd867cc7d8c2f"},
			Author:     "santiago.ceron@conekta.com",
		}
		bodyBytes, _ := json.Marshal(familyCompaniesRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familyCompaniesUri, "", string(bodyBytes))
		familyCompaniesServiceMock.Mock.On("Create",
			context.Request().Context(),
			mock.AnythingOfType("entities.FamilyCompanies")).
			Return(nil).Once()
		handler := familycom.NewFamilyCompaniesHandler(familyCompaniesServiceMock, logger)

		err := handler.Create(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		familyCompaniesServiceMock.AssertExpectations(t)
	})

	t.Run("when request company_ids has an not valid mongo id, then return BadRequest", func(t *testing.T) {
		familyCompaniesRequest := entities.FamilyCompaniesRequest{
			Name:       "Family Companies Name",
			CompanyIDs: []string{"62046bf28e99e83b0554fbae", "7112821neuehud121eduihuq"},
			Author:     "santiago.ceron@conekta.com",
		}
		expectedError := fmt.Sprintf(
			"company id [%s] is not a valid format of type mongo id",
			familyCompaniesRequest.CompanyIDs[1])
		bodyBytes, _ := json.Marshal(familyCompaniesRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familyCompaniesUri, "", string(bodyBytes))
		handler := familycom.NewFamilyCompaniesHandler(nil, logger)

		handler.Create(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when service fails, then return error", func(t *testing.T) {
		familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
		expectedError := errors.New("connection lost")
		familyCompaniesRequest := entities.FamilyCompaniesRequest{
			Name:       "Family Companies Name",
			CompanyIDs: []string{"61e4dd6da5997ad4d9e76945", "61e4dd7320fbfc5f0849fba5"},
			Author:     "santiago.ceron@conekta.com",
		}
		bodyBytes, _ := json.Marshal(familyCompaniesRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familyCompaniesUri, "", string(bodyBytes))
		familyCompaniesServiceMock.Mock.On("Create",
			context.Request().Context(),
			mock.AnythingOfType("entities.FamilyCompanies")).
			Return(expectedError).Once()
		handler := familycom.NewFamilyCompaniesHandler(familyCompaniesServiceMock, logger)

		handler.Create(context)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		familyCompaniesServiceMock.AssertExpectations(t)
	})

	t.Run("when family companies create proccesed ok, then return created", func(t *testing.T) {
		familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
		familyCompaniesRequest := entities.FamilyCompaniesRequest{
			Name:       "Family Companies Name",
			CompanyIDs: []string{"61e4dd6da5997ad4d9e76945", "61e4dd7320fbfc5f0849fba5"},
			Author:     "santiago.ceron@conekta.com",
		}
		bodyBytes, _ := json.Marshal(familyCompaniesRequest)

		context, rec := echo.SetupAsRecorder(http.MethodPost, familyCompaniesUri, "", string(bodyBytes))
		familyCompaniesServiceMock.Mock.On("Create",
			context.Request().Context(),
			mock.AnythingOfType("entities.FamilyCompanies")).
			Return(nil).Once()
		handler := familycom.NewFamilyCompaniesHandler(familyCompaniesServiceMock, logger)

		err := handler.Create(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		familyCompaniesServiceMock.AssertExpectations(t)
	})
}

func TestFamilyCompaniesHandler_Get(t *testing.T) {
	logger, _ := logs.New()

	t.Run("get all family companies with default pagination", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetFamilyCompanies(),
		}

		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		filter := entities.FamilyCompaniesFilter{}
		serviceMock := new(mocks.FamilyCompaniesServiceMock)

		handler := familycom.NewFamilyCompaniesHandler(serviceMock, logger)

		c, rec := echo.SetupAsRecorder(http.MethodGet, familyCompaniesUri, "?"+q.Encode(), "")
		serviceMock.On("Get",
			c.Request().Context(),
			entities.NewDefaultPagination(), filter).
			Return(serviceResponse, nil)

		handler.Get(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		serviceMock.AssertExpectations(t)
	})

	t.Run("when id param is not valid", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    []entities.FamilyCompanies{},
		}
		id := "61e9b9414e569c8bcbdc408b-x"

		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		q.Set("id", id)

		familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, familyCompaniesUri, "?"+q.Encode(), "")
		handler := familycom.NewFamilyCompaniesHandler(familyCompaniesServiceMock, logger)
		familyCompaniesServiceMock.On("Get", context.Request().Context(), entities.NewDefaultPagination(),
			entities.FamilyCompaniesFilter{
				ID: id,
			}).
			Return(entities.PagedResponse{}, nil)

		handler.Get(context)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, serviceResponse.Data)
		familyCompaniesServiceMock.AssertExpectations(t)
	})

	t.Run("when an error occurs in the service return nil", func(t *testing.T) {
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data:    testdata.GetFamilyCompanies(),
		}

		q := make(url.Values)
		q.Set("page", "1")
		q.Set("size", "25")
		filter := entities.FamilyCompaniesFilter{ID: ""}
		serviceMock := new(mocks.FamilyCompaniesServiceMock)

		expectedError := errors.New("service error")

		handler := familycom.NewFamilyCompaniesHandler(serviceMock, logger)

		c, rec := echo.SetupAsRecorder(http.MethodGet, familyCompaniesUri, "?"+q.Encode(), "")
		serviceMock.On("Get", c.Request().Context(), entities.NewDefaultPagination(), filter).
			Return(serviceResponse, expectedError)

		err := handler.Get(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		serviceMock.AssertExpectations(t)
	})
}

func TestFamilyCompaniesHandler_Update(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when id family companies is empty, then return BadRequest", func(t *testing.T) {
		expectedError := errors.New("id cannot be empty to update a family companies")
		request := testdata.GetFamilyCompaniesUpdateRequest()
		body, _ := json.Marshal(request)

		handler := familycom.NewFamilyCompaniesHandler(nil, logger)
		ctx, rec := echo.SetupAsRecorder(http.MethodPut, familyCompaniesUri, "", string(body))

		handler.Update(ctx)
		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, restError.Error(), expectedError.Error())
	})

	t.Run("when family companies json request is malformed, then return BadRequest", func(t *testing.T) {

		uriWithID := familyCompaniesUri

		req := testdata.GetFamilyCompaniesJsonRequestIsMalformed()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, "61eefa6d92eda66f6fb489cd", string(request))

		handler := familycom.NewFamilyCompaniesHandler(nil, logger)
		handler.Update(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
	})

	t.Run("when family companies request author is not send, then return BadRequest", func(t *testing.T) {
		uriWithID := familyCompaniesUri
		id := "611709bb70cbe3606baa3f8d"
		famiyRequest := `{
			"name": "Tiendas Electrónica",
			"company_ids": ["61e4dd6da5997ad4d9e76945","61eb21792e341c54221062b4","61eb217e66524deb95ad5143"]
		}`

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, id, famiyRequest)

		handler := familycom.NewFamilyCompaniesHandler(nil, logger)
		handler.Update(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
	})

	t.Run("when request company_ids is not a number, then return BadRequest", func(t *testing.T) {
		uriWithID := familyCompaniesUri

		familyRequest := testdata.GetFamilyCompaniesRequestWithCompanyIdsNotValid()
		expectedError := fmt.Sprintf(
			"company id [%s] is not a valid format of type mongo id",
			familyRequest.CompanyIDs[0])
		request, _ := json.Marshal(familyRequest)
		id := "611709bb70cbe3606baa3f8d"

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, id, string(request))

		handler := familycom.NewFamilyCompaniesHandler(nil, logger)
		handler.Update(context)

		restError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when service fails, then return error", func(t *testing.T) {
		familyServiceMock := new(mocks.FamilyCompaniesServiceMock)
		expectedError := errors.New("Service error")
		familyCompanies := testdata.GetFamilyCompaniesRequest()
		request, _ := json.Marshal(familyCompanies)

		context, _ := echo.SetupAsRecorder(http.MethodPost, familyCompaniesUri, "611709bb70cbe3606baa3f8d", string(request))
		familyServiceMock.On("Update", context.Request().Context(), "611709bb70cbe3606baa3f8d",
			mock.AnythingOfType("entities.FamilyCompanies")).
			Return(expectedError)

		handler := familycom.NewFamilyCompaniesHandler(familyServiceMock, logger)
		err := handler.Update(context)

		assert.NotNil(t, expectedError, err)
		familyServiceMock.AssertExpectations(t)
	})

	t.Run("when family companies update is success", func(t *testing.T) {
		familyID := "611709bb70cbe3606baa3f8d"

		familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
		uriWithID := familyCompaniesUri + "/"

		req := testdata.GetFamilyCompaniesRequest()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, familyID, string(request))
		familyCompaniesServiceMock.On("Update", context.Request().Context(), familyID,
			mock.AnythingOfType("entities.FamilyCompanies")).
			Return(nil)

		handler := familycom.NewFamilyCompaniesHandler(familyCompaniesServiceMock, logger)
		err := handler.Update(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
		familyCompaniesServiceMock.AssertExpectations(t)
	})
}

func TestFamilyCompaniesHandler_Delete(t *testing.T) {
	logger, _ := logs.New()
	familyID := "61e990e16d290da842dfbc62"

	t.Run("delete family companies when id is empty", func(t *testing.T) {
		familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
		expectedError := errors.New("empty id")

		context, recorder := echo.SetupAsRecorder(http.MethodPut, familyCompaniesUri, "", "")
		familyCompaniesServiceMock.On("Delete", context.Request().Context(), "").Return(expectedError)

		handler := familycom.NewFamilyCompaniesHandler(familyCompaniesServiceMock, logger)
		handler.Delete(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
	})

	t.Run("delete family companies service return error", func(t *testing.T) {
		familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
		expectedError := errors.New("Service error")

		context, rec := echo.SetupAsRecorder(http.MethodDelete, familyCompaniesUri, familyID, "")
		familyCompaniesServiceMock.On("Delete", context.Request().Context(), familyID).Return(expectedError).Once()

		handler := familycom.NewFamilyCompaniesHandler(familyCompaniesServiceMock, logger)
		handler.Delete(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, restError.Error(), expectedError.Error())
	})

	t.Run("delete family companies successful", func(t *testing.T) {
		familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, familyCompaniesUri, familyID, "")
		familyCompaniesServiceMock.On("Delete", context.Request().Context(), familyID).Return(nil)

		handler := familycom.NewFamilyCompaniesHandler(familyCompaniesServiceMock, logger)
		err := handler.Delete(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})
}
