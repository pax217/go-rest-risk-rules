package fields

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

const serviceName = "fields.service.%s"

type FieldService interface {
	AddField(ctx context.Context, field entities.Field) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, field entities.Field) error
	GetFields(ctx context.Context, filter entities.FieldsFilter, pagination entities.Pagination) (interface{}, error)
}

type fieldService struct {
	config     config.Config
	repository FieldsRepository
	logs       logs.Logger
	datadog    datadog.Metricer
}

func NewFieldsService(cfg config.Config, repository FieldsRepository, logger logs.Logger,
	metric datadog.Metricer) FieldService {
	return &fieldService{
		config:     cfg,
		repository: repository,
		logs:       logger,
		datadog:    metric,
	}
}

func (service *fieldService) AddField(ctx context.Context, field entities.Field) error {
	fieldsFound, err := service.repository.GetFields(ctx, field.GetFieldsFilter(false))
	if err != nil {
		return err
	}

	duplicateField := field.SearchSimilarField(fieldsFound)
	if !duplicateField.IsEmpty() {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("the field '%s' of type '%s' already exist", field.Name, field.Type))
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "AddField"))
		return err
	}

	err = service.repository.AddField(ctx, &field)
	metricData := metrics.NewMetricData(ctx, "AddField", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveFieldMetricName)
		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveFieldMetricName)
	return err
}

func (service *fieldService) Update(ctx context.Context, id string, field entities.Field) error {
	fieldsFound, err := service.repository.GetFields(ctx, field.GetFieldsFilter(false))
	if err != nil {
		return err
	}

	if service.existDuplicatedField(field, fieldsFound, id) {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("the field '%s' of type '%s' already exist", field.Name, field.Type))
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "AddField"))
		return err
	}

	err = service.repository.Update(ctx, id, &field)

	metricData := metrics.NewMetricData(ctx, "Update Field", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateFieldMetricName)
		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateFieldMetricName)
	return nil
}

func (service *fieldService) Delete(ctx context.Context, id string) error {
	err := service.repository.Delete(ctx, id)
	metricData := metrics.NewMetricData(context.TODO(), "Delete", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteFieldMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteFieldMetricName)
	return nil
}

func (service *fieldService) GetFields(ctx context.Context, filter entities.FieldsFilter,
	pagination entities.Pagination) (interface{}, error) {
	if filter.Paged {
		return service.repository.GetFieldsPaged(ctx, filter, pagination)
	} else {
		return service.repository.GetFields(ctx, filter)
	}
}

func (service *fieldService) existDuplicatedField(field entities.Field, fields []entities.Field, id string) bool {
	duplicateField := field.SearchSimilarField(fields)
	if duplicateField.IsEmpty() {
		return false
	}

	if duplicateField.IsTheSame(id) {
		return false
	}

	return true
}
