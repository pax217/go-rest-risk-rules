package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	White      TypeList = "Whitelist"
	Black      TypeList = "Blacklist"
	Gray       TypeList = "Graylist"
	EmailField          = "email"
)

type List struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CompanyID   string             `json:"company_id" bson:"company_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	Description string             `json:"description" bson:"description"`
	Decision    Decision           `json:"decision" bson:"decision"`
	Field       string             `json:"field" bson:"field"`
	IsGlobal    bool               `json:"is_global" bson:"is_global"`
	IsTest      bool               `json:"is_test" bson:"is_test"`
	Rule        string             `json:"rule" bson:"rule"`
	Type        string             `json:"type" bson:"type"`
	UpdatedAt   *time.Time         `json:"updated_at" bson:"updated_at"`
	UpdatedBy   *string            `json:"updated_by" bson:"updated_by"`
	Value       string             `json:"value" bson:"value"`
	TimeToLive  int64              `json:"time_to_live" bson:"time_to_live"`
	Expires     *time.Time         `json:"expires" bson:"expires"`
}

func (l *List) IsEmpty() bool { return l.CreatedBy == "" }

func (l *List) IsWhitelist() bool { return l.Type == White.String() }

func (l *List) IsBlacklist() bool { return l.Type == Black.String() }

func (l *List) IsGraylist() bool { return l.Type == Gray.String() }

type ListsSearch struct {
	Email     string `json:"email"`
	CardHash  string `json:"card_hash"`
	Phone     string `json:"phone"`
	CompanyID string `json:"company_id"`
}

type TypeList string

func (t TypeList) String() string { return string(t) }

type ListResponse struct {
	Decision      Decision `json:"decision"`
	TestDecision  Decision `json:"-" bson:"test_decision"`
	DecisionRules []List   `json:"decision_rules" bson:"decision_rules"`
	TestRules     []List   `json:"test_rules" bson:"test_rules"`
	Errors        []string `json:"errors"`
	Type          TypeList `json:"-"`
}

func NewListResponse() ListResponse {
	return ListResponse{
		Decision:      Undecided,
		TestDecision:  Undecided,
		DecisionRules: []List{},
		TestRules:     []List{},
		Errors:        []string{},
		Type:          "",
	}
}

func (d *ListResponse) GetResponses(typeList TypeList, desicion Decision) ListResponse {
	var listResponse = NewListResponse()
	d.getListResponse(&listResponse, typeList)
	d.getTestListResponse(&listResponse, typeList)
	listResponse.Decision = desicion
	listResponse.Type = typeList

	return listResponse
}

func (d *ListResponse) getListResponse(response *ListResponse, typeList TypeList) {
	for _, list := range d.DecisionRules {
		if list.Type == typeList.String() {
			response.DecisionRules = append(response.DecisionRules, list)
		}
	}
}

func (d *ListResponse) getTestListResponse(response *ListResponse, typeList TypeList) {
	for _, list := range d.TestRules {
		if list.Type == typeList.String() {
			response.TestRules = append(response.TestRules, list)
		}
	}
}

func (d *ListResponse) GetEntityName() string {
	return d.Type.String()
}

func (d *ListResponse) GetDecision() Decision {
	return d.Decision
}

func (d *ListResponse) GetTestDecision() Decision {
	return d.TestDecision
}

func (d *ListResponse) IsListResponseEmpty() bool {
	return len(d.DecisionRules) == 0 &&
		len(d.TestRules) == 0
}

func (l *List) IsValidListType() bool {
	return l.IsBlacklist() ||
		l.IsWhitelist() ||
		l.IsGraylist()
}
