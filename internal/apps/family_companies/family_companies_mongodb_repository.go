package familycom

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

const repositoryName = "familycom.Repository.mongo"

type familyCompaniesMongoDBRepository struct {
	config  config.Config
	mongodb mongodb.MongoDBier
	Logger  logs.Logger
}

type FamilyCompaniesRepository interface {
	AddFamilyCompanies(ctx context.Context, family *entities.FamilyCompanies) error
	GetFamilyCompanies(ctx context.Context, filter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error)
	Update(ctx context.Context, id string, familyCompanies *entities.FamilyCompanies) error
	Delete(ctx context.Context, id string) error
	SearchPaged(ctx context.Context,
		pag entities.Pagination,
		fil entities.FamilyCompaniesFilter) (entities.PagedResponse, error)
	Search(ctx context.Context, filter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error)
}

func NewFamilyCompaniesMongoDBRepository(cfg config.Config,
	mongoDBier mongodb.MongoDBier,
	logger logs.Logger) FamilyCompaniesRepository {
	return &familyCompaniesMongoDBRepository{
		config:  cfg,
		mongodb: mongoDBier,
		Logger:  logger,
	}
}

func (repository *familyCompaniesMongoDBRepository) AddFamilyCompanies(
	ctx context.Context,
	familyCompanies *entities.FamilyCompanies) error {
	result, err := repository.mongodb.
		Collection(repository.config.MongoDB.Collections.FamilyCompanies).
		InsertOne(ctx, familyCompanies)

	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod,
			fmt.Sprintf("%s.%s", repositoryName, "AddFamilyCompanies"), text.Family, familyCompanies)
		return err
	}

	familyCompanies.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (repository *familyCompaniesMongoDBRepository) GetFamilyCompanies(
	ctx context.Context,
	filter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error) {
	familyCompanies := make([]entities.FamilyCompanies, 0)
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.FamilyCompanies)
	opts := options.FindOptions{}
	opts.Collation = &options.Collation{
		Locale:   "en",
		Strength: 2,
	}

	query, err := getQueryByFilter(filter)
	if err != nil {
		repository.Logger.Error(ctx,
			err.Error(),
			text.LogTagMethod,
			fmt.Sprintf("%s.%s", repositoryName, "GetFamilyCompanies"))
		return nil, err
	}

	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.Logger.Error(ctx,
			err.Error(),
			text.LogTagMethod,
			fmt.Sprintf("%s.%s", repositoryName, "GetFamilyCompanies"))
		return nil, err
	}

	err = cursor.All(ctx, &familyCompanies)
	if err != nil {
		repository.Logger.Error(ctx,
			err.Error(),
			text.LogTagMethod,
			fmt.Sprintf("%s.%s", repositoryName, "GetFamilyCompanies"))
		return nil, err
	}

	return familyCompanies, nil
}

func (repository *familyCompaniesMongoDBRepository) SearchPaged(ctx context.Context,
	pag entities.Pagination,
	fil entities.FamilyCompaniesFilter) (entities.PagedResponse, error) {
	familiesCompanies := make([]entities.FamilyCompanies, 0)
	emptyPagedResponse := entities.PagedResponse{Data: familiesCompanies}
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.FamilyCompanies)
	opts := options.FindOptions{}
	opts.Collation = &options.Collation{
		Locale: "en",
	}

	query, err := getQueryByFilter(fil)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "SearchPaged"))
		return emptyPagedResponse, err
	}

	total, _ := collection.CountDocuments(ctx, query)
	hasMore := pag.HasMorePages(total)

	opts.SetLimit(pag.PageSize)
	opts.SetSkip(pag.GetPageStartIndex())
	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cursor, err := collection.Find(ctx, query, &opts)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "SearchPaged"))
		return emptyPagedResponse, err
	}

	err = cursor.All(ctx, &familiesCompanies)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "SearchPaged"))
		return emptyPagedResponse, err
	}

	return entities.NewPagedResponse(familiesCompanies, hasMore, total), ctx.Err()
}

func (repository *familyCompaniesMongoDBRepository) Search(ctx context.Context,
	filter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error) {
	familyCompanies := make([]entities.FamilyCompanies, 0)
	collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.FamilyCompanies)
	opt := options.FindOptions{}
	opt.Collation = &options.Collation{
		Locale: "en",
	}

	query, err := getQueryByFilter(filter)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Search"))
		return nil, err
	}

	opt.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cur, err := collection.Find(ctx, query, &opt)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Search"))
		return nil, err
	}

	err = cur.All(ctx, &familyCompanies)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Search"))
		return nil, err
	}

	return familyCompanies, nil
}

func (repository *familyCompaniesMongoDBRepository) Delete(ctx context.Context, id string) error {
	familyCompaniesID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.Logger.Error(ctx,
			err.Error(),
			text.LogTagMethod,
			fmt.Sprintf("%s.%s", repositoryName, "Delete"),
			text.FamilyCompanyID, id)
		return err
	}

	collections := repository.mongodb.Collection(repository.config.MongoDB.Collections.FamilyCompanies)
	filter := bson.M{"_id": familyCompaniesID}
	result, err := collections.DeleteOne(ctx, filter, nil)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Delete"))
	}

	if result.DeletedCount == 0 {
		err = exceptions.NewNotFoundException(fmt.Sprintf("error: family companies not found: '%s'", familyCompaniesID))
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Delete"))
		return err
	}

	return nil
}

func getQueryByFilter(filter entities.FamilyCompaniesFilter) (bson.M, error) {
	query := []bson.M{}
	findQuery := bson.M{}

	if filter.IsEmpty() {
		return findQuery, nil
	}

	if !strings.IsEmpty(filter.ID) {
		ID, err := primitive.ObjectIDFromHex(filter.ID)
		if err != nil {
			err = http.NewBadRequestError(fmt.Sprintf("error: invalid id of family_company: '%s'", filter.ID))
			return nil, err
		}
		query = append(query, bson.M{"_id": ID})
	}

	if filter.CompanyIDs != nil {
		query = append(query, bson.M{"company_ids": bson.M{"$in": filter.CompanyIDs}})
	}

	if !strings.IsEmpty(filter.Name) {
		query = append(query, bson.M{"name": filter.Name})
	}

	findQuery["$or"] = query
	return findQuery, nil
}

func (repository *familyCompaniesMongoDBRepository) Update(
	ctx context.Context,
	id string, familyCompanies *entities.FamilyCompanies) error {
	familyCompaniesID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repository.Logger.Error(ctx,
			err.Error(),
			text.LogTagMethod,
			fmt.Sprintf("%s.%s", repositoryName, "Update"),
			text.FamilyID,
			id)
		return err
	}

	Collection := repository.mongodb.Collection(repository.config.MongoDB.Collections.FamilyCompanies)
	filter := bson.M{"_id": familyCompaniesID}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "name", Value: familyCompanies.Name},
				primitive.E{Key: "company_ids", Value: familyCompanies.CompanyIDs},
				primitive.E{Key: "updated_at", Value: familyCompanies.UpdatedAt},
				primitive.E{Key: "updated_by", Value: familyCompanies.UpdatedBy},
			},
		},
	}
	result, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Update"))
		return err
	}

	if result.MatchedCount == 0 {
		err = fmt.Errorf("error, family companies document '%s' not found", familyCompaniesID)
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", repositoryName, "Update"))
		return err
	}

	return nil
}
