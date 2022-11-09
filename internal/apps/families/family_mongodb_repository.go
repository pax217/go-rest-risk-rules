package families

import (
	"context"
	"fmt"

	http "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	repositoryName = "family.familyRepository.mongo"
)

type familyMongoDBRepository struct {
	config  config.Config
	mongodb mongodb.MongoDBier
	Logger  logs.Logger
}

type FamilyRepository interface {
	AddFamily(ctx context.Context, family *entities.Family) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, family *entities.Family) error
	SearchPaged(ctx context.Context, pagination entities.Pagination,
		filter entities.FamilyFilter) (entities.PagedResponse, error)
	Search(ctx context.Context, filter entities.FamilyFilter) ([]entities.Family, error)
	SearchEvaluate(ctx context.Context, filter entities.FamilyFilter) ([]entities.Family, error)
}

func NewFamilyMongoDBRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logs.Logger) FamilyRepository {
	return &familyMongoDBRepository{
		config:  cfg,
		mongodb: mongoDBier,
		Logger:  logger,
	}
}

func (repository *familyMongoDBRepository) AddFamily(ctx context.Context, family *entities.Family) error {
	result, err := repository.mongodb.
		Collection(repository.config.MongoDB.Collections.Families).
		InsertOne(ctx, family)

	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod,
			fmt.Sprintf("%s.%s", repositoryName, "AddFamily"), text.Family, family)
		return err
	}

	family.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (repository *familyMongoDBRepository) Update(ctx context.Context, id string, family *entities.Family) error {
	familyID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Update"), text.FamilyID, id)
		return err
	}

	Collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Families)
	filter := bson.M{"_id": familyID}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "name", Value: family.Name},
				primitive.E{Key: "mccs", Value: family.Mccs},
				primitive.E{Key: "excluded_companies", Value: family.ExcludedCompanies},
				primitive.E{Key: "updated_at", Value: family.UpdatedAt},
				primitive.E{Key: "updated_by", Value: family.UpdatedBy},
			},
		},
	}
	result, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Update"))
		return err
	}

	if result.MatchedCount == 0 {
		err = fmt.Errorf("error, document '%s' not found", familyID)
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Update"))
		return err
	}

	return nil
}

func (repository *familyMongoDBRepository) Delete(ctx context.Context, id string) error {
	familyID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Delete"), text.FamilyID, id)
		return err
	}

	collections := repository.mongodb.Collection(repository.config.MongoDB.Collections.Families)
	filter := bson.M{"_id": familyID}
	result, err := collections.DeleteOne(ctx, filter, nil)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Delete"))
	}

	if result.DeletedCount == 0 {
		err = exceptions.NewNotFoundException(fmt.Sprintf("error: family not found: '%s'", familyID))
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Delete"))
		return err
	}

	return nil
}

func (repository *familyMongoDBRepository) SearchPaged(ctx context.Context, pagination entities.Pagination,
	filter entities.FamilyFilter) (entities.PagedResponse, error) {
	families := make([]entities.Family, 0)
	emptyPagedResponse := entities.PagedResponse{Data: families}
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Families)
	opts := options.FindOptions{}

	query, err := getQueryByFilter(filter)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "SearchPaged"))
		return emptyPagedResponse, err
	}

	total, _ := collection.CountDocuments(ctx, query)
	hasMore := pagination.HasMorePages(total)

	opts.SetLimit(pagination.PageSize)
	opts.SetSkip(pagination.GetPageStartIndex())
	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "SearchPaged"))
		return emptyPagedResponse, err
	}

	err = cursor.All(ctx, &families)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "SearchPaged"))
		return emptyPagedResponse, err
	}

	return entities.NewPagedResponse(families, hasMore, total), ctx.Err()
}

func (repository *familyMongoDBRepository) Search(ctx context.Context, filter entities.FamilyFilter) ([]entities.Family, error) {
	families := make([]entities.Family, 0)
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Families)
	opts := options.FindOptions{}

	query, err := getQueryByFilter(filter)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Search"))
		return nil, err
	}

	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Search"))
		return nil, err
	}

	err = cursor.All(ctx, &families)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Search"))
		return nil, err
	}

	return families, nil
}

func (repository *familyMongoDBRepository) SearchEvaluate(ctx context.Context, filter entities.FamilyFilter) (
	[]entities.Family, error) {
	families := make([]entities.Family, 0)
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.Families)
	opts := options.FindOptions{}
	var query []bson.M
	findQuery := bson.M{}

	if filter.Mccs != nil {
		query = append(query, bson.M{"mccs": bson.M{"$in": filter.Mccs}})
	}
	if filter.NotExcludedCompanies != nil {
		query = append(query, bson.M{"excluded_companies": bson.M{"$nin": filter.NotExcludedCompanies}})
	}
	findQuery["$and"] = query
	cursor, err := collection.Find(ctx, findQuery, &opts)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "SearchEvaluate"))
		return families, err
	}

	err = cursor.All(ctx, &families)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "SearchEvaluate"))
		return families, err
	}

	return families, nil
}

func getQueryByFilter(filter entities.FamilyFilter) (bson.M, error) {
	var query []bson.M
	findQuery := bson.M{}

	if filter.IsEmpty() {
		return findQuery, nil
	}

	if !strings.IsEmpty(filter.ID) {
		ID, err := primitive.ObjectIDFromHex(filter.ID)
		if err != nil {
			err = http.NewBadRequestError(fmt.Sprintf("error: invalid id of family: '%s'", filter.ID))
			return nil, err
		}
		query = append(query, bson.M{"_id": ID})
	}

	if filter.Mccs != nil {
		query = append(query, bson.M{"mccs": bson.M{"$in": filter.Mccs}})
	}

	if !strings.IsEmpty(filter.Name) {
		query = append(query, bson.M{"name": filter.Name})
	}

	findQuery["$or"] = query
	return findQuery, nil
}
