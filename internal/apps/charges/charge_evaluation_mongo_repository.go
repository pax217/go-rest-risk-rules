package charges

import (
	"context"
	"fmt"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/pkg/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const repositoryMethodName = "charge_evaluation.repository.mongo.%s"

type chargeEvaluationMongoDBRepository struct {
	logs    logs.Logger
	mongodb mongodb.MongoDBier
	config  config.Config
}

func NewChargeMongoDBRepository(conf config.Config, db mongodb.MongoDBier, logger logs.Logger) ChargeRepository {
	return &chargeEvaluationMongoDBRepository{
		logs:    logger,
		mongodb: db,
		config:  conf,
	}
}

func (repository *chargeEvaluationMongoDBRepository) Save(ctx context.Context,
	evaluation entities.EvaluationResponse) error {
	_, err := repository.mongodb.Collection(repository.config.MongoDB.Collections.ChargeEvaluations).
		InsertOne(ctx, evaluation, nil)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryMethodName, "Save"))
		return err
	}
	return nil
}

func (repository *chargeEvaluationMongoDBRepository) SaveOnlyRules(ctx context.Context,
	rulesEvaluationResponse entities.RulesEvaluationResponse) error {
	_, err := repository.mongodb.Collection(repository.config.MongoDB.Collections.ChargeEvaluationsOnlyRules).
		InsertOne(ctx, rulesEvaluationResponse, nil)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryMethodName, "SaveOnlyRules"))
		return err
	}
	return nil
}

func (repository *chargeEvaluationMongoDBRepository) Get(ctx context.Context, id string) (entities.EvaluationResponse, error) {
	var result entities.EvaluationResponse

	chargeEvaluationsCollection := repository.mongodb.Collection(repository.config.MongoDB.Collections.ChargeEvaluations)
	filter := bson.D{primitive.E{Key: "charge._id", Value: id}}

	err := chargeEvaluationsCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err.Error() == mongodb.NoResultsOnFind {
			err = exceptions.NewNotFoundException(fmt.Sprintf("error, charge_evaluation %s not found", id))
			repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryMethodName, "Get"))
			return entities.EvaluationResponse{}, err
		} else {
			repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryMethodName, "Get"))
			return entities.EvaluationResponse{}, err
		}
	}

	return result, nil
}

func (repository *chargeEvaluationMongoDBRepository) GetOnlyRules(ctx context.Context,
	id string) (entities.RulesEvaluationResponse, error) {
	var result entities.RulesEvaluationResponse

	chargeEvaluationsOnlyRulesCollection := repository.
		mongodb.Collection(repository.config.MongoDB.Collections.ChargeEvaluationsOnlyRules)
	filter := bson.D{primitive.E{Key: "charge._id", Value: id}}

	err := chargeEvaluationsOnlyRulesCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err.Error() == mongodb.NoResultsOnFind {
			err = exceptions.NewNotFoundException(fmt.Sprintf("error, charge_evaluation %s not found", id))
			repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryMethodName, "GetOnlyRules"))
			return entities.RulesEvaluationResponse{}, err
		} else {
			repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryMethodName, "GetOnlyRules"))
			return entities.RulesEvaluationResponse{}, err
		}
	}

	return result, nil
}
