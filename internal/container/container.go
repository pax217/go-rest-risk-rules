package container

import (
	"context"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"

	"github.com/conekta/risk-rules/internal/apps/chargebacks"
	"github.com/conekta/risk-rules/internal/apps/charges"
	"github.com/conekta/risk-rules/internal/apps/conditions"
	"github.com/conekta/risk-rules/internal/apps/families"
	familycom "github.com/conekta/risk-rules/internal/apps/family_companies"
	"github.com/conekta/risk-rules/internal/apps/fields"
	"github.com/conekta/risk-rules/internal/apps/lists"
	merchantsscore "github.com/conekta/risk-rules/internal/apps/merchants_score"
	"github.com/conekta/risk-rules/internal/apps/modules"
	"github.com/conekta/risk-rules/internal/apps/omniscores"
	"github.com/conekta/risk-rules/internal/apps/operators"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/apps/status"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/pkg/csv"
	"github.com/conekta/risk-rules/pkg/mongodb"
	"github.com/conekta/risk-rules/pkg/rest"
)

type Dependencies struct {
	StatusHandler          status.StatusHandler
	RulesHandler           rules.RuleHandler
	ChargeHandler          charges.ChargeHandler
	ModulesHandler         modules.ModuleHandler
	OperatorHandler        operators.OperatorHandler
	ConditionsHandler      conditions.ConditionHandler
	FieldsHandler          fields.FieldHandler
	FamilyHandler          families.FamilyHandler
	FamilyCompaniesHandler familycom.FamilyCompaniesHandler
	ChargebacksHandler     chargebacks.ChargebackHandler
	MerchantsScoreHandler  merchantsscore.MerchantsScoreHandler
	Config                 config.Config
	S3Reader               csv.S3Reader
	Logs                   logs.Logger
}

func Build() Dependencies {
	dependencies := Dependencies{}

	configs := config.NewConfig()

	logger, err := logs.New(logs.LoggerLevel(logs.Info))
	if err != nil {
		return dependencies
	}
	dependencies.Logs = logger

	mongoDB := mongodb.NewMongoDB(configs)
	metric := datadog.NewMetric(context.TODO(), logger, configs.Metrics.Host, configs.Metrics.Port)

	rulesValidator := rules.NewRulesValidator(dependencies.Logs)

	rulesMongoDBRepository := rules.NewRuleMongoDBRepository(configs, mongoDB, dependencies.Logs)
	modulesMongoDBRepository := modules.NewModulesMongoRepository(configs, mongoDB, dependencies.Logs)
	operatorMongoDBRepository := operators.NewOperatorMongoDBRepository(configs, mongoDB, dependencies.Logs)
	fieldsMongoDBRepository := fields.NewFieldsMongoDBRepository(configs, mongoDB, dependencies.Logs)
	conditionsMongoDBRepository := conditions.NewConditionsRepository(configs, mongoDB, dependencies.Logs)
	chargesMongoDBRepository := charges.NewChargeMongoDBRepository(configs, mongoDB, dependencies.Logs)
	familiesMongoDBRepository := families.NewFamilyMongoDBRepository(configs, mongoDB, dependencies.Logs)
	familyCompaniesMongoDBRepository := familycom.NewFamilyCompaniesMongoDBRepository(configs, mongoDB, dependencies.Logs)
	chargebacksMongoDBRepository := chargebacks.NewChargebacksMongoDBRepository(configs, mongoDB, dependencies.Logs)
	omniscoreRestClient := rest.NewOmniscoreClient(configs, dependencies.Logs)
	listsClient := rest.NewRkListsRestClient(configs, logger)
	merchantsScoreMongoDBRepository := merchantsscore.NewMerchantsMongoDBRepository(configs, mongoDB, dependencies.Logs)
	s3CsvReader := csv.NewS3Reader(configs, dependencies.Logs)
	merchantRepositoryS3 := merchantsscore.NewMerchantScoreS3Repository(configs, dependencies.Logs, s3CsvReader)

	modulesService := modules.NewModuleService(configs, modulesMongoDBRepository, dependencies.Logs, metric)
	operatorService := operators.NewOperatorService(configs, dependencies.Logs, operatorMongoDBRepository, metric)
	listsService := lists.NewListsService(configs, logger, metric, listsClient)
	fieldsService := fields.NewFieldsService(configs, fieldsMongoDBRepository, logger, metric)
	conditionsService := conditions.NewConditionsService(configs, conditionsMongoDBRepository, logger, metric)
	familiesService := families.NewFamilyService(configs, familiesMongoDBRepository, rulesMongoDBRepository, logger, metric)
	rulesService := rules.NewRulesService(configs, rulesValidator, rulesMongoDBRepository, logger, metric)
	familyCompaniesService := familycom.NewFamilyCompaniesService(configs, familyCompaniesMongoDBRepository,
		rulesMongoDBRepository, logger, metric)
	omniscoreService := omniscores.NewOmniscoreService(configs, logger, omniscoreRestClient)
	chargeService := charges.NewChargeService(configs, rulesValidator, rulesMongoDBRepository,
		listsService, chargesMongoDBRepository, familiesService, familyCompaniesService, chargebacksMongoDBRepository,
		omniscoreService, merchantsScoreMongoDBRepository, dependencies.Logs, metric)
	chargebackService := chargebacks.NewChargebacksService(configs, chargebacksMongoDBRepository, logger, metric)
	merchantsScoreService := merchantsscore.NewMerchantsScoreService(configs, logger, metric,
		merchantsScoreMongoDBRepository, merchantRepositoryS3)

	dependencies.StatusHandler = status.NewStatusHandler(configs, metric)
	dependencies.RulesHandler = rules.NewRulesHandler(configs, rulesService, logger)
	dependencies.RulesHandler = rules.NewRulesHandler(configs, rulesService, dependencies.Logs)
	dependencies.ChargeHandler = charges.NewChargeHandler(configs, chargeService, dependencies.Logs, metric)
	dependencies.OperatorHandler = operators.NewOperatorHandler(dependencies.Logs, operatorService)
	dependencies.ModulesHandler = modules.NewModuleHandler(modulesService, dependencies.Logs)
	dependencies.FieldsHandler = fields.NewFieldsHandler(configs, fieldsService, logger)
	dependencies.ConditionsHandler = conditions.NewConditionsHandler(conditionsService, dependencies.Logs)
	dependencies.FamilyHandler = families.NewFamilyHandler(familiesService, logger)
	dependencies.FamilyCompaniesHandler = familycom.NewFamilyCompaniesHandler(familyCompaniesService, logger)
	dependencies.ChargebacksHandler = chargebacks.NewChargebackHandler(chargebackService, configs, logger, metric)
	dependencies.MerchantsScoreHandler = merchantsscore.NewMerchantsScoreHandler(configs, logger, merchantsScoreService)
	dependencies.Config = configs

	return dependencies
}
