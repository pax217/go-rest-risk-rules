package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	rulesString "github.com/conekta/risk-rules/pkg/strings"
	"github.com/mitchellh/mapstructure"
)

const (
	SUM      = "SUM"
	SUBTRACT = "SUBTRACT"
	MLP      = "MLP"
	DIV      = "DIV"
)

const (
	minFormulaFieldsContent     = 2
	countFormulaSubtractContent = 2
)

var mathOperationValues = map[string]bool{
	SUM:      true,
	SUBTRACT: true,
	MLP:      true,
	DIV:      true,
}

type FormulaContent struct {
	Fields        *[]string `json:"fields" bson:"fields"  mapstructure:"fields" `
	MathOperation *string   `json:"math_operation" bson:"math_operation"  mapstructure:"math_operation" `
}

func ValidateFormulas(rules []RuleContent) error {
	for _, rule := range rules {
		if hasValidFormulaContent(rule) {
			return fmt.Errorf(
				"formula should have both fields and math operation. Fields: %q, MathOperation: %q",
				rulesString.StringArrPointerToStringArr(rule.Fields),
				rulesString.StringPointerToString(rule.MathOperation))
		}

		if !isFormulaRule(rule) {
			continue
		}

		if len(*rule.Fields) < minFormulaFieldsContent {
			return errors.New("formula rule should have at least two fields")
		}

		if !hasValidMathOperation(*rule.MathOperation) {
			return fmt.Errorf(
				"formula should have valid math operation. Requested math operation: %s", *rule.MathOperation)
		}

		if *rule.MathOperation == SUBTRACT && len(*rule.Fields) != countFormulaSubtractContent {
			return fmt.Errorf(
				"formula with SUBTRACT should have only 2 values. Formula fields length = %d", len(*rule.Fields))
		}

		ruleMap, err := rule.FormulaContent.ToMap()
		if err != nil {
			return errors.New("could not convert rule formula content to map[string]interface{}")
		}

		switch fields := ruleMap["fields"].(type) {
		case *[]string:
			for _, field := range *fields {
				_, err := strconv.Atoi(field)
				if err == nil {
					return fmt.Errorf("formula field with value %q must be non numeric string",
						field)
				}
			}

		default:
			return errors.New("formula fields type is not *[]string")
		}
	}
	return nil
}

func isFormulaRule(rule RuleContent) bool {
	return rule.Fields != nil && rule.MathOperation != nil
}

func hasValidFormulaContent(rule RuleContent) bool {
	return (rule.Fields != nil && rulesString.IsStringPointerEmpty(rule.MathOperation)) ||
		(rule.Fields == nil && !rulesString.IsStringPointerEmpty(rule.MathOperation))
}

func hasValidMathOperation(operation string) bool {
	_, ok := mathOperationValues[operation]
	return ok
}

func GenerateRuleFieldWithFormulaFields(rules *[]RuleContent) {
	for i, rule := range *(rules) {
		if !isFormulaRule(rule) {
			continue
		}

		fieldsString := fmt.Sprintf("(%s)", strings.Join(*rule.Fields, ","))
		(*rules)[i].Field = fmt.Sprintf("%s %s", *rule.MathOperation, fieldsString)
	}
}

func (c *FormulaContent) ToMap() (map[string]interface{}, error) {
	var mapCharge map[string]interface{}

	err := mapstructure.Decode(c, &mapCharge)
	return mapCharge, err
}
