package entities_test

import (
	"github.com/conekta/risk-rules/internal/entities"
	"testing"

	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestCharge_ValidatePriorities(t *testing.T) {
	t.Run("Validate priorities in request", func(t *testing.T) {
		charge := testdata.GetDefaultChargeWithoutConsole()
		charge.ValidateConsole()

		assert.NotEmpty(t, charge.Console)
		assert.True(t, len(charge.Console) > 0)
	})
}

func TestCharge_SetDefaultConsoleOnlyRules(t *testing.T) {
	t.Run("Validate set default console only rules exclude list", func(t *testing.T) {
		chargeRequest := testdata.GetChargeConsoleRules()

		chargeRequest.ValidateConsoleOnlyRules()

		assert.Equal(t, 4, len(chargeRequest.Console))
	})

	t.Run("Validate charge when not have console", func(t *testing.T) {
		chargeRequest := testdata.GetChargeWithOutConsoleRules()

		chargeRequest.ValidateConsoleOnlyRules()

		assert.Equal(t, 4, len(chargeRequest.Console))
	})
}

func TestChargeRequest_NotHaveSecondaryDecision(t *testing.T) {
	t.Run("Validate when component not have secondary decision", func(t *testing.T) {
		component := entities.Component{
			Name:     entities.IdentityModuleType,
			Priority: []entities.Decision{entities.Undecided},
		}
		assert.False(t, component.HaveSecondaryDecision())
	})
}
