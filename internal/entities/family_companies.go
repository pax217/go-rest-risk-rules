package entities

import (
	"fmt"
	"regexp"
	"time"

	customString "github.com/conekta/risk-rules/pkg/strings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RegexMongoID = `^[0-9a-fA-F]{24}$`

type FamilyCompanies struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	CompanyIDs []string           `json:"company_ids" bson:"company_ids"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	CreatedBy  string             `json:"created_by" bson:"created_by"`
	UpdatedAt  *time.Time         `json:"updated_at" bson:"updated_at"`
	UpdatedBy  *string            `json:"updated_by" bson:"updated_by"`
}

type FamilyCompaniesRequest struct {
	Name       string   `json:"name"   validate:"required"`
	CompanyIDs []string `json:"company_ids"   validate:"required"`
	Author     string   `json:"author" validate:"required"`
}
type FamilyCompaniesFilter struct {
	ID         string   `query:"id"`
	CompanyIDs []string `query:"company_ids"`
	Name       string   `query:"name"`
	Paged      bool     `query:"paged"`
}

func (familyCompanies *FamilyCompanies) IsTheSame(id string) bool {
	ID, _ := primitive.ObjectIDFromHex(id)
	return familyCompanies.ID == ID
}

func (famFilter *FamilyCompaniesFilter) IsEmpty() bool {
	return customString.IsEmpty(famFilter.ID) && customString.IsEmpty(famFilter.Name) && len(famFilter.CompanyIDs) == 0
}

func (request *FamilyCompaniesRequest) Validate() error {
	for _, companyID := range request.CompanyIDs {
		compile := regexp.MustCompile(RegexMongoID)

		isValidID := compile.MatchString(companyID)
		if !isValidID {
			return fmt.Errorf("company id [%s] is not a valid format of type mongo id", companyID)
		}
	}

	return nil
}

func (request *FamilyCompaniesRequest) NewFamilyCompaniesFromPostRequest() FamilyCompanies {
	now := time.Now().Truncate(time.Millisecond)
	return FamilyCompanies{
		ID:         primitive.NewObjectID(),
		Name:       request.Name,
		CreatedBy:  request.Author,
		CreatedAt:  now,
		CompanyIDs: request.CompanyIDs,
	}
}

func (request *FamilyCompaniesRequest) NewFamilyCompaniesFromPutRequest() FamilyCompanies {
	now := time.Now().Truncate(time.Millisecond)
	return FamilyCompanies{
		Name:       request.Name,
		CompanyIDs: request.CompanyIDs,
		UpdatedBy:  &request.Author,
		UpdatedAt:  &now,
	}
}
