package entities_test

import (
	"testing"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewModuleFromPutRequest(t *testing.T) {
	t.Run("check operator to update from ruleRequest", func(t *testing.T) {
		request := testdata.GetOperatorRequest()

		operatorToUpdate := request.NewModuleFromPutRequest()

		assert.False(t, operatorToUpdate.IsEmpty())
		assert.Equal(t, operatorToUpdate.UpdatedBy, &request.Author)
		assert.NotNil(t, operatorToUpdate.UpdatedAt)
		assert.Equal(t, operatorToUpdate.Type, request.Type)
		assert.Equal(t, operatorToUpdate.Title, request.Title)
		assert.Equal(t, operatorToUpdate.Name, request.Name)
		assert.NotNil(t, operatorToUpdate.Description, request.Description)
	})
}

func TestNewOperatorFromPostRequest(t *testing.T) {
	t.Run("check operator to create from ruleRequest", func(t *testing.T) {
		request := testdata.GetOperatorRequest()

		operatorToCreate := request.NewOperatorFromPostRequest()

		assert.False(t, operatorToCreate.IsEmpty())
		assert.NotNil(t, operatorToCreate.CreatedAt)
		assert.Equal(t, operatorToCreate.Type, request.Type)
		assert.Equal(t, operatorToCreate.Title, request.Title)
		assert.Equal(t, operatorToCreate.Name, request.Name)
		assert.NotNil(t, operatorToCreate.Description, request.Description)
	})
}

func TestValidateFilter(t *testing.T) {
	t.Run("check operator filter", func(t *testing.T) {
		filter := entities.OperatorFilter{
			ID:    "123",
			Name:  "+",
			Type:  "number",
			Title: "suma",
			Paged: true,
		}
		assert.Nil(t, filter.Validate())

		filter.Type = "invalid"
		assert.Error(t, filter.Validate())
	})
}

func TestValidateOperatorID(t *testing.T) {
	t.Run("check if is operator ID valid", func(t *testing.T) {
		operator := testdata.GetOperatorDefault()
		ID := operator.ID.Hex()
		assert.Nil(t, operator.ValidateOperatorID(ID))
		assert.True(t, operator.IsTheSame(ID))

		ID = "invalid"
		assert.Error(t, operator.ValidateOperatorID(ID))
	})
}

func TestSearchSimilarOperator(t *testing.T) {
	t.Run("check similar operator", func(t *testing.T) {
		operator := testdata.GetOperatorDefault()
		operators := make([]entities.Operator, 0)

		assert.Empty(t, operator.SearchSimilarOperator(operators))

		operators = append(operators, operator, operator)

		assert.NotEmpty(t, operator.SearchSimilarOperator(operators))
	})
}
