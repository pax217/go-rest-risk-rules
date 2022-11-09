package entities

import (
	customString "github.com/conekta/go_common/strings"
)

const (
	Accepted  Decision = "A"
	Declined  Decision = "D"
	Undecided Decision = "UN"
)

type EvaluationResult interface {
	GetDecision() Decision
	GetTestDecision() Decision
	GetEntityName() string
}

type EvaluationResults []EvaluationResult

func (eResult EvaluationResults) Len() int { return len(eResult) }

func (eResult EvaluationResults) Swap(i, j int) {
	eResult[i], eResult[j] = eResult[j], eResult[i]
}

type Decision string

func (decision Decision) String() string { return string(decision) }

func (decision *Decision) HasNoDecision() bool { return decision.String() == Undecided.String() }

func (decision Decision) ValidateDecision() Decision {
	if customString.IsEmpty(decision.String()) {
		return Undecided
	}

	return decision
}

type EvaluationResponse struct {
	Decision string          `json:"decision"`
	Modules  ModulesResponse `json:"modules"`
	Charge   ChargeRequest   `json:"charge"`
}

func NewUndecidedEvaluationResponse(charge ChargeRequest, evaluationOrder []string) EvaluationResponse {
	return EvaluationResponse{
		Decision: Undecided.String(),
		Modules: ModulesResponse{
			EvaluationOrder: evaluationOrder,
		},
		Charge: charge,
	}
}

func NewUndecidedEvaluationResponseOnlyRules(charge ChargeRequest) RulesEvaluationResponse {
	return RulesEvaluationResponse{
		Decision:      Undecided.String(),
		RulesModules:  RulesModulesResponse{},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        charge,
	}
}
