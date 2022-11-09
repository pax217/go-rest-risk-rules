package testdata

import (
	"time"

	"github.com/conekta/risk-rules/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDefaultField() entities.Field {
	now := time.Now().Truncate(time.Millisecond)
	return entities.Field{
		ID:          primitive.ObjectID{},
		CreatedAt:   time.Now(),
		CreatedBy:   "carlos.maldonado@conekta.com",
		UpdatedAt:   &now,
		UpdatedBy:   nil,
		Description: "Representa el campo email perteneciente al charge",
		Name:        "email",
		Type:        "string",
	}
}

func GetDefaultArrayFields() []entities.Field {
	fields := make([]entities.Field, 0)
	now := time.Now().Truncate(time.Millisecond)

	fields = append(fields, entities.Field{
		ID:          primitive.ObjectID{},
		CreatedAt:   time.Now(),
		CreatedBy:   "carlos.maldonado@conekta.com",
		UpdatedAt:   &now,
		UpdatedBy:   nil,
		Description: "Representa el campo email perteneciente al charge",
		Name:        "email",
		Type:        "string",
	})

	return fields
}

func GetFields() []entities.Field {
	return []entities.Field{
		{
			ID:          primitive.ObjectID{},
			CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
			CreatedBy:   "santiago@conekta.com",
			Description: "buyer email",
			Name:        "details.email",
			Type:        entities.OperatorTypeString,
		},
		{
			ID:          primitive.ObjectID{},
			CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
			CreatedBy:   "santiago@conekta.com",
			Description: "Cantidad de dinero en transacciones en los últimos 3 meses, asociando comprador y merchant",
			Name:        "aggregation.payer_company.charge.h2880.sum",
			Type:        entities.OperatorTypeNumber,
		},
	}
}

func GetFieldRequest() entities.FieldRequest {
	return entities.FieldRequest{
		Author:      "you@conekta.com",
		Description: "nothing",
		Name:        "email",
		Type:        "string",
	}
}

func GetFieldRequestNotValid() entities.FieldRequest {
	return entities.FieldRequest{
		Author:      "you@conekta.com",
		Description: "nothing",
		Name:        "",
		Type:        "string",
	}
}

func GetFielNotValid() entities.Field {
	return entities.Field{
		ID:          primitive.NewObjectID(),
		Description: "nothing",
		Name:        "",
		Type:        "string",
	}
}
