package charges

import (
	"context"
	"errors"
	"testing"

	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/apps/chargebacks"
	"github.com/conekta/risk-rules/internal/apps/families"
	familycom "github.com/conekta/risk-rules/internal/apps/family_companies"
	"github.com/conekta/risk-rules/internal/apps/lists"
	merchantsscore "github.com/conekta/risk-rules/internal/apps/merchants_score"
	"github.com/conekta/risk-rules/internal/apps/omniscores"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/rest"
	"github.com/conekta/risk-rules/test/mocks"
	"github.com/conekta/risk-rules/test/mocks/datadog"
	"github.com/conekta/risk-rules/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	familyID           = entities.Family{}.ID.Hex()
	familyCompaniesId  = entities.FamilyCompanies{}.ID.Hex()
	familyCompaniesIds = []string{familyCompaniesId}
	defaultScore       = float64(-1)
)

type fields struct {
	config                  config.Config
	rulesRepository         rules.RuleRepository
	listsService            lists.ListsService
	rulesValidatorService   rules.RuleValidator
	chargeRepository        ChargeRepository
	familyService           families.FamilyService
	familyCompaniesService  familycom.FamilyCompaniesService
	chargebackRepository    chargebacks.ChargebackRepository
	omniscoreService        omniscores.OmniscoreService
	merchantScoreRepository merchantsscore.MerchantsScoreRepository
}

type args struct {
	charge entities.ChargeRequest
}

type testItem struct {
	name    string
	fields  fields
	args    args
	want    entities.EvaluationResponse
	wantErr bool
}

type testItemOnlyRules struct {
	name    string
	fields  fields
	args    args
	want    entities.RulesEvaluationResponse
	wantErr bool
}

func TestChargeService_EvaluateChargeOnlyRulesMerchantScoreEnabled(t *testing.T) {
	log, _ := logs.New()
	cfg := config.NewConfig()
	cfg.MerchantScore.IsEnabled = true
	rulesRepositoryMockFirstCase, listServiceMockFirstCase, familyServiceMockFirstCase, familyCompaniesServiceMockFirstCase, chargebackRepositoryMockFirstCase, merchantScoreRepositoryMockFirstCase := getMockServiceFirstCaseOnlyRules()

	_, OmniscoreIsOff := getOmniscoreTestCases()

	chargeRepositoryOkMock := new(mocks.ChargeEvaluationRepositoryMock)
	chargeRepositoryOkMock.On("SaveOnlyRules", mock.Anything, mock.AnythingOfType("entities.RulesEvaluationResponse")).
		Return(nil)

	chargeRepositoryErrMock := new(mocks.ChargeEvaluationRepositoryMock)
	chargeRepositoryErrMock.On("SaveOnlyRules", mock.Anything, mock.AnythingOfType("entities.RulesEvaluationResponse")).
		Return(errors.New("connection to database lost"))

	testsCases := []testItemOnlyRules{
		{
			"the charge evaluation is undecided cause console is empty and merchant score is enabled",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFirstCase,
				listsService:            listServiceMockFirstCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyServiceMockFirstCase,
				familyCompaniesService:  familyCompaniesServiceMockFirstCase,
				chargebackRepository:    chargebackRepositoryMockFirstCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockFirstCase,
			},
			args{charge: testdata.GetChargeConsoleIsEmptyOnlyRules()},
			testdata.GetRulesEvaluationResponseUndecidedCauseConsoleIsEmptyOnlyRules(),
			false,
		},
	}

	for _, ttCase := range testsCases {
		t.Run(ttCase.name, func(t *testing.T) {
			r := NewChargeService(ttCase.fields.config, ttCase.fields.rulesValidatorService, ttCase.fields.rulesRepository,
				ttCase.fields.listsService, ttCase.fields.chargeRepository, ttCase.fields.familyService,
				ttCase.fields.familyCompaniesService, ttCase.fields.chargebackRepository, ttCase.fields.omniscoreService,
				ttCase.fields.merchantScoreRepository, log, new(datadog.MetricsDogMock),
			)
			got, err := r.EvaluateChargeOnlyRules(context.Background(), ttCase.args.charge)
			if (err != nil) != ttCase.wantErr {
				t.Errorf("EvaluateCharge() error = %v, wantErr %v , name %s", err, ttCase.wantErr, ttCase.name)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.EqualValues(t, ttCase.want, got)
		})
	}
}

func TestChargeService_EvaluateChargeOnlyRules(t *testing.T) {
	log, _ := logs.New()
	cfg := config.NewConfig()

	rulesRepositoryMockFirstCase, listServiceMockFirstCase, familyServiceMockFirstCase, familyCompaniesServiceMockFirstCase, chargebackRepositoryMockFirstCase, merchantScoreRepositoryMockFirstCase := getMockServiceFirstCaseOnlyRules()
	rulesRepositoryMockSecondCase, listServiceMockSecondCase, familyRepositoryMockSecondCase, familyCompaniesServiceMockSecondCase, chargebackRepositoryMockSecondCase, merchantScoreRepositoryMockSecondCase := getMockServiceSecondCaseOnlyRules()
	rulesRepositoryMockThirdCase, listServiceMockThirdCase, familyRepositoryMockThirdCase, familyCompaniesServiceMockThirdCase, chargebackRepositoryMockThirdCase, merchantScoreRepositoryMockThirdCase := getMockServiceThirdCaseOnlyRules()
	rulesRepositoryMockFourthCase, listServiceMockFourthCase, familyRepositoryMockFourthCase, familyCompaniesServiceMockFourthCase, chargebackRepositoryMockFourthCase, merchantScoreRepositoryMockFourthCase := getMockServiceFourthCaseOnlyRules()
	cfg.MerchantScore.IsEnabled = true
	rulesRepositoryMockFifthCase, listServiceMockFifthCase, familyRepositoryMockFifthCase, familyCompaniesServiceMockFifthCase, chargebackRepositoryMockFifthCase, merchantScoreRepositoryMockFifthCase := getMockServiceFifthCaseOnlyRules()

	rulesRepositoryMockSixthCase, listServiceMockSixthCase, familyRepositoryMockSixthCase, familyCompaniesServiceMockSixthCase, chargebackRepositoryMockSixthCase, merchantScoreRepositoryMockSixthCase := getMockServiceSixthCaseOnlyRules()
	cfg.MerchantScore.IsEnabled = false

	rulesRepositoryMockSevenCase, listServiceMockSevenCase, familyRepositoryMockSevenCase, familyCompaniesServiceMockSevenCase, chargebackRepositoryMockSevenCase, merchantScoreRepositoryMockSevenCase := getMockServiceSevenCaseOnlyRules()
	rulesRepositoryMockEighthCase, listServiceMockEighthCase, familyRepositoryMockEighthCase, familyCompaniesServiceMockEighthCase, chargebackRepositoryMockEighthCase, merchantScoreRepositoryMockEighthCase := getMockServiceEighthCaseOnlyRules()
	rulesRepositoryMockNinthCase, listServiceMockNinthCase, familyRepositoryMockNinthCase, familyCompaniesServiceMockNinthCase, chargebackRepositoryMockNinthCase, merchantScoreRepositoryMockNinthCase := getMockServiceNinthCaseOnlyRules()

	_, OmniscoreIsOff := getOmniscoreTestCases()

	chargeRepositoryOkMock := new(mocks.ChargeEvaluationRepositoryMock)
	chargeRepositoryOkMock.On("SaveOnlyRules", mock.Anything, mock.AnythingOfType("entities.RulesEvaluationResponse")).
		Return(nil)

	chargeRepositoryErrMock := new(mocks.ChargeEvaluationRepositoryMock)
	chargeRepositoryErrMock.On("SaveOnlyRules", mock.Anything, mock.AnythingOfType("entities.RulesEvaluationResponse")).
		Return(errors.New("connection to database lost"))

	testsCases := []testItemOnlyRules{
		{
			"the charge evaluation is undecided cause console is empty",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFirstCase,
				listsService:            listServiceMockFirstCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyServiceMockFirstCase,
				familyCompaniesService:  familyCompaniesServiceMockFirstCase,
				chargebackRepository:    chargebackRepositoryMockFirstCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockFirstCase,
			},
			args{charge: testdata.GetChargeConsoleIsEmptyOnlyRules()},
			testdata.GetRulesEvaluationResponseUndecidedCauseConsoleIsEmptyOnlyRules(),
			false,
		},
		{
			"the charge should be approved because it the company amount is in CompanyRules",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockSecondCase,
				listsService:            listServiceMockSecondCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockSecondCase,
				familyCompaniesService:  familyCompaniesServiceMockSecondCase,
				chargebackRepository:    chargebackRepositoryMockSecondCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockSecondCase,
			},
			args{charge: testdata.GetChargeConsoleCompanyRules()},
			testdata.GetRulesEvaluationResponseAcceptedCauseConsoleCompanyRules(),
			false,
		},
		{
			"the charge should be approved, family company found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockThirdCase,
				listsService:            listServiceMockThirdCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockThirdCase,
				familyCompaniesService:  familyCompaniesServiceMockThirdCase,
				chargebackRepository:    chargebackRepositoryMockThirdCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockThirdCase,
			},
			args{charge: testdata.GetChargeConsoleFamilyRules()},
			testdata.GetRulesEvaluationResponseAcceptedCauseConsoleFamilyRules(),
			false,
		},
		{
			"the charge should be approved, family mcc found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFourthCase,
				listsService:            listServiceMockFourthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockFourthCase,
				familyCompaniesService:  familyCompaniesServiceMockFourthCase,
				chargebackRepository:    chargebackRepositoryMockFourthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockFourthCase,
			},
			args{charge: testdata.GetChargeConsoleFamilyMccRules()},
			testdata.GetRulesEvaluationResponseAcceptedCauseConsoleFamilyMccRules(),
			false,
		},
		{
			"the charge should be Declined because the email is blocked globally",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFifthCase,
				listsService:            listServiceMockFifthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockFifthCase,
				familyCompaniesService:  familyCompaniesServiceMockFifthCase,
				chargebackRepository:    chargebackRepositoryMockFifthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockFifthCase,
			},
			args{charge: testdata.GetChargeConsoleGlobalRules()},
			testdata.GetRulesEvaluationResponseDeclinedCauseConsoleGlobalRules(),
			false,
		},
		{
			"the charge should be Declined because the email is blocked globally, with all components rules",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockSixthCase,
				listsService:            listServiceMockSixthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockSixthCase,
				familyCompaniesService:  familyCompaniesServiceMockSixthCase,
				chargebackRepository:    chargebackRepositoryMockSixthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockSixthCase,
			},
			args{charge: testdata.GetChargeConsoleRules()},
			testdata.GetRulesEvaluationResponseDeclinedCauseConsoleRules(),
			false,
		},
		{
			"the charge should be Declined because the email has email_proximity.stats.observations_count > 7",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockSevenCase,
				listsService:            listServiceMockSevenCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockSevenCase,
				familyCompaniesService:  familyCompaniesServiceMockSevenCase,
				chargebackRepository:    chargebackRepositoryMockSevenCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockSevenCase,
			},
			args{charge: testdata.GetChargeWithEmailProximity()},
			testdata.GetRulesEvaluationResponseDeclinedEmailProximityRules(),
			false,
		},
		{
			"the charge should be Undecided because the console has yellow flag module",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockEighthCase,
				listsService:            listServiceMockEighthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockEighthCase,
				familyCompaniesService:  familyCompaniesServiceMockEighthCase,
				chargebackRepository:    chargebackRepositoryMockEighthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockEighthCase,
			},
			args{charge: testdata.GetChargeYellowFlag()},
			testdata.GetRulesEvaluationResponseUndecidedYellowFlag(),
			false,
		},
		{
			"the charge should be Declined because found company rule and has yellow flag",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockNinthCase,
				listsService:            listServiceMockNinthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockNinthCase,
				familyCompaniesService:  familyCompaniesServiceMockNinthCase,
				chargebackRepository:    chargebackRepositoryMockNinthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockNinthCase,
			},
			args{charge: testdata.GetChargeYellowFlagAndGlobal()},
			testdata.GetRulesEvaluationResponseDeclinedYellowFlag(),
			false,
		},
	}

	for _, ttCase := range testsCases {
		t.Run(ttCase.name, func(t *testing.T) {
			r := NewChargeService(ttCase.fields.config, ttCase.fields.rulesValidatorService, ttCase.fields.rulesRepository,
				ttCase.fields.listsService, ttCase.fields.chargeRepository, ttCase.fields.familyService,
				ttCase.fields.familyCompaniesService, ttCase.fields.chargebackRepository, ttCase.fields.omniscoreService,
				ttCase.fields.merchantScoreRepository, log, new(datadog.MetricsDogMock),
			)
			got, err := r.EvaluateChargeOnlyRules(context.Background(), ttCase.args.charge)
			if (err != nil) != ttCase.wantErr {
				t.Errorf("EvaluateCharge() error = %v, wantErr %v , name %s", err, ttCase.wantErr, ttCase.name)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.EqualValues(t, ttCase.want, got)
		})
	}
}

func TestChargeService_EvaluateCharge(t *testing.T) {
	log, _ := logs.New()
	cfg := config.NewConfig()

	rulesRepositoryMockFirstCase, listServiceMockFirstCase, familyServiceMockFirstCase, familyCompaniesServiceMockFirstCase, chargebackRepositoryMockFirstCase, merchantScoreRepositoryMockFirstCase := getMockServiceFirstCase()
	rulesRepositoryMockSecondCase, listServiceMockSecondCase, familyRepositoryMockSecondCase, familyCompaniesServiceMockSecondCase, chargebackRepositoryMockSecondCase, merchantScoreRepositoryMockSecondCase := getMockServiceSecondCase()
	rulesRepositoryMockThirdCase, listServiceMockThirdCase, familyRepositoryMockThirdCase, familyCompaniesServiceMockThirdCase, chargebackRepositoryMockThirdCase, merchantScoreRepositoryMockThirdCase := getMockServiceThirdCase()
	rulesRepositoryMockFourthCase, listServiceMockFourthCase, familyRepositoryMockFourthCase, familyCompaniesServiceMockFourthCase, chargebackRepositoryMockFourthCase, merchantScoreRepositoryMockFourthCase := getMockServiceFourthCase()
	rulesRepositoryMockFifthCase, listServiceMockFifthCase, familyRepositoryMockFifthCase, familyCompaniesServiceMockFifthCase, chargebackRepositoryMockFifthCase, merchantScoreRepositoryMockFifthCase := getMockServiceFifthCase()
	rulesRepositoryMockSixthCase, listServiceMockSixthCase, familyRepositoryMockSixthCase, familyCompaniesServiceMockSixthCase, chargebackRepositoryMockSixthCase, merchantScoreRepositoryMockSixthCase := getMockServiceSixthCase()
	rulesRepositoryMockNinthCase, listServiceMockNinthCase, familyRepositoryMockNinthCase, familyCompaniesServiceMockNinthCase, chargebackRepositoryMockNinthCase, merchantScoreRepositoryMockNinthCase := getMockServiceNinthCase()
	rulesRepositoryMockEleventhCase, listServiceMockEleventhCase, familyRepositoryMockEleventhCase, familyCompaniesServiceMockEleventhCase, chargebackRepositoryMockEleventhCase, merchantScoreRepositoryMockEleventhCase := getMockServiceEleventhCase()
	rulesRepositoryMockThirteenthCase, listServiceMockThirteenthCase, familyRepositoryMockThirteenthCase, familyCompaniesServiceMockThirteenthCase, chargebackRepositoryMockThirteenthCase, merchantScoreRepositoryMockThirteenthCase := getMockServiceThirteenthCase()
	rulesRepositoryMockFourteenthCase, listServiceMockFourteenthCase, familyRepositoryMockFourteenthCase, familyCompaniesServiceMockFourteenthCase, chargebackRepositoryMockFourteenthCase, merchantScoreRepositoryFourteenthCase := getMockServiceFourteenthCase()
	rulesRepositoryMockFifteenthCase, listServiceMockFifteenthCase, familyRepositoryMockFifteenthCase, familyCompaniesServiceMockFifteenthCase, chargebackRepositoryMockFifteenthCase, merchantScoreRepositoryMockFifteenthCase := getMockServiceFifteenthCase()
	rulesRepositoryMockSixteenthCase, listServiceMockSixteenthCase, familyRepositoryMockSixteenthCase, familyCompaniesServiceMockSixteenthCase, chargebackRepositoryMockSixteenthCase, merchantScoreRepositorySixteenthCase := getMockServiceSixteenthCase()
	rulesRepositoryMockSeventeenthCase, listServiceMockSeventeenthCase, familyRepositoryMockSeventeenthCase, familyCompaniesServiceMockSeventeenthCase, chargebackRepositoryMockSeventeenthCase, merchantScoreRepositoryMockSeventeenthCase := getMockServiceSeventeenthCase()
	rulesRepositoryMockEighteenthCase, listServiceMockEighteenthCase, familyRepositoryMockEighteenthCase, familyCompaniesServiceMockEighteenthCase, chargebackRepositoryMockEighteenthCase, merchantScoreRepositoryMockEighteenthCase := getMockServiceSeventeenthCase()
	rulesRepositoryMockNineteenthCase, listServiceMockNineteenthCase, familyRepositoryMockNineteenthCase, familyCompaniesServiceMockNineteenthCase, chargebackRepositoryMockNineteenthCase, merchantScoreRepositoryMockNineteenthCase := getMockServiceNineteenthCase()
	rulesRepositoryMockTwentiethCase, listServiceMockTwentiethCase, familyRepositoryMockTwentiethCase, familyCompaniesServiceMockTwentiethCase, chargebackRepositoryMockTwentiethCase, merchantScoreRepositoryMockTwentiethCase := getMockServiceTwentiethCase()
	rulesRepositoryMockTwentyfirstCase, listServiceMockTwentyfirstCase, familyRepositoryMockTwentyfirstCase, familyCompaniesServiceMockTwentyfirstCase, chargebackRepositoryMockTwentyfirstCase, merchantScoreRepositoryMockTwentyfirstCase := getMockServiceTwentyfirstCase()
	rulesRepositoryMockTwentySecondCase, listServiceMockTwentySecondCase, familyRepositoryMockTwentySecondCase, familyCompaniesServiceMockTwentySecondCase, chargebackRepositoryMockTwentySecondCase, merchantScoreRepositoryMockTwentySecondCase := getMockServiceTwentySecondCase()

	OmniscoreIsOn, OmniscoreIsOff := getOmniscoreTestCases()

	chargeRepositoryOkMock := new(mocks.ChargeEvaluationRepositoryMock)
	chargeRepositoryOkMock.On("Save", mock.Anything, mock.AnythingOfType("entities.EvaluationResponse")).
		Return(nil)

	chargeRepositoryErrMock := new(mocks.ChargeEvaluationRepositoryMock)
	chargeRepositoryErrMock.On("Save", mock.Anything, mock.AnythingOfType("entities.EvaluationResponse")).
		Return(errors.New("connection to database lost"))

	tests := []testItem{
		{
			"the charge evaluation is undecided cause no match",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFirstCase,
				listsService:            listServiceMockFirstCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyServiceMockFirstCase,
				familyCompaniesService:  familyCompaniesServiceMockFirstCase,
				chargebackRepository:    chargebackRepositoryMockFirstCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockFirstCase,
			},
			args{charge: testdata.GetDefaultCharge()},
			testdata.GetEvaluationResponseSuccessful(),
			false,
		},
		{
			"the charge should be approved because it the company id is in the whitelist",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockSecondCase,
				listsService:            listServiceMockSecondCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockSecondCase,
				familyCompaniesService:  familyCompaniesServiceMockSecondCase,
				chargebackRepository:    chargebackRepositoryMockSecondCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockSecondCase,
			},
			args{charge: testdata.GetDefaultCharge()},
			testdata.GetEvaluationResponseAcceptedByWhiteListSuccessful(),
			false,
		},
		{
			"the charge should be approved because it the company id is in the whitelist even though it's also on the blacklist",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockThirdCase,
				listsService:            listServiceMockThirdCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockThirdCase,
				familyCompaniesService:  familyCompaniesServiceMockThirdCase,
				chargebackRepository:    chargebackRepositoryMockThirdCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockThirdCase,
			},
			args{charge: testdata.GetDefaultCharge()},
			testdata.GetEvaluationResponseAcceptedByWhiteListAndBlackListPresentSuccessful(),
			false,
		},
		{
			"the charge should be Declined because it the company id is in the blacklist",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFourthCase,
				listsService:            listServiceMockFourthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockFourthCase,
				familyCompaniesService:  familyCompaniesServiceMockFourthCase,
				chargebackRepository:    chargebackRepositoryMockFourthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockFourthCase,
			},
			args{charge: testdata.GetDefaultChargeBlacklist()},
			testdata.GetEvaluationResponseAcceptedByBlackListSuccessful(),
			false,
		},
		{
			"the charge should be Declined because it has a rule configured",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFifthCase,
				listsService:            listServiceMockFifthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockFifthCase,
				familyCompaniesService:  familyCompaniesServiceMockFifthCase,
				chargebackRepository:    chargebackRepositoryMockFifthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockFifthCase,
			},
			args{charge: testdata.GetChargeWithDeviceFingerprintBlocked()},
			testdata.GetEvaluationResponseDeclinedByRulesSuccessful(),
			false,
		},
		{
			"the charge should be Declined because the email is blocked globally",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockSixthCase,
				listsService:            listServiceMockSixthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockSixthCase,
				familyCompaniesService:  familyCompaniesServiceMockSixthCase,
				chargebackRepository:    chargebackRepositoryMockSixthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockSixthCase,
			},
			args{charge: testdata.GetChargeWithEmailBlockedGlobal()},
			testdata.GetEvaluationResponseDeclinedByGlobalRuleSuccessful(),
			false,
		},
		{
			"the charge should be approved, family found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockNinthCase,
				listsService:            listServiceMockNinthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockNinthCase,
				familyCompaniesService:  familyCompaniesServiceMockNinthCase,
				chargebackRepository:    chargebackRepositoryMockNinthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockNinthCase,
			},
			args{charge: testdata.GetDefaultChargeFamily()},
			testdata.GetEvaluationResponseSuccessfulFamily(),
			false,
		},
		{
			"the charge should be approved, family companies found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockEleventhCase,
				listsService:            listServiceMockEleventhCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockEleventhCase,
				familyCompaniesService:  familyCompaniesServiceMockEleventhCase,
				chargebackRepository:    chargebackRepositoryMockEleventhCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockEleventhCase,
			},
			args{charge: testdata.GetDefaultChargeFamilyMcc()},
			testdata.GetEvaluationResponseSuccessfulFamilyMcc(),
			false,
		},
		{
			"the charge should be approved, family companies found in family companies array and rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockThirteenthCase,
				listsService:            listServiceMockThirteenthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockThirteenthCase,
				familyCompaniesService:  familyCompaniesServiceMockThirteenthCase,
				chargebackRepository:    chargebackRepositoryMockThirteenthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockThirteenthCase,
			},
			args{charge: testdata.GetDefaultChargeFamilyMcc()},
			testdata.GetEvaluationResponseSuccessfulFamilyMcc(),
			false,
		},
		{
			"the charge should be undecided, exist graylist but not found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFourteenthCase,
				listsService:            listServiceMockFourteenthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockFourteenthCase,
				familyCompaniesService:  familyCompaniesServiceMockFourteenthCase,
				chargebackRepository:    chargebackRepositoryMockFourteenthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryFourteenthCase,
			},
			args{charge: testdata.GetDefaultChargeInGraylist()},
			testdata.GetEvaluationResponseUndecidedInGraylistWithoutRuleApplied(),
			false,
		},
		{
			"the charge should be Declined, exist graylist and found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockFifteenthCase,
				listsService:            listServiceMockFifteenthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockFifteenthCase,
				familyCompaniesService:  familyCompaniesServiceMockFifteenthCase,
				chargebackRepository:    chargebackRepositoryMockFifteenthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockFifteenthCase,
			},
			args{charge: testdata.GetDefaultChargeInGraylistAndRule()},
			testdata.GetEvaluationResponseDeclinedInGraylistWithRuleApplied(),
			false,
		},
		{
			"the charge should be Declined, exist chargeback and found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockSixteenthCase,
				listsService:            listServiceMockSixteenthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockSixteenthCase,
				familyCompaniesService:  familyCompaniesServiceMockSixteenthCase,
				chargebackRepository:    chargebackRepositoryMockSixteenthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositorySixteenthCase,
			},
			args{charge: testdata.GetDefaultChargeWithChargebacks()},
			testdata.GetEvaluationResponseDeclinedExistChargebackWithRuleApplied(),
			false,
		},
		{
			"the charge should be Declined, has omniscore and found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockSeventeenthCase,
				listsService:            listServiceMockSeventeenthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockSeventeenthCase,
				familyCompaniesService:  familyCompaniesServiceMockSeventeenthCase,
				chargebackRepository:    chargebackRepositoryMockSeventeenthCase,
				omniscoreService:        OmniscoreIsOn,
				merchantScoreRepository: merchantScoreRepositoryMockSeventeenthCase,
			},
			args{charge: testdata.GetChargeForOmniscoreRule()},
			testdata.GetEvaluationResponseDeclinedWithOmniscoreRuleApplied(),
			false,
		},
		{
			"the charge should be Undecided, because omniscore service is off/failed",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockEighteenthCase,
				listsService:            listServiceMockEighteenthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockEighteenthCase,
				familyCompaniesService:  familyCompaniesServiceMockEighteenthCase,
				chargebackRepository:    chargebackRepositoryMockEighteenthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockEighteenthCase,
			},
			args{charge: testdata.GetChargeRequestForOmniscoreRule()},
			testdata.GetEvaluationResponseUndecidedWithoutOmniscoreRuleApplied(),
			false,
		},
		{
			"the charge should be Undecided, because list type is incorrect in list",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockNineteenthCase,
				listsService:            listServiceMockNineteenthCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockNineteenthCase,
				familyCompaniesService:  familyCompaniesServiceMockNineteenthCase,
				chargebackRepository:    chargebackRepositoryMockNineteenthCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockNineteenthCase,
			},
			args{charge: testdata.GetChargeRequestForBlacklist()},
			testdata.GetEvaluationResponseUndecidedWithoutListTypeIncorrect(),
			false,
		},
		{
			"the charge should be approved, has merchant score and found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockTwentiethCase,
				listsService:            listServiceMockTwentiethCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockTwentiethCase,
				familyCompaniesService:  familyCompaniesServiceMockTwentiethCase,
				chargebackRepository:    chargebackRepositoryMockTwentiethCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockTwentiethCase,
			},
			args{charge: testdata.GetDefaultCharge()},
			testdata.GetEvaluationResponseSuccessfulMerchantScoreApproved(),
			false,
		},
		{
			"the charge should be approved, family companies found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockTwentyfirstCase,
				listsService:            listServiceMockTwentyfirstCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockTwentyfirstCase,
				familyCompaniesService:  familyCompaniesServiceMockTwentyfirstCase,
				chargebackRepository:    chargebackRepositoryMockTwentyfirstCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockTwentyfirstCase,
			},
			args{charge: testdata.GetDefaultCharge()},
			testdata.GetEvaluationResponseSuccessfulMerchantScoreDeclined(),
			false,
		},
		{
			"the charge should be approved, has merchant score and found rules applied",
			fields{
				config:                  cfg,
				rulesRepository:         rulesRepositoryMockTwentySecondCase,
				listsService:            listServiceMockTwentySecondCase,
				rulesValidatorService:   rules.NewRulesValidator(log),
				chargeRepository:        chargeRepositoryOkMock,
				familyService:           familyRepositoryMockTwentySecondCase,
				familyCompaniesService:  familyCompaniesServiceMockTwentySecondCase,
				chargebackRepository:    chargebackRepositoryMockTwentySecondCase,
				omniscoreService:        OmniscoreIsOff,
				merchantScoreRepository: merchantScoreRepositoryMockTwentySecondCase,
			},
			args{charge: testdata.GetDefaultCharge()},
			testdata.GetEvaluationResponseSuccessfulMarketSegmentApproved(),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewChargeService(tt.fields.config, tt.fields.rulesValidatorService, tt.fields.rulesRepository,
				tt.fields.listsService, tt.fields.chargeRepository, tt.fields.familyService,
				tt.fields.familyCompaniesService, tt.fields.chargebackRepository, tt.fields.omniscoreService,
				tt.fields.merchantScoreRepository, log, new(datadog.MetricsDogMock),
			)
			got, err := r.EvaluateCharge(context.Background(), tt.args.charge)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateCharge() error = %v, wantErr %v , name %s", err, tt.wantErr, tt.name)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestChargeService_Get(t *testing.T) {
	logger, _ := logs.New()
	t.Run("service returns repository response", func(t *testing.T) {
		chargeRepository := new(mocks.ChargeEvaluationRepositoryMock)
		chargeId := "charge-123"

		service := NewChargeService(
			config.Config{}, nil, nil, nil, chargeRepository, nil, nil, nil, nil, nil, logger, nil)
		chargeRepository.On("Get", nil, chargeId).Return(entities.EvaluationResponse{}, nil)

		response, err := service.Get(nil, chargeId)

		assert.Nil(t, err)
		assert.Empty(t, response)
	})
}

func TestChargeService_GetOnlyRules(t *testing.T) {
	logger, _ := logs.New()

	t.Run("service returns repository response", func(t *testing.T) {
		chargeRepository := new(mocks.ChargeEvaluationRepositoryMock)
		chargeId := "charge-123"

		service := NewChargeService(
			config.Config{}, nil, nil, nil, chargeRepository, nil, nil, nil, nil, nil, logger, nil)
		chargeRepository.On("GetOnlyRules", nil, chargeId).Return(entities.RulesEvaluationResponse{}, nil)

		response, err := service.GetOnlyRules(nil, chargeId)

		assert.Nil(t, err)
		assert.Empty(t, response)
	})
}

func getMockServiceFirstCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}
	charge := testdata.GetDefaultCharge()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().Times(3).
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{CompanyID: charge.CompanyID, Score: defaultScore}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID
	rulesRepositoryMock.On("GetRulesByFilters", context.Background(), entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return([]entities.Rule{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceSecondCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	white := testdata.GetDefaultWhiteList(false)
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	charge := testdata.GetDefaultCharge()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{white}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{CompanyID: charge.CompanyID, Score: defaultScore}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyID: familyID, FamilyCompaniesIDs: familyCompaniesIds}).
		Once().
		Return([]entities.Rule{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceThirdCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	white := testdata.GetDefaultWhiteList(false)
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	charge := testdata.GetDefaultCharge()
	charge.SetDefaultConsole()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{white}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: charge.Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{CompanyID: charge.CompanyID, Score: defaultScore}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID
	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyID: familyID, FamilyCompaniesIDs: familyCompaniesIds}).
		Once().
		Return([]entities.Rule{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceFourthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	charge := testdata.GetDefaultCharge()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{testdata.GetDefaultBlackList(false)}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultChargeBlacklist().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyID: familyID, FamilyCompaniesIDs: familyCompaniesIds}).
		Once().
		Return([]entities.Rule{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceFifthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleFingerprintBlocked(false)}
	familyFilter := entities.FamilyFilter{
		ID:                   "",
		Mccs:                 []string{testdata.GetChargeWithDeviceFingerprintBlocked().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetChargeWithDeviceFingerprintBlocked().CompanyID},
	}

	charge := testdata.GetChargeWithDeviceFingerprintBlocked()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetChargeWithDeviceFingerprintBlocked().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceSixthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleEmailBlockedGlobal(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyID},
	}

	charge := testdata.GetChargeWithEmailBlockedGlobal()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetChargeWithEmailBlockedGlobal().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetChargeWithEmailBlockedGlobal().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.GlobalRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{{CompanyIDs: familyCompaniesFilter.CompanyIDs}}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceNinthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	mcc := []string{testdata.GetDefaultCharge().CompanyMCC}
	rulesMock := []entities.Rule{testdata.GetDefaultRuleWithFamilyMccID(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 mcc,
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	charge := testdata.GetDefaultCharge()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID
	familyID := testdata.GetDefaultFamily().ID.Hex()
	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyID: familyID}, entities.FamilyCompanyRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(testdata.GetDefaultFamily(), nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceEleventhCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	companyIDs := testdata.GetFamilyCompaniesWithMatchingCompanyIDs().CompanyIDs
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleWithFamilyMccID(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultFamilyCompaniesIDCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultFamilyCompaniesIDCharge().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: companyIDs,
	}

	charge := testdata.GetDefaultFamilyCompaniesIDCharge()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultFamilyCompaniesIDCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetDefaultFamilyCompaniesIDCharge().CompanyID
	familyCompaniesID := testdata.GetDefaultFamilyCompanies().ID.Hex()
	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyCompaniesIDs: []string{familyCompaniesID}}, entities.FamilyMccRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{testdata.GetDefaultFamilyCompanies()}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceThirteenthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	companyIDs := testdata.GetFamilyCompaniesWithMatchingCompanyIDs().CompanyIDs

	rulesMock := []entities.Rule{testdata.GetDefaultRuleWithFamilyMccID(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultChargeFamilyMcc().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultChargeFamilyMcc().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: companyIDs,
	}

	charge := testdata.GetDefaultChargeFamilyMcc()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultFamilyCompaniesIDCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetDefaultFamilyCompaniesIDCharge().CompanyID
	familyCompaniesID := testdata.GetDefaultFamilyCompanies().ID.Hex()
	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyCompaniesIDs: []string{familyCompaniesID, familyCompaniesId, familyCompaniesId}}, entities.FamilyMccRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{testdata.GetDefaultFamilyCompanies(), {}, {}}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceFourteenthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	gray := testdata.GetDefaultGrayList(true)

	rulesMock := []entities.Rule{testdata.GetDefaultRule(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{""},
		NotExcludedCompanies: []string{""},
		ID:                   "",
		Name:                 "",
		Paged:                false,
	}
	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{""},
		ID:         "",
		Name:       "",
		Paged:      false,
	}

	charge := testdata.GetDefaultChargeInGraylist()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{gray}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultChargeInGraylist().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{FamilyID: familyID, FamilyCompaniesIDs: make([]string, 0)}).
		Once().
		Return(rulesMock, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceFifteenthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	gray := testdata.GetDefaultGrayList(false)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleEmailBlockedGlobalForGraylist(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{""},
		NotExcludedCompanies: []string{""},
		ID:                   "",
		Name:                 "",
		Paged:                false,
	}
	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{""},
		ID:         "",
		Name:       "",
		Paged:      false,
	}

	charge := testdata.GetDefaultChargeInGraylist()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{gray}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultChargeInGraylist().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{}, entities.GlobalRulesType).
		Once().
		Return(rulesMock, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceSixteenthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleEmailWithChargebacks(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{""},
		NotExcludedCompanies: []string{"7683457364"},
		ID:                   "",
		Name:                 "",
		Paged:                false,
	}
	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{"7683457364"},
		ID:         "",
		Name:       "",
		Paged:      false,
	}

	charge := testdata.GetDefaultChargeWithChargebacks()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultChargeWithChargebacks().Details.Email}).
		Once().
		Return(testdata.GetPayerDefeult(), nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: charge.CompanyID}, entities.GlobalRulesType).
		Once().
		Return(rulesMock, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceSeventeenthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleWithOmniscore(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{""},
		NotExcludedCompanies: []string{*testdata.GetDefaultRuleWithOmniscore(false).CompanyID},
		ID:                   "",
		Name:                 "",
		Paged:                false,
	}
	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{"7683457364"},
		ID:         "",
		Name:       "",
		Paged:      false,
	}

	charge := testdata.GetChargeForOmniscoreRule()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetChargeForOmniscoreRule().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: "7683457364"}, entities.GlobalRulesType).
		Once().
		Return(rulesMock, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceNineteenthCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	charge := testdata.GetDefaultCharge()
	blacklist := testdata.GetDefaultBlackList(false)
	blacklist.Type = "Blaklist"
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{blacklist}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultChargeBlacklist().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}).
		Once().
		Return([]entities.Rule{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceTwentiethCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	rulesMock := []entities.Rule{testdata.GetDefaultRuleMerchantScoreApproved(false)}

	charge := testdata.GetDefaultCharge()

	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultChargeBlacklist().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{CompanyID: charge.CompanyID, Score: defaultScore}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceTwentyfirstCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	rulesMock := []entities.Rule{testdata.GetDefaultRuleMerchantScoreDeclined(false)}

	charge := testdata.GetDefaultCharge()

	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultChargeBlacklist().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{CompanyID: charge.CompanyID, Score: defaultScore}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceTwentySecondCase() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	rulesMock := []entities.Rule{testdata.GetDefaultRuleMarketSegmentApproved(false)}

	charge := testdata.GetDefaultCharge()

	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{CompanyID: charge.CompanyID, Score: 0.1}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getOmniscoreTestCases() (omniscores.OmniscoreService, omniscores.OmniscoreService) {
	OmniscoreOnMock := new(mocks.OmniscoreServiceMock)
	OmniscoreOffMock := new(mocks.OmniscoreServiceMock)

	OmniscoreOnMock.On("GetScore", context.Background(), mock.AnythingOfType("ChargeRequest")).Return(0.4)
	OmniscoreOffMock.On("GetScore", context.Background(), mock.AnythingOfType("ChargeRequest")).Return(rest.DefaultScore)

	return OmniscoreOnMock, OmniscoreOffMock
}

func setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock *mocks.FamilyCompaniesServiceMock) {
	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{testdata.GetDefaultCharge().CompanyID},
	}
	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{{}}, nil)
}

func getMockServiceFirstCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}
	charge := testdata.GetDefaultCharge()

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{CompanyID: charge.CompanyID, Score: defaultScore}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID
	rulesRepositoryMock.On("GetRulesByFilters", context.Background(), entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return([]entities.Rule{}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceSecondCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}
	charge := testdata.GetDefaultCharge()

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{CompanyID: charge.CompanyID, Score: defaultScore}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID
	rulesRepositoryMock.On("GetRulesByFilters",
		context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return([]entities.Rule{testdata.GetDefaultRuleCompanyRuleAccepted(false)}, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceThirdCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	mcc := []string{testdata.GetDefaultCharge().CompanyMCC}
	rulesMock := []entities.Rule{testdata.GetDefaultRuleWithFamilyMccID(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 mcc,
		NotExcludedCompanies: []string{testdata.GetDefaultCharge().CompanyID},
	}

	charge := testdata.GetDefaultCharge()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetDefaultCharge().CompanyID
	familyID := testdata.GetDefaultFamily().ID.Hex()
	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyID: familyID}, entities.FamilyCompanyRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(testdata.GetDefaultFamily(), nil)

	setDefaultFilterFamilyCompaniesServiceMock(familyCompaniesServiceMock)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceFourthCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	companyIDs := testdata.GetFamilyCompaniesWithMatchingCompanyIDs().CompanyIDs
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleWithFamilyMccID(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetDefaultFamilyCompaniesIDCharge().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetDefaultFamilyCompaniesIDCharge().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: companyIDs,
	}

	charge := testdata.GetDefaultFamilyCompaniesIDCharge()
	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetDefaultFamilyCompaniesIDCharge().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetDefaultFamilyCompaniesIDCharge().CompanyID
	familyCompaniesID := testdata.GetDefaultFamilyCompanies().ID.Hex()
	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyCompaniesIDs: []string{familyCompaniesID}}, entities.FamilyMccRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{testdata.GetDefaultFamilyCompanies()}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceFifthCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleEmailBlockedGlobal(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyID},
	}

	charge := testdata.GetChargeWithEmailBlockedGlobal()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetChargeWithEmailBlockedGlobal().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetChargeWithEmailBlockedGlobal().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.GlobalRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{{CompanyIDs: familyCompaniesFilter.CompanyIDs}}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceSixthCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesMock := []entities.Rule{testdata.GetDefaultRuleEmailBlockedGlobal(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{testdata.GetChargeWithEmailBlockedGlobal().CompanyID},
	}

	charge := testdata.GetChargeWithEmailBlockedGlobal()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetChargeWithEmailBlockedGlobal().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetChargeConsoleRules().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return([]entities.Rule{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyID: "000000000000000000000000"},
		entities.FamilyCompanyRulesType).
		Once().
		Return([]entities.Rule{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyCompaniesIDs: []string{"000000000000000000000000"}},
		entities.FamilyMccRulesType).
		Once().
		Return([]entities.Rule{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.GlobalRulesType).
		Once().
		Return(rulesMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{{CompanyIDs: familyCompaniesFilter.CompanyIDs}}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceSevenCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesCompanyMock := []entities.Rule{testdata.GetDefaultRuleEmailGlobalUndefined(false)}
	rulesIdentityModuleMock := []entities.Rule{testdata.GetDefaultRuleEmailProximity(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetChargeWithEmailProximity().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetChargeWithEmailProximity().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{testdata.GetChargeWithEmailProximity().CompanyID},
	}

	charge := testdata.GetChargeWithEmailProximity()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetChargeWithEmailBlockedGlobal().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetChargeWithEmailProximity().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.CompanyRulesType).
		Once().
		Return(rulesCompanyMock, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyID: "000000000000000000000000"},
		entities.FamilyCompanyRulesType).
		Once().
		Return([]entities.Rule{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID, FamilyCompaniesIDs: []string{"000000000000000000000000"}},
		entities.FamilyMccRulesType).
		Once().
		Return([]entities.Rule{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.GlobalRulesType).
		Once().
		Return([]entities.Rule{}, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.IdentityModuleType).
		Once().
		Return(rulesIdentityModuleMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{{CompanyIDs: familyCompaniesFilter.CompanyIDs}}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceEighthCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesYellowFlagMock := []entities.Rule{testdata.GetDefaultRuleYellowFlag(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetChargeYellowFlag().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetChargeYellowFlag().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{testdata.GetChargeYellowFlag().CompanyID},
	}

	charge := testdata.GetChargeYellowFlag()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetChargeWithEmailBlockedGlobal().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetChargeYellowFlag().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.YellowFlagType).
		Once().
		Return(rulesYellowFlagMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{{CompanyIDs: familyCompaniesFilter.CompanyIDs}}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}

func getMockServiceNinthCaseOnlyRules() (rules.RuleRepository, lists.ListsService, families.FamilyService, familycom.FamilyCompaniesService, chargebacks.ChargebackRepository, merchantsscore.MerchantsScoreRepository) {
	rulesRepositoryMock := new(mocks.RulesRepositoryMock)
	listServiceMock := new(mocks.ListsServiceMock)
	familyServiceMock := new(mocks.FamilyServiceMock)
	familyCompaniesServiceMock := new(mocks.FamilyCompaniesServiceMock)
	chargebacksRepositoryMock := new(mocks.ChargebackRepositoryMock)
	merchantsScoreRepositoryMock := new(mocks.MerchantsScoreRepositoryMock)

	rulesYellowFlagMock := []entities.Rule{testdata.GetDefaultRuleYellowFlag(false)}
	rulesGlobalMock := []entities.Rule{testdata.GetDefaultRuleEmailBlockedGlobal(false)}
	familyFilter := entities.FamilyFilter{
		Mccs:                 []string{testdata.GetChargeYellowFlagAndGlobal().CompanyMCC},
		NotExcludedCompanies: []string{testdata.GetChargeYellowFlagAndGlobal().CompanyID},
	}

	familyCompaniesFilter := entities.FamilyCompaniesFilter{
		CompanyIDs: []string{testdata.GetChargeYellowFlagAndGlobal().CompanyID},
	}

	charge := testdata.GetChargeYellowFlagAndGlobal()
	listServiceMock.On("GetLists", context.Background(), charge.NewListsSearch()).
		Once().
		Return([]entities.List{}, nil)

	chargebacksRepositoryMock.On("Find", context.Background(), entities.Payer{Email: testdata.GetChargeWithEmailBlockedGlobal().Details.Email}).
		Once().
		Return(entities.Payer{}, nil)

	merchantsScoreRepositoryMock.On("FindByMerchantID", context.Background(), charge.CompanyID).
		Once().
		Return(entities.MerchantScore{}, nil)

	companyID := testdata.GetChargeYellowFlag().CompanyID

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.YellowFlagType).
		Once().
		Return(rulesYellowFlagMock, nil)

	rulesRepositoryMock.On("GetRulesByFilters", context.Background(),
		entities.RuleFilter{CompanyID: companyID}, entities.GlobalRulesType).
		Once().
		Return(rulesGlobalMock, nil)

	familyServiceMock.On("GetFamily", context.Background(), familyFilter).
		Once().
		Return(entities.Family{}, nil)

	familyCompaniesServiceMock.On("GetFamiliesCompaniesFromFilter", context.Background(), familyCompaniesFilter).
		Once().
		Return([]entities.FamilyCompanies{{CompanyIDs: familyCompaniesFilter.CompanyIDs}}, nil)

	return rulesRepositoryMock, listServiceMock, familyServiceMock, familyCompaniesServiceMock, chargebacksRepositoryMock, merchantsScoreRepositoryMock
}
