package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRuleRepository_AddRule(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when add rule successful", func(t *testing.T) {

		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		response, errAdd := repository.AddRule(ctx, rule)
		assert.NoError(t, errAdd)

		ruleDelete, err := repository.GetRulesByFilters(ctx, rule.GetRuleFilter(), entities.CompanyRulesType)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, ruleDelete)
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleDelete[0].ID)
	})

	t.Run("when add rule fail", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		rule.ID = primitive.NewObjectID()
		_, err := repository.AddRule(ctx, rule)
		assert.NoError(t, err)
		_, err = repository.AddRule(ctx, rule)

		assert.NotNil(t, err)
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, rule.ID)
	})

	t.Run("when add rule successful, validate decision ", func(t *testing.T) {
		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCreated, _ := repository.AddRule(ctx, rule)

		filter := rule.GetRuleFilter()
		filter.ID = rule.ID.Hex()
		response, err := repository.GetRulesByFilters(ctx, filter, entities.CompanyRulesType)
		RuleFounded := response

		assert.Nil(t, err)
		assert.Equal(t, ruleCreated.CompanyID, RuleFounded[0].CompanyID)
		assert.Equal(t, ruleCreated.Decision, RuleFounded[0].Decision)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleCreated.ID)
	})
}

func TestRuleRepository_GetRulesByFilters(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when find all rules successful", func(t *testing.T) {

		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCreated, _ := repository.AddRule(ctx, rule)
		response, err := repository.GetRulesByFilters(ctx, rule.GetRuleFilter(), entities.CompanyRulesType)
		RuleFounded := response

		assert.Nil(t, err)
		assert.Equal(t, ruleCreated.CompanyID, RuleFounded[0].CompanyID)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleCreated.ID)
	})

	t.Run("when find rules for identity module, id should be the same", func(t *testing.T) {
		rule := testdata.GetDefaultRuleEmailProximity(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCreated, _ := repository.AddRule(ctx, rule)
		response, err := repository.GetRulesByFilters(ctx, rule.GetRuleFilter(), entities.IdentityModuleType)
		RuleFounded := response

		assert.Nil(t, err)
		assert.Equal(t, ruleCreated.ID, RuleFounded[0].ID)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleCreated.ID)
	})
}

func TestRuleRepository_GlobalRuleType(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when find all rules successful", func(t *testing.T) {

		rule := testdata.GetDefaultRuleEmailBlockedGlobal(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCreated, _ := repository.AddRule(ctx, rule)
		response, err := repository.GetRulesByFilters(ctx, rule.GetRuleFilter(), entities.GlobalRulesType)
		RuleFounded := response

		assert.Nil(t, err)
		assert.Equal(t, ruleCreated.CompanyID, RuleFounded[0].CompanyID)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleCreated.ID)
	})
}

func TestRuleRepository_Find_Rules_Paged(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when find rules paged successful", func(t *testing.T) {

		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCreated, _ := repository.AddRule(ctx, rule)
		response, err := repository.FindRulesPaged(ctx, rule.GetRuleFilter(), entities.Pagination{})
		RuleFounded := response.Data.([]entities.Rule)

		assert.Nil(t, err)
		assert.Equal(t, ruleCreated.CompanyID, RuleFounded[0].CompanyID)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleCreated.ID)
	})

	t.Run("when filter has value on family_company_id", func(t *testing.T) {
		rule := testdata.GetDefaultRuleWithFamilyCompanyID(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCreated, _ := repository.AddRule(ctx, rule)
		response, err := repository.FindRulesPaged(ctx, rule.GetRuleFilter(), entities.Pagination{})
		RuleFounded := response.Data.([]entities.Rule)

		assert.Nil(t, err)
		assert.Equal(t, ruleCreated.CompanyID, RuleFounded[0].CompanyID)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleCreated.ID)
	})

	t.Run("when find rules and ID is invalid then return error", func(t *testing.T) {

		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCreated, _ := repository.AddRule(ctx, rule)
		modify := rule.GetRuleFilter()
		modify.ID = "invalid"
		response, err := repository.FindRulesPaged(ctx, modify, entities.Pagination{})

		assert.Error(t, err)
		assert.Empty(t, response)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleCreated.ID)
	})

	t.Run("when getting rules by valid ID, then return the rule with same ID", func(t *testing.T) {
		rule := testdata.GetDefaultRuleWithID(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCreated, _ := repository.AddRule(ctx, rule)
		filter := entities.RuleFilter{
			ID: rule.ID.Hex(),
		}
		response, err := repository.FindRulesPaged(ctx, filter, entities.Pagination{})
		ruleFound := response.Data.([]entities.Rule)

		assert.Nil(t, err)
		assert.Equal(t, ruleCreated.ID, ruleFound[0].ID)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleCreated.ID)
	})
}

func TestRuleRepository_DeleteRule(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when delete rule fail ID not found", func(t *testing.T) {

		expectedError := errors.New("error: record not found")
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.RemoveRule(ctx, "611709bb70cbe3606baa3f8d")

		assert.NotNil(t, err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("when delete rule fail ID incorrect", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		expectedError := errors.New("the provided hex string is not a valid ObjectID")
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.RemoveRule(ctx, "")

		assert.NotNil(t, err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("when delete rule successful", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		response, _ := repository.AddRule(ctx, rule)

		ruleCreated, _ := repository.GetRulesByFilters(ctx, response.GetRuleFilter(), entities.CompanyRulesType)

		err := repository.RemoveRule(ctx, ruleCreated[0].ID.Hex())
		assert.Nil(t, err)
	})
}

func TestRuleRepository_UpdateRule(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when rule_id is not valid, returns an error", func(t *testing.T) {

		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleID := "6154fb866f494f0769673fc3-x"

		err := repository.UpdateRule(ctx, ruleID, rule)

		assert.Error(t, err)
		assert.Contains(t, "the provided hex string is not a valid ObjectID", err.Error())
	})

	t.Run("when rule_id does not exist, returns an error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleID := "6154fb866f494f0769673fc3"

		err := repository.UpdateRule(ctx, ruleID, rule)

		assert.Error(t, err)
		assert.Contains(t, "error, document not found", err.Error())
	})

	t.Run("when the decision value changes successfully, then returns NoContent", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		famCompanyID := "61f168736db1497e893ffe75"
		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := repository.AddRule(ctx, rule)
		assert.NoError(t, err)

		ruleAdded, err := repository.GetRulesByFilters(ctx, rule.GetRuleFilter(), entities.CompanyRulesType)
		assert.NoError(t, err)
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleAdded[0].ID)

		updatedBy := "carlos.maldonado@conekta.com"
		updatedAt := time.Now().Truncate(time.Millisecond).UTC()
		ruleToUpdate := entities.Rule{
			CompanyID:       nil,
			FamilyCompanyID: &famCompanyID,
			ID:              ruleAdded[0].ID,
			UpdatedBy:       &updatedBy,
			UpdatedAt:       &updatedAt,
			IsYellowFlag:    false,
			Description:     "Description Updated",
			Decision:        "A",
		}
		err = repository.UpdateRule(ctx, ruleAdded[0].ID.Hex(), ruleToUpdate)

		var result entities.Rule
		err = mongoDB.Collection(cfg.MongoDB.Collections.Rules).
			FindOne(ctx, bson.M{"_id": ruleAdded[0].ID}).Decode(&result)

		assert.NoError(t, err)

		assert.NotEmpty(t, ruleToUpdate.ID)
		assert.EqualValues(t, ruleToUpdate.ID, result.ID)
		assert.EqualValues(t, ruleToUpdate.Decision, result.Decision)
	})

	t.Run("when update rule successful", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		famCompanyID := "61f168736db1497e893ffe75"
		rule := testdata.GetDefaultRule(true)
		repository := rules.NewRuleMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := repository.AddRule(ctx, rule)
		assert.NoError(t, err)

		ruleAdded, err := repository.GetRulesByFilters(ctx, rule.GetRuleFilter(), entities.CompanyRulesType)
		assert.NoError(t, err)
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Rules, ruleAdded[0].ID)

		updatedBy := "carlos.maldonado@conekta.com"
		updatedAt := time.Now().Truncate(time.Millisecond).UTC()
		ruleToUpdate := entities.Rule{
			CompanyID:       nil,
			FamilyCompanyID: &famCompanyID,
			ID:              ruleAdded[0].ID,
			UpdatedBy:       &updatedBy,
			IsYellowFlag:    false,
			UpdatedAt:       &updatedAt,
			Description:     "Description Updated",
		}
		err = repository.UpdateRule(ctx, ruleAdded[0].ID.Hex(), ruleToUpdate)

		var result entities.Rule
		err = mongoDB.Collection(cfg.MongoDB.Collections.Rules).
			FindOne(ctx, bson.M{"_id": ruleAdded[0].ID}).Decode(&result)

		assert.NoError(t, err)

		assert.NotEmpty(t, ruleToUpdate.ID)
		assert.EqualValues(t, ruleToUpdate.ID, result.ID)
		assert.EqualValues(t, ruleToUpdate.Description, result.Description)
		assert.EqualValues(t, ruleToUpdate.UpdatedAt, result.UpdatedAt)
		assert.EqualValues(t, ruleToUpdate.UpdatedBy, result.UpdatedBy)
		assert.EqualValues(t, ruleToUpdate.FamilyCompanyID, result.FamilyCompanyID)
		assert.EqualValues(t, ruleToUpdate.CompanyID, result.CompanyID)
	})
}
