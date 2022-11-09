package charges

import (
	"context"
	"fmt"

	familycom "github.com/conekta/risk-rules/internal/apps/family_companies"
	merchantsscore "github.com/conekta/risk-rules/internal/apps/merchants_score"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/chargebacks"
	"github.com/conekta/risk-rules/internal/apps/families"
	"github.com/conekta/risk-rules/internal/apps/lists"
	"github.com/conekta/risk-rules/internal/apps/omniscores"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/metrics"
	"github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
)

const (
	serviceMethodName = "charge.service.%s"
)

var evaluationOrder = []string{
	"blacklist",
	"whitelist",
	"rules",
}

type ChargeService interface {
	EvaluateCharge(ctx context.Context, charge entities.ChargeRequest) (entities.EvaluationResponse, error)
	Get(ctx context.Context, id string) (entities.EvaluationResponse, error)
	EvaluateChargeOnlyRules(ctx context.Context, charge entities.ChargeRequest) (entities.RulesEvaluationResponse, error)
	GetOnlyRules(ctx context.Context, id string) (entities.RulesEvaluationResponse, error)
}

type ChargeRepository interface {
	Save(ctx context.Context, evaluation entities.EvaluationResponse) error
	SaveOnlyRules(ctx context.Context, evaluation entities.RulesEvaluationResponse) error
	Get(ctx context.Context, id string) (entities.EvaluationResponse, error)
	GetOnlyRules(ctx context.Context, id string) (entities.RulesEvaluationResponse, error)
}

type chargeService struct {
	config                  config.Config
	rulesRepository         rules.RuleRepository
	listsService            lists.ListsService
	chargesRepository       ChargeRepository
	familyService           families.FamilyService
	familyCompaniesService  familycom.FamilyCompaniesService
	rulesValidatorService   rules.RuleValidator
	payerRepository         chargebacks.ChargebackRepository
	omniscoreService        omniscores.OmniscoreService
	merchantScoreRepository merchantsscore.MerchantsScoreRepository
	logs                    logs.Logger
	metrics                 datadog.Metricer
}

func NewChargeService(cfg config.Config, ruleValidator rules.RuleValidator, ruleRepository rules.RuleRepository,
	listsService lists.ListsService, chargeRepository ChargeRepository, familyService families.FamilyService,
	familyCompaniesService familycom.FamilyCompaniesService, payerRepository chargebacks.ChargebackRepository,
	omniscoreService omniscores.OmniscoreService, merchantScoreRepository merchantsscore.MerchantsScoreRepository,
	logger logs.Logger, metric datadog.Metricer) ChargeService {
	return &chargeService{
		config:                  cfg,
		rulesRepository:         ruleRepository,
		listsService:            listsService,
		rulesValidatorService:   ruleValidator,
		chargesRepository:       chargeRepository,
		familyService:           familyService,
		familyCompaniesService:  familyCompaniesService,
		payerRepository:         payerRepository,
		omniscoreService:        omniscoreService,
		merchantScoreRepository: merchantScoreRepository,
		logs:                    logger,
		metrics:                 metric,
	}
}

func (service *chargeService) EvaluateCharge(ctx context.Context,
	charge entities.ChargeRequest) (entities.EvaluationResponse, error) {
	result := entities.NewUndecidedEvaluationResponse(charge, evaluationOrder)

	charge.Payer.Chargebacks = service.FindChargebacks(ctx, charge.Details.Email)
	charge.Omniscore = service.omniscoreService.GetScore(ctx, charge)
	charge.MerchantScore = service.getScore(ctx, charge)

	definitiveDecision, testDecision, definitiveRulesResult, listResult := service.getDecisionByConsole(ctx, charge)

	result.Decision = definitiveDecision.ValidateDecision().String()
	result.Modules.WhiteList = listResult.GetResponses(entities.White, entities.Accepted)
	result.Modules.BlackList = listResult.GetResponses(entities.Black, entities.Declined)
	result.Modules.GrayList = listResult.GetResponses(entities.Gray, entities.Undecided)
	result.Modules.Rules = definitiveRulesResult
	result.Charge.IsGraylist = charge.IsGraylist
	result.Charge.Payer.Chargebacks = charge.Payer.Chargebacks
	result.Charge.Omniscore = charge.Omniscore
	result.Charge.MerchantScore = charge.MerchantScore
	result.Charge.MarketSegment = charge.MarketSegment

	go func() {
		ctxBg := context.Background()
		service.sendChargeMetrics(ctxBg, charge, definitiveDecision.ValidateDecision().String(),
			testDecision.ValidateDecision().String())

		err := service.chargesRepository.Save(ctxBg, result)
		if err != nil {
			service.logs.Error(ctxBg, err.Error())
		}
	}()

	return result, nil
}

func (service *chargeService) EvaluateChargeOnlyRules(ctx context.Context,
	charge entities.ChargeRequest) (entities.RulesEvaluationResponse, error) {
	result := entities.NewUndecidedEvaluationResponseOnlyRules(charge)

	charge.Payer.Chargebacks = service.FindChargebacks(ctx, charge.Details.Email)
	result.Omniscore = service.omniscoreService.GetScore(ctx, charge)
	result.MerchantScore = service.getScore(ctx, charge)

	definitiveDecision, testDecision, rulesModulesResponse := service.getDecisionByConsoleOnlyRules(ctx, charge)

	result.Decision = definitiveDecision.ValidateDecision().String()
	result.RulesModules = rulesModulesResponse
	go func() {
		ctxBg := context.Background()
		service.sendChargeMetrics(ctxBg, charge, definitiveDecision.ValidateDecision().String(),
			testDecision.ValidateDecision().String())

		err := service.chargesRepository.SaveOnlyRules(ctxBg, result)
		if err != nil {
			service.logs.Error(ctxBg, err.Error())
		}
	}()

	return result, nil
}

func (service *chargeService) getDecisionByConsole(ctx context.Context, charge entities.ChargeRequest,
) (definitiveDecision entities.Decision, testDecision entities.Decision,
	definitiveRulesResult entities.RulesResponse, listResult entities.ListResponse) {
	var decisionTaken, listDecisionTaken bool
	var rulesResult entities.RulesResponse
	var decision entities.Decision

	foundLists, err := service.listsService.GetLists(ctx, charge.NewListsSearch())
	if err != nil {
		listResult.Errors = append(listResult.Errors, err.Error())
	}

	for _, component := range charge.Console {
		if component.Name.IsList() {
			listResult, listDecisionTaken = service.getDecisionByList(ctx, charge, component, foundLists)
		} else {
			rulesResult = service.getDecisionByRule(ctx, charge, component)
		}

		if listResult.Type == entities.Gray && !listResult.IsListResponseEmpty() {
			charge.IsGraylist = true
		}

		evaluations := entities.EvaluationResults{&listResult, &rulesResult}
		decision, decisionTaken = calculateDecisionByEvaluation(evaluations, component, false)
		testDecision, _ = calculateDecisionByEvaluation(evaluations, component, true)

		if (listDecisionTaken || decisionTaken) && component.Name != entities.GraylistType {
			definitiveDecision = decision
			definitiveRulesResult = rulesResult
			return definitiveDecision, testDecision, definitiveRulesResult, listResult
		}

		if decision.ValidateDecision() != entities.Undecided {
			definitiveDecision = decision.ValidateDecision()
			definitiveRulesResult = rulesResult
		}
	}
	return definitiveDecision, testDecision, definitiveRulesResult, listResult
}

func (service *chargeService) getDecisionByConsoleOnlyRules(ctx context.Context, charge entities.ChargeRequest,
) (definitiveDecision entities.Decision, testDecision entities.Decision,
	rulesModulesResponse entities.RulesModulesResponse) {
	var decisionTaken bool
	var rulesResult entities.RulesResponse
	var decision entities.Decision

	for _, component := range charge.Console {
		rulesResult = service.getDecisionByRule(ctx, charge, component)
		rulesModulesResponse.SetRuleResponse(component, rulesResult)

		evaluations := entities.EvaluationResults{&rulesResult}
		decision, decisionTaken = calculateDecisionByEvaluation(evaluations, component, false)
		testDecision, _ = calculateDecisionByEvaluation(evaluations, component, true)

		if decisionTaken {
			definitiveDecision = decision
			return definitiveDecision, testDecision, rulesModulesResponse
		}

		if decision.ValidateDecision() != entities.Undecided {
			definitiveDecision = decision.ValidateDecision()
		}
	}

	return definitiveDecision, testDecision, rulesModulesResponse
}

func (service *chargeService) getDecisionByList(ctx context.Context, charge entities.ChargeRequest,
	component entities.Component, foundLists []entities.List) (entities.ListResponse, bool) {
	var listResult entities.ListResponse
	listDecisionTaken := service.EvaluateList(ctx, charge, component.Name, &listResult, foundLists)
	if listResult.Type == entities.Gray && !listResult.IsListResponseEmpty() {
		charge.IsGraylist = true
	}

	return listResult, listDecisionTaken
}

func (service *chargeService) getDecisionByRule(ctx context.Context, charge entities.ChargeRequest,
	component entities.Component) entities.RulesResponse {
	return service.EvaluateRules(ctx, charge, component)
}

func (service *chargeService) sendChargeMetrics(ctx context.Context, charge entities.ChargeRequest,
	decision, testDecision string) {
	metricData := metrics.NewMetricData(ctx, "EvaluateCharge", serviceMethodName, service.config.Env)
	metricData.AddCustomTags([]string{
		fmt.Sprintf(text.MetricTagTestRulesChangeDecision, decision == testDecision),
		fmt.Sprintf(text.MetricTagRulesDecision, decision),
		fmt.Sprintf(text.MetricTagTestRulesDecision, testDecision),
		fmt.Sprintf(text.MetricTagCompany, charge.CompanyID),
		fmt.Sprintf(text.MetricPaymentNetwork, charge.PaymentMethod.Brand),
		fmt.Sprintf(text.MetricCardType, charge.PaymentMethod.CardType),
		fmt.Sprintf(text.MetricCountry, charge.PaymentMethod.Country),
		fmt.Sprintf(text.MetricIssuer, charge.PaymentMethod.Issuer),
	})
	metricData.SetResult(true)
	metrics.SendAsyncMetrics(service.metrics, service.logs, metricData, text.EvaluateChargeMetricName)
}

func (service *chargeService) EvaluateRules(ctx context.Context, charge entities.ChargeRequest,
	component entities.Component) entities.RulesResponse {
	response := entities.NewRulesResponse()
	var familyID string
	var familyCompaniesIDs []string
	var totalApplied int64

	mapCharge, err := charge.ToMap()
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		return response
	}

	decisionRules := make([]entities.Rule, 0)
	testRules := make([]entities.Rule, 0)

	if component.Name == entities.FamilyCompanyRulesType {
		familyID = service.getFamilyIDFromCharge(ctx, charge)
	} else if component.Name == entities.FamilyMccRulesType {
		familyCompaniesIDs = service.getFamilyCompaniesIDsFromCharge(ctx, charge)
	}

	rulesFound, _ := service.rulesRepository.GetRulesByFilters(ctx,
		entities.RuleFilter{CompanyID: charge.CompanyID, FamilyID: familyID, FamilyCompaniesIDs: familyCompaniesIDs}, component.Name)

	totalApplied = 0
	for _, rule := range rulesFound {
		isApplied, err := service.rulesValidatorService.Evaluate(ctx, rule, mapCharge)
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
			continue
		}

		if rule.IsGlobal {
			response.EvaluatedGlobalRules++
		} else {
			response.EvaluatedNonGlobalRules++
		}

		if isApplied {
			totalApplied++
			if rule.IsTest {
				testRules = append(testRules, rule)
				if component.Priority[0] == rule.Decision {
					response.TestDecision = rule.Decision
				}
			} else {
				decisionRules = append(decisionRules, rule)
				if component.Priority[0] == rule.Decision {
					response.Decision = rule.Decision
				}
			}
		}
	}

	if component.HaveSecondaryDecision() {
		if shouldAssignSecondaryDecision(totalApplied, rulesFound, component, response) {
			response.Decision = component.Priority[1]
		}
	}

	response.DecisionRules = decisionRules
	response.TestRules = testRules

	return response
}

func shouldAssignSecondaryDecision(totalApplied int64, companyRules []entities.Rule,
	component entities.Component, response entities.RulesResponse) bool {
	return totalApplied > 0 && len(companyRules) > 0 && component.Priority[0] !=
		response.Decision || response.Decision == entities.Undecided
}

func (service *chargeService) getFamilyIDFromCharge(ctx context.Context, charge entities.ChargeRequest) string {
	family, err := service.familyService.GetFamily(ctx,
		entities.FamilyFilter{
			Mccs:                 []string{charge.CompanyMCC},
			NotExcludedCompanies: []string{charge.CompanyID},
		},
	)

	if err != nil {
		service.familyError(ctx, "getFamilyIDFromCharge", err)
		return strings.Empty
	}

	return family.ID.Hex()
}

func (service *chargeService) getFamilyCompaniesIDsFromCharge(ctx context.Context, charge entities.ChargeRequest) []string {
	familyCompaniesIDs := make([]string, 0)
	familiesCompanies, err := service.familyCompaniesService.GetFamiliesCompaniesFromFilter(ctx,
		entities.FamilyCompaniesFilter{
			CompanyIDs: []string{charge.CompanyID},
		},
	)

	if err != nil {
		service.familyError(ctx, "getFamilyCompaniesIDsFromCharge", err)
	}

	for _, familyCompanies := range familiesCompanies {
		familyCompaniesIDs = append(familyCompaniesIDs, familyCompanies.ID.Hex())
	}

	return familyCompaniesIDs
}

func (service *chargeService) familyError(ctx context.Context, methodName string, err error) {
	service.logs.Error(
		ctx, err.Error(),
		text.LogTagMethod, fmt.Sprintf(serviceMethodName, methodName),
	)
}

func (service *chargeService) EvaluateList(ctx context.Context, charge entities.ChargeRequest,
	component entities.ConsoleComponent, listResponse *entities.ListResponse, foundLists []entities.List) (decisionTaken bool) {
	listResponse.Type = entities.TypeList(component)

	for _, list := range foundLists {
		if !list.IsValidListType() {
			service.logs.Error(ctx,
				fmt.Sprintf(`chargeService:EvaluateList - No valid list type detected. 
				-component: %+v -listType: %s -company_id: %s -charge_id: %s `,
					component, list.Type, charge.CompanyID, charge.ID))
			continue
		}

		if component == entities.GraylistType && list.IsGraylist() {
			if list.IsTest {
				listResponse.TestDecision = entities.Undecided
				listResponse.TestRules = append(listResponse.TestRules, list)
			} else {
				listResponse.Decision = entities.Undecided
				listResponse.DecisionRules = append(listResponse.DecisionRules, list)
			}
		} else if component == entities.WhitelistType && list.IsWhitelist() {
			if list.IsTest {
				listResponse.TestDecision = entities.Accepted
				listResponse.TestRules = append(listResponse.TestRules, list)
			} else {
				listResponse.Decision = entities.Accepted
				listResponse.DecisionRules = append(listResponse.DecisionRules, list)
			}
			return true
		} else if component == entities.BlacklistType && list.IsBlacklist() {
			if list.IsTest {
				listResponse.TestDecision = entities.Declined
				listResponse.TestRules = append(listResponse.TestRules, list)
			} else {
				listResponse.Decision = entities.Declined
				listResponse.DecisionRules = append(listResponse.DecisionRules, list)
			}
			return true
		}
	}

	return false
}

func (service *chargeService) FindChargebacks(ctx context.Context, email string) int64 {
	foundPayer, err := service.payerRepository.Find(ctx, entities.Payer{Email: email})
	if err != nil {
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, "HasChargebacks", text.Email, email)
		return 0
	}

	return int64(len(foundPayer.Chargebacks))
}

func (service *chargeService) FindMerchantScore(ctx context.Context, companyID string) float64 {
	foundMerchantScore, err := service.merchantScoreRepository.FindByMerchantID(ctx, companyID)
	if err != nil {
		service.logs.Error(ctx, err.Error(), text.LogTagMethod, "HasMerchantScore", text.CompanyID, companyID)
		return 0
	}

	return foundMerchantScore.Score
}

func calculateDecisionByEvaluation(results entities.EvaluationResults,
	priorities entities.Component, isTest bool) (decision entities.Decision, decisionTaken bool) {
	decision = entities.Undecided

	for _, result := range results {
		if isTest {
			if result.GetTestDecision() != entities.Undecided {
				decision = result.GetTestDecision()
			}
		} else {
			if result.GetDecision() != entities.Undecided {
				decision = result.GetDecision()
			}
		}

		if shouldTakeDecision(priorities, decision, isTest) {
			return decision, true
		}
	}

	return decision, false
}

func shouldTakeDecision(priorities entities.Component,
	decision entities.Decision,
	isTest bool) bool {
	return priorities.Priority[0] == decision &&
		!strings.IsEmpty(decision.String()) &&
		!isTest && priorities.HaveSecondaryDecision()
}

func (service *chargeService) getScore(ctx context.Context, charge entities.ChargeRequest) float64 {
	if service.config.MerchantScore.IsEnabled {
		return service.FindMerchantScore(ctx, charge.CompanyID)
	}

	return -1
}

func (service *chargeService) Get(ctx context.Context, id string) (entities.EvaluationResponse, error) {
	return service.chargesRepository.Get(ctx, id)
}

func (service *chargeService) GetOnlyRules(ctx context.Context, id string) (entities.RulesEvaluationResponse, error) {
	return service.chargesRepository.GetOnlyRules(ctx, id)
}
