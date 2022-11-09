package merchantsscore_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/conekta/go_common/logs"
	merchantsscore "github.com/conekta/risk-rules/internal/apps/merchants_score"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	formatFileName = "%s/sc_merchant_%s_000"
	dateFormat     = "2006-01-02"
)

func Test_MerchantsScore_FileProcessing(t *testing.T) {
	logger, _ := logs.New()
	configs := config.NewConfig()

	t.Run("when only add merchants score is success", func(t *testing.T) {
		repository := new(mocks.MerchantsScoreRepositoryMock)
		s3Repository := new(mocks.MerchantScoreS3RepositoryMock)
		service := merchantsscore.NewMerchantsScoreService(configs, logger, new(datadog.MetricsDogMock), repository, s3Repository)

		scores := testdata.GetDefaultMerchantScoreData()

		s3Repository.Mock.On("GetFileContent", context.TODO(), createFileName(configs.MerchantScore.S3PrefixFile)).
			Return(scores, nil).Once()

		repository.Mock.On("WriteMerchantsScore", context.TODO(), scores).
			Return(nil).Once()

		err := service.MerchantScoreProcessing(context.TODO())

		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("adding merchants score returns error", func(t *testing.T) {
		repository := new(mocks.MerchantsScoreRepositoryMock)
		s3Repository := new(mocks.MerchantScoreS3RepositoryMock)
		service := merchantsscore.NewMerchantsScoreService(configs, logger, new(datadog.MetricsDogMock), repository, s3Repository)
		expectedError := errors.New("database connection lost")

		scores := testdata.GetDefaultMerchantScoreData()

		s3Repository.Mock.On("GetFileContent", context.TODO(), createFileName(configs.MerchantScore.S3PrefixFile)).
			Return(scores, nil).Once()

		repository.Mock.On("WriteMerchantsScore", context.TODO(), scores).
			Return(expectedError).Once()

		err := service.MerchantScoreProcessing(context.TODO())

		assert.NotNil(t, expectedError, err)
		repository.AssertExpectations(t)
	})

	t.Run("when get file content return error", func(t *testing.T) {
		repository := new(mocks.MerchantsScoreRepositoryMock)
		s3Repository := new(mocks.MerchantScoreS3RepositoryMock)
		service := merchantsscore.NewMerchantsScoreService(configs, logger, new(datadog.MetricsDogMock), repository, s3Repository)
		expectedError := errors.New("file error")

		s3Repository.Mock.On("GetFileContent", context.TODO(), createFileName(configs.MerchantScore.S3PrefixFile)).
			Return([]entities.MerchantScore{}, expectedError).Once()

		err := service.MerchantScoreProcessing(context.TODO())

		assert.NotNil(t, expectedError, err)
		repository.AssertExpectations(t)
	})
}

func createFileName(prefix string) string {
	now := time.Now().UTC()
	return fmt.Sprintf(formatFileName, prefix, now.Format(dateFormat))
}
