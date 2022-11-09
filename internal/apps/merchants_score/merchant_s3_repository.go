package merchantsscore

import (
	"context"
	"fmt"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/csv"
	"github.com/conekta/risk-rules/pkg/text"
	"github.com/gocarina/gocsv"
)

const (
	merchantScoreS3RepositoryName = "merchantsscore.repository.s3"
)

type merchantScoreS3Repository struct {
	config   config.Config
	s3Reader csv.S3Reader
	Logger   logs.Logger
}

type MerchantScoreS3Repository interface {
	GetFileContent(ctx context.Context, fileName string) ([]entities.MerchantScore, error)
}

func NewMerchantScoreS3Repository(cfg config.Config, logger logs.Logger, s3Reader csv.S3Reader) MerchantScoreS3Repository {
	return &merchantScoreS3Repository{
		config:   cfg,
		Logger:   logger,
		s3Reader: s3Reader,
	}
}

func (repository *merchantScoreS3Repository) GetFileContent(ctx context.Context, fileName string) ([]entities.MerchantScore, error) {
	scores := make([]entities.MerchantScore, 0)
	download, err := repository.s3Reader.ReadS3File(ctx, repository.config.MerchantScore.S3Bucket, fileName)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", merchantScoreS3RepositoryName, "GetFileContent"))
		return scores, err
	}

	err = gocsv.UnmarshalBytes(download, &scores)
	if err != nil {
		repository.Logger.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", merchantScoreS3RepositoryName, "GetFileContent"))
		return scores, err
	}

	return scores, nil
}
