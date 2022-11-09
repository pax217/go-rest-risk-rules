package rules_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/conekta/risk-rules/pkg/echo"

	customHttp "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const rulesUri = "/risk-rules/v1/rules"

func Test_ruleHandler_List(t *testing.T) {
	logger, _ := logs.New()
	companyID := "611c02f729b0258dfcc84cc8"
	ruleID := "611709bb70cbe3606baa3f8d"

	t.Run("when list service fails", func(t *testing.T) {
		expectedErr := errors.New("error: database connection lost")
		uriWithID := rulesUri + "?company_id=" + companyID + "&id=" + ruleID
		ruleFilter := entities.RuleFilter{ID: ruleID, CompanyID: companyID}
		rulesServiceMock := new(mocks.RuleServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, uriWithID, "", "")
		handler := rules.NewRulesHandler(config.Config{}, rulesServiceMock, logger)
		rulesServiceMock.On("ListRules", context.Request().Context(), ruleFilter,
			entities.NewDefaultPagination()).Return(entities.PagedResponse{}, expectedErr)

		handler.GetPaged(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedErr.Error())
	})

	t.Run("when list service is ok", func(t *testing.T) {
		uriWithID := rulesUri + "?company_id=" + companyID + "&id=" + ruleID
		serviceResponse := entities.PagedResponse{
			HasMore: false,
			Total:   1,
			Object:  "list",
			Data: []entities.Rule{
				{
					Rule: "amount > 8 and amount < 52",
				},
			},
		}
		ruleFilter := entities.RuleFilter{ID: ruleID, CompanyID: companyID}

		rulesServiceMock := new(mocks.RuleServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, uriWithID, "", "")
		handler := rules.NewRulesHandler(config.Config{}, rulesServiceMock, logger)
		rulesServiceMock.On("ListRules", context.Request().Context(), ruleFilter,
			entities.NewDefaultPagination()).Return(serviceResponse, nil)

		handler.GetPaged(context)

		var pagedResponse entities.PagedResponse

		json.Unmarshal(rec.Body.Bytes(), &pagedResponse)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.EqualValues(t, serviceResponse.Data.([]entities.Rule)[0].Rule, interfaceToRules(pagedResponse.Data)[0].Rule)
	})

	t.Run("when rule_id is invalid", func(t *testing.T) {
		invalidRuleID := "611c02f729b0258dfcc84cc8-x"
		expectedErr := errors.New("invalid id")
		uriWithID := rulesUri + "?id=" + invalidRuleID
		ruleFilter := entities.RuleFilter{ID: invalidRuleID, CompanyID: companyID}

		rulesServiceMock := new(mocks.RuleServiceMock)
		context, rec := echo.SetupAsRecorder(http.MethodGet, uriWithID, "", "")
		handler := rules.NewRulesHandler(config.Config{}, rulesServiceMock, logger)
		rulesServiceMock.On("ListRules", context.Request().Context(), ruleFilter,
			entities.NewDefaultPagination()).Return(entities.PagedResponse{}, expectedErr)

		handler.GetPaged(context)

		httpErrorResponse, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpErrorResponse.Status())
		assert.Contains(t, httpErrorResponse.Message(), expectedErr.Error())
	})
}

func Test_ruleHandler_AddRule(t *testing.T) {
	logger, _ := logs.New()

	t.Run("create rule global without companyID", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		var isGlobal = true

		rule := testdata.GetDefaultRuleRequestWithAmount()
		rule.IsGlobal = &isGlobal
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("create rule json malformed", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		request := testdata.GetJsonMalformed()

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", request)
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("create rule validation fail", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule := testdata.GetDefaultRuleRequestFailValidation()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("create rule service fail", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule, expectedError := testdata.GetDefaultRuleRequestReturnError()
		request, _ := json.Marshal(rule)

		context, _ := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, expectedError)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		err := handler.AddRule(context)

		assert.NotNil(t, expectedError, err)
	})

	t.Run("global case BadRequest -> non global && no value", func(t *testing.T) {
		var isGlobal = false
		const msgErr = "non global rule, one option: [family_mcc - company_id - family_companies] have to be passed"
		rule := testdata.GetDefaultRuleNotGlobalWithoutFamilyIDAndCompanyID()
		rule.IsGlobal = &isGlobal
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		handler := rules.NewRulesHandler(config.Config{}, nil, logger)
		handler.AddRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.Equal(t, msgErr, httpError.Message())
	})

	t.Run("family case BadRequest -> non global && family_id && company_id", func(t *testing.T) {
		const msgExp = "non global rule, only one option: [family_mcc - company_id - family_companies] could be set at same time"
		rule := testdata.GetDefaultRuleNotGlobalWithFamilyIDAndCompanyID()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		handler := rules.NewRulesHandler(config.Config{}, nil, logger)
		handler.AddRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.Equal(t, msgExp, httpError.Message())
	})

	t.Run("global case Ok-> global == false , family_id=nil, company_id=617ae8bc92d0e243227eee9d", func(t *testing.T) {
		rule := testdata.GetDefaultRuleNotGlobalWithCompanyID()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))

		ruleService := new(mocks.RuleServiceMock)
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, "", httpError.Message())
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("global case OK-> global == false , family_id=617ae8a4649e59500b7cd54d, company_id=nil", func(t *testing.T) {
		rule := testdata.GetDefaultRuleNotGlobalWithFamilyID()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))

		ruleService := new(mocks.RuleServiceMock)
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, "", httpError.Message())
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("global case BadRequest-> global == true , family_id=617ae8a4649e59500b7cd54d, company_id=nil", func(t *testing.T) {
		rule := testdata.GetDefaultRuleGlobalWithFamilyID()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		handler := rules.NewRulesHandler(config.Config{}, nil, logger)
		handler.AddRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.Equal(t, "family_mcc should not be passed, the rule is configured as Global", httpError.Message())
	})

	t.Run("global case BadRequest-> global == true , family_id=nil, company_id=617ae8bc92d0e243227eee9d", func(t *testing.T) {
		rule := testdata.GetDefaultRuleGlobalWithCompanyID()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		handler := rules.NewRulesHandler(config.Config{}, nil, logger)
		handler.AddRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.Equal(t, "company_id should not be passed, the rule is configured as Global", httpError.Message())
	})

	t.Run("global case OK-> global == true , family_id=nil, company_id=nil", func(t *testing.T) {
		rule := testdata.GetDefaultRuleGlobal()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))

		ruleService := new(mocks.RuleServiceMock)
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		err := handler.AddRule(context)

		assert.NoError(t, err)
		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, "", httpError.Message())
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("create rule with formula validation fail", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule := testdata.GetDefaultRuleRequestFailFormulaValidation()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("when creating rule with formula fields and rule field, it should fail", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule := testdata.GetRuleRequestWithFieldAndFormulaFields()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("decision value is not valid, then return BadRequest", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule := testdata.GetDefaultRuleRequestWithNotValueDecision()
		expectedErrMessage := fmt.Sprintf("decision value [%s], is not a valid value", rule.Decision)
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, expectedErrMessage, httpError.Message())
	})

	t.Run("decision value valid, then return StatusOK", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule := testdata.GetDefaultRuleRequestWithValidValueDecision()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("create rule with yellow flag success", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule := testdata.GetDefaultRuleRequestWithYellowFlag()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("create rule with yellow flag return bad request for incorrect decision", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		expectedError := errors.New("yellow flag rules must have the undecided decision")

		rule := testdata.GetDefaultRuleRequestWithYellowFlagIncorrectDecision()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "", string(request))
		ruleService.On("AddRule", context.Request().Context(), mock.Anything).Return(entities.Rule{}, nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.AddRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, expectedError.Error(), httpError.Message())
	})
}

func Test_ruleHandler_UpdateRule(t *testing.T) {
	configs := config.NewConfig()
	logger, _ := logs.New()

	t.Run("update rule successful", func(t *testing.T) {
		ruleID := "611709bb70cbe3606baa3f8d"

		ruleService := new(mocks.RuleServiceMock)
		uriWithID := rulesUri + "/"

		req := testdata.GetDefaultRuleRequestWithAmount()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, ruleID, string(request))
		ruleService.On("UpdateRule", context.Request().Context(), ruleID, mock.AnythingOfType("entities.Rule")).
			Return(nil)

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		err := handler.UpdateRule(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("update rule when id is empty", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		uriWithID := rulesUri

		req := testdata.GetDefaultRuleRequestWithAmount()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, "", string(request))
		ruleService.On("UpdateRule", context.Request().Context(), "", req.NewRuleFromPutRequest()).
			Return(errors.New("empty id"))

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		handler.UpdateRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
	})

	t.Run("update rule when json request is malformed", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		uriWithID := rulesUri

		req := testdata.GetJsonRequestIsMalformed()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, "611709bb70cbe3606baa3f8d", string(request))
		ruleService.On("UpdateRule", context.Request().Context(), "611709bb70cbe3606baa3f8d", req).
			Return(nil)

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		handler.UpdateRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
	})

	t.Run("update rule When Json Request Without Module", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		uriWithID := rulesUri

		req := testdata.GetJsonRequestWithoutModule()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, "611709bb70cbe3606baa3f8d", string(request))
		ruleService.On("UpdateRule", context.Request().Context(), "611709bb70cbe3606baa3f8d", req.NewRuleFromPutRequest()).
			Return(nil)

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		handler.UpdateRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
	})

	t.Run("update rule not global without companyid should be return error", func(t *testing.T) {
		var isGlobal = false
		const msgErr = "non global rule, one option: [family_mcc - company_id - family_companies] have to be passed"
		rule := testdata.GetDefaultRuleRequestWithAmount()
		rule.IsGlobal = &isGlobal
		request, _ := json.Marshal(rule)

		ruleService := new(mocks.RuleServiceMock)

		context, recorder := echo.SetupAsRecorder(
			http.MethodPut, "/risk-rules/v1/rules/", "611709bb70cbe3606baa3f8d", string(request))
		handler := rules.NewRulesHandler(configs, ruleService, logger)
		handler.UpdateRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
		assert.Equal(t, msgErr, httpError.Message())
	})

	t.Run("update rule service return error", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule, expectedError := testdata.GetDefaultRuleRequestServiceReturnError()
		request, _ := json.Marshal(rule)

		context, _ := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "611709bb70cbe3606baa3f8d", string(request))
		ruleService.On("UpdateRule", context.Request().Context(), "611709bb70cbe3606baa3f8d", mock.AnythingOfType("entities.Rule")).
			Return(expectedError)

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		err := handler.UpdateRule(context)

		assert.NotNil(t, expectedError, err)
	})

	t.Run("update rule when decision value is not valid, then return BadRequest", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		rule := testdata.GetDefaultRuleRequestWithNotValueDecision()
		expectedErrMessage := fmt.Sprintf("decision value [%s], is not a valid value", rule.Decision)
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPost, "/risk-rules/v1/rules", "611709bb70cbe3606baa3f8d", string(request))
		ruleService.On("UpdateRule", context.Request().Context(), "611709bb70cbe3606baa3f8d", mock.AnythingOfType("entities.Rule")).
			Return(expectedErrMessage)

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		handler.UpdateRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, expectedErrMessage, httpError.Message())
	})

	t.Run("update rule when decision value is valid, then return StatusOK", func(t *testing.T) {
		ruleID := "611709bb70cbe3606baa3f8d"

		ruleService := new(mocks.RuleServiceMock)
		uriWithID := rulesUri + "/"

		req := testdata.GetDefaultRuleRequestWithValidValueDecision()
		request, _ := json.Marshal(req)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, ruleID, string(request))
		ruleService.On("UpdateRule", context.Request().Context(), ruleID, mock.AnythingOfType("entities.Rule")).
			Return(nil)

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		err := handler.UpdateRule(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("update rule with yellow flag success", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		ruleID := "611709bb70cbe3606baa3f8d"

		rule := testdata.GetDefaultRuleRequestWithYellowFlag()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, "/", ruleID, string(request))
		ruleService.On("UpdateRule", context.Request().Context(), ruleID, mock.Anything).Return(nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.UpdateRule(context)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("create rule with yellow flag return bad request for incorrect decision", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		ruleID := "611709bb70cbe3606baa3f8d"
		expectedError := errors.New("yellow flag rules must have the undecided decision")

		rule := testdata.GetDefaultRuleRequestWithYellowFlagIncorrectDecision()
		request, _ := json.Marshal(rule)

		context, recorder := echo.SetupAsRecorder(http.MethodPut, "/", ruleID, string(request))
		ruleService.On("UpdateRule", context.Request().Context(), ruleID, mock.Anything).Return(nil)

		handler := rules.NewRulesHandler(config.Config{}, ruleService, logger)
		handler.UpdateRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, expectedError.Error(), httpError.Message())
	})
}

func Test_ruleHandler_DeleteRule(t *testing.T) {
	configs := config.NewConfig()
	logger, _ := logs.New()

	t.Run("delete rule when id is empty", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		uriWithID := rulesUri
		expectedError := errors.New("empty id")

		context, recorder := echo.SetupAsRecorder(http.MethodPut, uriWithID, "", "")
		ruleService.On("RemoveRule", context.Request().Context(), "").Return(expectedError)

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		handler.RemoveRule(context)

		httpError, _ := customHttp.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, httpError.Status())
	})

	t.Run("delete rule service return error", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)
		expectedError := errors.New("Service error")

		context, rec := echo.SetupAsRecorder(http.MethodDelete, "/risk-rules/v1/rules/", "611709bb70cbe3606baa3f8d", "")
		ruleService.On("RemoveRule", context.Request().Context(), "611709bb70cbe3606baa3f8d").Return(expectedError).Once()

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		handler.RemoveRule(context)

		restError, _ := customHttp.NewRestErrorFromBytes(rec.Body.Bytes())

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, restError.Error(), expectedError.Error())
	})

	t.Run("delete rule successful", func(t *testing.T) {
		ruleService := new(mocks.RuleServiceMock)

		context, recorder := echo.SetupAsRecorder(http.MethodDelete, "/risk-rules/v1/rules/", "611709bb70cbe3606baa3f8d", "")
		ruleService.On("RemoveRule", context.Request().Context(), "611709bb70cbe3606baa3f8d").Return(nil)

		handler := rules.NewRulesHandler(configs, ruleService, logger)
		err := handler.RemoveRule(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})
}

func interfaceToRules(in interface{}) []entities.Rule {
	interfaceArray := in.([]interface{})
	var ruleItem entities.Rule
	var ruleJson []byte
	rules := make([]entities.Rule, 0)

	for _, ruleMap := range interfaceArray {
		ruleJson, _ = json.Marshal(ruleMap)
		json.Unmarshal(ruleJson, &ruleItem)
		rules = append(rules, ruleItem)
	}
	return rules
}
