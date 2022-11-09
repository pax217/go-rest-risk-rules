package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/charges"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EvaluationDelete struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

func TestChargeRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when save charge successful", func(t *testing.T) {
		charge := testdata.GetEvaluationResponseAcceptedByWhiteListSuccessful()
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		charge.Charge.ID = primitive.NewObjectID().Hex()
		defer cancel()

		err := repository.Save(ctx, charge)
		ID := get(ctx, mongoDB, cfg, charge.Charge.ID)

		assert.Nil(t, err)
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.ChargeEvaluations, ID)
	})
}

func TestChargeRepository_SaveOnlyRules(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when save charge successful", func(t *testing.T) {
		rulesEvaluationResponse := testdata.GetRulesEvaluationResponseAcceptedCauseConsoleCompanyRules()
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		rulesEvaluationResponse.Charge.ID = primitive.NewObjectID().Hex()
		defer cancel()

		err := repository.SaveOnlyRules(ctx, rulesEvaluationResponse)
		ID := getOnlyRules(ctx, mongoDB, cfg, rulesEvaluationResponse.Charge.ID)

		assert.Nil(t, err)
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.ChargeEvaluationsOnlyRules, ID)
	})
}

func TestChargeRepository_Get(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when get charge successful", func(t *testing.T) {

		charge := testdata.GetEvaluationResponseAcceptedByWhiteListSuccessful()
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, charge)
		assert.Nil(t, err)

		ID := get(ctx, mongoDB, cfg, charge.Charge.ID)

		assert.Nil(t, err)
		assert.NotNil(t, ID)
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.ChargeEvaluations, ID)
	})

	t.Run("when get charge unsuccessful", func(t *testing.T) {

		charge := testdata.GetEvaluationResponseAcceptedByWhiteListSuccessful()
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		response, err := repository.Get(ctx, charge.Charge.ID)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("error, charge_evaluation %s not found", charge.Charge.ID))
		assert.Empty(t, response.Decision)
	})

	t.Run("when get charge unsuccessful invalid collection", func(t *testing.T) {

		charge := testdata.GetEvaluationResponseAcceptedByWhiteListSuccessful()
		cfg.MongoDB.Collections.ChargeEvaluations = ""
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		response, err := repository.Get(ctx, charge.Charge.ID)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("(InvalidNamespace) Invalid namespace specified '%s.'", cfg.MongoDB.Database))
		assert.Empty(t, response.Decision)
	})
}

func TestChargeRepository_GetOnlyRules(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when get charge successful", func(t *testing.T) {

		charge := testdata.GetRulesEvaluationResponseAcceptedCauseConsoleCompanyRules()
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.SaveOnlyRules(ctx, charge)
		assert.Nil(t, err)

		ID := getOnlyRules(ctx, mongoDB, cfg, charge.Charge.ID)

		assert.Nil(t, err)
		assert.NotNil(t, ID)
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.ChargeEvaluationsOnlyRules, ID)
	})

	t.Run("when get charge unsuccessful", func(t *testing.T) {

		charge := testdata.GetRulesEvaluationResponseAcceptedCauseConsoleCompanyRules()
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		response, err := repository.GetOnlyRules(ctx, charge.Charge.ID)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("error, charge_evaluation %s not found", charge.Charge.ID))
		assert.Empty(t, response.Decision)
	})

	t.Run("when get charge unsuccessful invalid collection", func(t *testing.T) {
		charge := testdata.GetRulesEvaluationResponseAcceptedCauseConsoleCompanyRules()
		cfg.MongoDB.Collections.ChargeEvaluationsOnlyRules = ""
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		response, err := repository.GetOnlyRules(ctx, charge.Charge.ID)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("(InvalidNamespace) Invalid namespace specified '%s.'", cfg.MongoDB.Database))
		assert.Empty(t, response.Decision)
	})
}

func get(ctx context.Context, mongoDB mongodb.MongoDBier, cfg config.Config, id string) primitive.ObjectID {
	var result EvaluationDelete
	moduleCollection := mongoDB.Collection(cfg.MongoDB.Collections.ChargeEvaluations)
	filter := bson.D{primitive.E{Key: "charge._id", Value: id}}

	moduleCollection.FindOne(ctx, filter).Decode(&result)
	return result.ID
}

func getOnlyRules(ctx context.Context, mongoDB mongodb.MongoDBier, cfg config.Config, id string) primitive.ObjectID {
	var result EvaluationDelete
	moduleCollection := mongoDB.Collection(cfg.MongoDB.Collections.ChargeEvaluationsOnlyRules)
	filter := bson.D{primitive.E{Key: "charge._id", Value: id}}

	moduleCollection.FindOne(ctx, filter).Decode(&result)
	return result.ID
}

func TestSave(t *testing.T) {

	logger, _ := logs.New()

	t.Run("error saving a evaluation", func(t *testing.T) {
		cfg := config.NewConfig()
		cfg.MongoDB.Collections.ChargeEvaluations = ""
		mongoDB := mongodb.NewMongoDB(cfg)
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)

		err := repository.Save(context.TODO(), testdata.GetEvaluationResponseSuccessful())

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("(InvalidNamespace) Invalid namespace specified '%s.'", cfg.MongoDB.Database))
	})
	t.Run("successfully", func(t *testing.T) {
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		repository := charges.NewChargeMongoDBRepository(cfg, mongoDB, logger)

		err := repository.Save(context.TODO(), testdata.GetEvaluationResponseSuccessful())

		assert.NoError(t, err)
	})
}
