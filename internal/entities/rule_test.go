package entities_test

import (
	"testing"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewRuleFromPutRequest(t *testing.T) {
	t.Run("check rule to update from ruleRequest", func(t *testing.T) {
		ruleRequest := testdata.GetDefaultRuleRequestGlobalToFamily()

		ruleToUpdate := ruleRequest.NewRuleFromPutRequest()

		assert.Equal(t, ruleToUpdate.UpdatedBy, &ruleRequest.Author)
		assert.NotNil(t, ruleToUpdate.UpdatedAt)
		assert.Equal(t, ruleToUpdate.IsTest, *ruleRequest.IsTest)
		assert.Equal(t, ruleToUpdate.Module, ruleRequest.Module)
		assert.False(t, ruleToUpdate.IsGlobal, &ruleRequest.IsGlobal)
		assert.NotNil(t, ruleToUpdate.Description, ruleRequest.Description)
		assert.Nil(t, ruleToUpdate.CompanyID)
		assert.Equal(t, ruleToUpdate.FamilyMccID, &ruleRequest.FamilyID)
		assert.Equal(t, ruleToUpdate.Rules, ruleRequest.Rules)
	})
}

func TestNewRuleFromPostRequest(t *testing.T) {
	t.Run("check rule to create from ruleRequest", func(t *testing.T) {
		ruleRequest := testdata.GetDefaultRuleRequestGlobalToFamily()

		ruleToUpdate := ruleRequest.NewRuleFromPostRequest()

		assert.Nil(t, ruleRequest.Validate())
		assert.Equal(t, ruleToUpdate.Module, ruleRequest.Module)
		assert.False(t, ruleToUpdate.IsGlobal, &ruleRequest.IsGlobal)
		assert.NotNil(t, ruleToUpdate.Description, ruleRequest.Description)
		assert.Nil(t, ruleToUpdate.CompanyID)
		assert.Equal(t, ruleToUpdate.FamilyMccID, &ruleRequest.FamilyID)
		assert.Equal(t, ruleToUpdate.Rules, ruleRequest.Rules)
		ruleRequest.Decision = entities.Declined
		ruleRequest.IsYellowFlag = true
		assert.True(t, ruleRequest.IsYellowFlag)
		assert.Error(t, ruleRequest.Validate())
	})
}

func TestList_RuleResponse(t *testing.T) {
	t.Run("Check if is empty White List Response", func(t *testing.T) {
		ruleResponse := entities.NewRulesResponse()
		decision := ruleResponse.GetDecision()
		testDecision := ruleResponse.GetTestDecision()
		entityName := ruleResponse.GetEntityName()

		assert.Equal(t, entities.Undecided, decision)
		assert.Equal(t, entities.Undecided, testDecision)
		assert.NotEmpty(t, entityName)
		assert.Empty(t, ruleResponse.DecisionRules)
	})
}

func TestList_RuleFilter(t *testing.T) {
	t.Run("Check if is valid ID in RuleFilter", func(t *testing.T) {
		ruleFilter, _ := testdata.GetDefaultRuleFilter()

		assert.True(t, ruleFilter.IsIDValid())

		ruleFilter.ID = "ivalid"
		assert.False(t, ruleFilter.IsIDValid())

	})
}

func TestList_RuleRequest(t *testing.T) {
	t.Run("Check if is valid RuleRequest", func(t *testing.T) {
		request := testdata.GetDefaultRuleRequestWithAmount()
		assert.Nil(t, request.Validate())

		isGlobal := true
		request.IsGlobal = &isGlobal
		request.CompanyID = "2"
		assert.Error(t, request.Validate())

		isGlobal = false
		request.CompanyID = ""
		request.FamilyID = ""
		request.FamilyCompanyID = ""
		assert.Error(t, request.Validate())

		isGlobal = true
		request.IsGlobal = &isGlobal
		request.CompanyID = "2"
		request.FamilyID = "22"
		assert.Error(t, request.Validate())

		isGlobal = false
		request.IsGlobal = &isGlobal
		request.CompanyID = ""
		request.FamilyCompanyID = "22"
		assert.Error(t, request.Validate())
	})
}

func TestList_IsContained(t *testing.T) {
	t.Run("Check if rule is contained", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		rules := make([]entities.Rule, 0)

		assert.False(t, rule.IsContained(rules))

		rules = append(rules, rule)
		assert.True(t, rule.IsContained(rules))
	})
}

func TestRule_SetRuleResponse(t *testing.T) {
	t.Run("Check set rule response to rulesModulesResponse",
		func(t *testing.T) {
			consoleDefaultOnlyRules := testdata.SetDefaultConsoleOnlyRules()
			rulesModulesResponse := entities.RulesModulesResponse{}

			rulesResponse := entities.RulesResponse{
				DecisionRules: []entities.Rule{
					testdata.GetDefaultRule(false),
				},
			}

			for _, component := range consoleDefaultOnlyRules {
				rulesModulesResponse.SetRuleResponse(component, rulesResponse)
			}
		})
}
