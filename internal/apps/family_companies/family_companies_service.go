package familycom

import (
	"context"
	"fmt"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/pkg/metrics"
	"github.com/conekta/risk-rules/pkg/text"
)

const serviceName = "familycom.service.%s"

type FamilyCompaniesService interface {
	Create(ctx context.Context, familyCompaniesRequest entities.FamilyCompanies) error
	GetFamiliesCompaniesFromFilter(ctx context.Context, filter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, familyCompanies entities.FamilyCompanies) error
	Get(ctx context.Context, pagination entities.Pagination,
		filter entities.FamilyCompaniesFilter) (interface{}, error)
}

type familyCompaniesService struct {
	logs                      logs.Logger
	datadog                   datadog.Metricer
	config                    config.Config
	familyCompaniesRepository FamilyCompaniesRepository
	ruleRepository            rules.RuleRepository
}

func NewFamilyCompaniesService(cfg config.Config,
	familyCompaniesRepository FamilyCompaniesRepository,
	ruleRepository rules.RuleRepository,
	logger logs.Logger,
	datadogMetric datadog.Metricer) FamilyCompaniesService {
	return &familyCompaniesService{
		config:                    cfg,
		familyCompaniesRepository: familyCompaniesRepository,
		ruleRepository:            ruleRepository,
		logs:                      logger,
		datadog:                   datadogMetric,
	}
}

func (service *familyCompaniesService) Create(ctx context.Context, familyCompanies entities.FamilyCompanies) error {
	metricData := metrics.NewMetricData(ctx, "CreateFamilyCompanies", serviceName, service.config.Env)

	filter := entities.FamilyCompaniesFilter{
		Name: familyCompanies.Name,
	}

	familyCompaniesFound, err := service.familyCompaniesRepository.GetFamilyCompanies(ctx, filter)
	if err != nil {
		return err
	}
	if len(familyCompaniesFound) > 0 {
		var causes exceptions.Causes

		errMessage := fmt.Sprintf("family companies with name [%s] is already created", familyCompanies.Name)
		causes.Code = exceptions.FamilyCompaniesNameDuplicated
		err = exceptions.NewDuplicatedExceptionWithCause(errMessage, causes)
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(serviceName, "Create"))
		return err
	}

	err = service.familyCompaniesRepository.AddFamilyCompanies(ctx, &familyCompanies)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveListMetricName)
		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveListMetricName)

	return err
}

func (service *familyCompaniesService) Get(ctx context.Context,
	pagination entities.Pagination,
	filter entities.FamilyCompaniesFilter) (interface{}, error) {
	if filter.Paged {
		return service.familyCompaniesRepository.SearchPaged(ctx, pagination, filter)
	} else {
		return service.familyCompaniesRepository.Search(ctx, filter)
	}
}

func (service *familyCompaniesService) Update(ctx context.Context, id string, familyCompanies entities.FamilyCompanies) error {
	metricData := metrics.NewMetricData(ctx, "Update", serviceName, service.config.Env)
	filter := entities.FamilyCompaniesFilter{
		Name:  familyCompanies.Name,
		Paged: false,
	}

	familyCompaniesFound, err := service.familyCompaniesRepository.
		GetFamilyCompanies(ctx, filter)

	if err != nil {
		return err
	}

	filteredFamilies := service.excludeFamily(familyCompaniesFound, id)

	if len(filteredFamilies) > 0 {
		var causes exceptions.Causes

		errMessage := fmt.Sprintf("family companies name: [%s] is duplicated", familyCompanies.Name)
		causes.Code = exceptions.FamilyCompaniesNameDuplicated
		err = exceptions.NewDuplicatedExceptionWithCause(errMessage, causes)
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, repositoryName, text.Functionality, "Update")

		return err
	}

	err = service.familyCompaniesRepository.Update(ctx, id, &familyCompanies)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateListMetricName)
		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateListMetricName)

	return nil
}

func (service *familyCompaniesService) excludeFamily(
	families []entities.FamilyCompanies,
	id string) []entities.FamilyCompanies {
	result := make([]entities.FamilyCompanies, 0)
	for _, op := range families {
		if !op.IsTheSame(id) {
			result = append(result, op)
		}
	}
	return result
}

func (service *familyCompaniesService) GetFamiliesCompaniesFromFilter(ctx context.Context,
	filter entities.FamilyCompaniesFilter) ([]entities.FamilyCompanies, error) {
	familiesFound, err := service.familyCompaniesRepository.GetFamilyCompanies(ctx, filter)
	if err != nil {
		return []entities.FamilyCompanies{}, err
	}

	return familiesFound, nil
}

func (service *familyCompaniesService) Delete(ctx context.Context, id string) error {
	pagination := entities.Pagination{}
	pagedRules, err := service.ruleRepository.
		FindRulesPaged(ctx, entities.RuleFilter{FamilyCompanyID: id}, pagination)
	if err != nil {
		return err
	}
	associatedRules := pagedRules.Data.([]entities.Rule)
	if len(associatedRules) > 0 {
		var causes exceptions.Causes
		message := fmt.Sprintf("family companies with id [%s], is associated with the rule [%s]",
			id,
			associatedRules[0].Description)
		causes.Code = exceptions.FamilyCompaniesAssociatedWithARule
		err = exceptions.NewAssociatedExceptionWithCause(message, causes)

		service.logs.Error(ctx, err.Error(), text.LogTagMethod, repositoryName, text.Functionality, "Delete")
		return err
	}

	err = service.familyCompaniesRepository.Delete(ctx, id)
	metricData := metrics.NewMetricData(ctx, "Delete", serviceName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteListMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteListMetricName)

	return nil
}
