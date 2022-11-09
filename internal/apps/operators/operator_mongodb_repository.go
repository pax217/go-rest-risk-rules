package operators

import (
	"context"
	"fmt"

	http "github.com/conekta/go_common/http/resterror"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/pkg/text"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const repositoryName = "operator.repository.mongo.%s"

type OperatorRepository interface {
	Save(ctx context.Context, operator *entities.Operator) error
	Get(ctx context.Context, operatorFilter entities.OperatorFilter) ([]entities.Operator, error)
	GetPaged(ctx context.Context, operatorFilter entities.OperatorFilter,
		pagination entities.Pagination) (entities.PagedResponse, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, operator entities.Operator) error
}
type operatorMongoDBRepository struct {
	logs    logs.Logger
	config  config.Config
	mongodb mongodb.MongoDBier
}

func NewOperatorMongoDBRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logs.Logger) OperatorRepository {
	return &operatorMongoDBRepository{
		config:  cfg,
		mongodb: mongoDBier,
		logs:    logger,
	}
}

func (repository *operatorMongoDBRepository) Save(ctx context.Context, operator *entities.Operator) error {
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Operators)
	result, err := collection.InsertOne(ctx, operator)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "AddList"))
		return err
	}

	operator.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (repository *operatorMongoDBRepository) Update(ctx context.Context, id string, operator entities.Operator) error {
	operatorID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Update"))
		return err
	}

	Collection := repository.mongodb.Collection(config.Configs.MongoDB.Collections.Operators)
	filter := bson.M{"_id": operatorID}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "updated_at", Value: operator.UpdatedAt},
				primitive.E{Key: "updated_by", Value: operator.UpdatedBy},
				primitive.E{Key: "name", Value: operator.Name},
				primitive.E{Key: "title", Value: operator.Title},
				primitive.E{Key: "description", Value: operator.Description},
				primitive.E{Key: "type", Value: operator.Type},
			},
		},
	}
	result, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Update"))
		return err
	}

	if result.MatchedCount == 0 {
		err = http.NewNotFoundError(fmt.Sprintf("error: record not found: %s", operatorID))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Update"))
		return err
	}

	return nil
}

func (repository *operatorMongoDBRepository) Delete(ctx context.Context, id string) error {
	operatorID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = http.NewBadRequestError(fmt.Sprintf("error: invalid id: %s", id))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Delete"))
		return err
	}

	Collection := repository.mongodb.Collection(config.Configs.MongoDB.Collections.Operators)
	filter := bson.M{"_id": operatorID}

	result, err := Collection.DeleteOne(ctx, filter)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Delete"))
		return err
	}

	if result.DeletedCount == 0 {
		err = http.NewNotFoundError(fmt.Sprintf("error: record not found: %s", operatorID))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Delete"))
		return err
	}

	return nil
}

func (repository *operatorMongoDBRepository) Get(ctx context.Context,
	filter entities.OperatorFilter) ([]entities.Operator, error) {
	operators := make([]entities.Operator, 0)
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Operators)
	opts := options.FindOptions{}

	query, err := repository.getQueryByFilter(filter)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Get"))
		return nil, err
	}

	opts.SetSort(bson.D{primitive.E{Key: "title", Value: 1}})
	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Get"))
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &operators)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Get"))
		return nil, err
	}

	return operators, nil
}

func (repository *operatorMongoDBRepository) GetPaged(ctx context.Context, filter entities.OperatorFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	operators := make([]entities.Operator, 0)
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Operators)
	opts := options.FindOptions{}
	opts.SetLimit(pagination.PageSize)
	opts.SetSkip(pagination.GetPageStartIndex())
	opts.SetSort(bson.D{primitive.E{Key: "title", Value: 1}})

	query, err := repository.getQueryByFilter(filter)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetPaged"))
		return entities.PagedResponse{}, err
	}

	total, _ := collection.CountDocuments(ctx, query)
	hasMore := pagination.HasMorePages(total)

	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetPaged"))
		return entities.PagedResponse{}, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &operators)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetPaged"))
		return entities.PagedResponse{}, err
	}

	return entities.NewPagedResponse(operators, hasMore, total), ctx.Err()
}

func (repository *operatorMongoDBRepository) getQueryByFilter(filter entities.OperatorFilter) (bson.M, error) {
	query := bson.M{}
	var ID primitive.ObjectID
	var err error

	if filter.IsEmptyOperatorFilter() {
		return query, nil
	}

	if filter.ID != "" {
		ID, err = primitive.ObjectIDFromHex(filter.ID)
		if err != nil {
			err = http.NewBadRequestError(fmt.Sprintf("error: invalid id: %s", filter.ID))
			return query, err
		}
	}

	query = bson.M{
		"$or": bson.A{
			bson.M{"_id": ID},
			bson.M{"type": filter.Type},
			bson.M{"name": filter.Name},
			bson.M{"title": filter.Title},
		},
	}

	return query, nil
}
