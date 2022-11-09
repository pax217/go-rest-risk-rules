package conditions

import (
	"context"
	"fmt"

	http "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const repositoryName = "condition.repository.mongo.%s"

type conditionMongoRepository struct {
	logs          logs.Logger
	mongodb       mongodb.MongoDBier
	configuration config.Config
}

func NewConditionsRepository(c config.Config, dBier mongodb.MongoDBier, logger logs.Logger) ConditionRepository {
	return &conditionMongoRepository{
		logs:          logger,
		mongodb:       dBier,
		configuration: c,
	}
}

func (repository *conditionMongoRepository) Add(ctx context.Context, condition *entities.Condition) error {
	moduleCollection := repository.mongodb.Collection(repository.configuration.MongoDB.Collections.Conditions)

	result, err := moduleCollection.InsertOne(ctx, condition)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Add"))
		return err
	}

	condition.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (repository *conditionMongoRepository) FindByName(ctx context.Context, condition entities.Condition) (
	entities.Condition, error) {
	var result entities.Condition
	moduleCollection := repository.mongodb.Collection(repository.configuration.MongoDB.Collections.Conditions)
	filter := bson.D{primitive.E{Key: "name", Value: condition.Name}}

	err := moduleCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil && (err.Error() != mongodb.NoResultsOnFind) {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "FindByName"))
		return entities.Condition{}, err
	}

	return result, nil
}

func (repository *conditionMongoRepository) GetAll(
	ctx context.Context,
	conditionsFilter entities.ConditionsFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	conditions := make([]entities.Condition, 0)
	searchOptions := options.FindOptions{}
	searchOptions.Collation = &options.Collation{
		Locale: "en",
	}
	collection := repository.mongodb.Collection(repository.configuration.MongoDB.Collections.Conditions)
	searchOptions.SetLimit(pagination.PageSize)
	searchOptions.SetSkip(pagination.GetPageStartIndex())
	filter := bson.M{}

	if !strings.IsEmpty(conditionsFilter.ID) {
		oID, err := primitive.ObjectIDFromHex(conditionsFilter.ID)
		if err != nil {
			repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetAll"))
			return entities.NewPagedResponse(conditions, false, 0), nil
		}
		filter = bson.M{"_id": oID}
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetAll"))
		return entities.PagedResponse{}, err
	}

	hasMore := pagination.HasMorePages(total)
	searchOptions.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cursor, err := collection.Find(ctx, filter, &searchOptions)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetAll"))
		return entities.PagedResponse{}, err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, &conditions)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetAll"))
		return entities.PagedResponse{}, err
	}

	return entities.NewPagedResponse(conditions, hasMore, total), ctx.Err()
}

func (repository *conditionMongoRepository) Update(ctx context.Context, id string, condition entities.Condition) error {
	conditionID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Update"))
		return err
	}

	Collection := repository.mongodb.Collection(config.Configs.MongoDB.Collections.Conditions)
	filter := bson.M{"_id": conditionID}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "updated_at", Value: condition.UpdatedAt},
				primitive.E{Key: "updated_by", Value: condition.UpdatedBy},
				primitive.E{Key: "name", Value: condition.Name},
				primitive.E{Key: "description", Value: condition.Description},
			},
		},
	}

	result, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Update"))
		return err
	}

	if result.MatchedCount == 0 {
		err = fmt.Errorf("error, document %s not found", conditionID)
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Update"))
		return err
	}

	return nil
}

func (repository *conditionMongoRepository) Delete(ctx context.Context, id string) error {
	conditionID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Delete"))
		return err
	}

	Collection := repository.mongodb.Collection(config.Configs.MongoDB.Collections.Conditions)
	filter := bson.M{"_id": conditionID}

	result, err := Collection.DeleteOne(ctx, filter)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Delete"))
		return err
	}

	if result.DeletedCount == 0 {
		err = http.NewNotFoundError(fmt.Sprintf("error: record not found: %s", conditionID))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Delete"))
		return err
	}

	return nil
}
