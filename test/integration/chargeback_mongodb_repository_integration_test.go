package integration

import (
	"context"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/chargebacks"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestPayerChargebackRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	logger, _ := logs.New()
	configs := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(configs)

	t.Run("when payer save is success", func(t *testing.T) {
		repository := chargebacks.NewChargebacksMongoDBRepository(configs, mongoDB, logger)
		payer := testdata.GetDefaultPayer().NewPayerFromPostRequest()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, payer)

		assert.Nil(t, err)

		payerFound, err := repository.Find(ctx, payer)

		assert.Nil(t, err)
		assert.Equal(t, payer.Email, payerFound.Email)

		defer mongoDB.CleanCollectionByIds(ctx, configs.MongoDB.Collections.Payers, payer.ID)
	})

	t.Run("when family update is success", func(t *testing.T) {
		repository := chargebacks.NewChargebacksMongoDBRepository(configs, mongoDB, logger)
		payer := testdata.GetDefaultPayer().NewPayerFromPostRequest()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repository.Save(ctx, payer)

		assert.Nil(t, err)

		payerFound, err := repository.Find(ctx, payer)

		assert.Nil(t, err)
		assert.Equal(t, payer.Email, payerFound.Email)

		err = repository.Update(ctx, payer)

		assert.Nil(t, err)

		defer mongoDB.CleanCollectionByIds(ctx, configs.MongoDB.Collections.Payers, payer.ID)
	})

	t.Run("when find a payer is success", func(t *testing.T) {
		repository := chargebacks.NewChargebacksMongoDBRepository(configs, mongoDB, logger)
		payer := testdata.GetDefaultPayer().NewPayerFromPostRequest()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := repository.Save(ctx, payer)

		assert.Nil(t, err)

		payerFound, err := repository.Find(ctx, payer)

		assert.Nil(t, err)
		assert.Equal(t, payer.Email, payerFound.Email)
		assert.Equal(t, payer.Chargebacks[0].Amount, payerFound.Chargebacks[0].Amount)
		assert.Equal(t, payer.Chargebacks[0].ChargeID, payerFound.Chargebacks[0].ChargeID)
		assert.Equal(t, payer.Chargebacks[0].ChargebackID, payerFound.Chargebacks[0].ChargebackID)
		assert.Equal(t, payer.Chargebacks[0].Currency, payerFound.Chargebacks[0].Currency)
		assert.Equal(t, payer.Chargebacks[0].Reason, payerFound.Chargebacks[0].Reason)
		assert.Equal(t, payer.Chargebacks[0].Status, payerFound.Chargebacks[0].Status)

		defer mongoDB.CleanCollectionByIds(ctx, configs.MongoDB.Collections.Payers, payer.ID)
	})
}
