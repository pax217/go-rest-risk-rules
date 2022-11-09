package entities

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	customString "github.com/conekta/risk-rules/pkg/strings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	EntityName = "Rules"
	OperatorIN = "in"
)

type RuleContent struct {
	Field     string `json:"field" bson:"field" validate:"required_without=Fields,excluded_with=Fields"`
	Operator  string `json:"operator" bson:"operator" validate:"required"`
	Value     string `json:"value" bson:"value" validate:"required"`
	Condition string `json:"condition" bson:"condition" validate:"required"`
	Not       bool   `json:"not"  bson:"not"`
	FormulaContent
}

func (rc *RuleContent) IsOperatorIn() bool {
	return rc.Operator == OperatorIN
}

func (rc *RuleContent) FormatRuleValue() string {
	if rc.IsOperatorIn() {
		return rc.Value
	}
	hasLetters, _ := regexp.MatchString("[a-zA-Z]", rc.Value)
	if hasLetters && !customString.IsBoolean(rc.Value) {
		rc.Value = fmt.Sprintf("%q", rc.Value)
	}
	return rc.Value
}

func (rc *RuleContent) RuleAsString(isFirstOne bool) string {
	rc.Value = rc.FormatRuleValue()
	if isFirstOne {
		return strings.Trim(fmt.Sprintf("%s %s %s %s", rc.GetNot(), rc.Field, rc.Operator, rc.Value), " ")
	}
	if rc.Not {
		return fmt.Sprintf(" %s %s %s %s %s", rc.Condition, rc.GetNot(), rc.Field, rc.Operator, rc.Value)
	}
	return fmt.Sprintf(" %s %s %s %s", rc.Condition, rc.Field, rc.Operator, rc.Value)
}

func (rc *RuleContent) GetNot() string {
	if rc.Not {
		return "not"
	}
	return customString.Empty
}

type Rule struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       *time.Time         `json:"updated_at" bson:"updated_at"`
	CreatedBy       string             `json:"created_by" bson:"created_by"`
	UpdatedBy       *string            `json:"updated_by" bson:"updated_by"`
	IsTest          bool               `json:"is_test" bson:"is_test"`
	Module          string             `json:"module" bson:"module"`
	IsGlobal        bool               `json:"is_global" bson:"is_global"`
	Description     string             `json:"description" bson:"description"`
	CompanyID       *string            `json:"company_id" bson:"company_id"`
	FamilyMccID     *string            `json:"family_id" bson:"family_id"`
	FamilyCompanyID *string            `json:"family_company_id" bson:"family_company_id"`
	Rule            string             `json:"rule" bson:"rule"`
	Rules           []RuleContent      `json:"rules" bson:"rules"`
	Decision        Decision           `json:"decision" bson:"decision"`
	IsYellowFlag    bool               `json:"is_yellow_flag" bson:"is_yellow_flag"`
}

type RuleRequest struct {
	Decision        Decision      `json:"decision" validate:"required"`
	IsTest          *bool         `json:"is_test" validate:"required"`
	Module          string        `json:"module" validate:"required"`
	IsGlobal        *bool         `json:"is_global" validate:"required"`
	Description     string        `json:"description" validate:"required"`
	CompanyID       string        `json:"company_id"`
	FamilyID        string        `json:"family_id"`
	FamilyCompanyID string        `json:"family_company_id"`
	Rules           []RuleContent `json:"rules" validate:"required,gt=0,dive,required"`
	Author          string        `json:"author" validate:"required"`
	IsYellowFlag    bool          `json:"is_yellow_flag"`
}

func (rReq *RuleRequest) NewRuleFromPostRequest() Rule {
	now := time.Now().UTC()

	orderRulesWithFormulas(&rReq.Rules)
	GenerateRuleFieldWithFormulaFields(&rReq.Rules)

	return Rule{
		ID:              primitive.NewObjectID(),
		CreatedBy:       rReq.Author,
		CreatedAt:       now,
		IsTest:          *rReq.IsTest,
		Module:          rReq.Module,
		IsGlobal:        *rReq.IsGlobal,
		Description:     rReq.Description,
		CompanyID:       customString.StringToStringPointer(rReq.CompanyID),
		FamilyMccID:     customString.StringToStringPointer(rReq.FamilyID),
		FamilyCompanyID: customString.StringToStringPointer(rReq.FamilyCompanyID),
		Rules:           rReq.Rules,
		Decision:        rReq.Decision,
		IsYellowFlag:    rReq.IsYellowFlag,
	}
}

func (rReq *RuleRequest) NewRuleFromPutRequest() Rule {
	now := time.Now().UTC()

	orderRulesWithFormulas(&rReq.Rules)
	GenerateRuleFieldWithFormulaFields(&rReq.Rules)

	return Rule{
		UpdatedBy:       &rReq.Author,
		UpdatedAt:       &now,
		IsTest:          *rReq.IsTest,
		Module:          rReq.Module,
		IsGlobal:        *rReq.IsGlobal,
		Description:     rReq.Description,
		CompanyID:       customString.StringToStringPointer(rReq.CompanyID),
		FamilyCompanyID: customString.StringToStringPointer(rReq.FamilyCompanyID),
		FamilyMccID:     &rReq.FamilyID,
		Rules:           rReq.Rules,
		Decision:        rReq.Decision,
		IsYellowFlag:    rReq.IsYellowFlag,
	}
}

func (rReq *RuleRequest) HasMultipleValues() bool {
	return len(rReq.CompanyID) > 0 && len(rReq.FamilyID) > 0 ||
		len(rReq.CompanyID) > 0 && len(rReq.FamilyCompanyID) > 0 ||
		len(rReq.FamilyCompanyID) > 0 && len(rReq.FamilyID) > 0
}

func (rReq *RuleRequest) HasNoValue() bool {
	return customString.IsEmpty(rReq.CompanyID) && customString.IsEmpty(rReq.FamilyID) && customString.IsEmpty(rReq.FamilyCompanyID)
}

func (rReq *RuleRequest) Validate() error {
	err := ValidateIsGlobal(rReq)
	if err != nil {
		return err
	}

	err = ValidateDecision(rReq.Decision)
	if err != nil {
		return err
	}

	err = ValidateFormulas(rReq.Rules)
	if err != nil {
		return err
	}

	err = ValidateIsYellowFlag(rReq)
	if err != nil {
		return err
	}

	return nil
}

func ValidateIsYellowFlag(rReq *RuleRequest) error {
	if rReq.IsYellowFlag && rReq.Decision != Undecided {
		return errors.New("yellow flag rules must have the undecided decision")
	}
	return nil
}

func ValidateIsGlobal(rReq *RuleRequest) error {
	if !*rReq.IsGlobal {
		if rReq.HasMultipleValues() {
			return errors.New("non global rule, only one option: [family_mcc - company_id - family_companies] could be set at same time")
		}
		if rReq.HasNoValue() {
			return errors.New("non global rule, one option: [family_mcc - company_id - family_companies] have to be passed")
		}
	} else {
		if len(rReq.CompanyID) > 0 {
			return errors.New("company_id should not be passed, the rule is configured as Global")
		}

		if len(rReq.FamilyID) > 0 {
			return errors.New("family_mcc should not be passed, the rule is configured as Global")
		}

		if len(rReq.FamilyCompanyID) > 0 {
			return errors.New("family_companies should not be passed, the rule is configured as Global")
		}
	}

	return nil
}

func ValidateDecision(decision Decision) error {
	if !hasValidRuleDecision(decision) {
		return fmt.Errorf("decision value [%s], is not a valid value", decision)
	}

	return nil
}

func hasValidRuleDecision(decision Decision) bool {
	_, ok := ruleDecisionValues[decision]
	return ok
}

var ruleDecisionValues = map[Decision]bool{
	Accepted:  true,
	Declined:  true,
	Undecided: true,
}

type RuleUpdate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email" validate:"required,email"`
	IsTest      bool               `json:"is_test" bson:"is_test" validate:"required"`
	Module      string             `json:"module" bson:"module" validate:"required"`
	IsGlobal    bool               `json:"is_global" bson:"is_global" validate:"required"`
	Description string             `json:"description" bson:"description" validate:"required"`
	Rule        string             `json:"rule" bson:"rule"`
	Rules       []RuleContent      `json:"rules" bson:"rules" validate:"required,gt=0,dive,required"`
	Decision    string             `json:"decision" bson:"decision" validate:"required"`
}

type RulesResponse struct {
	Decision                Decision `json:"decision"`
	TestDecision            Decision `json:"-"`
	DecisionRules           []Rule   `json:"decision_rules" bson:"decision_rules"`
	TestRules               []Rule   `json:"test_rules" bson:"test_rules"`
	EvaluatedGlobalRules    int64    `json:"evaluated_global_rules" bson:"evaluated_global_rules"`
	EvaluatedNonGlobalRules int64    `json:"evaluated_non_global_rules" bson:"evaluated_non_global_rules"`
	Errors                  []string `json:"errors"`
}

type RuleFilter struct {
	ID                 string   `json:"id" query:"id"`
	IsTest             bool     `json:"is_test" query:"is_test"`
	Module             string   `json:"module" query:"module"`
	IsGlobal           bool     `json:"is_global" query:"is_global"`
	CompanyID          string   `json:"company_id" query:"company_id"`
	FamilyID           string   `json:"family_id" query:"family_id"`
	FamilyCompanyID    string   `json:"family_company_id" query:"family_company_id"`
	FamilyCompaniesIDs []string `json:"family_companies_ids" query:"family_companies_ids"`
	Rule               string   `json:"rule" query:"rule"`
}

func NewRulesResponse() RulesResponse {
	return RulesResponse{
		Decision:                Undecided,
		TestDecision:            Undecided,
		DecisionRules:           []Rule{},
		TestRules:               []Rule{},
		EvaluatedGlobalRules:    0,
		EvaluatedNonGlobalRules: 0,
		Errors:                  []string{},
	}
}

func (rResp *RulesResponse) GetDecision() Decision {
	if len(rResp.DecisionRules) > 0 {
		return rResp.Decision
	}
	return Undecided
}

func (rResp *RulesResponse) GetTestDecision() Decision {
	if len(rResp.TestRules) > 0 {
		return rResp.Decision
	}
	return Undecided
}

func (rResp *RulesResponse) GetEntityName() string {
	return EntityName
}

func (r *Rule) GetRuleFilter() RuleFilter {
	familyCompaniesIds := make([]string, 0)
	var companyID, familyID, familyCompanyID string

	if !customString.IsStringPointerEmpty(r.CompanyID) {
		companyID = *r.CompanyID
	}
	if !customString.IsStringPointerEmpty(r.FamilyMccID) {
		familyID = *r.FamilyMccID
	}
	if !customString.IsStringPointerEmpty(r.FamilyCompanyID) {
		familyCompanyID = *r.FamilyCompanyID
		familyCompaniesIds = append(familyCompaniesIds, *r.FamilyCompanyID)
	}
	return RuleFilter{
		Rule:               r.Rule,
		CompanyID:          companyID,
		FamilyID:           familyID,
		IsTest:             r.IsTest,
		IsGlobal:           r.IsGlobal,
		Module:             r.Module,
		FamilyCompanyID:    familyCompanyID,
		FamilyCompaniesIDs: familyCompaniesIds,
	}
}

func (r *Rule) IsContained(rules []Rule) bool {
	for _, rule := range rules {
		if rule.Rule == r.Rule {
			return true
		}
	}

	return false
}

func (s *RuleFilter) IsEmptyCompanyID() bool {
	return customString.IsEmpty(s.CompanyID)
}

func (s *RuleFilter) IsEmptyFamilyID() bool {
	return customString.IsEmpty(s.FamilyID)
}

func (s *RuleFilter) IsIDValid() bool {
	if !customString.IsEmpty(s.ID) {
		_, err := primitive.ObjectIDFromHex(s.ID)
		if err != nil {
			return false
		}
	}

	return true
}

func (s *RuleFilter) IsEmptyFamilyCompaniesID() bool {
	return customString.IsEmpty(s.FamilyCompanyID)
}

func orderRulesWithFormulas(rules *[]RuleContent) {
	sort.Slice(*rules, func(i, j int) bool {
		if (*rules)[i].FormulaContent.Fields != nil && (*rules)[i].FormulaContent.MathOperation != nil {
			return true
		}
		return false
	})
}

type RulesEvaluationResponse struct {
	Decision      string               `json:"decision"`
	RulesModules  RulesModulesResponse `json:"modules"`
	Omniscore     float64              `json:"omniscore"`
	MerchantScore float64              `json:"merchant_score"`
	Charge        ChargeRequest        `json:"charge"`
}

type RulesModulesResponse struct {
	CompanyRules       *RulesResponse `json:"company_rules,omitempty"`
	FamilyCompanyRules *RulesResponse `json:"family_company_rules,omitempty"`
	FamilyMccRules     *RulesResponse `json:"family_mcc_rules,omitempty"`
	GlobalRules        *RulesResponse `json:"global_rules,omitempty"`
	IdentityModule     *RulesResponse `json:"identity_module,omitempty"`
	YellowFlagModule   *RulesResponse `json:"yellow_flag_module,omitempty"`
}

func (rulesModulesResponse *RulesModulesResponse) SetRuleResponse(component Component, rulesResponse RulesResponse) {
	if len(rulesResponse.DecisionRules) > 0 {
		switch component.Name {
		case CompanyRulesType:
			rulesModulesResponse.CompanyRules = &rulesResponse
		case FamilyCompanyRulesType:
			rulesModulesResponse.FamilyCompanyRules = &rulesResponse
		case FamilyMccRulesType:
			rulesModulesResponse.FamilyMccRules = &rulesResponse
		case GlobalRulesType:
			rulesModulesResponse.GlobalRules = &rulesResponse
		case IdentityModuleType:
			rulesModulesResponse.IdentityModule = &rulesResponse
		case YellowFlagType:
			rulesModulesResponse.YellowFlagModule = &rulesResponse
		}
	}
}
