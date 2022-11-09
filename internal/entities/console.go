package entities

import "github.com/conekta/risk-rules/pkg/strings"

type ConsoleComponent string

const (
	WhitelistType          ConsoleComponent = "Whitelist"
	BlacklistType          ConsoleComponent = "Blacklist"
	GraylistType           ConsoleComponent = "Graylist"
	CompanyRulesType       ConsoleComponent = "CompanyRules"
	FamilyCompanyRulesType ConsoleComponent = "FamilyCompanyRules"
	FamilyMccRulesType     ConsoleComponent = "FamilyMccRules"
	GlobalRulesType        ConsoleComponent = "GlobalRules"
	ScoreRulesType         ConsoleComponent = "ScoreRules"
	IdentityModuleType     ConsoleComponent = "IdentityModule"
	YellowFlagType         ConsoleComponent = "YellowFlag"
)

var ComponentsWithOutSecondaryDecision = []ConsoleComponent{
	IdentityModuleType,
	YellowFlagType,
}

func (c ConsoleComponent) IsList() bool {
	return strings.Contains(string(c), "list")
}

func (c ConsoleComponent) IsIdentityModule() bool {
	return strings.Contains(string(c), "identity")
}
