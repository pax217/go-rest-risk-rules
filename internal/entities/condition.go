package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Condition struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at" `
	UpdatedAt   *time.Time         `json:"updated_at" bson:"updated_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	UpdatedBy   *string            `json:"updated_by" bson:"updated_by"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
}

type ConditionRequest struct {
	Author      string `json:"author" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type ConditionsFilter struct {
	ID string `query:"id"`
}

func (c ConditionRequest) NewConditionFromPostRequest() Condition {
	nowStr := time.Now().Truncate(time.Millisecond)

	return Condition{
		ID:          primitive.NewObjectID(),
		CreatedAt:   nowStr,
		CreatedBy:   c.Author,
		Name:        c.Name,
		Description: c.Description,
		UpdatedBy:   nil,
		UpdatedAt:   nil,
	}
}

func (c ConditionRequest) NewConditionFromPutRequest() Condition {
	nowStr := time.Now().Truncate(time.Millisecond)

	return Condition{
		ID:          primitive.NewObjectID(),
		Name:        c.Name,
		Description: c.Description,
		UpdatedBy:   &c.Author,
		UpdatedAt:   &nowStr,
	}
}
func (f *Condition) SetID(id primitive.ObjectID) {
	f.ID = id
}
