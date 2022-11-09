package entities

import (
	"fmt"
	"time"

	"github.com/conekta/risk-rules/internal/entities/exceptions"

	http "github.com/conekta/go_common/http/resterror"
	customString "github.com/conekta/go_common/strings"
	"github.com/conekta/risk-rules/pkg/strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const mccLen = 4

type Family struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	Name              string             `json:"name" bson:"name"`
	Mccs              []string           `json:"mccs" bson:"mccs"`
	ExcludedCompanies []string           `json:"excluded_companies" bson:"excluded_companies"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	CreatedBy         string             `json:"created_by" bson:"created_by"`
	UpdatedAt         *time.Time         `json:"updated_at" bson:"updated_at"`
	UpdatedBy         *string            `json:"updated_by" bson:"updated_by"`
}

type FamilyRequest struct {
	Name              string   `json:"name"   validate:"required"`
	Mccs              []string `json:"mccs"   validate:"required"`
	ExcludedCompanies []string `json:"excluded_companies"`
	Author            string   `json:"author" validate:"required"`
}

type FamilyFilter struct {
	ID                   string   `query:"id"`
	Mccs                 []string `query:"mccs"`
	Name                 string   `query:"name"`
	Paged                bool     `query:"paged"`
	NotExcludedCompanies []string
}

func (famFilter *FamilyFilter) IsEmpty() bool {
	return customString.IsEmpty(famFilter.ID) && customString.IsEmpty(famFilter.Name) && len(famFilter.Mccs) == 0
}

func (f *Family) IsEmpty() bool {
	return customString.IsEmpty(f.Name) && len(f.Mccs) == 0
}

func (f *Family) SearchDuplicatedMcc(otherFamily Family) []string {
	duplicates := make([]string, 0)
	for _, mcc := range f.Mccs {
		for _, otherMcc := range otherFamily.Mccs {
			if mcc == otherMcc {
				duplicates = append(duplicates, mcc)
			}
		}
	}
	return duplicates
}

func (f *Family) IsTheSame(id string) bool {
	ID, _ := primitive.ObjectIDFromHex(id)
	return f.ID == ID
}

func (request *FamilyRequest) NewFamilyFromPostRequest() Family {
	now := time.Now().Truncate(time.Millisecond)
	return Family{
		ID:                primitive.NewObjectID(),
		Name:              request.Name,
		Mccs:              request.Mccs,
		ExcludedCompanies: request.ExcludedCompanies,
		CreatedBy:         request.Author,
		CreatedAt:         now,
	}
}

func (request *FamilyRequest) NewFamilyFromPutRequest() Family {
	now := time.Now().Truncate(time.Millisecond)
	return Family{
		Name:              request.Name,
		Mccs:              request.Mccs,
		ExcludedCompanies: request.ExcludedCompanies,
		UpdatedBy:         &request.Author,
		UpdatedAt:         &now,
	}
}

func (request *FamilyRequest) Validate() error {
	for _, s := range request.Mccs {
		_, err := strings.ToInt64(s, 0)
		if err != nil {
			return http.NewBadRequestError(fmt.Sprintf("family mcc [%s] is not number value", s))
		}

		if len(s) > mccLen {
			return http.NewBadRequestError(fmt.Sprintf("family mcc [%s] length must be 4 positions", s))
		}
	}
	if request.ExcludedCompanies != nil && len(request.ExcludedCompanies) > 0 {
		for i := range request.ExcludedCompanies {
			err := strings.IsHex(request.ExcludedCompanies[i])
			if err != nil {
				message := "a valid company id should have a hexadecimal format like a  622fb6f934089500011e270f"
				return exceptions.NewInvalidRequestWithCauses(fmt.Sprintf("company %s is not valid",
					request.ExcludedCompanies[i]), exceptions.Causes{Code: "001", Message: message})
			}
		}
	}
	return nil
}
