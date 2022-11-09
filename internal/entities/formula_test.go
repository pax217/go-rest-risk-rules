package entities_test

import (
	"fmt"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormula_ValidateFormula(t *testing.T) {
	t.Run("When rule is a valid formula, it should return no error", func(t *testing.T) {
		rules := testdata.GetDefaultFormulaRules()
		err := entities.ValidateFormulas(rules)

		assert.NoError(t, err)
	})

	t.Run("When rule is not formula, it should return no error", func(t *testing.T) {
		rules := testdata.GetNoFormulaRules()
		err := entities.ValidateFormulas(rules)

		assert.NoError(t, err)
	})

	t.Run("When formula has no MathOperation, it should return error", func(t *testing.T) {
		rules := testdata.GetNoMathOperationFormula()
		err := entities.ValidateFormulas(rules)

		assert.Error(t, err)
	})

	t.Run("When formula has no Fields, it should return error", func(t *testing.T) {
		rules := testdata.GetNoFieldsFormula()
		err := entities.ValidateFormulas(rules)

		assert.Error(t, err)
	})

	t.Run("When formula has only one Field, it should return error", func(t *testing.T) {
		rules := testdata.GetOneFieldFormula()
		err := entities.ValidateFormulas(rules)

		assert.Error(t, err)
	})

	t.Run("When formula has invalid MathOperation, it should return error", func(t *testing.T) {
		rules := testdata.GetInvalidMathOperationFormula()
		err := entities.ValidateFormulas(rules)

		assert.Error(t, err)
	})

	t.Run("When formula MathOperation is SUBTRACT and Fields size != 2, it should return error", func(t *testing.T) {
		rules := testdata.GetInvalidSubtractFormula()
		err := entities.ValidateFormulas(rules)

		assert.Error(t, err)
	})

	t.Run("When formula fields have numeric string value, it should return error", func(t *testing.T) {
		rules := testdata.GetNumericFieldsFormula()
		err := entities.ValidateFormulas(rules)

		assert.Error(t, err)
	})
}

func TestFormula_GenerateRuleFieldWithFormulaFields(t *testing.T) {
	t.Run("When formula has two fields, it should return a valid formula field", func(t *testing.T) {
		rules := testdata.GetDefaultFormulaRules()
		expectedField := fmt.Sprintf("%s (%s,%s)", *rules[0].MathOperation, (*rules[0].Fields)[0], (*rules[0].Fields)[1])
		entities.GenerateRuleFieldWithFormulaFields(&rules)

		assert.Equal(t, expectedField, rules[0].Field)
	})

	t.Run("When formula has three fields or more, it should return a valid formula field", func(t *testing.T) {
		rules := testdata.GetDefaultFormulaWithThreeFields()
		expectedField := fmt.Sprintf("%s (%s,%s,%s)", *rules[0].MathOperation, (*rules[0].Fields)[0], (*rules[0].Fields)[1], (*rules[0].Fields)[2])
		entities.GenerateRuleFieldWithFormulaFields(&rules)

		assert.Equal(t, expectedField, rules[0].Field)
	})
}
