package merchantsscore

import (
	"context"
	"fmt"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/pkg/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const repositoryName = "merchant.repository.mongo.%s"

type MerchantsScoreRepository interface {
	WriteMerchantsScore(ctx context.Context, merchant []entities.MerchantScore) error
	FindByMerchantID(ctx context.Context, companyID string) (entities.MerchantScore, error)
}

type merchantsScoreMongoDBRepository struct {
	config  config.Config
	mongodb mongodb.MongoDBier
	logs    logs.Logger
}

func NewMerchantsMongoDBRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logs.Logger) MerchantsScoreRepository {
	return &merchantsScoreMongoDBRepository{
		config:  cfg,
		mongodb: mongoDBier,
		logs:    logger,
	}
}

func (repository *merchantsScoreMongoDBRepository) WriteMerchantsScore(ctx context.Context, merchant []entities.MerchantScore) error {
	merchantCollection := repository.mongodb.Collection(repository.config.MongoDB.Collections.MerchantsScore)

	for _, row := range merchant {
		filter := bson.M{"company_id": row.CompanyID}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "score", Value: row.Score}}}}
		opts := options.Update().SetUpsert(true)

		_, err := merchantCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "WriteMerchantsScore"))
			return err
		}
	}

	return nil
}

func (repository *merchantsScoreMongoDBRepository) FindByMerchantID(ctx context.Context,
	companyID string) (entities.MerchantScore, error) {
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.MerchantsScore)
	merchantScore := entities.MerchantScore{}
	query := bson.M{"company_id": companyID}

	err := collection.FindOne(ctx, query).Decode(&merchantScore)
	if err != nil && (err.Error() != mongodb.NoResultsOnFind) {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "FindByMerchantID"))
		return merchantScore, err
	}

	return merchantScore, nil
}
