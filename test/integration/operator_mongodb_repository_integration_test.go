package integration

import (
	"context"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/operators"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestOperatorRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)
	operator := testdata.GetOperators()[0]
	repository := operators.NewOperatorMongoDBRepository(cfg, mongoDB, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := repository.Save(ctx, &operator)

	assert.NoError(t, err)
	assert.NotEmpty(t, operator.ID)

	defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Operators, operator.ID)

	var op entities.Operator
	err = mongoDB.Collection(cfg.MongoDB.Collections.Operators).FindOne(ctx, bson.M{"_id": operator.ID}).Decode(&op)

	assert.NoError(t, err)
	assert.EqualValues(t, operator, op)
}

func TestOperatorRepository_Get(t *testing.T) {
	t.Run("When operator is found", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		operator := testdata.GetOperators()[0]
		repository := operators.NewOperatorMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Operators, operator)

		operatorsFounded, err := repository.Get(ctx, entities.OperatorFilter{
			Type: operator.Type,
			Name: operator.Name,
		})

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Operators, operatorsFounded[0].ID)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(operatorsFounded), "operators len should be 1")
		assert.Equal(t, operator.Name, operatorsFounded[0].Name)
		assert.Equal(t, operator.CreatedBy, operatorsFounded[0].CreatedBy)
		assert.EqualValues(t, operator.CreatedAt, operatorsFounded[0].CreatedAt)

	})
}

func TestOperatorRepository_GetPaged(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)
	operatorsData := testdata.GetOperators()
	repository := operators.NewOperatorMongoDBRepository(cfg, mongoDB, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Operators, operatorsData[0], operatorsData[1])
	operatorIDS := mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Fields, operatorsData[0], operatorsData[1])
	operatorsData[0].SetID(operatorIDS[0])
	operatorsData[1].SetID(operatorIDS[1])

	pagination := entities.Pagination{
		PageNumber: 0,
		PageSize:   500,
	}
	operatorsFounded, err := repository.GetPaged(
		ctx,
		entities.OperatorFilter{
			Type: "type",
		},
		pagination)
	data := operatorsFounded.Data.([]entities.Operator)

	assert.NoError(t, err)
	assert.Equalf(t, len(operatorsData), len(data), "operators len should be %d", len(operatorsData))
	assert.Equal(t, operatorsData[0].Name, data[1].Name)
	assert.Equal(t, operatorsData[0].Type, data[1].Type)
	assert.Equal(t, operatorsData[0].Title, data[1].Title)
	assert.Equal(t, operatorsData[0].Description, data[1].Description)

	assert.Equal(t, operatorsData[1].Name, data[0].Name)
	assert.Equal(t, operatorsData[1].Type, data[0].Type)
	assert.Equal(t, operatorsData[1].Title, data[0].Title)
	assert.Equal(t, operatorsData[1].Description, data[0].Description)

	for i := range data {
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Operators, data[i].ID)
	}
}

func TestOperatorRepository_GetPaged_error(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)
	repository := operators.NewOperatorMongoDBRepository(cfg, mongoDB, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pagination := entities.Pagination{
		PageNumber: 0,
		PageSize:   500,
	}
	operatorsFounded, err := repository.GetPaged(
		ctx,
		entities.OperatorFilter{
			ID: "invalid",
		},
		pagination)

	assert.Error(t, err)
	assert.Empty(t, operatorsFounded)
}

func TestOperatorRepository_Delete(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)
	repository := operators.NewOperatorMongoDBRepository(cfg, mongoDB, logger)

	t.Run("on operator_id not found, then return error", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Delete(ctx, "615238fcc1b177e01aefccf8")

		assert.Error(t, err)
		assert.Contains(t,
			err.Error(),
			"error: record not found: ObjectID(\"615238fcc1b177e01aefccf8\")")
	})

	t.Run("on operator deleted is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		operator := testdata.GetOperatorDefault()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, &operator)

		assert.NoError(t, err)
		assert.NotEmpty(t, operator.ID)

		err = repository.Delete(ctx, operator.ID.Hex())

		assert.NoError(t, err)
	})

	t.Run("when ID is invalid, then return error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Delete(ctx, "x")

		assert.Error(t, err)
	})
}

func TestOperatorRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)
	repository := operators.NewOperatorMongoDBRepository(cfg, mongoDB, logger)

	t.Run("on condition_id not found, then return error", func(t *testing.T) {

		condition := testdata.GetOperatorDefault()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Update(ctx, "615238fcc1b177e01aefccf8", condition)

		assert.Error(t, err)
	})

	t.Run("on condition updated ok, then return nil", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		operator := testdata.GetOperatorDefault()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, &operator)

		now := time.Now().UTC().Truncate(time.Millisecond)
		updatedBy := "carlos.maldonado@conekta.com"

		operatorToUpdate := entities.Operator{
			ID:          operator.ID,
			Name:        "name updated",
			Description: "description updated",
			Type:        "string updated",
			UpdatedBy:   &updatedBy,
			UpdatedAt:   &now,
		}

		err = repository.Update(
			ctx,
			operatorToUpdate.ID.Hex(),
			operatorToUpdate)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Operators, operatorToUpdate.ID)

		var opr entities.Operator
		err = mongoDB.Collection(cfg.MongoDB.Collections.Operators).FindOne(ctx, bson.M{"_id": operatorToUpdate.ID}).Decode(&opr)

		assert.NoError(t, err)
		assert.NotEmpty(t, operatorToUpdate.ID)
		assert.EqualValues(t, operatorToUpdate.ID, opr.ID)
		assert.EqualValues(t, operatorToUpdate.Name, opr.Name)
		assert.EqualValues(t, operatorToUpdate.Description, opr.Description)
		assert.EqualValues(t, operatorToUpdate.UpdatedAt, opr.UpdatedAt)
		assert.EqualValues(t, operatorToUpdate.UpdatedBy, opr.UpdatedBy)
	})

	t.Run("when ID is invalid, then return error", func(t *testing.T) {

		condition := testdata.GetOperatorDefault()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Update(ctx, "x", condition)

		assert.Error(t, err)
	})
}
