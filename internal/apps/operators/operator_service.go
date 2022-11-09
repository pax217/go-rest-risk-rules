package operators

import (
	"context"
	"fmt"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/pkg/metrics"
	"github.com/conekta/risk-rules/pkg/text"
)

const serviceName = "operator.service.%s"

type OperatorService interface {
	AddOperator(ctx context.Context, operator entities.Operator) error
	Get(ctx context.Context, operatorFilter entities.OperatorFilter,
		pagination entities.Pagination) (interface{}, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, operator entities.Operator) error
}

type operatorService struct {
	logs               logs.Logger
	operatorRepository OperatorRepository
	datadog            datadog.Metricer
	config             config.Config
}

func NewOperatorService(configs config.Config, logger logs.Logger,
	operatorRepository OperatorRepository, metric datadog.Metricer) OperatorService {
	return &operatorService{
		logs:               logger,
		operatorRepository: operatorRepository,
		datadog:            metric,
		config:             configs,
	}
}

func (service *operatorService) AddOperator(ctx context.Context, operator entities.Operator) error {
	filter := entities.OperatorFilter{
		Type: operator.Type,
		Name: operator.Name,
	}

	operatorsFound, err := service.operatorRepository.Get(ctx, filter)
	if err != nil {
		return err
	}

	duplicateOperator := operator.SearchSimilarOperator(operatorsFound)
	if !duplicateOperator.IsEmpty() {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("the operator '%s' of type '%s' already exist",
			operator.Name, operator.Type))
		service.logs.Error(ctx, err.Error())
		return err
	}

	err = service.operatorRepository.Save(ctx, &operator)
	metricData := metrics.NewMetricData(ctx, "AddOperator", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveOperatorMetricName)
		return err
	}

	metricData.SetResult(false)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveOperatorMetricName)

	return nil
}

func (service *operatorService) Get(ctx context.Context, filter entities.OperatorFilter,
	pagination entities.Pagination) (interface{}, error) {
	if filter.Paged {
		return service.operatorRepository.GetPaged(ctx, filter, pagination)
	} else {
		return service.operatorRepository.Get(ctx, filter)
	}
}

func (service *operatorService) Update(ctx context.Context, id string, operator entities.Operator) error {
	metricData := metrics.NewMetricData(ctx, "Update Operator", serviceName, service.config.Env)
	err := operator.ValidateOperatorID(id)
	if err != nil {
		err := exceptions.NewInvalidRequest(err.Error())
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(handlerName, "Update"))
		return err
	}

	filter := entities.OperatorFilter{
		Type: operator.Type,
		Name: operator.Name,
	}

	operatorsFound, err := service.operatorRepository.Get(ctx, filter)
	if err != nil {
		return err
	}

	if service.existDuplicatedOperator(operator, operatorsFound, id) {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("operator: [%s] of type [%s] is already created",
			operator.Name, operator.Type))
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "Update"))
		return err
	}

	err = service.operatorRepository.Update(ctx, id, operator)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateOperatorMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateOperatorMetricName)

	return nil
}

func (service *operatorService) Delete(ctx context.Context, id string) error {
	err := service.operatorRepository.Delete(ctx, id)

	metricData := metrics.NewMetricData(ctx, "Delete Operator", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteOperatorMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteOperatorMetricName)

	return nil
}

func (service *operatorService) existDuplicatedOperator(operator entities.Operator,
	operators []entities.Operator, id string) bool {
	duplicateOperator := operator.SearchSimilarOperator(operators)
	if duplicateOperator.IsEmpty() {
		return false
	}

	if duplicateOperator.IsTheSame(id) {
		return false
	}

	return true
}
