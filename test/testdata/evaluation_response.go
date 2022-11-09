package testdata

import (
	"github.com/conekta/risk-rules/internal/entities"
)

var evaluationOrder = []string{
	"blacklist",
	"whitelist",
	"rules",
}

func GetEvaluationResponseSuccessful() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Undecided.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{},
		},
		Charge: GetDefaultCharge(),
	}
}

func GetRulesEvaluationResponseUndecidedCauseConsoleIsEmptyOnlyRules() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Undecided.String(),
		RulesModules: entities.RulesModulesResponse{
			CompanyRules:       nil,
			FamilyCompanyRules: nil,
			FamilyMccRules:     nil,
			GlobalRules:        nil,
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeConsoleIsEmptyOnlyRules(),
	}
}

func GetRulesEvaluationResponseAcceptedCauseConsoleCompanyRules() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Accepted.String(),
		RulesModules: entities.RulesModulesResponse{
			CompanyRules: &entities.RulesResponse{
				Decision:                entities.Accepted,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleCompanyRuleAccepted(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{},
			},
			FamilyCompanyRules: nil,
			FamilyMccRules:     nil,
			GlobalRules:        nil,
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeConsoleCompanyRules(),
	}
}

func GetRulesEvaluationResponseAcceptedCauseConsoleFamilyRules() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Accepted.String(),
		RulesModules: entities.RulesModulesResponse{
			CompanyRules: nil,
			FamilyCompanyRules: &entities.RulesResponse{
				Decision:                entities.Accepted,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleWithFamilyMccID(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{},
			},
			FamilyMccRules: nil,
			GlobalRules:    nil,
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeConsoleFamilyRules(),
	}
}

func GetRulesEvaluationResponseAcceptedCauseConsoleFamilyMccRules() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Accepted.String(),
		RulesModules: entities.RulesModulesResponse{
			CompanyRules:       nil,
			FamilyCompanyRules: nil,
			FamilyMccRules: &entities.RulesResponse{
				Decision:                entities.Accepted,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleWithFamilyMccID(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{},
			},
			GlobalRules: nil,
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeConsoleFamilyMccRules(),
	}
}

func GetRulesEvaluationResponseDeclinedCauseConsoleGlobalRules() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Declined.String(),
		RulesModules: entities.RulesModulesResponse{
			CompanyRules:       nil,
			FamilyCompanyRules: nil,
			FamilyMccRules:     nil,
			GlobalRules: &entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleEmailBlockedGlobal(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    1,
				EvaluatedNonGlobalRules: 0,
				Errors:                  []string{},
			},
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeConsoleGlobalRules(),
	}
}

func GetRulesEvaluationResponseDeclinedCauseConsoleRules() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Declined.String(),
		RulesModules: entities.RulesModulesResponse{
			CompanyRules:       nil,
			FamilyCompanyRules: nil,
			FamilyMccRules:     nil,
			GlobalRules: &entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleEmailBlockedGlobal(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    1,
				EvaluatedNonGlobalRules: 0,
				Errors:                  []string{},
			},
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeConsoleRules(),
	}
}
func GetRulesEvaluationResponseDeclinedEmailProximityRules() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Undecided.String(),
		RulesModules: entities.RulesModulesResponse{
			CompanyRules:       nil,
			FamilyCompanyRules: nil,
			FamilyMccRules:     nil,
			GlobalRules:        nil,
			IdentityModule: &entities.RulesResponse{
				Decision:                entities.Undecided,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleEmailProximity(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{},
			},
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeWithEmailProximity(),
	}
}

func GetRulesEvaluationResponseUndecidedYellowFlag() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Undecided.String(),
		RulesModules: entities.RulesModulesResponse{
			YellowFlagModule: &entities.RulesResponse{
				Decision:                entities.Undecided,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleYellowFlag(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{},
			},
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeYellowFlag(),
	}
}

func GetRulesEvaluationResponseDeclinedYellowFlag() entities.RulesEvaluationResponse {
	return entities.RulesEvaluationResponse{
		Decision: entities.Declined.String(),
		RulesModules: entities.RulesModulesResponse{
			YellowFlagModule: &entities.RulesResponse{
				Decision:                entities.Undecided,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleYellowFlag(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{},
			},
			GlobalRules: &entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleEmailBlockedGlobal(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    1,
				EvaluatedNonGlobalRules: 0,
				Errors:                  []string{},
			},
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Charge:        GetChargeYellowFlagAndGlobal(),
	}
}

func GetEvaluationResponseSuccessfulFamily() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Accepted.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Accepted,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleWithFamilyMccID(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{}},
		},
		Charge: GetDefaultChargeFamily(),
	}
}

func GetEvaluationResponseSuccessfulFamilyMcc() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Accepted.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Accepted,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleWithFamilyMccID(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{}},
		},
		Charge: GetDefaultChargeFamilyMcc(),
	}
}

func GetEvaluationResponseAcceptedByWhiteListSuccessful() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Accepted.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{GetDefaultWhiteList(false)},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{},
		},
		Charge: GetDefaultCharge(),
	}
}

func GetEvaluationResponseAcceptedByWhiteListAndBlackListPresentSuccessful() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Accepted.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{GetDefaultWhiteList(false)},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{},
		},
		Charge: GetDefaultCharge(),
	}
}

func GetEvaluationResponseAcceptedByBlackListSuccessful() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Declined.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{GetDefaultBlackList(false)},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{},
		},
		Charge: GetDefaultChargeBlacklist(),
	}
}

func GetEvaluationResponseDeclinedByRulesSuccessful() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Declined.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleFingerprintBlocked(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{}},
		},
		Charge: GetChargeWithDeviceFingerprintBlocked(),
	}
}

func GetEvaluationResponseDeclinedByGlobalRuleSuccessful() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Declined.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleEmailBlockedGlobal(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    1,
				EvaluatedNonGlobalRules: 0,
				Errors:                  []string{}},
		},
		Charge: GetChargeWithEmailBlockedGlobal(),
	}
}

func GetEvaluationResponseSuccessfulWithoutAggregation() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Undecided.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White},
			BlackList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Undecided,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{}},
		},
		Charge: GetDefaultChargeWithoutAggregation(),
	}
}

func GetEvaluationResponseUndecidedInGraylistWithoutRuleApplied() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Undecided.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{GetDefaultGrayList(true)},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{},
		},
		Charge: GetDefaultChargeInGraylist(),
	}
}

func GetEvaluationResponseDeclinedInGraylistWithRuleApplied() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Declined.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{GetDefaultGrayList(false)},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleEmailBlockedGlobalForGraylist(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    1,
				EvaluatedNonGlobalRules: 0,
				Errors:                  []string{}},
		},
		Charge: GetDefaultChargeInGraylistAndRule(),
	}
}

func GetEvaluationResponseDeclinedExistChargebackWithRuleApplied() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Declined.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleEmailWithChargebacks(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    1,
				EvaluatedNonGlobalRules: 0,
				Errors:                  []string{}},
		},
		Charge: GetDefaultChargeWithChargebacks(),
	}
}

func GetEvaluationResponseDeclinedWithOmniscoreRuleApplied() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Declined.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleWithOmniscore(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    1,
				EvaluatedNonGlobalRules: 0,
				Errors:                  []string{}},
		},
		Charge: GetChargeForOmniscoreRule(),
	}
}

func GetEvaluationResponseUndecidedWithoutOmniscoreRuleApplied() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Undecided.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{},
		},
		Charge: GetChargeForOmniscoreRuleNotApplied(),
	}
}

func GetEvaluationResponseUndecidedWithoutListTypeIncorrect() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Undecided.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			Rules: entities.RulesResponse{},
		},
		Charge: GetChargeRequestForBlacklist(),
	}
}

func GetEvaluationResponseSuccessfulMerchantScoreApproved() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Accepted.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Accepted,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleMerchantScoreApproved(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{}},
		},
		Charge: GetDefaultCharge(),
	}
}

func GetEvaluationResponseSuccessfulMerchantScoreDeclined() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Declined.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Declined,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleMerchantScoreDeclined(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{}},
		},
		Charge: GetDefaultCharge(),
	}
}

func GetEvaluationResponseSuccessfulMarketSegmentApproved() entities.EvaluationResponse {
	return entities.EvaluationResponse{
		Decision: entities.Accepted.String(),
		Modules: entities.ModulesResponse{
			EvaluationOrder: evaluationOrder,
			WhiteList: entities.ListResponse{
				Decision:      entities.Accepted,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.White,
			},
			GrayList: entities.ListResponse{
				Decision:      entities.Undecided,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Gray,
			},
			BlackList: entities.ListResponse{
				Decision:      entities.Declined,
				TestDecision:  entities.Undecided,
				DecisionRules: []entities.List{},
				TestRules:     []entities.List{},
				Errors:        []string{},
				Type:          entities.Black,
			},
			Rules: entities.RulesResponse{
				Decision:                entities.Accepted,
				TestDecision:            entities.Undecided,
				DecisionRules:           []entities.Rule{GetDefaultRuleMarketSegmentApproved(false)},
				TestRules:               []entities.Rule{},
				EvaluatedGlobalRules:    0,
				EvaluatedNonGlobalRules: 1,
				Errors:                  []string{}},
		},
		Charge: GetDefaultCharge(),
	}
}
