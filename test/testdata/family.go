package testdata

import (
	"time"

	"github.com/conekta/risk-rules/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDefaultFamily() entities.Family {
	now := time.Now().UTC().Truncate(time.Millisecond)
	id, _ := primitive.ObjectIDFromHex("615324eb5bc1dea9ce66068f")
	return entities.Family{
		ID:                id,
		Name:              "Name Family",
		Mccs:              []string{"1111", "2222", "3333"},
		ExcludedCompanies: []string{"6262cb610211a64464781a5f", "6262cc2f14c692731e992f47"},
		CreatedAt:         now,
		CreatedBy:         "carlos.maldonado@conekta.com",
		UpdatedAt:         nil,
		UpdatedBy:         nil,
	}
}

func GetFamilyRequest() entities.FamilyRequest {
	return entities.FamilyRequest{
		Name:   "Name Family",
		Mccs:   []string{"1111", "2222", "3333"},
		Author: "carlos.maldonado@conekta.com",
	}
}

func GetFamilyUpdateRequest() entities.FamilyRequest {
	return entities.FamilyRequest{
		Name:   "Family Name Updated",
		Mccs:   []string{"4444", "5555", "6666"},
		Author: "carlos.maldonado@conekta.com",
	}
}

func GetFamilies() []entities.Family {
	now := time.Now().UTC().Truncate(time.Millisecond)
	families := []entities.Family{
		{
			ID:                primitive.NewObjectID(),
			Name:              "proveedores de tecnología",
			Mccs:              []string{"222", "333"},
			ExcludedCompanies: []string{"6262cb610211a64464781a5f"},
			CreatedAt:         now.Truncate(time.Millisecond),
			CreatedBy:         "carlos.maldonado@conekta.com",
		},
		{
			ID:                primitive.NewObjectID(),
			Name:              "familia de vendedoras de frutas",
			Mccs:              []string{"1111", "2222", "3333"},
			ExcludedCompanies: []string{"6262cc2f14c692731e992f47"},
			CreatedAt:         now.Truncate(time.Millisecond),
			CreatedBy:         "santiago.ceron@conekta.com",
		},
	}

	return families
}

func GetFamilyRequestWithMccCharacters() entities.FamilyRequest {
	return entities.FamilyRequest{
		Name:   "Name Family",
		Mccs:   []string{"ABCD"},
		Author: "carlos.maldonado@conekta.com",
	}
}

func GetFamilyRequestWithMccLengthIsDifferentForm4() entities.FamilyRequest {
	return entities.FamilyRequest{
		Name:   "Name Family",
		Mccs:   []string{"12345"},
		Author: "carlos.maldonado@conekta.com",
	}
}
