package entities

import (
	"time"

	"github.com/conekta/risk-rules/pkg/strings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Field struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   *time.Time         `json:"updated_at" bson:"updated_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	UpdatedBy   *string            `json:"updated_by" bson:"updated_by"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Description string             `json:"description" bson:"description" validate:"required"`
	Type        string             `json:"type" bson:"type" validate:"required"`
}

type FieldRequest struct {
	Name        string `json:"name" bson:"name" validate:"required"`
	Description string `json:"description" bson:"description" validate:"required"`
	Type        string `json:"type" bson:"type" validate:"required"`
	Author      string `json:"author" bson:"author" validate:"required"`
}

type FieldsFilter struct {
	ID    string `query:"id"`
	Name  string `query:"name"`
	Type  string `query:"type"`
	Paged bool   `query:"paged"`
}

func (request *FieldRequest) NewFieldFromPostRequest() Field {
	now := time.Now().Truncate(time.Millisecond)

	return Field{
		CreatedAt:   now,
		CreatedBy:   request.Author,
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		UpdatedBy:   nil,
		UpdatedAt:   nil,
	}
}

func (request *FieldRequest) NewFieldFromPutRequest() Field {
	now := time.Now().Truncate(time.Millisecond)

	return Field{
		UpdatedAt:   &now,
		UpdatedBy:   &request.Author,
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
	}
}

func (f *Field) IsEmpty() bool {
	return f.Name == "" || f.Type == ""
}

func (f *Field) SetID(id primitive.ObjectID) {
	f.ID = id
}

func (f *Field) GetFieldsFilter(paged bool) FieldsFilter {
	var ID string
	if !f.ID.IsZero() {
		ID = f.ID.Hex()
	}

	return FieldsFilter{
		ID:    ID,
		Name:  f.Name,
		Type:  f.Type,
		Paged: paged,
	}
}

func (f *Field) IsTheSame(id string) bool {
	ID, _ := primitive.ObjectIDFromHex(id)

	return f.ID == ID
}

func (f *FieldsFilter) IsEmptyFieldsFilter() bool {
	return strings.IsEmpty(f.ID) && strings.IsEmpty(f.Name) && strings.IsEmpty(f.Type)
}

func (f *Field) SearchSimilarField(fields []Field) Field {
	for _, field := range fields {
		if f.Name == field.Name && f.Type == field.Type {
			return field
		}
	}

	return Field{}
}
