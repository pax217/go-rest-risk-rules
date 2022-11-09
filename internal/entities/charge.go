package entities

import (
	"github.com/mitchellh/mapstructure"
)

type ChargeRequest struct {
	ID                  string                  `json:"_id" mapstructure:"id" bson:"_id" validate:"required"`
	Amount              float64                 `json:"amount" mapstructure:"amount" validate:"required"`
	DeviceFingerprint   string                  `json:"device_fingerprint" mapstructure:"device_fingerprint" bson:"device_fingerprint"`
	OrderID             string                  `json:"order_id" mapstructure:"order_id" bson:"order_id" `
	Status              string                  `json:"status" mapstructure:"status"`
	CompanyID           string                  `json:"company_id" mapstructure:"company_id" bson:"company_id" validate:"required"`
	CompanyMCC          string                  `json:"company_mcc" mapstructure:"company_mcc" bson:"company_mcc" validate:"required"`
	MonthlyInstallments int                     `json:"monthly_installments" mapstructure:"monthly_installments" bson:"monthly_installments"`
	LiveMode            bool                    `json:"live_mode" mapstructure:"live_mode" bson:"live_mode"`
	Details             DetailsRequest          `json:"details" mapstructure:"details" validate:"required"`
	PaymentMethod       PaymentMethodRequest    `json:"payment_method" mapstructure:"payment_method" bson:"payment_method" validate:"required"`
	Aggregation         Aggregation             `json:"aggregation" mapstructure:"aggregation"`
	Payer               PayerRequest            `json:"payer" mapstructure:"payer" bson:"payer"`
	IsGraylist          bool                    `json:"is_graylist" mapstructure:"is_graylist" bson:"is_graylist"`
	Omniscore           float64                 `json:"omniscore" mapstructure:"omniscore" bson:"omniscore"`
	Console             []Component             `json:"console" mapstructure:"console" bson:"console"`
	MerchantScore       float64                 `json:"merchant_score" mapstructure:"merchant_score" bson:"merchant_score"`
	EmailProximity      EmailEvaluationResponse `json:"email_proximity" mapstructure:"email_proximity" bson:"email_proximity,omitempty"`
	MarketSegment       string                  `json:"market_segment" mapstructure:"market_segment" bson:"market_segment"`
	IsYellowFlag        bool                    `json:"is_yellow_flag" mapstructure:"is_yellow_flag" bson:"is_yellow_flag"`
}

type Component struct {
	Name     ConsoleComponent `json:"name" mapstructure:"name" bson:"name"`
	Priority []Decision       `json:"priority" mapstructure:"priority" bson:"priority,omitempty"`
}

func (cm *Component) HaveSecondaryDecision() bool {
	for _, value := range ComponentsWithOutSecondaryDecision {
		if value == cm.Name {
			return false
		}
	}

	return true
}

func (c *ChargeRequest) SetDefaultConsole() {
	c.Console = []Component{
		{
			Name:     WhitelistType,
			Priority: []Decision{Accepted},
		},
		{
			Name:     GraylistType,
			Priority: []Decision{Undecided},
		},
		{
			Name:     BlacklistType,
			Priority: []Decision{Declined},
		},
		{
			Name:     CompanyRulesType,
			Priority: []Decision{Accepted, Declined, Undecided},
		},
		{
			Name:     FamilyCompanyRulesType,
			Priority: []Decision{Accepted, Declined, Undecided},
		},
		{
			Name:     FamilyMccRulesType,
			Priority: []Decision{Accepted, Declined, Undecided},
		},
		{
			Name:     GlobalRulesType,
			Priority: []Decision{Accepted, Declined, Undecided},
		},
	}
}

func (c *ChargeRequest) SetDefaultConsoleOnlyRules() {
	c.Console = []Component{
		{
			Name:     CompanyRulesType,
			Priority: []Decision{Accepted, Declined, Undecided},
		},
		{
			Name:     FamilyCompanyRulesType,
			Priority: []Decision{Accepted, Declined, Undecided},
		},
		{
			Name:     FamilyMccRulesType,
			Priority: []Decision{Accepted, Declined, Undecided},
		},
		{
			Name:     GlobalRulesType,
			Priority: []Decision{Accepted, Declined, Undecided},
		},
	}
}

func (c *ChargeRequest) ValidateConsole() {
	if len(c.Console) == 0 {
		c.SetDefaultConsole()
	}
}

func (c *ChargeRequest) ValidateConsoleOnlyRules() {
	if len(c.Console) == 0 {
		c.SetDefaultConsoleOnlyRules()
	} else {
		c.excludeListComponents()
	}
}

func (c *ChargeRequest) excludeListComponents() {
	rulesList := make([]Component, 0)
	for _, component := range c.Console {
		if !component.Name.IsList() {
			rulesList = append(rulesList, component)
		}
	}
	c.Console = rulesList
}

func (c *ChargeRequest) ToMap() (map[string]interface{}, error) {
	var mapCharge map[string]interface{}

	err := mapstructure.Decode(c, &mapCharge)
	if err != nil {
		panic(err)
	}

	return mapCharge, nil
}

func (c *ChargeRequest) NewListsSearch() ListsSearch {
	return ListsSearch{
		CompanyID: c.CompanyID,
		Email:     c.Details.Email,
		CardHash:  c.PaymentMethod.CardHash,
		Phone:     c.Details.Phone,
	}
}

type PaymentMethodRequest struct {
	BinNumber string `json:"bin_number" mapstructure:"bin_number" bson:"bin_number"`
	Brand     string `json:"brand" mapstructure:"brand"`
	CardType  string `json:"card_type" mapstructure:"card_type" bson:"card_type"`
	CardHash  string `json:"card_hash" mapstructure:"card_hash" bson:"-"`
	Country   string `json:"country" mapstructure:"country"`
	Issuer    string `json:"issuer" mapstructure:"issuer"`
}

type DetailsRequest struct {
	Email     string `json:"email" mapstructure:"email" validate:"required"`
	IPAddress string `json:"ip_address" mapstructure:"ip_address" bson:"ip_address"`
	Phone     string `json:"phone" mapstructure:"phone" bson:"phone"`
	Name      string `json:"name" mapstructure:"name"`
}

type Aggregation struct {
	PayerAggregation AggregationAttribute `json:"payer" mapstructure:"payer"`
	PayerCompany     AggregationAttribute `json:"payer_company" mapstructure:"payer_company" bson:"payer_company"`
	BinNumber        AggregationAttribute `json:"bin_number" mapstructure:"bin_number" bson:"bin_number"`
	CardHash         AggregationAttribute `json:"card_hash" mapstructure:"card_hash" bson:"card_hash"`
}

type AggregationAttribute struct {
	Charge AggregationEvent `json:"charge" mapstructure:"charge"`
}

type AggregationEvent struct {
	M1   AggregationEventProperties `json:"m1" mapstructure:"m1" bson:"m_1"`
	M5   AggregationEventProperties `json:"m5" mapstructure:"m5" bson:"m_5"`
	H1   AggregationEventProperties `json:"h1" mapstructure:"h1" bson:"h_1"`
	H2   AggregationEventProperties `json:"h2" mapstructure:"h2" bson:"h_2"`
	H3   AggregationEventProperties `json:"h3" mapstructure:"h3" bson:"h_3"`
	H6   AggregationEventProperties `json:"h6" mapstructure:"h6" bson:"h_6"`
	H9   AggregationEventProperties `json:"h9" mapstructure:"h9" bson:"h_9"`
	H12  AggregationEventProperties `json:"h12" mapstructure:"h12" bson:"h_12"`
	D1   AggregationEventProperties `json:"d1" mapstructure:"d1" bson:"d_1"`
	D2   AggregationEventProperties `json:"d2" mapstructure:"d2" bson:"d_2"`
	D3   AggregationEventProperties `json:"d3" mapstructure:"d3" bson:"d_3"`
	D4   AggregationEventProperties `json:"d4" mapstructure:"d4" bson:"d_4"`
	D7   AggregationEventProperties `json:"d7" mapstructure:"d7" bson:"d_7"`
	D15  AggregationEventProperties `json:"d15" mapstructure:"d15" bson:"d_15"`
	D30  AggregationEventProperties `json:"d30" mapstructure:"d30" bson:"d_30"`
	D60  AggregationEventProperties `json:"d60" mapstructure:"d60" bson:"d_60"`
	D120 AggregationEventProperties `json:"d120" mapstructure:"d120" bson:"d_120"`
}

type AggregationEventProperties struct {
	Sum   int `json:"sum" mapstructure:"sum"`
	Count int `json:"count" mapstructure:"count"`
}

type PayerRequest struct {
	Chargebacks int64 `json:"chargebacks" mapstructure:"chargebacks" bson:"chargebacks"`
}
