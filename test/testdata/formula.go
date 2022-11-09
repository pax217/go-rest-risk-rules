package testdata

import (
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/strings"
)

func GetDefaultFormulaRules() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"amount.h9", "amount.h12"},
				MathOperation: strings.StringToStringPointer(entities.SUM),
			},
		},
	}
}

func GetNoFormulaRules() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Field:     "amount",
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        nil,
				MathOperation: nil,
			},
		},
	}
}

func GetNoMathOperationFormula() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"amount.h9", "amount.h12"},
				MathOperation: nil,
			},
		},
	}
}

func GetNoFieldsFormula() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				MathOperation: strings.StringToStringPointer("SUM"),
			},
		},
	}
}

func GetOneFieldFormula() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"amount.h9"},
				MathOperation: strings.StringToStringPointer(entities.SUM),
			},
		},
	}
}

func GetInvalidMathOperationFormula() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"amount.h9"},
				MathOperation: strings.StringToStringPointer("SUMS"),
			},
		},
	}
}

func GetInvalidSubtractFormula() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"amount.h6", "amount.h9", "amount.h12"},
				MathOperation: strings.StringToStringPointer(entities.SUBTRACT),
			},
		},
	}
}

func GetNumericFieldsFormula() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"3", "9"},
				MathOperation: strings.StringToStringPointer(entities.SUBTRACT),
			},
		},
	}
}

func GetDefaultFormulaWithThreeFields() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"amount.h6", "amount.h9", "amount.h12"},
				MathOperation: strings.StringToStringPointer(entities.SUM),
			},
		},
	}
}

func GetRulesWithFormulaFieldsAndRuleField() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Field:     "amount",
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"amount.h6", "amount.h9", "amount.h12"},
				MathOperation: strings.StringToStringPointer(entities.SUM),
			},
		},
	}
}

func GetDefaultFormulaRulesWithProcessedField() []entities.RuleContent {
	return []entities.RuleContent{
		{
			Field:     "SUM (amount.h9,amount.h12)",
			Operator:  "==",
			Value:     "1",
			Condition: "and",
			Not:       false,
			FormulaContent: entities.FormulaContent{
				Fields:        &[]string{"amount.h9", "amount.h12"},
				MathOperation: strings.StringToStringPointer(entities.SUM),
			},
		},
	}
}
