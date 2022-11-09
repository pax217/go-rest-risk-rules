package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/conekta/go_common/logs"
	familycom "github.com/conekta/risk-rules/internal/apps/family_companies"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFamilyCompaniesRepository_AddFamilyCompanies(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("on family_company_id duplicated then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		familyCompanies := testdata.GetDefaultFamilyCompanies()
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		expectedError := fmt.Sprintf(
			`write exception: write errors: [E11000 duplicate key error collection: rules.family_companies index: _id_ dup key: { _id: ObjectId('%s') }]`,
			familyCompanies.ID.Hex())
		defer cancel()

		err := repository.AddFamilyCompanies(ctx, &familyCompanies)

		assert.NoError(t, err)
		assert.NotEmpty(t, familyCompanies.ID)

		var familyCompaniesFound entities.FamilyCompanies
		mongoDB.Collection(cfg.MongoDB.Collections.FamilyCompanies).FindOne(ctx, bson.M{"_id": familyCompanies.ID}).Decode(&familyCompaniesFound)

		err = repository.AddFamilyCompanies(ctx, &familyCompaniesFound)

		mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.FamilyCompanies, familyCompaniesFound.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err.Error())
	})

	t.Run("on family_companies added is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		familyCompanies := testdata.GetDefaultFamilyCompanies()
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		err := repository.AddFamilyCompanies(ctx, &familyCompanies)

		assert.NoError(t, err)
		assert.NotEmpty(t, familyCompanies.ID)

		var familyCompaniesFound entities.FamilyCompanies
		mongoDB.Collection(cfg.MongoDB.Collections.FamilyCompanies).FindOne(ctx, bson.M{"_id": familyCompanies.ID}).Decode(&familyCompaniesFound)

		assert.NoError(t, err)
		assert.NotNil(t, familyCompanies.ID)
		assert.Equal(t, familyCompanies, familyCompaniesFound)

		mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.FamilyCompanies, familyCompaniesFound.ID)
	})
}

func TestFamilyCompaniesRepository_Update(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(configs)

	t.Run("when family companies id is incorrect, then return error", func(t *testing.T) {
		expectedError := "the provided hex string is not a valid ObjectID"
		request := testdata.GetFamilyCompaniesRequest()
		familyCompanies := request.NewFamilyCompaniesFromPutRequest()

		repository := familycom.NewFamilyCompaniesMongoDBRepository(configs, mongoDB, logger)

		err := repository.Update(context.TODO(), "", &familyCompanies)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), expectedError)
	})

	t.Run("when family companies id no exist, then return error", func(t *testing.T) {
		request := testdata.GetFamilyCompaniesRequest()
		family := request.NewFamilyCompaniesFromPutRequest()
		id := primitive.NewObjectID()
		repository := familycom.NewFamilyCompaniesMongoDBRepository(configs, mongoDB, logger)

		err := repository.Update(context.TODO(), id.Hex(), &family)

		assert.NotNil(t, err)
	})

	t.Run("when family comopanies update is success", func(t *testing.T) {
		repository := familycom.NewFamilyCompaniesMongoDBRepository(configs, mongoDB, logger)
		familyCompanies := testdata.GetDefaultFamilyCompanies()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.AddFamilyCompanies(ctx, &familyCompanies)
		assert.NoError(t, err)

		err = repository.Update(ctx, familyCompanies.ID.Hex(), &familyCompanies)
		var familyCompaniesUpdated entities.FamilyCompanies

		mongoDB.Collection(configs.MongoDB.Collections.FamilyCompanies).
			FindOne(ctx, bson.M{"_id": familyCompanies.ID}).Decode(&familyCompaniesUpdated)

		assert.NoError(t, err)
		assert.EqualValues(t, familyCompanies.ID, familyCompaniesUpdated.ID)
		assert.EqualValues(t, familyCompanies.Name, familyCompaniesUpdated.Name)
		assert.EqualValues(t, familyCompanies.CompanyIDs, familyCompaniesUpdated.CompanyIDs)
		assert.EqualValues(t, familyCompanies.CreatedAt, familyCompaniesUpdated.CreatedAt)
		assert.EqualValues(t, familyCompanies.CreatedBy, familyCompaniesUpdated.CreatedBy)
		assert.EqualValues(t, familyCompanies.UpdatedAt, familyCompaniesUpdated.UpdatedAt)
		assert.EqualValues(t, familyCompanies.UpdatedBy, familyCompaniesUpdated.UpdatedBy)

		defer mongoDB.CleanCollectionByIds(ctx, configs.MongoDB.Collections.FamilyCompanies, familyCompanies.ID)
	})
}

func TestFamilyCompaniesRepository_SearchPaged(t *testing.T) {
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when filter objectID is not valid", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		invalidId := "invalid-format-abc"
		expectedErr := fmt.Sprintf(
			"message: error: invalid id of family_company: '%s' - status: 400 - error: bad_request - causes: []",
			invalidId)

		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foundFamilies, err := repository.SearchPaged(ctx, entities.NewDefaultPagination(),
			entities.FamilyCompaniesFilter{
				ID: invalidId,
			})

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err.Error())
		assert.EqualValues(t, foundFamilies.Data, []entities.FamilyCompanies{})
	})

	t.Run("when family companies data is found", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		familyCompaniesData := testdata.GetFamilyCompanies()
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.FamilyCompanies)
		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.FamilyCompanies, familyCompaniesData[0], familyCompaniesData[1])

		families, err := repository.SearchPaged(ctx, entities.NewDefaultPagination(),
			entities.FamilyCompaniesFilter{
				CompanyIDs: []string{familyCompaniesData[0].CompanyIDs[0], familyCompaniesData[1].CompanyIDs[0]},
			})
		data := families.Data.([]entities.FamilyCompanies)

		assert.NoError(t, err)
		assert.Equalf(t, len(familyCompaniesData), len(data), "family companies len should be %d", len(familyCompaniesData))

		assert.Equal(t, familyCompaniesData[0], data[1])
		assert.Equal(t, familyCompaniesData[1], data[0])

		for i := range data {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.FamilyCompanies, data[i].ID)
		}
	})

	t.Run("when family companies data is found, it should be in alphabetical order", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		familyCompaniesData := testdata.GetFamilyCompaniesToAlphabeticalOrder()
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.FamilyCompanies)
		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.FamilyCompanies,
			familyCompaniesData[0], familyCompaniesData[1], familyCompaniesData[2], familyCompaniesData[3])

		families, err := repository.SearchPaged(ctx,
			entities.NewDefaultPagination(),
			entities.FamilyCompaniesFilter{})
		data := families.Data.([]entities.FamilyCompanies)

		assert.NoError(t, err)
		assert.Equalf(t, len(familyCompaniesData), len(data), "family companies len should be %d", len(familyCompaniesData))

		assert.Equal(t, familyCompaniesData[0], data[2])
		assert.Equal(t, familyCompaniesData[1], data[0])
		assert.Equal(t, familyCompaniesData[2], data[3])
		assert.Equal(t, familyCompaniesData[3], data[1])

		for i := range data {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.FamilyCompanies, data[i].ID)
		}
	})
}

func TestFamilyCompaniesRepository_Delete(t *testing.T) {

	t.Run("on family_company_id not found, then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		expectedError := errors.New(`error: family companies not found: 'ObjectID("61e9977bb5cce893d0aa9e45")'`)
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Delete(ctx, "61e9977bb5cce893d0aa9e45")

		assert.Error(t, err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("on family companies deleted is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		familyCompanies := testdata.GetDefaultFamilyCompanies()
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		familyCompaniesCollection := mongoDB.Collection(cfg.MongoDB.Collections.FamilyCompanies)
		familyCompaniesCollection.InsertOne(ctx, familyCompanies)

		assert.NotEmpty(t, familyCompanies.ID)

		err := repository.Delete(ctx, familyCompanies.ID.Hex())

		assert.NoError(t, err)

		defer mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.Rules)
	})

	t.Run("when ID is invalid, then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Delete(ctx, "x")

		assert.Error(t, err)
	})
}

func TestFamilyCompaniesRepository_GetFamilyCompanies(t *testing.T) {
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when searching by name value with either lower and upper cases characters, find all values successfully", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		familyCompaniesData := testdata.GetFamilyCompanies()
		name := "tiendas Deportivas"
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.FamilyCompanies)
		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.FamilyCompanies, familyCompaniesData[0], familyCompaniesData[1])

		familyCompaniesFound, err := repository.GetFamilyCompanies(ctx,
			entities.FamilyCompaniesFilter{
				Name: name,
			})

		assert.NoError(t, err)
		assert.Equalf(t, 1, len(familyCompaniesFound), "family companies len should be %d", 1)
		assert.Equal(t, familyCompaniesData[0].Name, familyCompaniesFound[0].Name)
		assert.Equal(t, familyCompaniesData[0].CompanyIDs, familyCompaniesFound[0].CompanyIDs)

		for i := range familyCompaniesData {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.FamilyCompanies, familyCompaniesData[i].ID)
		}
	})
}

func TestFamilyCompaniesRepository_Search(t *testing.T) {
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when family companies exists", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		familyCompaniesData := testdata.GetFamilyCompanies()
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.FamilyCompanies)
		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.FamilyCompanies, familyCompaniesData[0], familyCompaniesData[1])

		families, err := repository.Search(ctx,
			entities.FamilyCompaniesFilter{
				CompanyIDs: []string{familyCompaniesData[0].CompanyIDs[0], familyCompaniesData[1].CompanyIDs[0]},
			})

		assert.NoError(t, err)
		assert.Equalf(t, len(familyCompaniesData), len(families), "family companies len should be %d", len(familyCompaniesData))

		assert.Equal(t, familyCompaniesData[0], families[1])
		assert.Equal(t, familyCompaniesData[1], families[0])

		for i := range families {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.FamilyCompanies, families[i].ID)
		}
	})

	t.Run("when family companies exists, it should be in alphabetical order", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		familyCompaniesData := testdata.GetFamilyCompaniesToAlphabeticalOrder()
		repository := familycom.NewFamilyCompaniesMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.FamilyCompanies)
		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.FamilyCompanies,
			familyCompaniesData[0], familyCompaniesData[1], familyCompaniesData[2], familyCompaniesData[3])

		families, err := repository.Search(ctx,
			entities.FamilyCompaniesFilter{})

		assert.NoError(t, err)
		assert.Equalf(t, len(familyCompaniesData), len(families), "family companies len should be %d", len(familyCompaniesData))

		assert.Equal(t, familyCompaniesData[0], families[2])
		assert.Equal(t, familyCompaniesData[1], families[0])
		assert.Equal(t, familyCompaniesData[2], families[3])
		assert.Equal(t, familyCompaniesData[3], families[1])

		for i := range families {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.FamilyCompanies, families[i].ID)
		}
	})
}
