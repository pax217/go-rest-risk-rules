package rules_test

import (
	"context"
	"testing"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/test/testdata"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestEvaluation(t *testing.T) {
	logger, _ := logs.New()

	charge := entities.ChargeRequest{
		Amount:              28,
		DeviceFingerprint:   "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "60ad5c44926c8400016cbfdc",
		MonthlyInstallments: 10,
		LiveMode:            false,
		Details: entities.DetailsRequest{
			Email:     "eliosf27@gmail.com12",
			IPAddress: "127.0.0.1",
			Phone:     "+52477266334212",
			Name:      "de M12",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "cGfNEDJZjyj",
			Country:  "US",
		},
		Aggregation: testdata.GetDefaultAggregation(),
	}
	chargeMap, err := charge.ToMap()
	if err != nil {
		panic(err)
	}

	rulesValidator := rules.NewRulesValidator(logger)
	tests := []struct {
		name string
		rule string
		want bool
	}{
		{
			name: "not equal",
			rule: "not monthly_installments eq 1500",
			want: true,
		},
		{
			name: "equal",
			rule: "monthly_installments eq 10",
			want: true,
		},
		{
			name: "equal",
			rule: "monthly_installments == 10",
			want: true,
		},
		{
			name: "not less than",
			rule: "not monthly_installments lt 12",
			want: false,
		},
		{
			name: "less than",
			rule: "monthly_installments lt 12",
			want: true,
		},
		{
			name: "less than",
			rule: "monthly_installments < 12",
			want: true,
		},
		{
			name: "not less than equal to",
			rule: "not monthly_installments le 12",
			want: false,
		},
		{
			name: "less than equal to",
			rule: "monthly_installments le 12",
			want: true,
		},
		{
			name: "less than equal to",
			rule: "monthly_installments <= 12",
			want: true,
		},
		{
			name: "not greater than",
			rule: "not monthly_installments gt 12",
			want: true,
		},
		{
			name: "greater than",
			rule: "monthly_installments gt 12",
			want: false,
		},
		{
			name: "greater than",
			rule: "monthly_installments > 9",
			want: true,
		},
		{
			name: "not greater than equal to",
			rule: "not monthly_installments ge 12",
			want: true,
		},
		{
			name: "greater than equal to",
			rule: "monthly_installments ge 10",
			want: true,
		},
		{
			name: "greater than equal to",
			rule: "monthly_installments >= 10",
			want: true,
		},
		{
			name: "not contains",
			rule: "not company_id co \"test\"",
			want: true,
		},
		{
			name: "contains",
			rule: "company_id co \"test\"",
			want: false,
		},
		{
			name: "not ends with",
			rule: "not company_id ew \"test\"",
			want: true,
		},
		{
			name: "ends with",
			rule: "company_id ew \"test\"",
			want: false,
		},
		{
			name: "not in a list",
			rule: "not monthly_installments in [12]",
			want: true,
		},
		{
			name: "in a list",
			rule: "monthly_installments in [12]",
			want: false,
		},
		{
			name: "string in a list",
			rule: `payment_method.country in ["mx", "us"]`,
			want: false,
		},
		{
			name: "is not live mode",
			rule: "live_mode eq false",
			want: true,
		},
		{
			name: "payer has more than 20 MXN in 1 hours",
			rule: "aggregation.payer.charge.h1.sum > 20",
			want: true,
		},
		{
			name: "payer has more than 20 charges with the same merchant in 1 hours",
			rule: "aggregation.payer_company.charge.h1.count > 20",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := rulesValidator.Evaluate(context.TODO(), entities.Rule{Rule: tt.rule}, chargeMap)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, isValid, tt.name)
		})
	}
}

func TestWrongEvaluation(t *testing.T) {
	logger, _ := logs.New()

	charge := entities.ChargeRequest{
		Amount:              28,
		DeviceFingerprint:   "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "60ad5c44926c8400016cbfdc",
		MonthlyInstallments: 10,
		Details: entities.DetailsRequest{
			Email:     "eliosf27@gmail.com12",
			IPAddress: "127.0.0.1",
			Phone:     "+52477266334212",
			Name:      "de M12",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "cGfNEDJZjyj",
			Country:  "US",
		},
	}
	chargeMap, err := charge.ToMap()
	if err != nil {
		panic(err)
	}

	rulesValidator := rules.NewRulesValidator(logger)
	tests := []struct {
		name string
		rule string
		want string
	}{
		{
			name: "not equal",
			rule: "not xxx eq 1500",
			want: "Eval operand missing in input object",
		},
		{
			name: "not equal",
			rule: "not xxx etrrtrtrt 1500",
			want: "invalid memory address or nil pointer dereference",
		},
		{
			name: "not equal",
			rule: "not xxx eq \"1500\"",
			want: "Operand not present",
		},
		{
			name: "not equal",
			rule: "",
			want: "empty rule",
		},
		{
			name: "not equal",
			rule: "aasasasasasasasasasasasasasasasassasas",
			want: "invalid rule",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rulesValidator.Evaluate(context.TODO(), entities.Rule{Rule: tt.rule}, chargeMap)
			assert.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestMathOperations(t *testing.T) {
	logger, _ := logs.New()
	t.Run("subtract operation", func(t *testing.T) {
		charge := testdata.GetDefaultCharge()

		chargeMap, err := charge.ToMap()
		if err != nil {
			panic(err)
		}

		rulesValidator := rules.NewRulesValidator(logger)
		rule := "SUBTRACT (aggregation.payer_company.charge.h1.count,aggregation.payer_company.charge.h2.count) EQ -2"
		result, err := rulesValidator.Evaluate(context.TODO(), entities.Rule{Rule: rule}, chargeMap)

		assert.True(t, result, "the result should be -2")
		assert.NoError(t, err)
	})

	t.Run("division operation", func(t *testing.T) {
		charge := testdata.GetDefaultCharge()

		chargeMap, err := charge.ToMap()
		if err != nil {
			panic(err)
		}

		rulesValidator := rules.NewRulesValidator(logger)
		rule := "DIV (aggregation.payer.charge.h12.count,aggregation.payer.charge.h1.count) EQ 3"
		result, err := rulesValidator.Evaluate(context.TODO(), entities.Rule{Rule: rule}, chargeMap)

		assert.True(t, result, "the result should be 3")
		assert.NoError(t, err)
	})

	t.Run("division by zero operation", func(t *testing.T) {
		charge := testdata.GetDefaultCharge()
		charge.Aggregation = testdata.GetZeroDivisionAggregation()

		chargeMap, err := charge.ToMap()
		if err != nil {
			panic(err)
		}

		rulesValidator := rules.NewRulesValidator(logger)
		rule := "DIV (aggregation.payer_company.charge.h1.count,aggregation.payer_company.charge.h12.count) EQ 0"
		result, err := rulesValidator.Evaluate(context.TODO(), entities.Rule{Rule: rule}, chargeMap)

		assert.True(t, result, "the result should be 0")
		assert.NoError(t, err)
	})
}
