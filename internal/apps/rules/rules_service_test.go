package rules_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_ruleService_List(t *testing.T) {

	t.Run("test when repository passes directly data to service", func(t *testing.T) {
		pagination := entities.NewDefaultPagination()
		ruleFilter, _ := testdata.GetDefaultRuleFilter()

		notRulesFound := entities.PagedResponse{
			HasMore: false,
			Total:   0,
			Object:  "",
			Data:    nil,
		}
		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", nil, ruleFilter, pagination).Return(notRulesFound, nil)
		service := rules.NewRulesService(config.Config{}, nil, ruleRepository, nil, nil)

		pagedRules, err := service.ListRules(nil, ruleFilter, pagination)

		assert.Nil(t, err)
		assert.Equal(t, notRulesFound, pagedRules)
	})
}

func Test_ruleService_Build(t *testing.T) {
	t.Run("create tow numeric rules", func(t *testing.T) {
		rulesAsString := "amount > 8 and amount < 802.1"
		rulesContent := []entities.RuleContent{{
			Field:     "amount",
			Operator:  ">",
			Value:     "8",
			Condition: "and",
		}, {
			Field:     "amount",
			Operator:  "<",
			Value:     "802.1",
			Condition: "and",
		},
		}

		service := rules.NewRulesService(config.Config{}, nil, nil, nil, nil)
		ruleBuilt := service.BuildRule(rulesContent)

		assert.Equal(t, rulesAsString, ruleBuilt)
	})

	t.Run("create a string rule", func(t *testing.T) {
		ruleExpected := `company_id eq "60ad5c44926c8400016cbfdc"`
		rulesContent := []entities.RuleContent{{
			Field:     "company_id",
			Operator:  "eq",
			Value:     "60ad5c44926c8400016cbfdc",
			Condition: "and",
		},
		}

		service := rules.NewRulesService(config.Config{}, nil, nil, nil, nil)
		ruleBuilt := service.BuildRule(rulesContent)

		assert.Equal(t, ruleExpected, ruleBuilt)
	})

	t.Run("create a bool rule", func(t *testing.T) {
		ruleExpected := "live_mode eq true"
		rulesContent := []entities.RuleContent{{
			Field:     "live_mode",
			Operator:  "eq",
			Value:     "true",
			Condition: "and",
		},
		}

		service := rules.NewRulesService(config.Config{}, nil, nil, nil, nil)
		ruleBuilt := service.BuildRule(rulesContent)

		assert.Equal(t, ruleExpected, ruleBuilt)
	})

	t.Run("create a NOT IN  strings rule", func(t *testing.T) {
		rulesContent := []entities.RuleContent{
			{
				Field:     "payment_method.country",
				Operator:  "in",
				Value:     `["MX","US"]`,
				Condition: "and",
				Not:       true,
			},
			{
				Field:     "live_mode",
				Operator:  "eq",
				Value:     `true`,
				Condition: "and",
				Not:       false,
			},
		}
		ruleExpected := fmt.Sprintf("not payment_method.country in %s and live_mode eq true", rulesContent[0].Value)
		service := rules.NewRulesService(config.Config{}, nil, nil, nil, nil)
		ruleBuilt := service.BuildRule(rulesContent)

		assert.Equalf(t, ruleExpected, ruleBuilt, "The rule should be %s", ruleExpected)
	})
	t.Run("create a  IN  strings rule", func(t *testing.T) {
		rulesContent := []entities.RuleContent{{
			Field:     "payment_method.country",
			Operator:  "in",
			Value:     `["MX","US"]`,
			Condition: "and",
			Not:       false,
		},
		}
		ruleExpected := fmt.Sprintf("payment_method.country in %s", rulesContent[0].Value)
		service := rules.NewRulesService(config.Config{}, nil, nil, nil, nil)
		ruleBuilt := service.BuildRule(rulesContent)

		assert.Equalf(t, ruleExpected, ruleBuilt, "The rule should be %s", ruleExpected)
	})
}

func Test_ruleService_CreateRule(t *testing.T) {
	logger, _ := logs.New()
	rulesList := make([]entities.Rule, 0)
	rulesList = append(rulesList, testdata.GetDefaultRuleEmailBlockedGlobal(true))

	t.Run("test when validate rule fail", func(t *testing.T) {
		cxt := context.TODO()
		rule := testdata.GetRuleWithRuleEmpty(true)
		expectedError := fmt.Errorf("[RuleValidator.evaluate] empty rule: [%v]", rule)
		rulesValidator := rules.NewRulesValidator(logger)

		ruleRepository := new(mocks.RulesRepositoryMock)

		service := rules.NewRulesService(config.Config{}, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		_, err := service.AddRule(cxt, rule)

		assert.Error(t, err)
		assert.EqualValues(t, err, expectedError)
	})

	t.Run("test when add rule fail", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		expectedError := errors.New("error")
		rulesValidator := rules.NewRulesValidator(logger)

		ruleRepository := new(mocks.RulesRepositoryMock)

		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{Data: rulesList}, nil)
		ruleRepository.On("AddRule", rule, context.TODO()).Return(entities.Rule{}, expectedError)

		service := rules.NewRulesService(config.Config{}, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		_, err := service.AddRule(context.TODO(), rule)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("test when add rule successful", func(t *testing.T) {

		rule := testdata.GetDefaultRule(true)
		rulesValidator := rules.NewRulesValidator(nil)
		rulesList := make([]entities.Rule, 0)
		rulesList = append(rulesList, testdata.GetDefaultRuleEmailBlockedGlobal(true))
		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{Data: rulesList}, nil)
		ruleRepository.On("AddRule", rule, context.TODO()).Return(rule, nil)

		ruleRepository.On("GetFamilyCompaniesFromFilter", context.TODO(), entities.FamilyFilter{}).
			Return(entities.Family{}, nil)

		service := rules.NewRulesService(config.Config{}, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		response, err := service.AddRule(context.TODO(), rule)

		assert.NoError(t, err)
		assert.Equal(t, rule, response)
	})

	t.Run("test when add rule fail per rule duplicated", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		rulesValidator := rules.NewRulesValidator(nil)
		expectedError := errors.New("the rule 'device_fingerprint == \"w45345\"' already exist")
		ruleRepository := new(mocks.RulesRepositoryMock)

		ruleRepository.On("GetFamilyCompaniesFromFilter", context.TODO(), entities.FamilyFilter{}).
			Return(entities.Family{}, nil)
		rulesReturn := make([]entities.Rule, 0)
		rulesReturn = append(rulesReturn, rule)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).
			Return(entities.PagedResponse{Data: rulesReturn}, nil)
		ruleRepository.On("AddRule", rule, context.TODO()).Return(rule, nil)

		service := rules.NewRulesService(config.Config{}, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		_, err := service.AddRule(context.TODO(), rule)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("test when add rule successful and operator is IN", func(t *testing.T) {
		rule := testdata.GetDefaultRuleIn(true)
		rulesValidator := rules.NewRulesValidator(logger)

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{Data: rulesList}, nil)
		ruleRepository.On("AddRule", rule, context.TODO()).Return(rule, nil)
		service := rules.NewRulesService(config.Config{}, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		response, err := service.AddRule(context.TODO(), rule)

		assert.NoError(t, err)
		assert.Equal(t, rule, response)
	})

	t.Run("test when add rule successful and operator is IN and value are a couple of number", func(t *testing.T) {

		rule := testdata.GetDefaultRuleInNumber(true)
		rulesValidator := rules.NewRulesValidator(logger)

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{Data: rulesList}, nil)
		ruleRepository.On("AddRule", rule, context.TODO()).Return(rule, nil)
		service := rules.NewRulesService(config.Config{}, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		response, err := service.AddRule(context.TODO(), rule)

		assert.NoError(t, err)
		assert.Equal(t, rule, response)
	})

	t.Run("test when add rule fail where GetRulesByFilters return error", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		rulesValidator := rules.NewRulesValidator(nil)
		expectedError := errors.New("error")

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{}, expectedError)
		ruleRepository.On("AddRule", context.TODO(), rule).Return(nil)

		service := rules.NewRulesService(config.Config{}, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		_, err := service.AddRule(context.TODO(), rule)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("test when add rule successful", func(t *testing.T) {

		rule := testdata.GetDefaultRule(true)
		rulesValidator := rules.NewRulesValidator(nil)
		rulesList := make([]entities.Rule, 0)
		rulesList = append(rulesList, testdata.GetDefaultRuleEmailBlockedGlobal(true))
		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{Data: rulesList}, nil)
		ruleRepository.On("AddRule", rule, context.TODO()).Return(rule, nil)

		ruleRepository.On("GetFamilyCompaniesFromFilter", context.TODO(), entities.FamilyFilter{}).
			Return(entities.Family{}, nil)

		service := rules.NewRulesService(config.Config{}, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		response, err := service.AddRule(context.TODO(), rule)

		assert.NoError(t, err)
		assert.Equal(t, rule, response)
	})
}

func Test_ruleService_UpdateRule(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()
	rulesList := make([]entities.Rule, 0)
	rulesList = append(rulesList, testdata.GetDefaultRuleEmailBlockedGlobal(true))

	t.Run("test when validate rule fail", func(t *testing.T) {
		rule := testdata.GetRuleWithRuleEmptyServiceError(true)
		expectedError := fmt.Errorf("[RuleValidator.evaluate] empty rule: [%v]", rule)
		rulesValidator := rules.NewRulesValidator(logger)

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("UpdateRule", rule, context.TODO()).Return(expectedError)

		service := rules.NewRulesService(configs, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		err := service.UpdateRule(context.TODO(), "611709bb70cbe3606baa3f8d", rule)

		assert.Error(t, err)
		assert.EqualValues(t, expectedError, err)
	})

	t.Run("test when update rule fail", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		expectedError := errors.New("error")
		rulesValidator := rules.NewRulesValidator(nil)

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{}, expectedError)
		ruleRepository.On("UpdateRule", context.TODO(), "611709bb70cbe3606baa3f8d", rule).Return(expectedError)

		service := rules.NewRulesService(configs, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		err := service.UpdateRule(context.TODO(), "611709bb70cbe3606baa3f8d", rule)

		assert.NotNil(t, err)
		assert.Equal(t, err, expectedError)
	})

	t.Run("test when update rule successful", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		rulesValidator := rules.NewRulesValidator(nil)

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{Data: rulesList}, nil)
		ruleRepository.On("UpdateRule", context.TODO(), "611709bb70cbe3606baa3f8d", rule).Return(nil)

		service := rules.NewRulesService(configs, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		err := service.UpdateRule(context.TODO(), "611709bb70cbe3606baa3f8d", rule)

		assert.Nil(t, err)
	})

	t.Run("test when update rule successful", func(t *testing.T) {
		rule := testdata.GetDefaultRuleWithFormula(true)
		rulesValidator := rules.NewRulesValidator(nil)

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{Data: rulesList}, nil)
		ruleRepository.On("UpdateRule", context.TODO(), "611709bb70cbe3606baa3f8d", rule).Return(nil)

		service := rules.NewRulesService(configs, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		err := service.UpdateRule(context.TODO(), "611709bb70cbe3606baa3f8d", rule)

		assert.Nil(t, err)
	})

	t.Run("test when update rule fail per rule duplicated", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		rulesValidator := rules.NewRulesValidator(nil)
		expectedError := errors.New("the rule 'device_fingerprint == \"w45345\"' already exist")
		rulesReturn := make([]entities.Rule, 0)
		rulesReturn = append(rulesReturn, rule)

		ruleRepository := new(mocks.RulesRepositoryMock)

		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{Data: rulesReturn}, expectedError)
		ruleRepository.On("UpdateRule", context.TODO(), "611709bb70cbe3606baa3f8d", rule).Return(nil)

		service := rules.NewRulesService(configs, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		err := service.UpdateRule(context.TODO(), "611709bb70cbe3606baa3f8d", rule)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("test when update rule fail where GetRulesByFilters return error", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		rulesValidator := rules.NewRulesValidator(nil)
		expectedError := errors.New("error")

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("FindRulesPaged", context.TODO(), rule.GetRuleFilter(), entities.Pagination{}).Return(entities.PagedResponse{}, expectedError)
		ruleRepository.On("UpdateRule", context.TODO(), "611709bb70cbe3606baa3f8d", rule).Return(nil)

		service := rules.NewRulesService(configs, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		err := service.UpdateRule(context.TODO(), "611709bb70cbe3606baa3f8d", rule)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})
}

func Test_ruleService_DeleteRule(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("test when delete rule fail", func(t *testing.T) {
		expectedError := errors.New("error")
		rulesValidator := rules.NewRulesValidator(nil)

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("RemoveRule", context.TODO(), "611709bb70cbe3606baa3f8d").Return(expectedError)

		service := rules.NewRulesService(configs, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		err := service.RemoveRule(context.TODO(), "611709bb70cbe3606baa3f8d")

		assert.NotNil(t, err)
		assert.Equal(t, err, expectedError)
	})

	t.Run("test when delete rule successful", func(t *testing.T) {
		rulesValidator := rules.NewRulesValidator(nil)

		ruleRepository := new(mocks.RulesRepositoryMock)
		ruleRepository.On("RemoveRule", context.TODO(), "611709bb70cbe3606baa3f8d").Return(nil)

		service := rules.NewRulesService(configs, rulesValidator, ruleRepository, logger, new(datadog.MetricsDogMock))

		err := service.RemoveRule(context.TODO(), "611709bb70cbe3606baa3f8d")

		assert.Nil(t, err)
	})
}
