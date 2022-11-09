package integration

import (
	"context"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	merchantsscore "github.com/conekta/risk-rules/internal/apps/merchants_score"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMerchantsScoreRepository_Add(t *testing.T) {
	t.Run("on add merchants is success", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		merchants := testdata.GetDefaultMerchantScoreData()
		repository := merchantsscore.NewMerchantsMongoDBRepository(cfg, mongoDB, logger)
		defer cancel()

		err := repository.WriteMerchantsScore(ctx, merchants)

		assert.Nil(t, err)

		var merchant entities.MerchantScore
		err = mongoDB.Collection(cfg.MongoDB.Collections.MerchantsScore).FindOne(ctx, bson.M{"company_id": merchants[0].CompanyID}).Decode(&merchant)
		assert.Nil(t, err)
		assert.NotEmpty(t, merchant)
		assert.Equal(t, merchants[0].CompanyID, merchant.CompanyID)
		assert.Equal(t, merchants[0].Score, merchant.Score)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.MerchantsScore, merchant.ID)

		err = mongoDB.Collection(cfg.MongoDB.Collections.MerchantsScore).FindOne(ctx, bson.M{"company_id": merchants[1].CompanyID}).Decode(&merchant)
		assert.Nil(t, err)
		assert.NotEmpty(t, merchant)
		assert.Equal(t, merchants[1].CompanyID, merchant.CompanyID)
		assert.Equal(t, merchants[1].Score, merchant.Score)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.MerchantsScore, merchant.ID)
	})
}

func TestMerchantsScoreRepository_Find(t *testing.T) {
	t.Run("on find merchant score is success", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}

		logger, _ := logs.New()
		cfg := config.NewConfig()
		mongoDB := mongodb.NewMongoDB(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		merchants := testdata.GetDefaultMerchantScoreData()
		repository := merchantsscore.NewMerchantsMongoDBRepository(cfg, mongoDB, logger)
		defer cancel()

		err := repository.WriteMerchantsScore(ctx, merchants)

		assert.Nil(t, err)

		var merchant entities.MerchantScore
		merchant, err = repository.FindByMerchantID(ctx, merchants[0].CompanyID)
		assert.Nil(t, err)
		assert.NotEmpty(t, merchant)
		assert.Equal(t, merchants[0].CompanyID, merchant.CompanyID)
		assert.Equal(t, merchants[0].Score, merchant.Score)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.MerchantsScore, merchant.ID)

		merchant, err = repository.FindByMerchantID(ctx, merchants[1].CompanyID)
		assert.Nil(t, err)
		assert.NotEmpty(t, merchant)
		assert.Equal(t, merchants[1].CompanyID, merchant.CompanyID)
		assert.Equal(t, merchants[1].Score, merchant.Score)

		defer mongoDB.CleanCollectionByIds(ctx, cfg.MongoDB.Collections.MerchantsScore, merchant.ID)
	})
}
