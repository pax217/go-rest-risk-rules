package chargebacks

import (
	"context"
	"errors"
	"fmt"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/pkg/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	repositoryName        = "chargeback.repository"
	documentNotFoundError = "error, document not found"
)

type ChargebackRepository interface {
	Save(ctx context.Context, payer entities.Payer) error
	Update(ctx context.Context, payer entities.Payer) error
	Find(ctx context.Context, payer entities.Payer) (entities.Payer, error)
}
type chargebackRepository struct {
	logs    logs.Logger
	mongodb mongodb.MongoDBier
	config  config.Config
}

func NewChargebacksMongoDBRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logs.Logger) ChargebackRepository {
	return &chargebackRepository{
		logs:    logger,
		mongodb: mongoDBier,
		config:  cfg,
	}
}

func (repository *chargebackRepository) Save(ctx context.Context, payer entities.Payer) error {
	_, err := repository.mongodb.Collection(repository.config.MongoDB.Collections.Payers).InsertOne(ctx, payer)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Save"),
			text.Email, payer.Email, text.Payer, payer)
		return err
	}

	return nil
}

func (repository *chargebackRepository) Update(ctx context.Context, payer entities.Payer) error {
	_id, err := primitive.ObjectIDFromHex(payer.ID.Hex())
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Update"))
		return err
	}

	rulesCollection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Payers)
	filter := bson.M{"_id": _id}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "chargebacks", Value: payer.Chargebacks},
			},
		},
	}

	result, err := rulesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Update"))
		return err
	}

	if result.MatchedCount == 0 {
		err := errors.New(documentNotFoundError)
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Update"),
			text.PayerID, payer.ID, text.Payer, payer)
		return err
	}

	return nil
}

func (repository *chargebackRepository) Find(ctx context.Context, payer entities.Payer) (entities.Payer, error) {
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Payers)
	payerFound := entities.Payer{}
	query := bson.M{"email": payer.Email}

	err := collection.FindOne(ctx, query).Decode(&payerFound)
	if err != nil && (err.Error() != mongodb.NoResultsOnFind) {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Find"),
			text.Email, payer.Email, text.Payer, payer)
		return payerFound, err
	}

	return payerFound, nil
}
