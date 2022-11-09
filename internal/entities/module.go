package entities

import (
	"time"

	"github.com/conekta/risk-rules/pkg/strings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Module struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   *time.Time         `json:"updated_at" bson:"updated_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	UpdatedBy   *string            `json:"updated_by" bson:"updated_by"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
}

type ModuleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Author      string `json:"author" validate:"required"`
}

type ModulesResponse struct {
	EvaluationOrder []string      `json:"evaluation_order"`
	WhiteList       ListResponse  `json:"whitelist"`
	BlackList       ListResponse  `json:"blacklist"`
	GrayList        ListResponse  `json:"graylist"`
	Rules           RulesResponse `json:"rules"`
}

type ModuleFilter struct {
	ID    string `query:"id"`
	Name  string `query:"name"`
	Paged bool   `query:"paged"`
}

func (m *ModuleFilter) IsEmptyModuleFilter() bool {
	return strings.IsEmpty(m.ID) && strings.IsEmpty(m.Name)
}

func (m *Module) SetID(id primitive.ObjectID) {
	m.ID = id
}

func (m *Module) GetModuleFilter(paged bool) ModuleFilter {
	var ID string
	if !m.ID.IsZero() {
		ID = m.ID.Hex()
	}

	return ModuleFilter{
		ID:    ID,
		Name:  m.Name,
		Paged: paged,
	}
}

func (m *Module) SearchSimilarModule(modules []Module) Module {
	for _, module := range modules {
		if m.Name == module.Name {
			return module
		}
	}

	return Module{}
}

func (request *ModuleRequest) NewModuleFromPostRequest() Module {
	now := time.Now().Truncate(time.Millisecond)

	return Module{
		CreatedBy:   request.Author,
		Name:        request.Name,
		Description: request.Description,
		CreatedAt:   now,
		UpdatedAt:   nil,
		UpdatedBy:   nil,
	}
}

func (request *ModuleRequest) NewModuleFromPutRequest() Module {
	now := time.Now().Truncate(time.Millisecond)

	return Module{
		ID:          primitive.NewObjectID(),
		UpdatedBy:   &request.Author,
		Name:        request.Name,
		Description: request.Description,
		UpdatedAt:   &now,
	}
}

func (m *Module) IsTheSame(id string) bool {
	ID, _ := primitive.ObjectIDFromHex(id)

	return m.ID == ID
}
