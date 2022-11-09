package modules

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

const serviceName = "modules.service.%s"

type ModuleService interface {
	Add(ctx context.Context, module entities.Module) error
	GetAll(ctx context.Context, pagination entities.Pagination, filter entities.ModuleFilter) (interface{}, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, module entities.Module) error
}

type moduleService struct {
	repository ModuleRepository
	logs       logs.Logger
	metrics    datadog.Metricer
	config     config.Config
}

func (service *moduleService) Add(ctx context.Context, module entities.Module) error {
	modulesFound, err := service.repository.Get(ctx, module.GetModuleFilter(false))
	if err != nil {
		return err
	}

	duplicateModule := module.SearchSimilarModule(modulesFound)
	if !strings.IsEmpty(duplicateModule.Name) {
		err = exceptions.NewDuplicatedException(
			fmt.Sprintf("module with name '%s' is already created", module.Name))
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "Add"))
		return err
	}

	err = service.repository.Save(ctx, &module)
	metricData := metrics.NewMetricData(ctx, "Add", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.metrics, service.logs, metricData, text.SaveModuleMetricName)
		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.metrics, service.logs, metricData, text.SaveModuleMetricName)
	return nil
}

func (service *moduleService) Update(ctx context.Context, id string, module entities.Module) error {
	modulesFound, err := service.repository.Get(ctx, module.GetModuleFilter(false))
	if err != nil {
		return err
	}

	if service.existDuplicatedModule(module, modulesFound, id) {
		err = exceptions.NewDuplicatedException(
			fmt.Sprintf("module with name '%s' is already created", module.Name))
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "Add"))
		return err
	}

	err = service.repository.Update(ctx, id, module)

	metricData := metrics.NewMetricData(ctx, "Update Module", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.metrics, service.logs, metricData, text.UpdateModuleMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.metrics, service.logs, metricData, text.UpdateModuleMetricName)

	return nil
}

func (service *moduleService) Delete(ctx context.Context, id string) error {
	err := service.repository.Delete(ctx, id)

	metricData := metrics.NewMetricData(ctx, "Delete Module", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.metrics, service.logs, metricData, text.DeleteModuleMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.metrics, service.logs, metricData, text.DeleteModuleMetricName)

	return nil
}

func NewModuleService(conf config.Config, repo ModuleRepository, logger logs.Logger, metric datadog.Metricer) ModuleService {
	return &moduleService{
		repository: repo,
		logs:       logger,
		metrics:    metric,
		config:     conf,
	}
}

func (service *moduleService) GetAll(ctx context.Context, pagination entities.Pagination,
	filter entities.ModuleFilter) (interface{}, error) {
	if filter.Paged {
		return service.repository.GetPaged(ctx, filter, pagination)
	} else {
		return service.repository.Get(ctx, filter)
	}
}

func (service *moduleService) existDuplicatedModule(module entities.Module,
	modules []entities.Module, id string) bool {
	duplicateModule := module.SearchSimilarModule(modules)
	if strings.IsEmpty(duplicateModule.Name) {
		return false
	}

	if duplicateModule.IsTheSame(id) {
		return false
	}

	return true
}
