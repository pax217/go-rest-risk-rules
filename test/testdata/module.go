package testdata

import (
	"time"

	"github.com/conekta/risk-rules/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDefaultModule() entities.Module {
	now := time.Now().UTC().Truncate(time.Millisecond)
	updatedBy := "santiago.ceron@conekta.com"
	return entities.Module{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now,
		UpdatedAt:   &now,
		CreatedBy:   "carlos.maldonado@conekta.com",
		UpdatedBy:   &updatedBy,
		Name:        "policy_compliance",
		Description: "Regla para validar contratos con Bancomer",
	}
}

func GetModules() []entities.Module {
	now := time.Now().UTC().Truncate(time.Millisecond)
	updatedAt := now.Truncate(time.Millisecond)
	modules := make([]entities.Module, 0)

	modules = append(modules, entities.Module{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now.Truncate(time.Millisecond),
		CreatedBy:   "carlos.maldonado@conekta.com",
		Description: "Regla para validar contratos con OXXO",
		Name:        "policy_compliance",
		UpdatedAt:   &updatedAt,
	}, entities.Module{
		ID:          primitive.NewObjectID(),
		CreatedAt:   now.Truncate(time.Millisecond),
		CreatedBy:   "carlos.maldonado@conekta.com",
		Description: "Regla para validar contratos con Bancomer",
		Name:        "country",
		UpdatedAt:   &updatedAt,
	})

	return modules
}
