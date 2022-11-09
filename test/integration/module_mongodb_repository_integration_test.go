package integration

import (
	"context"
	"testing"
	"time"

	"github.com/conekta/risk-rules/internal/apps/modules"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
)

func TestModuleRepository_Save(t *testing.T) {
	t.Run("on module_id duplicated then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		module := testdata.GetDefaultModule()
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, &module)
		assert.NoError(t, err)

		err = repository.Save(ctx, &module)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Modules, module.ID)

		assert.Error(t, err)
	})

	t.Run("on module added is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		module := testdata.GetDefaultModule()
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, &module)
		assert.NoError(t, err)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Modules, module.ID)

		var mdl entities.Module
		err = mongoDB.Collection(cfg.MongoDB.Collections.Modules).FindOne(ctx, bson.M{"_id": module.ID}).Decode(&mdl)

		assert.NoError(t, err)
		assert.NotEmpty(t, module.ID)
		assert.EqualValues(t, module.ID, mdl.ID)
		assert.EqualValues(t, module.CreatedAt, mdl.CreatedAt)
		assert.EqualValues(t, module.UpdatedAt, mdl.UpdatedAt)
		assert.EqualValues(t, module.CreatedBy, mdl.CreatedBy)
		assert.EqualValues(t, module.UpdatedBy, mdl.UpdatedBy)
		assert.EqualValues(t, module.Name, mdl.Name)
		assert.EqualValues(t, module.Description, mdl.Description)
	})
}

func TestModuleRepository_FindByName(t *testing.T) {
	t.Run("on module not found", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		module := testdata.GetDefaultModule()
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := repository.Get(ctx, module.GetModuleFilter(false))

		assert.Nil(t, err)
		assert.Empty(t, result)
	})

	t.Run("on module found", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		module := testdata.GetDefaultModule()
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, &module)
		assert.NoError(t, err)

		moduleFounded, err := repository.Get(ctx, module.GetModuleFilter(false))

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Modules, module.ID)

		assert.NoError(t, err)
		assert.Equal(t, moduleFounded[0].ID, module.ID)
	})
}

func TestModuleRepository_GetPaged(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}
	logger, _ := logs.New()
	cfg := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(cfg)
	mockData := testdata.GetModules()
	repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	modulesID := mongoDB.PrepareData(ctx, cfg.MongoDB.Collections.Modules, mockData[0], mockData[1])
	mockData[0].SetID(modulesID[0])
	mockData[1].SetID(modulesID[1])

	pagination := entities.NewDefaultPagination()
	found, err := repository.GetPaged(ctx, entities.ModuleFilter{Paged: true}, pagination)
	data := found.Data.([]entities.Module)

	assert.NoError(t, err)
	assert.Equal(t, len(mockData), len(data))
	assert.Equal(t, mockData[0], data[1])
	assert.Equal(t, mockData[1], data[0])

	for i := range data {
		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Modules, data[i].ID)
	}
}

func TestModuleRepository_Update(t *testing.T) {
	t.Run("on module_id not found, then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		module := testdata.GetDefaultModule()
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Update(ctx, "615238fcc1b177e01aefccf8", module)

		assert.Error(t, err)
	})

	t.Run("on module updated ok, then return nil", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		module := testdata.GetDefaultModule()
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, &module)

		now := time.Now().UTC().Truncate(time.Millisecond)
		updatedBy := "santiago.ceron@conekta.com"
		moduleToUpdate := entities.Module{
			ID:          module.ID,
			CreatedAt:   module.CreatedAt,
			UpdatedAt:   &now,
			CreatedBy:   "carlos.maldonado@conekta.com",
			UpdatedBy:   &updatedBy,
			Name:        "name updated",
			Description: "description updated",
		}

		err = repository.Update(
			ctx,
			moduleToUpdate.ID.Hex(),
			moduleToUpdate)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.Modules, moduleToUpdate.ID)

		var mdl entities.Module
		err = mongoDB.Collection(cfg.MongoDB.Collections.Modules).FindOne(ctx, bson.M{"_id": moduleToUpdate.ID}).Decode(&mdl)

		assert.NoError(t, err)
		assert.NotEmpty(t, moduleToUpdate.ID)
		assert.EqualValues(t, moduleToUpdate.ID, mdl.ID)
		assert.EqualValues(t, moduleToUpdate.CreatedAt, mdl.CreatedAt)
		assert.EqualValues(t, moduleToUpdate.UpdatedAt, mdl.UpdatedAt)
		assert.EqualValues(t, moduleToUpdate.CreatedBy, mdl.CreatedBy)
		assert.EqualValues(t, moduleToUpdate.UpdatedBy, mdl.UpdatedBy)
		assert.EqualValues(t, moduleToUpdate.Name, mdl.Name)
		assert.EqualValues(t, moduleToUpdate.Description, mdl.Description)
	})

	t.Run("when ID is invalid return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		module := testdata.GetDefaultModule()
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Update(ctx, "x", module)

		assert.Error(t, err)
	})
}

func TestModuleRepository_Delete(t *testing.T) {
	t.Run("on module_id not found, then return error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Delete(ctx, "615238fcc1b177e01aefccf8")

		assert.Error(t, err)
	})

	t.Run("on module deleted is ok", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		module := testdata.GetDefaultModule()
		repository := modules.NewModulesMongoRepository(cfg, mongoDB, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, &module)

		assert.NoError(t, err)
		assert.NotEmpty(t, module.ID)

		err = repository.Delete(ctx, module.ID.Hex())

		assert.NoError(t, err)
	})
}
