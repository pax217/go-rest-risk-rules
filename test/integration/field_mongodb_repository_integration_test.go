package integration

import (
	"context"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/fields"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFieldRepository_AddField(t *testing.T) {
	t.Run("when add field on repository error id duplicated then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		field := testdata.GetFielNotValid()
		repository := fields.NewFieldsMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.AddField(ctx, &field)
		assert.NoError(t, err)

		err = repository.AddField(ctx, &field)

		mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Fields, field.ID)

		assert.NotNil(t, err)
	})

	t.Run("when add field on repository ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		field := testdata.GetDefaultField()
		repository := fields.NewFieldsMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.AddField(ctx, &field)

		assert.NoError(t, err)
		assert.NotEmpty(t, field.ID)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Fields, field.ID)

		var fld entities.Field
		err = mongoDB.Collection(cfg.MongoDB.Collections.Fields).FindOne(ctx, bson.M{"_id": field.ID}).Decode(&fld)

		assert.NoError(t, err)
		assert.EqualValues(t, field.ID, fld.ID)
	})
}

func TestFieldRepository_GetFieldsPaged(t *testing.T) {

	t.Run("disordered database result is sorted on get request", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		fieldsData := testdata.GetFields()
		repository := fields.NewFieldsMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.Fields)

		objectIDS := mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Fields, fieldsData[0], fieldsData[1])
		fieldsData[0].SetID(objectIDS[0])
		fieldsData[1].SetID(objectIDS[1])

		pagination := entities.NewDefaultPagination()

		operatorsFounded, err := repository.GetFieldsPaged(
			ctx,
			entities.FieldsFilter{},
			pagination)
		data := operatorsFounded.Data.([]entities.Field)

		assert.NoError(t, err)
		assert.Equalf(t, len(fieldsData), len(data), "operators len should be %d", len(fieldsData))
		assert.Equal(t, fieldsData[0], data[1])
		assert.Equal(t, fieldsData[1], data[0])

		for i := range data {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Fields, data[i].ID)
		}
	})
}

func TestFieldRepository_GetFields(t *testing.T) {

	t.Run("disordered database result is sorted on get request", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		fieldsData := testdata.GetFields()
		repository := fields.NewFieldsMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoDB.ClearCollection(ctx, cfg.MongoDB.Collections.Fields)

		objectIDS := mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Fields, fieldsData[0], fieldsData[1])
		fieldsData[0].SetID(objectIDS[0])
		fieldsData[1].SetID(objectIDS[1])

		operatorsFounded, err := repository.GetFields(
			ctx,
			entities.FieldsFilter{},
		)

		assert.NoError(t, err)
		assert.Equalf(t, len(fieldsData), len(operatorsFounded), "operators len should be %d", len(fieldsData))
		assert.Equal(t, fieldsData[0], operatorsFounded[1])
		assert.Equal(t, fieldsData[1], operatorsFounded[0])

		for i := range operatorsFounded {
			mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Fields, operatorsFounded[i].ID)
		}
	})
}

func TestFieldRepository_Update(t *testing.T) {
	t.Run("when id is not valid then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		id := "507f1f77bcf86cd799439011-x"
		field := testdata.GetDefaultField()
		repository := fields.NewFieldsMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Update(ctx, id, &field)
		assert.NotNil(t, err)
	})

	t.Run("when update field on database is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		field := testdata.GetDefaultField()
		repository := fields.NewFieldsMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.AddField(ctx, &field)
		assert.NoError(t, err)

		now := time.Now().UTC().Truncate(time.Millisecond)
		updatedBy := "carlos.maldonado@conekta.com"
		fieldToUpdate := entities.Field{
			ID:          field.ID,
			Name:        "name updated",
			Description: "description updated",
			Type:        "type updated",
			UpdatedBy:   &updatedBy,
			UpdatedAt:   &now,
		}

		err = repository.Update(
			ctx,
			fieldToUpdate.ID.Hex(),
			&fieldToUpdate)

		assert.NoError(t, err)
		assert.NotEmpty(t, field.ID)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Fields, fieldToUpdate.ID)

		var fld entities.Field
		err = mongoDB.Collection(cfg.MongoDB.Collections.Fields).FindOne(ctx, bson.M{"_id": fieldToUpdate.ID}).Decode(&fld)

		assert.NoError(t, err)
		assert.EqualValues(t, fieldToUpdate.ID, fld.ID)
		assert.EqualValues(t, fieldToUpdate.Name, fld.Name)
		assert.EqualValues(t, fieldToUpdate.Description, fld.Description)
		assert.EqualValues(t, fieldToUpdate.Type, fld.Type)
		assert.EqualValues(t, fieldToUpdate.UpdatedAt, fld.UpdatedAt)
		assert.EqualValues(t, fieldToUpdate.UpdatedBy, fld.UpdatedBy)
	})
}

func TestFieldRepository_Delete(t *testing.T) {
	t.Run("when id is valid then delete field", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		field := testdata.GetDefaultField()
		repository := fields.NewFieldsMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.AddField(ctx, &field)
		assert.NoError(t, err)

		ID := field.ID.Hex()
		err = repository.Delete(ctx, ID)
		assert.NoError(t, err)

		fieldsFounded, err := repository.GetFields(
			ctx,
			entities.FieldsFilter{
				Name: field.Name,
				Type: field.Type,
			},
		)

		assert.NoError(t, err)
		assert.Empty(t, fieldsFounded)
	})

	t.Run("when ID is invalid", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		field := testdata.GetDefaultField()
		repository := fields.NewFieldsMongoDBRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.AddField(ctx, &field)
		assert.NoError(t, err)

		ID := "x"
		err = repository.Delete(ctx, ID)
		assert.NotNil(t, err)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Fields, field.ID)
	})
}
