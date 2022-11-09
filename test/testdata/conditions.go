package testdata

import (
	"time"

	"github.com/conekta/risk-rules/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDefaultCondition() entities.Condition {
	now := time.Now().UTC().Truncate(time.Millisecond)
	return entities.Condition{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now.Truncate(time.Millisecond),
		CreatedBy:   "carlos.maldonado@conekta.com",
		Name:        "and",
		Description: "checks if two or more conditions are true",
	}
}

func GetConditions() []entities.Condition {
	now := time.Now().UTC().Truncate(time.Millisecond)
	conditions := make([]entities.Condition, 0)
	conditions = append(conditions, entities.Condition{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now.Truncate(time.Millisecond),
		CreatedBy:   "carlos.maldonado@conekta.com",
		Name:        "and",
		Description: "checks if two or more conditions are true",
	}, entities.Condition{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now.Truncate(time.Millisecond),
		CreatedBy:   "carlos.maldonado@conekta.com",
		Name:        "or",
		Description: "check if two or more conditions are met",
	}, entities.Condition{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now.Truncate(time.Millisecond),
		CreatedBy:   "carlos.maldonado@conekta.com",
		Name:        "equal",
		Description: "check if two or more conditions are equal",
	})

	return conditions
}
