package merchantsscore

import (
	"context"
	"fmt"
	"time"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/pkg/metrics"
	"github.com/conekta/risk-rules/pkg/text"
)

const (
	serviceName    = "merchantsscore.service.%s"
	formatFileName = "%s/sc_merchant_%s_000"
	dateFormat     = "2006-01-02"
)

type MerchantsScoreService interface {
	MerchantScoreProcessing(ctx context.Context) error
}

type merchantsScoreService struct {
	logs         logs.Logger
	datadog      datadog.Metricer
	config       config.Config
	repository   MerchantsScoreRepository
	repositoryS3 MerchantScoreS3Repository
}

func NewMerchantsScoreService(cfg config.Config, logger logs.Logger, metric datadog.Metricer, repository MerchantsScoreRepository,
	repositoryS3 MerchantScoreS3Repository) MerchantsScoreService {
	return &merchantsScoreService{
		config:       cfg,
		logs:         logger,
		repository:   repository,
		datadog:      metric,
		repositoryS3: repositoryS3,
	}
}

func (service *merchantsScoreService) MerchantScoreProcessing(ctx context.Context) error {
	merchantsScore, err := service.repositoryS3.GetFileContent(ctx, createFileName(service.config.MerchantScore.S3PrefixFile))
	if err != nil {
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "MerchantScoreProcessing"))
		return err
	}

	err = service.repository.WriteMerchantsScore(ctx, merchantsScore)
	metricData := metrics.NewMetricData(ctx, "MerchantScoreProcessing", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveMerchantsScoreMetricName)

		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveMerchantsScoreMetricName)

	return nil
}

func createFileName(prefix string) string {
	now := time.Now().UTC()
	return fmt.Sprintf(formatFileName, prefix, now.Format(dateFormat))
}
