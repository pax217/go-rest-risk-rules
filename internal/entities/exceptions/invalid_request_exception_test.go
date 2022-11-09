package exceptions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInvalidRequestWithCauses(t *testing.T) {

	exc := NewInvalidRequestWithCauses("a valid company id should have a hexadecimal format like a "+
		"622fb6f934089500011e270a", Causes{
		Code:    "001",
		Message: "reason",
	})

	assert.NotNil(t, exc)
	assert.NotEmpty(t, exc.Error())
	assert.True(t, exc.IsInvalidRequestException())
	assert.Equal(t, "001", exc.Causes().Code)
	assert.Equal(t, "reason", exc.Causes().Message)
	assert.Equal(t, "a valid company id should have a hexadecimal format like a "+
		"622fb6f934089500011e270a", exc.Error())
}
func TestNewInvalidRequest(t *testing.T) {

	exc := NewInvalidRequest("a valid company id should have a hexadecimal format like a " +
		"622fb6f934089500011e270a")

	assert.NotNil(t, exc)
	assert.NotEmpty(t, exc.Error())
	assert.True(t, exc.IsInvalidRequestException())
	assert.Empty(t, exc.Causes().Message)
	assert.Empty(t, exc.Causes().Code)
	assert.Equal(t, "a valid company id should have a hexadecimal format like a "+
		"622fb6f934089500011e270a", exc.Error())
}
