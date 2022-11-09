package families

import (
	"context"
	"fmt"
	"strings"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/internal/entities/exceptions"
	"github.com/conekta/risk-rules/pkg/metrics"
	"github.com/conekta/risk-rules/pkg/text"
)

const serviceMethodName = "families.service.%s"

type FamilyService interface {
	Create(ctx context.Context, family entities.Family) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, family entities.Family) error
	Get(ctx context.Context, pagination entities.Pagination,
		filter entities.FamilyFilter) (interface{}, error)
	GetFamilyFromFilter(ctx context.Context, filter entities.FamilyFilter) (entities.Family, error)
	GetFamily(ctx context.Context, filter entities.FamilyFilter) (entities.Family, error)
}

type familyService struct {
	logs             logs.Logger
	datadog          datadog.Metricer
	config           config.Config
	familyRepository FamilyRepository
	ruleRepository   rules.RuleRepository
}

func NewFamilyService(cfg config.Config,
	familyRepository FamilyRepository,
	ruleRepository rules.RuleRepository,
	logger logs.Logger,
	datadogMetric datadog.Metricer) FamilyService {
	return &familyService{
		config:           cfg,
		familyRepository: familyRepository,
		ruleRepository:   ruleRepository,
		logs:             logger,
		datadog:          datadogMetric,
	}
}

func (service *familyService) Create(ctx context.Context, family entities.Family) error {
	metricData := metrics.NewMetricData(ctx, "AddFamily", serviceMethodName, service.config.Env)
	filter := entities.FamilyFilter{
		Mccs: family.Mccs,
		Name: family.Name,
	}

	familyFound, err := service.GetFamilyFromFilter(ctx, filter)
	if err != nil {
		return err
	}
	if !familyFound.IsEmpty() {
		return service.BuildExistingFamilyError(ctx, []entities.Family{familyFound}, family)
	}

	err = service.familyRepository.AddFamily(ctx, &family)

	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveListMetricName)
		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveListMetricName)

	return err
}

func (service *familyService) Update(ctx context.Context, id string, family entities.Family) error {
	metricData := metrics.NewMetricData(ctx, "Update", serviceMethodName, service.config.Env)
	filter := entities.FamilyFilter{
		Mccs: family.Mccs,
		Name: family.Name,
	}

	pagedResponseFamilies, err := service.familyRepository.SearchPaged(ctx, entities.NewDefaultPagination(), filter)
	if err != nil {
		return err
	}

	families := pagedResponseFamilies.Data.([]entities.Family)
	filteredFamilies := service.excludeFamily(families, id)

	if len(filteredFamilies) > 0 {
		return service.BuildExistingFamilyError(ctx, filteredFamilies, family)
	}

	err = service.familyRepository.Update(ctx, id, &family)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateListMetricName)
		return err
	}

	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateListMetricName)

	return nil
}

func (service *familyService) excludeFamily(families []entities.Family, id string) []entities.Family {
	result := make([]entities.Family, 0)
	for _, op := range families {
		if !op.IsTheSame(id) {
			result = append(result, op)
		}
	}
	return result
}

func (service *familyService) GetFamilyFromFilter(ctx context.Context, filter entities.FamilyFilter) (entities.Family, error) {
	familyFound, err := service.familyRepository.SearchPaged(ctx, entities.NewDefaultPagination(), filter)
	if err != nil {
		return entities.Family{}, err
	}

	entityFamilies := familyFound.Data.([]entities.Family)
	if len(entityFamilies) > 0 {
		return entityFamilies[0], nil
	}

	return entities.Family{}, nil
}
func (service *familyService) GetFamily(ctx context.Context, filter entities.FamilyFilter) (entities.Family, error) {
	families, err := service.familyRepository.SearchEvaluate(ctx, filter)
	if err != nil {
		return entities.Family{}, err
	}

	if len(families) > 0 {
		return families[0], nil
	}

	return entities.Family{}, nil
}
func (service *familyService) BuildExistingFamilyError(ctx context.Context,
	familiesFound []entities.Family, family entities.Family) error {
	var errMessage string
	var causes exceptions.Causes

	if len(familiesFound) == 1 {
		duplicatedMccs := family.SearchDuplicatedMcc(familiesFound[0])

		if len(duplicatedMccs) > 0 {
			errMessage = fmt.Sprintf("family: [%s] has MCCs [%s] duplicated",
				familiesFound[0].Name, strings.Join(duplicatedMccs, ","))
			causes.Code = exceptions.FamiliesMccsDuplicated
		} else {
			errMessage = fmt.Sprintf("family name: [%s] is duplicated", family.Name)
			causes.Code = exceptions.FamiliesNameDuplicated
		}
	} else {
		errMessage = fmt.Sprintf("there is more than one family with the name[%s] or the same mcc[%s]",
			family.Name, family.Mccs)
	}

	err := exceptions.NewDuplicatedExceptionWithCause(errMessage, causes)
	service.logs.Error(ctx, err.Error(), text.LogTagMethod, repositoryName, text.Functionality, "BuildExistingFamilyError")
	return err
}

func (service *familyService) Get(ctx context.Context, pagination entities.Pagination,
	filter entities.FamilyFilter) (interface{}, error) {
	if filter.Paged {
		return service.familyRepository.SearchPaged(ctx, pagination, filter)
	} else {
		return service.familyRepository.Search(ctx, filter)
	}
}

func (service *familyService) Delete(ctx context.Context, id string) error {
	pagination := entities.Pagination{}
	pagedRules, err := service.ruleRepository.FindRulesPaged(ctx, entities.RuleFilter{FamilyID: id}, pagination)
	if err != nil {
		return err
	}

	associatedRules := pagedRules.Data.([]entities.Rule)
	if len(associatedRules) > 0 {
		var causes exceptions.Causes

		message := fmt.Sprintf("family with id [%s], is associated with the rule [%s]",
			id,
			associatedRules[0].Description)
		causes.Code = exceptions.FamilyAssociatedWithARule
		err = exceptions.NewAssociatedExceptionWithCause(message, causes)

		service.logs.Error(ctx, err.Error(), text.LogTagMethod, repositoryName, text.Functionality, "Delete")
		return err
	}

	err = service.familyRepository.Delete(ctx, id)
	metricData := metrics.NewMetricData(ctx, "Delete", serviceMethodName, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteListMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteListMetricName)

	return nil
}
