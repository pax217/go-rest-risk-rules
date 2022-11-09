package lists

import (
	"context"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/rest"
)

type ListsService interface {
	GetLists(ctx context.Context, listsSearch entities.ListsSearch) ([]entities.List, error)
}

type listsService struct {
	config      config.Config
	logs        logs.Logger
	datadog     datadog.Metricer
	listsClient rest.RkListsClient
}

func NewListsService(cfg config.Config, logger logs.Logger, datadogMetric datadog.Metricer,
	client rest.RkListsClient) ListsService {
	return &listsService{
		config:      cfg,
		logs:        logger,
		datadog:     datadogMetric,
		listsClient: client,
	}
}

func (service *listsService) GetLists(ctx context.Context, listsSearch entities.ListsSearch) ([]entities.List, error) {
	return service.listsClient.ListsSearch(ctx, listsSearch)
}
