package entities

import (
	"fmt"
	"time"

	customString "github.com/conekta/go_common/strings"

	http "github.com/conekta/go_common/http/resterror"
	"github.com/conekta/risk-rules/pkg/strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	OperatorTypeNumber  = "number"
	OperatorTypeString  = "string"
	OperatorTypeBoolean = "boolean"
)

type Operator struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	Description string             `json:"description"  bson:"description"`
	Name        string             `json:"name" bson:"name"`
	Title       string             `json:"title" bson:"title"`
	Type        string             `json:"type" bson:"type"`
	UpdatedAt   *time.Time         `json:"updated_at" bson:"updated_at"`
	UpdatedBy   *string            `json:"updated_by" bson:"updated_by"`
}

type OperatorRequest struct {
	Author      string `json:"author" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Type        string `json:"type" validate:"required"`
}

type OperatorFilter struct {
	ID    string `query:"id"`
	Type  string `query:"type"`
	Name  string `query:"name"`
	Title string `query:"title"`
	Paged bool   `query:"paged"`
}

func (o *Operator) SetID(id primitive.ObjectID) {
	o.ID = id
}

func (operatorFilter *OperatorFilter) IsEmptyOperatorFilter() bool {
	return customString.IsEmpty(operatorFilter.ID) &&
		customString.IsEmpty(operatorFilter.Type) &&
		customString.IsEmpty(operatorFilter.Name) &&
		customString.IsEmpty(operatorFilter.Title)
}

func (request *OperatorRequest) NewOperatorFromPostRequest() Operator {
	now := time.Now().Truncate(time.Second)

	return Operator{
		CreatedAt:   now,
		CreatedBy:   request.Author,
		Description: request.Description,
		Name:        request.Name,
		Title:       request.Title,
		Type:        request.Type,
	}
}

func (request *OperatorRequest) NewModuleFromPutRequest() Operator {
	now := time.Now().Truncate(time.Second)

	return Operator{
		UpdatedBy:   &request.Author,
		Name:        request.Name,
		Title:       request.Title,
		Description: request.Description,
		UpdatedAt:   &now,
		Type:        request.Type,
	}
}

func (operatorFilter *OperatorFilter) Validate() error {
	if !strings.IsEmpty(operatorFilter.Type) {
		if !(operatorFilter.Type == OperatorTypeString ||
			operatorFilter.Type == OperatorTypeNumber ||
			operatorFilter.Type == OperatorTypeBoolean) {
			return http.NewBadRequestError(fmt.Sprintf("operator type value not found: %s", operatorFilter.Type))
		}
	}

	return nil
}

func (o *Operator) IsTheSame(id string) bool {
	ID, _ := primitive.ObjectIDFromHex(id)
	return o.ID == ID
}

func (o *Operator) SearchSimilarOperator(operators []Operator) Operator {
	for _, operator := range operators {
		if o.Name == operator.Name && o.Type == operator.Type {
			return operator
		}
	}

	return Operator{}
}

func (o *Operator) IsEmpty() bool {
	return o.Name == "" || o.Type == ""
}

func (o *Operator) ValidateOperatorID(id string) error {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = fmt.Errorf("error: invalid id: %s", id)
		return err
	}

	return nil
}
