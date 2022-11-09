package rules

import (
	"context"
	"fmt"

	"github.com/conekta/risk-rules/internal/entities/exceptions"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/metrics"
	"github.com/conekta/risk-rules/pkg/text"
)

const ruleServiceMethod = "rules.service.%s"

type RuleService interface {
	AddRule(ctx context.Context, rule entities.Rule) (entities.Rule, error)
	UpdateRule(ctx context.Context, ruleID string, ruleReq entities.Rule) error
	RemoveRule(ctx context.Context, ID string) error
	ListRules(ctx context.Context, ruleFilter entities.RuleFilter, pagination entities.Pagination) (entities.PagedResponse, error)
	BuildRule(ruleContent []entities.RuleContent) string
}

type ruleService struct {
	config         config.Config
	ruleRepository RuleRepository
	rules          RuleValidator
	logs           logs.Logger
	datadog        datadog.Metricer
}

func NewRulesService(cfg config.Config,
	rules RuleValidator,
	ruleRepository RuleRepository,
	logger logs.Logger,
	metric datadog.Metricer) RuleService {
	return &ruleService{
		config:         cfg,
		ruleRepository: ruleRepository,
		rules:          rules,
		logs:           logger,
		datadog:        metric,
	}
}

func (service *ruleService) AddRule(ctx context.Context, rule entities.Rule) (entities.Rule, error) {
	rule.Rule = service.BuildRule(rule.Rules)
	err := service.validate(ctx, rule)
	if err != nil {
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(ruleServiceMethod, "AddRule"))
		return rule, err
	}

	rulesFound, err := service.ruleRepository.FindRulesPaged(ctx, rule.GetRuleFilter(), entities.Pagination{})
	if err != nil {
		return entities.Rule{}, err
	}

	if rule.IsContained(rulesFound.Data.([]entities.Rule)) {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("the rule '%s' already exist", rule.Rule))
		service.logs.Error(ctx, err.Error())

		return entities.Rule{}, err
	}

	ruleResp, err := service.ruleRepository.AddRule(ctx, rule)
	metricData := metrics.NewMetricData(ctx, "Add", ruleServiceMethod, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveRuleMetricName)
		return rule, err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.SaveRuleMetricName)
	return ruleResp, err
}

func (service *ruleService) UpdateRule(ctx context.Context, ruleID string, rule entities.Rule) error {
	rule.Rule = service.BuildRule(rule.Rules)
	err := service.validate(ctx, rule)
	if err != nil {
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf(ruleServiceMethod, "UpdateRule"))
		return err
	}

	rulesFound, err := service.ruleRepository.FindRulesPaged(ctx, rule.GetRuleFilter(), entities.Pagination{})
	if err != nil {
		return err
	}

	if rule.IsContained(rulesFound.Data.([]entities.Rule)) && !service.isTheSame(rulesFound.Data.([]entities.Rule), ruleID) {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("the rule '%s' already exist", rule.Rule))
		service.logs.Error(ctx, err.Error())
		return err
	}

	err = service.ruleRepository.UpdateRule(ctx, ruleID, rule)
	metricData := metrics.NewMetricData(ctx, "Update", ruleServiceMethod, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateRuleMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.UpdateRuleMetricName)
	return nil
}

func (service *ruleService) RemoveRule(ctx context.Context, ruleID string) error {
	err := service.ruleRepository.RemoveRule(ctx, ruleID)
	metricData := metrics.NewMetricData(context.TODO(), "Remove", ruleServiceMethod, service.config.Env)
	if err != nil {
		metricData.SetResult(false)
		metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteRuleMetricName)
		return err
	}
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.datadog, service.logs, metricData, text.DeleteRuleMetricName)
	return nil
}

func (service *ruleService) ListRules(ctx context.Context, ruleFilter entities.RuleFilter,
	pagination entities.Pagination) (entities.PagedResponse, error) {
	return service.ruleRepository.FindRulesPaged(ctx, ruleFilter, pagination)
}

func (service *ruleService) validate(ctx context.Context, rule entities.Rule) error {
	entity := entities.ChargeRequest{}
	data, _ := entity.ToMap()

	_, err := service.rules.Evaluate(ctx, rule, data)

	return err
}

func (service *ruleService) BuildRule(ruleContent []entities.RuleContent) string {
	var ruleResult string
	for idx, rule := range ruleContent {
		isFirstOne := idx < 1
		ruleResult += rule.RuleAsString(isFirstOne)
	}

	return ruleResult
}

func (service *ruleService) isTheSame(rules []entities.Rule, ruleID string) bool {
	for _, rule := range rules {
		if rule.ID.Hex() == ruleID {
			return true
		}
	}

	return false
}
