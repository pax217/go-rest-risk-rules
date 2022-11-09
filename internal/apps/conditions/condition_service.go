package conditions

import (
	"context"
	"fmt"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/pkg/metrics"
	"github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
)

const serviceName = "conditions.service.%s"

type ConditionRepository interface {
	Add(ctx context.Context, condition *entities.Condition) error
	FindByName(ctx context.Context, condition entities.Condition) (entities.Condition, error)
	GetAll(ctx context.Context, filter entities.ConditionsFilter, pagination entities.Pagination) (entities.PagedResponse, error)
	Update(ctx context.Context, id string, condition entities.Condition) error
	Delete(ctx context.Context, id string) error
}

type ConditionService interface {
	Add(ctx context.Context, condition entities.Condition) error
	GetAll(ctx context.Context, filter entities.ConditionsFilter, pagination entities.Pagination) (entities.PagedResponse, error)
	Update(ctx context.Context, id string, condition entities.Condition) error
	Delete(ctx context.Context, id string) error
}

type conditionsService struct {
	repository ConditionRepository
	log        logs.Logger
	datadog    datadog.Metricer
	config     config.Config
}

func (service *conditionsService) Add(ctx context.Context, condition entities.Condition) error {
	conditionFound, err := service.repository.FindByName(ctx, condition)
	if err != nil {
		return err
	}

	if !strings.IsEmpty(conditionFound.Name) {
		err = exceptions.NewDuplicatedException(
			fmt.Sprintf("condition with name '%s' is already created", condition.Name))
		service.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "Add"))
		return err
	}

	err = service.repository.Add(ctx, &condition)
	metricData := metrics.NewMetricData(ctx, "Add", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.log, metricData, text.SaveConditionsMetricName)
		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.log, metricData, text.SaveConditionsMetricName)

	return nil
}

func (service *conditionsService) GetAll(ctx context.Context, filter entities.ConditionsFilter, pagination entities.Pagination) (
	entities.PagedResponse, error) {
	return service.repository.GetAll(ctx, filter, pagination)
}

func (service *conditionsService) Update(ctx context.Context, id string, condition entities.Condition) error {
	conditionFound, err := service.repository.FindByName(ctx, condition)
	if err != nil {
		return err
	}

	if !strings.IsEmpty(conditionFound.Name) && conditionFound.ID.Hex() != id {
		err = exceptions.NewDuplicatedException(
			fmt.Sprintf("condition with name '%s' is already created", condition.Name))
		service.log.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "Update"))
		return err
	}

	err = service.repository.Update(ctx, id, condition)

	metricData := metrics.NewMetricData(ctx, "Update Condition", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.log, metricData, text.UpdateConditionMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.log, metricData, text.UpdateConditionMetricName)

	return nil
}

func (service *conditionsService) Delete(ctx context.Context, id string) error {
	err := service.repository.Delete(ctx, id)

	metricData := metrics.NewMetricData(ctx, "Delete Condition", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.log, metricData, text.DeleteConditionMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.log, metricData, text.DeleteConditionMetricName)

	return nil
}

func NewConditionsService(configs config.Config, repository ConditionRepository,
	logger logs.Logger, metricer datadog.Metricer) ConditionService {
	return &conditionsService{
		repository: repository,
		log:        logger,
		datadog:    metricer,
		config:     configs,
	}
}
