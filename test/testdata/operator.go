package testdata

import (
	"time"

	"github.com/conekta/risk-rules/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOperatorDefault() entities.Operator {
	now := time.Now().UTC().Truncate(time.Millisecond)
	return entities.Operator{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now,
		CreatedBy:   "you",
		Description: "nothing else",
		Name:        "*",
		Title:       "Por",
		UpdatedAt:   &now,
		Type:        "type",
	}
}

func GetOperators() []entities.Operator {
	now := time.Now().UTC().Truncate(time.Millisecond)
	operators := make([]entities.Operator, 0)
	operators = append(operators, entities.Operator{
		ID:          primitive.ObjectID{},
		CreatedAt:   now,
		CreatedBy:   "you",
		Description: "nothing else",
		Name:        "*",
		Title:       "Más",
		UpdatedAt:   &now,
		Type:        "type",
	}, entities.Operator{
		ID:          primitive.ObjectID{},
		CreatedAt:   now,
		CreatedBy:   "me",
		Description: "nothing in",
		Name:        ">",
		Title:       "Mayor que",
		UpdatedAt:   &now,
		Type:        "type",
	})
	return operators
}

func GetOperatorNotValid() entities.Operator {
	return entities.Operator{
		CreatedBy:   "you",
		Description: "nothing else",
		Name:        "",
		Type:        "type",
	}
}

func GetOperatorRequest() entities.OperatorRequest {
	return entities.OperatorRequest{
		Author:      "carlos.maldonado@conekta.com",
		Name:        "+",
		Title:       "Más",
		Description: "indicates the sum",
		Type:        "string",
	}
}
