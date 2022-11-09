package entities_test

import (
	"fmt"
	"testing"

	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestValidateFamily(t *testing.T) {
	t.Run("excluded companies are not valid", func(t *testing.T) {
		companyID := "622fb6f934089500011e270"
		familyRequest := entities.FamilyRequest{
			Name:              "family test",
			Mccs:              []string{"7874"},
			ExcludedCompanies: []string{companyID},
			Author:            "me@gmail.com",
		}

		assert.Error(t, familyRequest.Validate())
		assert.EqualError(t, familyRequest.Validate(), fmt.Sprintf("company %s is not valid", companyID))
	})

	t.Run("excluded companies are not valid", func(t *testing.T) {
		companyID := "1234"
		familyRequest := entities.FamilyRequest{
			Name:              "family test",
			Mccs:              []string{"7874"},
			ExcludedCompanies: []string{companyID},
			Author:            "me@gmail.com",
		}

		assert.Error(t, familyRequest.Validate())
		assert.EqualError(t, familyRequest.Validate(), fmt.Sprintf("company %s is not valid", companyID))
	})
	t.Run("excluded companies are valid", func(t *testing.T) {
		companyID := "622fb6f934089500011e270f"
		familyRequest := entities.FamilyRequest{
			Name:              "family test",
			Mccs:              []string{"7874"},
			ExcludedCompanies: []string{companyID},
			Author:            "me@gmail.com",
		}

		assert.NoError(t, familyRequest.Validate())
	})
	t.Run("excluded companies are empty", func(t *testing.T) {
		familyRequest := entities.FamilyRequest{
			Name:              "family test",
			Mccs:              []string{"7874"},
			ExcludedCompanies: []string{},
			Author:            "me@gmail.com",
		}

		assert.NoError(t, familyRequest.Validate())
	})
	t.Run("excluded companies are nil", func(t *testing.T) {
		familyRequest := entities.FamilyRequest{
			Name:              "family test",
			Mccs:              []string{"7874"},
			ExcludedCompanies: nil,
			Author:            "me@gmail.com",
		}

		assert.NoError(t, familyRequest.Validate())
	})
	t.Run("invalid mccs", func(t *testing.T) {
		familyRequest := entities.FamilyRequest{
			Name:              "family test",
			Mccs:              []string{"a", "5468"},
			ExcludedCompanies: nil,
			Author:            "me@gmail.com",
		}

		assert.Error(t, familyRequest.Validate())
		assert.Contains(t, familyRequest.Validate().Error(), fmt.Sprintf("family mcc [%s] is not number value", familyRequest.Mccs[0]))
	})

	t.Run("invalid len of mccs", func(t *testing.T) {
		familyRequest := entities.FamilyRequest{
			Name:              "family test",
			Mccs:              []string{"41532", "5468"},
			ExcludedCompanies: nil,
			Author:            "me@gmail.com",
		}

		assert.Error(t, familyRequest.Validate())
		assert.Contains(t, familyRequest.Validate().Error(), fmt.Sprintf("family mcc [%s] length must be 4 positions", familyRequest.Mccs[0]))
	})
}

func TestIsTheSame(t *testing.T) {
	t.Run("when is the same ID", func(t *testing.T) {
		family := testdata.GetDefaultFamily()
		assert.True(t, family.IsTheSame(family.ID.Hex()))
	})
}

func TestFamilyIsEmpty(t *testing.T) {
	t.Run("when family is empty", func(t *testing.T) {
		family := entities.Family{}
		assert.True(t, family.IsEmpty())
	})

	t.Run("when family has a name", func(t *testing.T) {
		family := entities.Family{Name: "pglow"}
		assert.False(t, family.IsEmpty())
	})
}

func TestNewFamilyFromPostRequest(t *testing.T) {
	request := testdata.GetFamilyRequest()
	response := request.NewFamilyFromPostRequest()

	assert.NotEmpty(t, response)
	assert.NotNil(t, response.ID)
	assert.Equal(t, request.Name, response.Name)
	assert.Equal(t, request.Mccs, response.Mccs)
	assert.Equal(t, request.ExcludedCompanies, response.ExcludedCompanies)
}

func TestSearchDuplicatedMcc(t *testing.T) {
	request := testdata.GetDefaultFamily()

	assert.NotEmpty(t, request.SearchDuplicatedMcc(request))

	request = entities.Family{}
	assert.Empty(t, request.SearchDuplicatedMcc(request))
}
