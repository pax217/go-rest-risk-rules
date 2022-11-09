package fields

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

const repositoryName = "field.repository.mongo.%s"

type FieldsRepository interface {
	AddField(ctx context.Context, request *entities.Field) error
	Update(ctx context.Context, id string, field *entities.Field) error
	Delete(ctx context.Context, id string) error
	GetFieldsPaged(ctx context.Context, filter entities.FieldsFilter, pagination entities.Pagination) (entities.PagedResponse, error)
	GetFields(ctx context.Context, filter entities.FieldsFilter) ([]entities.Field, error)
}

type fieldsMongoDBRepository struct {
	config  config.Config
	mongodb mongodb.MongoDBier
	logs    logs.Logger
}

func NewFieldsMongoDBRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logs.Logger) FieldsRepository {
	return &fieldsMongoDBRepository{
		config:  cfg,
		mongodb: mongoDBier,
		logs:    logger,
	}
}

func (repository *fieldsMongoDBRepository) AddField(ctx context.Context, field *entities.Field) error {
	result, err := repository.mongodb.Collection(repository.config.MongoDB.Collections.Fields).InsertOne(ctx, field, nil)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "AddField"))
		return err
	}
	field.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (repository *fieldsMongoDBRepository) Update(
	ctx context.Context,
	id string,
	field *entities.Field) error {
	fieldID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "UpdateField"))
		return err
	}

	Collection := repository.mongodb.Collection(config.Configs.MongoDB.Collections.Fields)
	filter := bson.M{"_id": fieldID}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "updated_at", Value: field.UpdatedAt},
				primitive.E{Key: "updated_by", Value: field.UpdatedBy},
				primitive.E{Key: "name", Value: field.Name},
				primitive.E{Key: "description", Value: field.Description},
				primitive.E{Key: "type", Value: field.Type},
			},
		},
	}
	result, err := Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "UpdateField"))
		return err
	}

	if result.MatchedCount == 0 {
		err = http.NewNotFoundError(fmt.Sprintf("record not found: %s", fieldID))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "UpdateField"))
		return err
	}

	field.ID = fieldID

	return nil
}

func (repository *fieldsMongoDBRepository) Delete(ctx context.Context, id string) error {
	fieldID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "DeleteField"))
		return err
	}

	rulesCollection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Fields)
	filter := bson.M{"_id": fieldID}
	result, err := rulesCollection.DeleteOne(ctx, filter, nil)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "DeleteField"))
		return err
	}

	if result.DeletedCount == 0 {
		err = http.NewNotFoundError(fmt.Sprintf("record not found: %s", fieldID))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "DeleteField"))
		return err
	}

	return nil
}

func (repository *fieldsMongoDBRepository) GetFieldsPaged(ctx context.Context, fieldsFilter entities.FieldsFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	fields := make([]entities.Field, 0)
	opts := options.FindOptions{}
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Fields)

	opts.SetLimit(pagination.PageSize)
	opts.SetSkip(pagination.GetPageStartIndex())
	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})

	query := repository.getQueryByFilter(fieldsFilter)

	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetFieldsPaged"))
		return entities.PagedResponse{}, err
	}

	hasMore := pagination.HasMorePages(total)

	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetFieldsPaged"))
		return entities.PagedResponse{}, err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, &fields)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetFieldsPaged"))
		return entities.PagedResponse{}, err
	}

	return entities.NewPagedResponse(fields, hasMore, total), ctx.Err()
}

func (repository *fieldsMongoDBRepository) GetFields(ctx context.Context,
	fieldsFilter entities.FieldsFilter) ([]entities.Field, error) {
	fields := make([]entities.Field, 0)
	opts := options.FindOptions{}
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Fields)

	query := repository.getQueryByFilter(fieldsFilter)

	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetFields"))
		return nil, err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, &fields)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetFields"))
		return nil, err
	}

	return fields, nil
}

func (repository *fieldsMongoDBRepository) getQueryByFilter(filter entities.FieldsFilter) bson.M {
	query := bson.M{}
	var ID primitive.ObjectID

	if filter.IsEmptyFieldsFilter() {
		return query
	}

	if filter.ID != "" {
		ID, _ = primitive.ObjectIDFromHex(filter.ID)
	}

	query = bson.M{
		"$or": bson.A{
			bson.M{"_id": ID},
			bson.M{"name": filter.Name},
			bson.M{"type": filter.Type},
		},
	}

	return query
}
