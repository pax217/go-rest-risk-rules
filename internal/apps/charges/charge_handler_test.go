package charges_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/charges"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/echo"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

var (
	metrics = new(datadog.MetricsDogMock)
)

func TestChargeHandler_Get(t *testing.T) {
	logger, _ := logs.New()

	t.Run("no id is sent", func(t *testing.T) {
		expectedErrorMsg := "empty id"

		context, rec := echo.SetupAsRecorder(http.MethodGet, "/charges", "", "")
		handler := charges.NewChargeHandler(config.Config{}, nil, logger, metrics)

		err := handler.GetEvaluation(context)

		assert.NoError(t, err)
		httpError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.Equal(t, expectedErrorMsg, httpError.Message())
	})

	t.Run("service fails", func(t *testing.T) {
		expectedError := errors.New("datbase connection lost")
		id := "abc-123"

		context, rec := echo.SetupAsRecorder(http.MethodGet, "/charges", id, "")
		service := new(mocks.ChargeServiceMock)
		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)
		service.On("Get", context.Request().Context(), id).Return(
			entities.EvaluationResponse{}, expectedError)

		err := handler.GetEvaluation(context)

		assert.NoError(t, err)
		httpError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpError.Status())
		assert.Equal(t, expectedError.Error(), httpError.Message())
	})

	t.Run("service works", func(t *testing.T) {
		id := "abc-123"
		decision := "A"
		result := entities.EvaluationResponse{
			Decision: decision,
			Modules:  entities.ModulesResponse{},
			Charge:   entities.ChargeRequest{},
		}

		context, rec := echo.SetupAsRecorder(http.MethodGet, "/charges", id, "")
		service := new(mocks.ChargeServiceMock)
		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)
		service.On("Get", context.Request().Context(), id).Return(result, nil)

		err := handler.GetEvaluation(context)

		assert.NoError(t, err)
		var response entities.EvaluationResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, decision, response.Decision)
	})
}

func TestChargeHandler_Evaluate(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when the evaluation is ok", func(t *testing.T) {
		charge := testdata.GetDefaultCharge()
		request, _ := json.Marshal(charge)

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate", "", string(request))
		response := testdata.GetEvaluationResponseSuccessful()
		service.On("EvaluateCharge", context.Request().Context(), charge).Return(response, nil)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.Evaluate(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("charge has aggregation fields and evaluation is ok", func(t *testing.T) {
		request := testdata.GetDefaultCharge()
		body, _ := json.Marshal(request)

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate", "", string(body))
		response := testdata.GetEvaluationResponseSuccessful()
		service.On("EvaluateCharge", context.Request().Context(),
			mock.AnythingOfType("entities.ChargeRequest")).Return(response, nil)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.Evaluate(context)

		var evaluation entities.EvaluationResponse
		json.Unmarshal(rec.Body.Bytes(), &evaluation)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, request.Aggregation.BinNumber.Charge.H1, evaluation.Charge.Aggregation.BinNumber.Charge.H1)
		assert.Equal(t, request.Aggregation.BinNumber.Charge.H2, evaluation.Charge.Aggregation.BinNumber.Charge.H2)
		assert.Equal(t, request.Aggregation.BinNumber.Charge.H12, evaluation.Charge.Aggregation.BinNumber.Charge.H12)
	})

	t.Run("charge without aggregation fields and evaluation is ok", func(t *testing.T) {
		request := testdata.GetChargeWithoutAggregationRequest()

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate", "", request)
		response := testdata.GetEvaluationResponseSuccessfulWithoutAggregation()
		service.On("EvaluateCharge", context.Request().Context(),
			mock.AnythingOfType("entities.ChargeRequest")).Return(response, nil)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.Evaluate(context)

		var evaluation entities.EvaluationResponse
		json.Unmarshal(rec.Body.Bytes(), &evaluation)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, entities.AggregationAttribute{}, evaluation.Charge.Aggregation.BinNumber)
	})

	t.Run("when request marshal ok", func(t *testing.T) {
		request := testdata.GetChargeAggregationRequest()

		body, err := json.Marshal(request)

		assert.Nil(t, err)
		assert.NotNil(t, body)
	})

	t.Run("when request is malformed", func(t *testing.T) {
		request := testdata.GetChargeMalformedRequest()

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate", "", request)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.Evaluate(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when validation fail", func(t *testing.T) {
		request := testdata.GetChargeInvalidRequest()

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate", "", request)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.Evaluate(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when service return error", func(t *testing.T) {
		charge := testdata.GetDefaultCharge()
		request, _ := json.Marshal(charge)
		expectedError := errors.New("service error")
		metrics := new(datadog.MetricsDogMock)

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate", "", string(request))
		service.On("EvaluateCharge", context.Request().Context(), charge).
			Return(entities.EvaluationResponse{}, expectedError)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.Evaluate(context)

		assert.NoError(t, err)
		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, restError.Error(), expectedError.Error())
	})

	t.Run("when company_mcc in request is not sent, then return BadRequest", func(t *testing.T) {
		expectedError := "Key: 'ChargeRequest.CompanyMCC' Error:Field validation for 'CompanyMCC' failed on the 'required' tag"
		charge := testdata.GetChargeWithoutCompanyMccInRequest()
		request, _ := json.Marshal(charge)

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate", "", string(request))
		response := testdata.GetEvaluationResponseSuccessful()
		service.On("EvaluateCharge", context.Request().Context(), charge).Return(response, nil)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		handler.Evaluate(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(restError.Message(), expectedError))
	})

	t.Run("when the evaluation is ok, validate structure charge -> aggregation -> card hash", func(t *testing.T) {
		charge := testdata.GetDefaultCharge()
		request, _ := json.Marshal(charge)

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate", "", string(request))
		response := testdata.GetEvaluationResponseSuccessful()
		service.On("EvaluateCharge", context.Request().Context(), charge).Return(response, nil)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.Evaluate(context)

		var evaluationResponse entities.EvaluationResponse

		json.Unmarshal(rec.Body.Bytes(), &evaluationResponse)

		assert.Nil(t, err)
		assert.NotNil(t, evaluationResponse.Charge.Aggregation)
		assert.Equal(t, charge.Aggregation.CardHash.Charge.H1, evaluationResponse.Charge.Aggregation.CardHash.Charge.H1)
		assert.Equal(t, charge.Aggregation.CardHash.Charge.H2, evaluationResponse.Charge.Aggregation.CardHash.Charge.H2)
		assert.Equal(t, charge.Aggregation.CardHash.Charge.H12, evaluationResponse.Charge.Aggregation.CardHash.Charge.H12)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

}

func TestChargeHandler_EvaluateOnlyRules(t *testing.T) {
	logger, _ := logs.New()

	t.Run("when charge request is malformed, the return StatusBadRequest",
		func(t *testing.T) {
			chargeRequest := testdata.GetChargeMalformedRequest()
			request, _ := json.Marshal(chargeRequest)

			context, rec := echo.SetupAsRecorder(http.MethodPost,
				"/charges/evaluate_only_rules",
				"",
				string(request))
			handler := charges.NewChargeHandler(config.Config{}, nil, logger, metrics)

			err := handler.EvaluateOnlyRules(context)

			var pagedResponse entities.PagedResponse

			json.Unmarshal(rec.Body.Bytes(), &pagedResponse)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})

	t.Run("when charge request is not valid, the return StatusBadRequest",
		func(t *testing.T) {
			chargeRequest := testdata.GetChargeInvalidRequestCompanyIDRequired()
			request, _ := json.Marshal(chargeRequest)

			context, rec := echo.SetupAsRecorder(http.MethodPost,
				"/charges/evaluate_only_rules",
				"",
				string(request))
			handler := charges.NewChargeHandler(config.Config{}, nil, logger, metrics)

			err := handler.EvaluateOnlyRules(context)

			var pagedResponse entities.PagedResponse

			json.Unmarshal(rec.Body.Bytes(), &pagedResponse)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})

	t.Run("when service fails, then return error", func(t *testing.T) {
		charge := testdata.GetChargeConsoleCompanyRules()
		request, _ := json.Marshal(charge)
		expectedError := errors.New("service error")
		metrics = new(datadog.MetricsDogMock)

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost,
			"/charges/evaluate_only_rules",
			"", string(request))
		service.On("EvaluateChargeOnlyRules",
			context.Request().Context(),
			charge).
			Return(entities.RulesEvaluationResponse{}, expectedError)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.EvaluateOnlyRules(context)

		assert.NoError(t, err)
		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, restError.Error(), expectedError.Error())
	})

	t.Run("when the evaluation is ok", func(t *testing.T) {
		charge := testdata.GetChargeConsoleCompanyRules()
		request, _ := json.Marshal(charge)

		service := new(mocks.ChargeServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodPost, "/charges/evaluate_only_rules", "", string(request))

		service.On("EvaluateChargeOnlyRules", context.Request().Context(), charge).
			Return(entities.RulesEvaluationResponse{}, nil)

		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)

		err := handler.EvaluateOnlyRules(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestChargeEvaluationMongoDBRepository_GetOnlyRules(t *testing.T) {
	logger, _ := logs.New()

	t.Run("no id is sent", func(t *testing.T) {
		expectedErrorMsg := "empty id"

		context, rec := echo.SetupAsRecorder(http.MethodGet, "/charges", "", "")
		handler := charges.NewChargeHandler(config.Config{}, nil, logger, metrics)

		err := handler.GetEvaluationOnlyRules(context)

		assert.NoError(t, err)
		httpError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.Equal(t, expectedErrorMsg, httpError.Message())
	})

	t.Run("service fails", func(t *testing.T) {
		expectedError := errors.New("datbase connection lost")
		id := "abc-123"

		context, rec := echo.SetupAsRecorder(http.MethodGet, "/charges", id, "")
		service := new(mocks.ChargeServiceMock)
		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)
		service.On("GetOnlyRules", context.Request().Context(), id).Return(
			entities.RulesEvaluationResponse{}, expectedError)

		err := handler.GetEvaluationOnlyRules(context)

		assert.NoError(t, err)
		httpError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpError.Status())
		assert.Equal(t, expectedError.Error(), httpError.Message())
	})

	t.Run("service works", func(t *testing.T) {
		id := "abc-123"
		decision := "A"
		result := testdata.GetRulesEvaluationResponseAcceptedCauseConsoleCompanyRules()

		context, rec := echo.SetupAsRecorder(http.MethodGet, "/charges", id, "")
		service := new(mocks.ChargeServiceMock)
		handler := charges.NewChargeHandler(config.Config{}, service, logger, metrics)
		service.On("GetOnlyRules", context.Request().Context(), id).Return(result, nil)

		err := handler.GetEvaluationOnlyRules(context)

		assert.NoError(t, err)
		var response entities.EvaluationResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, decision, response.Decision)
	})
}
