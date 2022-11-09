package testdata

import (
	"github.com/conekta/risk-rules/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GetDefaultFamilyCompanies() entities.FamilyCompanies {
	now := time.Now().UTC().Truncate(time.Millisecond)
	id, _ := primitive.ObjectIDFromHex("61e60c44097c21df49d3f50b")
	return entities.FamilyCompanies{
		ID:         id,
		Name:       "Name Family Companies",
		CompanyIDs: []string{"61e4dd6da5997ad4d9e76945", "61e4dd7320fbfc5f0849fba5"},
		CreatedAt:  now,
		CreatedBy:  "carlos.maldonado@conekta.com",
		UpdatedAt:  nil,
		UpdatedBy:  nil,
	}
}

func GetFamilyCompaniesRequest() entities.FamilyCompaniesRequest {
	return entities.FamilyCompaniesRequest{
		Name:       "Name Family",
		CompanyIDs: []string{"61e4dd6da5997ad4d9e76945", "61e4dd7320fbfc5f0849fba5"},
		Author:     "carlos.maldonado@conekta.com",
	}
}

func GetFamilyCompaniesRequestWithOutAuthor() entities.FamilyCompaniesRequest {
	return entities.FamilyCompaniesRequest{
		Name:       "Name Family",
		CompanyIDs: []string{"61e4dd6da5997ad4d9e76945", "61e4dd7320fbfc5f0849fba5"},
	}
}

func GetFamilyCompaniesRequestWithBadCompanyId() entities.FamilyCompaniesRequest {
	return entities.FamilyCompaniesRequest{
		Name:       "Name Family",
		CompanyIDs: []string{"7112821neuehud121eduihuq"},
		Author:     "carlos.maldonado@conekta.com",
	}
}

func GetFamilyCompaniesWithMatchingCompanyIDs() entities.FamilyCompanies {
	now := time.Now().UTC().Truncate(time.Millisecond)
	id, _ := primitive.ObjectIDFromHex("61e60c44097c21df49d3f50b")
	return entities.FamilyCompanies{
		ID:         id,
		Name:       "Name Family Companies",
		CompanyIDs: []string{"7683457364"},
		CreatedAt:  now,
		CreatedBy:  "carlos.maldonado@conekta.com",
		UpdatedAt:  nil,
		UpdatedBy:  nil,
	}
}

func GetDefaultFamilyCompaniesFilter() entities.FamilyCompaniesFilter {
	return entities.FamilyCompaniesFilter{
		ID:         "b3607ae794df433c6e39bfc8",
		CompanyIDs: []string{"7683457364", "7683457365", "7683457366"},
		Name:       "Test Family Companies",
		Paged:      false,
	}
}

func GetFamilyCompaniesUpdateRequest() entities.FamilyCompaniesRequest {
	return entities.FamilyCompaniesRequest{
		Name:       "Family Companies Name Updated",
		CompanyIDs: []string{"61eef7f302687a228c7cf24f", "61eef7fae492d9ba56d74e9f"},
		Author:     "carlos.maldonado@conekta.com",
	}
}

func GetFamilyCompaniesJsonRequestIsMalformed() string {
	request := `{
    "name": "Tiendas Electrónica",
    "company_ids": ["61e4dd6da5997ad4d9e76945","61eb21792e341c54221062b4","61eb217e66524deb95ad5143"],
    "author": "santiago.ceron@conekta.com",
	}`

	return request
}

func GetFamilyCompaniesRequestWithCompanyIdsNotValid() entities.FamilyCompaniesRequest {
	return entities.FamilyCompaniesRequest{
		Name:       "Name Family Companies",
		CompanyIDs: []string{"abc", "def"},
		Author:     "carlos.maldonado@conekta.com",
	}
}

func GetFamilyCompanies() []entities.FamilyCompanies {
	now := time.Now().UTC().Truncate(time.Millisecond)
	familyCompanies := []entities.FamilyCompanies{
		{
			ID:         primitive.NewObjectID(),
			Name:       "Tiendas Deportivas",
			CompanyIDs: []string{"61e4dd6da5997ad4d9e76945", "61e4dd7320fbfc5f0849fba5"},
			CreatedAt:  now.Truncate(time.Millisecond),
			CreatedBy:  "carlos.maldonado@conekta.com",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Tiendas Departamentales",
			CompanyIDs: []string{"61e9b80d6ab36bef5dc41da6", "61e9b811d94dd65161bb2a8d", "61e9b816eeba5ae7462d3a8d"},
			CreatedAt:  now.Truncate(time.Millisecond),
			CreatedBy:  "santiago.ceron@conekta.com",
		},
	}

	return familyCompanies
}

func GetFamilyCompaniesToAlphabeticalOrder() []entities.FamilyCompanies {
	now := time.Now().UTC().Truncate(time.Millisecond)
	familyCompanies := []entities.FamilyCompanies{
		{
			ID:         primitive.NewObjectID(),
			Name:       "Ropa y Calzado",
			CompanyIDs: []string{"61e4dd6da5997ad4d9e76945", "61e4dd7320fbfc5f0849fba5"},
			CreatedAt:  now.Truncate(time.Millisecond),
			CreatedBy:  "carlos.maldonado@conekta.com",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Farmacias",
			CompanyIDs: []string{"61e9b80d6ab36bef5dc41da6", "61e9b811d94dd65161bb2a8d", "61e9b816eeba5ae7462d3a8d"},
			CreatedAt:  now.Truncate(time.Millisecond),
			CreatedBy:  "santiago.ceron@conekta.com",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Tiendas Departamentales",
			CompanyIDs: []string{"61fd5006592b091af8d2de9c"},
			CreatedAt:  now.Truncate(time.Millisecond),
			CreatedBy:  "carlos.maldonado@conekta.com",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Joyerias",
			CompanyIDs: []string{"61fd503706aef018dacaf3ca"},
			CreatedAt:  now.Truncate(time.Millisecond),
			CreatedBy:  "carlos.maldonado@conekta.com",
		},
	}

	return familyCompanies
}
