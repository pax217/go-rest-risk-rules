package omniscores

import (
	"context"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/rest"
)

type OmniscoreService interface {
	GetScore(ctx context.Context, charge entities.ChargeRequest) float64
}

type omniscoreService struct {
	config          config.Config
	logs            logs.Logger
	omniscoreClient rest.OmniscoreClient
}

func NewOmniscoreService(cfg config.Config, logger logs.Logger, omniscoreClient rest.OmniscoreClient) OmniscoreService {
	return &omniscoreService{
		config:          cfg,
		logs:            logger,
		omniscoreClient: omniscoreClient,
	}
}

func (o *omniscoreService) GetScore(ctx context.Context, charge entities.ChargeRequest) float64 {
	score, _ := o.omniscoreClient.GetScore(ctx, charge)
	return score
}
