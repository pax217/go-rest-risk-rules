package config

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		ProjectName    string `default:"risk-rules"`
		ProjectVersion string `envconfig:"PROJECT_VERSION" default:"1.9.0"`
		Port           string `envconfig:"PORT" default:"8000" required:"true"`
		Env            string `envconfig:"ENV" default:"local"`
		MongoDB        struct {
			Collections struct {
				Rules                      string `envconfig:"RULES" default:"rules"`
				Lists                      string `envconfig:"LISTS" default:"lists"`
				Conditions                 string `envconfig:"CONDITIONS" default:"conditions"`
				Operators                  string `envconfig:"OPERATORS" default:"operators"`
				Modules                    string `envconfig:"MODULES" default:"modules"`
				Fields                     string `envconfig:"FIELDS" default:"fields"`
				ChargeEvaluations          string `envconfig:"CHARGE_EVALUATIONS" default:"charge_evaluations"`
				ChargeEvaluationsOnlyRules string `envconfig:"CHARGE_EVALUATIONS_ONLY_RULES" default:"charge_evaluations_only_rules"`
				Families                   string `envconfig:"FAMILIES" default:"families"`
				FamilyCompanies            string `envconfig:"FAMILY_COMPANIES" default:"family_companies"`
				Payers                     string `envconfig:"PAYERS" default:"payers"`
				MerchantsScore             string `envconfig:"MERCHANTS_SCORE" default:"merchants_score"`
			}
			Database string `envconfig:"MONGODB_DATABASE" default:"rules"`
			URI      string `envconfig:"MONGODB_URI" default:"mongodb://localhost:27017"`
		}
		Metrics struct {
			Host string `envconfig:"DATADOG_METRICS_HOST" default:"localhost"`
			Port string `envconfig:"DATADOG_METRICS_PORT" default:"8125"`
		}
		RequestHeaderToken string `envconfig:"REQUEST_HEADER_TOKEN"`
		EventBus           struct {
			Chargebacks struct {
				BoostrapServers         string `envconfig:"KAFKA_CHARGEBACK_BOOSTRAP_SERVERS" default:"localhost:19094"`
				Topic                   string `envconfig:"KAFKA_CHARGEBACK_TOPIC" default:"risk.chargebacks.created"`
				GroupID                 string `envconfig:"KAFKA_CHARGEBACK_GROUP_ID" default:"risk_rules_group1"`
				EnabledAuth             bool   `envconfig:"KAFKA_CHARGEBACK_ENABLED_AUTH" default:"true"`
				EnabledSslCertification bool   `envconfig:"KAFKA_CHARGEBACK_ENABLED_SSL_CERTIFICATION" default:"false"`
				Mechanism               string `envconfig:"KAFKA_CHARGEBACK_MECHANISM" default:"SCRAM-SHA-512"`
				SecurityProtocol        string `envconfig:"KAFKA_CHARGEBACK_SECURITY_PROTOCOL" default:"SASL_SSL"`
				Password                string `envconfig:"KAFKA_CHARGEBACK_PASSWORD" default:"password"`
				User                    string `envconfig:"KAFKA_CHARGEBACK_USER" default:"metricsreporter"`
			}
		}
		Omniscore struct {
			IsEnabled           bool   `envconfig:"IS_OMNISCORE_ENABLED" default:"false"`
			Host                string `envconfig:"OMNISCORE_HOST" default:"http://localhost:3000"`
			TimeoutMilliseconds int    `envconfig:"OMNISCORE_TIMEOUT_MILLISECONDS" default:"5000"`
		}
		InternalService struct {
			Host                string `envconfig:"INTERNAL_SERVICE_HOST" default:"http://localhost:3000"`
			TimeoutMilliseconds int    `envconfig:"INTERNAL_SERVICE_TIMEOUT_MILLISECONDS" default:"5000"`
		}
		MerchantScore struct {
			IsEnabled    bool   `envconfig:"IS_MERCHANT_SCORE_ENABLED" default:"false"`
			S3Bucket     string `envconfig:"S3_BUCKET" default:"testbucket"`
			S3PrefixFile string `envconfig:"S3_PREFIX_FILE" default:"merchant_score"`
			Region       string `envconfig:"AWS_REGION" default:"us-east-1"`
		}
	}
)

var (
	Configs Config
)

func NewConfig() Config {
	if err := envconfig.Process("", &Configs); err != nil {
		panic(err.Error())
	}

	return Configs
}
