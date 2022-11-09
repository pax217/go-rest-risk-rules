package modules

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

const repositoryName = "module.repository.mongo.%s"

type ModuleRepository interface {
	Get(ctx context.Context, filter entities.ModuleFilter) ([]entities.Module, error)
	Save(ctx context.Context, module *entities.Module) error
	GetPaged(ctx context.Context, filter entities.ModuleFilter, pagination entities.Pagination) (entities.PagedResponse, error)
	Update(ctx context.Context, id string, module entities.Module) error
	Delete(ctx context.Context, id string) error
}
type modulesMongoRepository struct {
	mongodb mongodb.MongoDBier
	config  config.Config
	logs    logs.Logger
}

func NewModulesMongoRepository(cfg config.Config, dBier mongodb.MongoDBier, logger logs.Logger) ModuleRepository {
	return &modulesMongoRepository{
		mongodb: dBier,
		config:  cfg,
		logs:    logger,
	}
}

func (repository *modulesMongoRepository) Save(ctx context.Context, module *entities.Module) error {
	result, err := repository.mongodb.Collection(repository.config.MongoDB.Collections.Modules).InsertOne(ctx, module, nil)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Save"))
		return err
	}

	module.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (repository *modulesMongoRepository) Update(ctx context.Context, id string, module entities.Module) error {
	moduleID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = http.NewBadRequestError(fmt.Sprintf("error, invalid id in modules: '%s'", id))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "UpdateModule"), text.ModuleID, id)
		return err
	}

	Collection := repository.mongodb.Collection(config.Configs.MongoDB.Collections.Modules)
	filter := bson.M{"_id": moduleID}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "updated_at", Value: module.UpdatedAt},
				primitive.E{Key: "updated_by", Value: module.UpdatedBy},
				primitive.E{Key: "name", Value: module.Name},
				primitive.E{Key: "description", Value: module.Description},
			},
		},
	}
	result, err := Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "UpdateModule"))
		return err
	}

	if result.MatchedCount == 0 {
		err = fmt.Errorf("error, document '%s' not found", moduleID)
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "UpdateModule"))
		return err
	}

	return nil
}

func (repository *modulesMongoRepository) Delete(ctx context.Context, id string) error {
	moduleID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = http.NewBadRequestError(fmt.Sprintf("error, invalid id '%s' in modules", id))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "DeleteModule"), text.ModuleID, id)
		return err
	}

	Collection := repository.mongodb.Collection(config.Configs.MongoDB.Collections.Modules)
	filter := bson.M{"_id": moduleID}

	result, err := Collection.DeleteOne(ctx, filter)

	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "DeleteModule"))
		return err
	}

	if result.DeletedCount == 0 {
		err = http.NewNotFoundError(fmt.Sprintf("error: record not found in modules: '%s'", moduleID))
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "DeleteModule"), text.ModuleID, id)
		return err
	}

	return nil
}

func (repository *modulesMongoRepository) Get(ctx context.Context, filter entities.ModuleFilter) ([]entities.Module, error) {
	modules := make([]entities.Module, 0)
	moduleCollection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Modules)
	opts := options.FindOptions{}

	query, err := repository.getQueryByFilter(filter)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Get"))
		return nil, err
	}

	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cursor, err := moduleCollection.Find(ctx, query, &opts)
	if err != nil && (err.Error() != mongodb.NoResultsOnFind) {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Get"))
		return nil, err
	}

	err = cursor.All(ctx, &modules)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "Get"))
		return nil, err
	}

	return modules, nil
}

func (repository *modulesMongoRepository) GetPaged(ctx context.Context, filter entities.ModuleFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	modules := make([]entities.Module, 0)
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Modules)
	opts := options.FindOptions{}

	query, err := repository.getQueryByFilter(filter)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetPaged"))
		return entities.PagedResponse{}, err
	}

	total, _ := collection.CountDocuments(ctx, query)
	hasMore := pagination.HasMorePages(total)

	opts.SetLimit(pagination.PageSize)
	opts.SetSkip(pagination.GetPageStartIndex())
	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetPaged"))
		return entities.PagedResponse{}, err
	}

	err = cursor.All(ctx, &modules)
	if err != nil {
		repository.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(repositoryName, "GetPaged"))
		return entities.PagedResponse{}, err
	}

	return entities.NewPagedResponse(modules, hasMore, total), ctx.Err()
}

func (repository *modulesMongoRepository) getQueryByFilter(filter entities.ModuleFilter) (bson.M, error) {
	query := []bson.M{}
	findQuery := bson.M{}

	if filter.IsEmptyModuleFilter() {
		return findQuery, nil
	}

	if !strings.IsEmpty(filter.ID) {
		ID, err := primitive.ObjectIDFromHex(filter.ID)
		if err != nil {
			err = http.NewBadRequestError(fmt.Sprintf("error: invalid id in modules: '%s'", filter.ID))
			return nil, err
		}
		query = append(query, bson.M{"_id": ID})
	}

	if !strings.IsEmpty(filter.Name) {
		query = append(query, bson.M{"name": filter.Name})
	}

	findQuery["$or"] = query
	return findQuery, nil
}
