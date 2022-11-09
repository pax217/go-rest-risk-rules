package integration

import (
	"context"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/conditions"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConditionRepository_Add(t *testing.T) {
	t.Run("on condition_id duplicated then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		condition := testdata.GetDefaultCondition()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Add(ctx, &condition)

		assert.NoError(t, err)
		assert.NotEmpty(t, condition.ID)

		var cdt entities.Condition
		err = mongoDB.Collection(cfg.MongoDB.Collections.Conditions).FindOne(ctx, bson.M{"_id": condition.ID}).Decode(&cdt)

		err = repository.Add(ctx, &condition)

		mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Conditions, cdt.ID)

		assert.Error(t, err)
	})

	t.Run("on condition added is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		condition := testdata.GetDefaultCondition()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Add(ctx, &condition)

		assert.NoError(t, err)
		assert.NotEmpty(t, condition.ID)

		var cdt entities.Condition
		err = mongoDB.Collection(cfg.MongoDB.Collections.Conditions).FindOne(ctx, bson.M{"_id": condition.ID}).Decode(&cdt)

		mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Conditions, cdt.ID)
	})
}

func TestConditionRepository_FindByName(t *testing.T) {

	t.Run("on condition not found", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		condition := testdata.GetDefaultCondition()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conditionFound, err := repository.FindByName(ctx, condition)

		assert.NoError(t, err)
		assert.Equal(t, entities.Condition{}, conditionFound)
	})

	t.Run("on condition found", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		condition := testdata.GetDefaultCondition()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Add(ctx, &condition)

		assert.NoError(t, err)

		conditionFound, err := repository.FindByName(ctx, condition)

		assert.NoError(t, err)
		assert.NotEqual(t, entities.Condition{}, conditionFound)

		var cdt entities.Condition
		err = mongoDB.Collection(cfg.MongoDB.Collections.Conditions).FindOne(ctx, bson.M{"_id": condition.ID}).Decode(&cdt)

		mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Conditions, cdt.ID)
	})
}

func TestConditionRepository_GetAll(t *testing.T) {

	t.Run("on condition_id no valid, then return array empty conditions", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conditionsFilter := entities.ConditionsFilter{
			ID: "615238fcc1b177e01aefccf8-x",
		}
		pagination := entities.Pagination{
			PageNumber: 0,
			PageSize:   500,
		}

		conditionsFounded, err := repository.GetAll(
			ctx,
			conditionsFilter,
			pagination)

		assert.Nil(t, err)
		assert.EqualValues(t, conditionsFounded.Data, []entities.Condition{})
	})

	t.Run("on condition filter valid, then return array of conditions", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		conditionsData := testdata.GetConditions()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conditionsFilter := entities.ConditionsFilter{}
		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Conditions, conditionsData[0], conditionsData[1])
		pagination := entities.Pagination{
			PageNumber: 0,
			PageSize:   500,
		}

		conditionsFounded, err := repository.GetAll(
			ctx,
			conditionsFilter,
			pagination)

		data := conditionsFounded.Data.([]entities.Condition)

		assert.Nil(t, err)
		assert.NotEqual(t, conditionsFounded.Data, []entities.Condition{})

		for i := range data {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Conditions, data[i].ID)
		}
	})

	t.Run("on get all successfully, then return array of conditions sort by name", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		conditionsData := testdata.GetConditions()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conditionsFilter := entities.ConditionsFilter{}

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.Conditions)

		objectIDS := mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Conditions, conditionsData[0], conditionsData[1], conditionsData[2])
		conditionsData[0].SetID(objectIDS[0])
		conditionsData[1].SetID(objectIDS[1])
		conditionsData[2].SetID(objectIDS[2])

		pagination := entities.Pagination{
			PageNumber: 0,
			PageSize:   500,
		}

		conditionsFounded, err := repository.GetAll(
			ctx,
			conditionsFilter,
			pagination)

		data := conditionsFounded.Data.([]entities.Condition)

		assert.NoError(t, err)
		assert.Equalf(t, len(conditionsData), len(data), "conditions len should be %d", len(conditionsData))
		assert.Equal(t, conditionsData[0], data[0])
		assert.Equal(t, conditionsData[1], data[2])
		assert.Equal(t, conditionsData[2], data[1])

		for i := range data {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Conditions, data[i].ID)
		}
	})
}

func TestConditionRepository_Update(t *testing.T) {

	t.Run("on condition_id not found, then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		condition := testdata.GetDefaultCondition()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Update(ctx, "615238fcc1b177e01aefccf8", condition)

		assert.Error(t, err)
	})

	t.Run("on condition updated ok, then return nil", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		condition := testdata.GetDefaultCondition()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Add(ctx, &condition)

		now := time.Now().UTC().Truncate(time.Millisecond)
		updatedBy := "carlos.maldonado@conekta.com"

		conditionToUpdate := entities.Condition{
			ID:          condition.ID,
			Name:        "name updated",
			Description: "description updated",
			UpdatedBy:   &updatedBy,
			UpdatedAt:   &now,
		}

		err = repository.Update(
			ctx,
			conditionToUpdate.ID.Hex(),
			conditionToUpdate)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Conditions, conditionToUpdate.ID)

		var cdt entities.Condition
		err = mongoDB.Collection(cfg.MongoDB.Collections.Conditions).FindOne(ctx, bson.M{"_id": conditionToUpdate.ID}).Decode(&cdt)

		assert.NoError(t, err)
		assert.NotEmpty(t, conditionToUpdate.ID)
		assert.EqualValues(t, conditionToUpdate.ID, cdt.ID)
		assert.EqualValues(t, conditionToUpdate.Name, cdt.Name)
		assert.EqualValues(t, conditionToUpdate.Description, cdt.Description)
		assert.EqualValues(t, conditionToUpdate.UpdatedAt, cdt.UpdatedAt)
		assert.EqualValues(t, conditionToUpdate.UpdatedBy, cdt.UpdatedBy)
	})

}

func TestConditionRepository_Delete(t *testing.T) {
	t.Run("on condition_id not found, then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Delete(ctx, "615238fcc1b177e01aefccf8")

		assert.Error(t, err)
	})

	t.Run("on condition deleted is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		condition := testdata.GetDefaultCondition()
		repository := conditions.NewConditionsRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Add(ctx, &condition)

		assert.NoError(t, err)
		assert.NotEmpty(t, condition.ID)

		err = repository.Delete(ctx, condition.ID.Hex())

		assert.NoError(t, err)
	})
}
