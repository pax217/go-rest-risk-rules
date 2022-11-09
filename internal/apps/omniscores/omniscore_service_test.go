package omniscores

import (
	"context"
	"errors"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ServiceOmniscore_GetScore(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()
	ctx := context.TODO()
	charge := testdata.GetDefaultCharge()

	t.Run("When omniscore client retrieves score successfully, return score", func(t *testing.T) {
		expectedScore := 0.4
		omniscoreClientMock := new(mocks.OmniscoreClientMock)
		omniscoreClientMock.On("GetScore", ctx, charge).Return(expectedScore, nil)
		omniscoreService := NewOmniscoreService(configs, logger, omniscoreClientMock)

		score := omniscoreService.GetScore(ctx, charge)

		assert.Equal(t, expectedScore, score)
	})

	t.Run("When omniscore client retrieves score unsuccessfully, return default score value", func(t *testing.T) {
		expectedScore := float64(-1)
		omniscoreClientMock := new(mocks.OmniscoreClientMock)
		omniscoreClientMock.On("GetScore", ctx, charge).Return(
			expectedScore, errors.New("omniscore rest error"))
		omniscoreService := NewOmniscoreService(configs, logger, omniscoreClientMock)

		score := omniscoreService.GetScore(ctx, charge)

		assert.Equal(t, expectedScore, score)
	})
}
