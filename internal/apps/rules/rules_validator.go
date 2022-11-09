package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/text"

	"github.com/conekta/Conekta-Golang-Rules-Engine/parser"
)

const validatorServiceMethod = "rules_validator.service.%s"

type RuleValidator interface {
	Evaluate(ctx context.Context, rule entities.Rule, info map[string]interface{}) (bool, error)
}

type rulesValidator struct {
	logs logs.Logger
}

func NewRulesValidator(logger logs.Logger) RuleValidator {
	return &rulesValidator{
		logs: logger,
	}
}

func (v *rulesValidator) Evaluate(ctx context.Context, rule entities.Rule, info map[string]interface{}) (bool, error) {
	if rule.Rule == "" {
		return false, fmt.Errorf("[RuleValidator.evaluate] empty rule: [%v]", rule)
	}

	newRuleString := strings.ReplaceAll(rule.Rule, "\\", "")
	ev, err := parser.NewEvaluator(newRuleString)
	if err != nil {
		er := fmt.Errorf(
			"[RuleValidator.evaluate] On creating new evaluator for rule=[%+v] id=[%s] error=[%v]",
			rule.Rule, rule.ID, err)
		v.logs.Error(ctx, er.Error(), text.LogTagMethod, fmt.Sprintf(validatorServiceMethod, "evaluate"))

		return false, er
	}
	ans, err := ev.Process(info)
	if err != nil {
		er := fmt.Errorf(
			"[RuleValidator.evaluate] On process rule=[%v] id=[%s] error=[%v]",
			rule.Rule, rule.ID.Hex(), err)
		v.logs.Error(ctx, er.Error(), text.LogTagMethod, fmt.Sprintf(validatorServiceMethod, "process"))

		return false, er
	}
	if err = ev.LastDebugErr(); err != nil {
		er := fmt.Errorf(
			"[RuleValidator.evaluate] On Last debug error rule=[%v] id=[%s] error=[%v]",
			rule.Rule, rule.ID, err)
		v.logs.Error(ctx, er.Error(), text.LogTagMethod, fmt.Sprintf(validatorServiceMethod, "lastDebugErr"))

		return false, er
	}

	return ans, nil
}
