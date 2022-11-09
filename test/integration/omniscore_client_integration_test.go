package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/pkg/rest"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestOmniscoreRest_GetScore(t *testing.T) {
	t.Run("when omniscore service responds successfully", func(t *testing.T) {

		expectedScore := 0.4
		logger, _ := logs.New()
		cfg := config.NewConfig()

		charge := testdata.GetDefaultCharge()
		charge.ID = "615324eb5bc1dea9ce66068c"
		restClient := rest.NewOmniscoreClient(cfg, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		score, err := restClient.GetScore(ctx, charge)

		assert.NoError(t, err)
		assert.Equal(t, expectedScore, score)
	})

	t.Run("when omniscore service timeouts", func(t *testing.T) {

		expectedScore := float64(-1)
		logger, _ := logs.New()
		cfg := config.NewConfig()

		charge := testdata.GetDefaultCharge()
		charge.ID = "615324eb5bc1dea9ce66068a"
		restClient := rest.NewOmniscoreClient(cfg, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		score, err := restClient.GetScore(ctx, charge)

		assert.Error(t, err, fmt.Sprintf("should return an error with charge_id %s", charge.ID))
		assert.Equal(t, expectedScore, score, fmt.Sprintf("should return an error with charge_id %s", charge.ID))
	})

	t.Run("when omniscore service responds an http error", func(t *testing.T) {

		expectedScore := float64(-1)
		logger, _ := logs.New()
		cfg := config.NewConfig()

		charge := testdata.GetDefaultCharge()
		charge.ID = "615324eb5bc1dea9ce66068b"
		restClient := rest.NewOmniscoreClient(cfg, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		score, err := restClient.GetScore(ctx, charge)

		assert.Error(t, err, fmt.Sprintf("should return an error with charge_id %s", charge.ID))
		assert.Equal(t, expectedScore, score, fmt.Sprintf("should return an error with charge_id %s", charge.ID))
	})

	t.Run("when omniscore service is not enabled", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration tests in short mode.")
		}
		expectedScore := float64(-1)
		logger, _ := logs.New()
		cfg := config.NewConfig()
		cfg.Omniscore.IsEnabled = false

		charge := testdata.GetDefaultCharge()
		restClient := rest.NewOmniscoreClient(cfg, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		score, err := restClient.GetScore(ctx, charge)

		assert.NoError(t, err)
		assert.Equal(t, expectedScore, score)
	})
}
