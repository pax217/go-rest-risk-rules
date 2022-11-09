package testdata

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
)

var (
	expectedError = errors.New("Service error")
)

func GetDefaultRule(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "2"

	rule := entities.Rule{
		CreatedAt:       now,
		CreatedBy:       "me",
		IsTest:          isATest,
		Module:          "policy_compliance",
		IsGlobal:        false,
		Description:     "empty",
		CompanyID:       &companyId,
		IsYellowFlag:    false,
		FamilyCompanyID: nil,
		FamilyMccID:     nil,
		Rule:            "",
		Rules:           []entities.RuleContent{{Field: "device_fingerprint", Operator: "==", Value: "w45345", Condition: "and"}},
		Decision:        "D",
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleWithID(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "2"

	rule := entities.Rule{
		ID:              primitive.NewObjectID(),
		CreatedAt:       now,
		CreatedBy:       "me",
		IsTest:          isATest,
		Module:          "policy_compliance",
		IsGlobal:        false,
		Description:     "empty",
		IsYellowFlag:    false,
		CompanyID:       &companyId,
		FamilyCompanyID: nil,
		FamilyMccID:     nil,
		Rule:            "",
		Rules:           []entities.RuleContent{{Field: "device_fingerprint", Operator: "==", Value: "w45345", Condition: "and"}},
		Decision:        "D",
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleWithApprovedDecision(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := GetDefaultCharge().CompanyID

	rule := entities.Rule{
		ID:              primitive.NewObjectID(),
		CreatedAt:       now,
		CreatedBy:       "me@gmail.com",
		IsTest:          isATest,
		Module:          "policy_compliance",
		IsGlobal:        false,
		IsYellowFlag:    false,
		Description:     "empty",
		CompanyID:       &companyId,
		FamilyCompanyID: nil,
		FamilyMccID:     nil,
		Rule:            "",
		Rules:           []entities.RuleContent{{Field: "amount", Operator: "==", Value: fmt.Sprintf("%f", GetDefaultCharge().Amount), Condition: "and"}},
		Decision:        "A",
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleWithFamilyMccID(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "2"
	familyID := "61e4dd6da5997ad4d9e76945"

	rule := entities.Rule{
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     false,
		Description:  "empty",
		IsYellowFlag: false,
		CompanyID:    &companyId,
		FamilyMccID:  &familyID,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "device_fingerprint", Operator: "==", Value: "fingerblockeed", Condition: "and"}},
		Decision:     "A",
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleWithFamilyCompanyID(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	familyCompanyID := "61e991ad1214eac062ada43d"

	rule := entities.Rule{
		CreatedAt:       now,
		CreatedBy:       "me",
		IsTest:          isATest,
		Module:          "policy_compliance",
		IsGlobal:        false,
		IsYellowFlag:    false,
		Description:     "empty",
		CompanyID:       nil,
		FamilyMccID:     nil,
		FamilyCompanyID: &familyCompanyID,
		Rule:            "",
		Rules:           []entities.RuleContent{{Field: "device_fingerprint", Operator: "==", Value: "fingerblockeed", Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleFingerprintBlocked(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "2"

	rule := entities.Rule{
		Decision:     "D",
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		IsYellowFlag: false,
		Module:       "policy_compliance",
		IsGlobal:     false,
		Description:  "empty",
		CompanyID:    &companyId,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "device_fingerprint", Operator: "==", Value: "fingerblockeed", Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleEmailBlockedGlobal(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyID := GetChargeWithEmailBlocked().CompanyID

	rule := entities.Rule{
		Decision:     "D",
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     true,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    &companyID,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "details.email", Operator: "==", Value: "block@hotmail.com", Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleEmailGlobalUndefined(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyID := GetChargeWithEmailBlocked().CompanyID

	rule := entities.Rule{
		Decision:     "UN",
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     false,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    &companyID,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "details.email", Operator: "==", Value: "block_diferent@hotmail.com", Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleEmailProximity(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyID := GetChargeWithEmailBlocked().CompanyID

	rule := entities.Rule{
		Decision:     "D",
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "identity_module",
		IsGlobal:     false,
		IsYellowFlag: false,
		Description:  "Some description identity module",
		CompanyID:    &companyID,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "email_proximity.stats.observations_count", Operator: ">", Value: "7", Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleYellowFlag(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyID := GetChargeYellowFlag().CompanyID

	rule := entities.Rule{
		Decision:     "UN",
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "identity_module",
		IsGlobal:     false,
		IsYellowFlag: true,
		Description:  "Some description identity module",
		CompanyID:    &companyID,
		Rule:         "",
		Rules: []entities.RuleContent{
			{Field: "amount", Operator: ">", Value: "800", Condition: "and"},
			{Field: "amount", Operator: "<", Value: "1000", Condition: "or"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleRequestWithAmount() entities.RuleRequest {
	trueValue := true

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleRequestWithNotValueDecision() entities.RuleRequest {
	trueValue := true

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "not_valid_value",
	}
}

func GetDefaultRuleRequestWithValidValueDecision() entities.RuleRequest {
	trueValue := true

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "A",
	}
}

func GetDefaultRuleNotGlobalWithoutFamilyIDAndCompanyID() entities.RuleRequest {
	trueValue := false

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleNotGlobalWithFamilyIDAndCompanyID() entities.RuleRequest {
	trueValue := false

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		CompanyID:   "617ae8a4649e59500b7cd54d",
		FamilyID:    "617ae8bc92d0e243227eee9d",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleNotGlobalWithCompanyID() entities.RuleRequest {
	trueValue := false

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		CompanyID:   "617ae8a4649e59500b7cd54d",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleNotGlobalWithFamilyID() entities.RuleRequest {
	trueValue := false
	familyID := "617ae8a4649e59500b7cd54d"

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		FamilyID:    familyID,
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleGlobalWithFamilyID() entities.RuleRequest {
	trueValue := true
	familyID := "617ae8a4649e59500b7cd54d"

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		FamilyID:    familyID,
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleGlobal() entities.RuleRequest {
	trueValue := true

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleGlobalWithCompanyID() entities.RuleRequest {
	trueValue := true
	CompanyID := "617ae8bc92d0e243227eee9d"

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		CompanyID:   CompanyID,
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetJsonMalformed() string {
	request := `{
		"is_test": false,
		"module": "policy_compliance",
		"is_global": true,
		"description": "empty",
		"company_id": "2",
		"author": "me@conekta.com",
		"rules": [
			{
				"field": "device_fingerprint",
				"operator": "==",
				"value": "fingerblockeed",
				"condition": "and",
			}
		]
	}`
	return request
}

func GetDefaultRuleRequestFailValidation() entities.RuleRequest {
	trueValue := true
	return entities.RuleRequest{
		IsTest:      &trueValue,
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleRequestReturnError() (entities.RuleRequest, error) {
	trueValue := true

	request := entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}

	return request, expectedError
}

func GetDefaultRuleRequestGlobalToFamily() entities.RuleRequest {
	falseValue := false
	isTest := true

	request := entities.RuleRequest{
		IsTest:      &isTest,
		Module:      "policy_compliance",
		IsGlobal:    &falseValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		CompanyID:   "",
		FamilyID:    "6171ad238730441d5ec9537b",
		Decision:    "D",
	}

	return request
}

func GetRuleWithRuleEmpty(isATest bool) entities.Rule {
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "2"

	rule := entities.Rule{
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     false,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    &companyId,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "", Operator: "", Value: "", Condition: ""}},
	}

	return rule
}
func GetDefaultRuleIn(isATest bool) entities.Rule {
	companyId := "2"
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)

	rule := entities.Rule{
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     false,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    &companyId,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "device_fingerprint", Operator: "in", Value: `["mx","ux"]`, Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}
func GetDefaultRuleInNumber(isATest bool) entities.Rule {
	companyId := "2"
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)

	rule := entities.Rule{
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     false,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    &companyId,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "amount", Operator: "in", Value: `[450.0,180.0]`, Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetJsonRequestIsMalformed() string {
	request := `{
		"is_test": true,
		"module": "policy_compliance",
		"is_global": true,
		"description": "charge is too low",
		"company_id": "",
		"rules": [
			{
				"field": "amount",
				"operator": "\u003c",
				"value": "8",
				"condition": "and"
			}
		],
		"author": "sfc@conekta.com",
	}`

	return request
}

func GetJsonRequestWithoutModule() entities.RuleRequest {
	trueValue := true

	return entities.RuleRequest{
		IsTest:      &trueValue,
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleRequestServiceReturnError() (entities.RuleRequest, error) {
	trueValue := true

	request := entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		Rules:       []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}

	return request, expectedError
}

func GetRuleWithRuleEmptyServiceError(isATest bool) entities.Rule {
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "2"

	rule := entities.Rule{
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     false,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    &companyId,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "", Operator: "", Value: "", Condition: ""}},
	}

	return rule
}

func GetDefaultRuleRequestFailFormulaValidation() entities.RuleRequest {
	trueValue := false

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		CompanyID:   "617ae8a4649e59500b7cd54d",
		Rules:       GetNoMathOperationFormula(),
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetRuleRequestWithFieldAndFormulaFields() entities.RuleRequest {
	trueValue := false

	return entities.RuleRequest{
		IsTest:      &trueValue,
		Module:      "policy_compliance",
		IsGlobal:    &trueValue,
		Description: "charge is too low",
		CompanyID:   "617ae8a4649e59500b7cd54d",
		Rules:       GetRulesWithFormulaFieldsAndRuleField(),
		Author:      "sfc@conekta.com",
		Decision:    "D",
	}
}

func GetDefaultRuleFilter() (entities.RuleFilter, error) {
	request := entities.RuleFilter{
		ID:        "617ae8a4649e59500b7cd54d",
		CompanyID: "2",
		IsTest:    false,
		Module:    "",
		IsGlobal:  false,
		Rule:      "",
		FamilyID:  "",
	}

	return request, expectedError
}

func GetDefaultRuleEmailBlockedGlobalForGraylist(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)

	rule := entities.Rule{
		Decision:     "D",
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     true,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    nil,
		Rule:         "",
		Rules:        []entities.RuleContent{{Field: "details.email", Operator: "==", Value: "mail@hotmail.com", Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleWithFormula(isATest bool) entities.Rule {
	rule := GetDefaultRule(isATest)
	rule.Rule = "SUM (amount.h9,amount.h12) == 1 and device_fingerprint == \"w45345\""
	rule.Rules = append(GetDefaultFormulaRulesWithProcessedField(), rule.Rules...)
	return rule
}

func GetDefaultRuleEmailWithChargebacks(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyID := "7683457364"

	rule := entities.Rule{
		Decision:     "D",
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     true,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    &companyID,
		Rule:         "payer.chargebacks > 0",
		Rules:        []entities.RuleContent{{Field: "payer.chargebacks", Operator: ">", Value: "0", Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleWithOmniscore(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyID := "7683457364"

	rule := entities.Rule{
		Decision:     "D",
		CreatedAt:    now,
		CreatedBy:    "me",
		IsTest:       isATest,
		Module:       "policy_compliance",
		IsGlobal:     true,
		IsYellowFlag: false,
		Description:  "empty",
		CompanyID:    &companyID,
		Rule:         "omniscore > 0.3",
		Rules:        []entities.RuleContent{{Field: "omniscore", Operator: ">", Value: "0.3", Condition: "and"}},
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleMerchantScoreApproved(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "7683457364"

	rule := entities.Rule{
		CreatedAt:       now,
		CreatedBy:       "me",
		IsTest:          isATest,
		Module:          "policy_compliance",
		IsGlobal:        false,
		IsYellowFlag:    false,
		Description:     "empty",
		CompanyID:       &companyId,
		FamilyCompanyID: nil,
		FamilyMccID:     nil,
		Rule:            "",
		Rules:           []entities.RuleContent{{Field: "merchant_score", Operator: "==", Value: "-1", Condition: "and"}},
		Decision:        "A",
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleMerchantScoreDeclined(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "7683457364"

	rule := entities.Rule{
		CreatedAt:       now,
		CreatedBy:       "me",
		IsTest:          isATest,
		Module:          "policy_compliance",
		IsGlobal:        false,
		IsYellowFlag:    false,
		Description:     "empty",
		CompanyID:       &companyId,
		FamilyCompanyID: nil,
		FamilyMccID:     nil,
		Rule:            "",
		Rules:           []entities.RuleContent{{Field: "merchant_score", Operator: "==", Value: "-1", Condition: "and"}},
		Decision:        "D",
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleCompanyRuleAccepted(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "7683457364"

	rule := entities.Rule{
		CreatedAt:       now,
		CreatedBy:       "me",
		IsTest:          isATest,
		Module:          "policy_compliance",
		IsGlobal:        false,
		Description:     "empty",
		CompanyID:       &companyId,
		FamilyCompanyID: nil,
		FamilyMccID:     nil,
		Rule:            "",
		Rules: []entities.RuleContent{
			{Field: "amount", Operator: ">", Value: "800", Condition: "and"},
			{Field: "amount", Operator: "<", Value: "100", Condition: "or"},
		},
		Decision: "A",
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}

func GetDefaultRuleRequestWithYellowFlag() entities.RuleRequest {
	trueValue := true

	return entities.RuleRequest{
		IsTest:       &trueValue,
		Module:       "policy_compliance",
		IsGlobal:     &trueValue,
		Description:  "charge is too low",
		Rules:        []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:       "sfc@conekta.com",
		Decision:     "UN",
		IsYellowFlag: true,
	}
}

func GetDefaultRuleRequestWithYellowFlagIncorrectDecision() entities.RuleRequest {
	trueValue := true

	return entities.RuleRequest{
		IsTest:       &trueValue,
		Module:       "policy_compliance",
		IsGlobal:     &trueValue,
		Description:  "charge is too low",
		Rules:        []entities.RuleContent{{Field: "amount", Operator: "<", Value: "8", Condition: "and"}},
		Author:       "sfc@conekta.com",
		Decision:     "A",
		IsYellowFlag: true,
	}
}

func GetDefaultRuleMarketSegmentApproved(isATest bool) entities.Rule {
	ruleService := rules.NewRulesService(config.NewConfig(), nil, nil, nil, nil)
	now := time.Date(2021, 07, 24, 12, 30, 00, 00, time.UTC).Truncate(time.Millisecond)
	companyId := "7683457364"

	rule := entities.Rule{
		CreatedAt:       now,
		CreatedBy:       "me",
		IsTest:          isATest,
		Module:          "policy_compliance",
		IsGlobal:        false,
		Description:     "empty",
		CompanyID:       &companyId,
		FamilyCompanyID: nil,
		FamilyMccID:     nil,
		Rule:            "",
		Rules:           []entities.RuleContent{{Field: "market_segment", Operator: "==", Value: "long tail SMB", Condition: "and"}},
		Decision:        "A",
	}

	rule.Rule = ruleService.BuildRule(rule.Rules)
	return rule
}
