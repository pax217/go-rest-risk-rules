package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/families"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFamilyRepository_Add(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("on family_id duplicated then return error", func(t *testing.T) {

		family := testdata.GetDefaultFamily()
		repository := families.NewFamilyMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		expectedError := fmt.Sprintf(
			`write exception: write errors: [E11000 duplicate key error collection: rules.families index: _id_ dup key: { _id: ObjectId('%s') }]`,
			family.ID.Hex())
		defer cancel()

		err := repository.AddFamily(ctx, &family)

		assert.NoError(t, err)
		assert.NotEmpty(t, family.ID)

		var familyFound entities.Family
		mongoDB.Collection(cfg.MongoDB.Collections.Families).FindOne(ctx, bson.M{"_id": family.ID}).Decode(&familyFound)

		err = repository.AddFamily(ctx, &familyFound)

		mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Families, familyFound.ID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err.Error())
	})

	t.Run("on family added is ok", func(t *testing.T) {

		family := testdata.GetDefaultFamily()
		repository := families.NewFamilyMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.AddFamily(ctx, &family)

		assert.NoError(t, err)
		assert.NotEmpty(t, family.ID)

		var familyFound entities.Family
		mongoDB.Collection(cfg.MongoDB.Collections.Families).FindOne(ctx, bson.M{"_id": family.ID}).Decode(&familyFound)

		assert.NoError(t, err)
		assert.NotNil(t, family.ID)
		assert.Equal(t, family, familyFound)
		assert.Equal(t, 2, len(familyFound.ExcludedCompanies))

		mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Families, familyFound.ID)
	})
}

func TestFamilyRepository_SearchPaged(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)

	t.Run("when filter objectID is not valid", func(t *testing.T) {

		invalidId := "invalid-format-abc"
		expectedErr := fmt.Sprintf(
			"message: error: invalid id of family: '%s' - status: 400 - error: bad_request - causes: []",
			invalidId)

		repository := families.NewFamilyMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foundFamilies, err := repository.SearchPaged(ctx, entities.NewDefaultPagination(),
			entities.FamilyFilter{
				ID: invalidId,
			})

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err.Error())
		assert.EqualValues(t, foundFamilies.Data, []entities.Family{})
	})

	t.Run("when families data is found", func(t *testing.T) {

		familiesData := testdata.GetFamilies()
		repository := families.NewFamilyMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.Families)
		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Families, familiesData[0], familiesData[1])

		response, err := repository.SearchPaged(ctx, entities.NewDefaultPagination(),
			entities.FamilyFilter{
				Mccs: []string{familiesData[0].Mccs[0], familiesData[1].Mccs[0]},
			})
		data := response.Data.([]entities.Family)

		assert.NoError(t, err)
		assert.Equalf(t, len(familiesData), len(data), "response len should be %d", len(familiesData))

		assert.Equal(t, familiesData[0], data[1])
		assert.Equal(t, familiesData[1], data[0])

		for i := range data {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Families, data[i].ID)
		}
	})

	t.Run("when families data is not found because family contains excluded companies", func(t *testing.T) {

		familiesData := testdata.GetFamilies()
		repository := families.NewFamilyMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.Families)
		mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Families, familiesData[0], familiesData[1])

		data, err := repository.SearchEvaluate(ctx,
			entities.FamilyFilter{
				Mccs:                 []string{familiesData[0].Mccs[0], familiesData[1].Mccs[0]},
				NotExcludedCompanies: []string{familiesData[0].ExcludedCompanies[0]},
			})

		assert.NoError(t, err)
		assert.Equal(t, 1, len(data), "families len should be 1")

		assert.Equal(t, familiesData[1], data[0])
		assert.NotEqual(t, familiesData[1].ID, familiesData[0].ID)

		for i := range familiesData {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Families, familiesData[i].ID)
		}
	})
}

func TestFamilyRepository_Update(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	logger, _ := logs.New()
	configs := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(configs)

	t.Run("when family id is incorrect, then return error", func(t *testing.T) {
		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPutRequest()

		repository := families.NewFamilyMongoDBRepository(configs, mongoDB, logger)

		err := repository.Update(context.TODO(), "", &family)

		assert.NotNil(t, err)
	})

	t.Run("when family id no exist, then return error", func(t *testing.T) {
		request := testdata.GetFamilyRequest()
		family := request.NewFamilyFromPutRequest()
		id := primitive.NewObjectID()
		repository := families.NewFamilyMongoDBRepository(configs, mongoDB, logger)

		err := repository.Update(context.TODO(), id.Hex(), &family)

		assert.NotNil(t, err)
	})

	t.Run("when family update is success", func(t *testing.T) {
		repository := families.NewFamilyMongoDBRepository(configs, mongoDB, logger)
		family := testdata.GetDefaultFamily()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.AddFamily(ctx, &family)
		assert.NoError(t, err)

		defer mongoDB.CleanCollectionByIds(ctx, configs.MongoDB.Collections.Families, family.ID)
		family.ExcludedCompanies = []string{"6108753dd8567400011cdc00"}
		err = repository.Update(ctx, family.ID.Hex(), &family)

		assert.Nil(t, err)

		responses, err := repository.Search(ctx, entities.FamilyFilter{
			ID: family.ID.Hex(),
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, len(responses))
		assert.Equal(t, responses[0].ExcludedCompanies, family.ExcludedCompanies)
		assert.Equal(t, responses[0].ID, family.ID)
	})
}

func TestFamilyRepository_Delete(t *testing.T) {

	t.Run("on family_id not found, then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		expectedError := errors.New(`error: family not found: 'ObjectID("615238fcc1b177e01aefccf8")'`)
		repository := families.NewFamilyMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Delete(ctx, "615238fcc1b177e01aefccf8")

		assert.Error(t, err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("on family deleted with rule isGlobal in true, should delete ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		family := testdata.GetDefaultFamily()
		ruleGlobal := testdata.GetDefaultRuleGlobal()
		repository := families.NewFamilyMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ruleCollection := mongoDB.Collection(cfg.MongoDB.Collections.Rules)
		ruleCollection.InsertOne(ctx, ruleGlobal)

		familyCollection := mongoDB.Collection(cfg.MongoDB.Collections.Families)
		familyCollection.InsertOne(ctx, family)

		assert.NotEmpty(t, family.ID)

		err := repository.Delete(ctx, family.ID.Hex())

		assert.NoError(t, err)

		defer mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.Rules)
	})

	t.Run("on family deleted is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		family := testdata.GetDefaultFamily()
		repository := families.NewFamilyMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		familyCollection := mongoDB.Collection(cfg.MongoDB.Collections.Families)
		familyCollection.InsertOne(ctx, family)

		assert.NotEmpty(t, family.ID)

		err := repository.Delete(ctx, family.ID.Hex())

		assert.NoError(t, err)
	})
}
