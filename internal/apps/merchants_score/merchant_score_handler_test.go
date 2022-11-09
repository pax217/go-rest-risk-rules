package merchantsscore_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/conekta/go_common/logs"
	merchantsscore "github.com/conekta/risk-rules/internal/apps/merchants_score"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/pkg/echo"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/stretchr/testify/assert"
)

const uri = "/risk-rules/v1/merchants_score"

func Test_MerchantsScoreHandler_FileProcessing(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("When merchant score processing is success", func(t *testing.T) {
		service := new(mocks.MerchantsScoreServiceMock)

		context, recorder := echo.SetupAsRecorder(http.MethodGet, uri, "", "")
		service.On("MerchantScoreProcessing", context.Request().Context()).Return(nil)

		handler := merchantsscore.NewMerchantsScoreHandler(configs, logger, service)
		handler.MerchantScoreProcessing(context)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("When merchant score return error", func(t *testing.T) {
		service := new(mocks.MerchantsScoreServiceMock)
		expectedError := errors.New("Syntax error")

		context, _ := echo.SetupAsRecorder(http.MethodPost, uri, "", "")
		service.On("MerchantScoreProcessing", context.Request().Context()).Return(expectedError)

		handler := merchantsscore.NewMerchantsScoreHandler(configs, logger, service)
		err := handler.MerchantScoreProcessing(context)

		assert.NotNil(t, expectedError, err)
	})
}
