package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/conekta/Conekta-Golang-Rules-Engine/parser"
	"github.com/conekta/risk-rules/internal/apps/rules"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/mongodb"
)

type Rule struct {
	BeginSeparator string `json:"begin_separator"`
	EndSeparator   string `json:"end_separator"`
	Operator       string `json:"operator"`
	Condition      string `json:"condition"`
	Field          string `json:"field"`
	Value          string `json:"value"`
}

func Build(rule Rule, isLastOne bool) string {
	if isLastOne {
		return fmt.Sprintf("%s %s %s", rule.Field, rule.Operator, rule.Value)
	}

	return fmt.Sprintf("%s %s %s %s ", rule.Field, rule.Operator, rule.Value, rule.Condition)
}

func BuildRules(r []Rule) string {
	ruleResult := ""
	for idx, rule := range r {
		lastOne := len(r)-1 == idx
		ruleResult += Build(rule, lastOne)
	}
	return ruleResult
}

func FormatRules() {
	r := []Rule{
		{
			BeginSeparator: "(",
			Operator:       ">",
			Condition:      "and",
			Field:          "age",
			Value:          "27",
		},
		{
			EndSeparator: ")",
			Operator:     "<",
			Condition:    "or",
			Field:        "age",
			Value:        "58",
		},
		{
			Operator:  "<",
			Condition: "",
			Field:     "age",
			Value:     "58",
		},
	}

	fmt.Println("building rule: ", BuildRules(r))

	company := entities.Company{
		ID: "0",
		Providers: []entities.Provider{
			{
				ID:         "3",
				Name:       "3",
				Precedence: 3,
				IsOn:       false,
			},
			{
				ID:         "1",
				Name:       "1",
				Precedence: 1,
				IsOn:       false,
			},
			{
				ID:         "2",
				Name:       "2",
				Precedence: 2,
				IsOn:       false,
			},
		},
	}

	fmt.Println("not sorting providers: ", company)
	company.Sort()
	fmt.Println("sorting providers: ", company)
}

func EvaluateRule() {
	charge := entities.ChargeRequest{
		Amount:              28,
		DeviceFingerprint:   "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "60ad5c44926c8400016cbfdc",
		MonthlyInstallments: 12,
		Details: entities.DetailsRequest{
			Email:     "eliosf27@gmail.com12",
			IPAddress: "127.0.0.1",
			Phone:     "+52477266334212",
			Name:      "de M12",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "cGfNEDJZjyj",
			Country:  "US",
		},
	}
	chargeMap, err := charge.ToMap()
	if err != nil {
		panic(err)
	}

	r := []string{
		"not monthly_installments eq 1500", // not equal
		"monthly_installments eq 12",       // equal
		"monthly_installments == 12",       // equal
		"not monthly_installments lt 12",   // not less than
		"monthly_installments lt 12",       // less than
		"monthly_installments < 12",        // less than
		"not monthly_installments le 12",   // not less than equal to
		"monthly_installments le 12",       // less than equal to
		"monthly_installments <= 12",       // less than equal to
		"not monthly_installments gt 12",   // not greater than
		"monthly_installments gt 12",       // greater than
		"monthly_installments > 12",        // greater than
		"not monthly_installments ge 12",   // not greater than equal to
		"monthly_installments ge 12",       // greater than equal to
		"monthly_installments >= 12",       // greater than equal to
		"not company_id co \"test\"",       // not contains
		"company_id co \"test\"",           // contains
		"not company_id ew \"test\"",       // not ends with
		"company_id ew \"test\"",           // ends with
		"not monthly_installments in [12]", // not in a list
		"monthly_installments in [12]",     // in a list
	}
	EvaluateList(r, chargeMap)
}

func EvaluateList(r []string, charge map[string]interface{}) {
	for _, rule := range r {
		result, err := Evaluate(rule, charge)
		fmt.Println("rule: ", rule, " result: ", result, " err: ", err)
	}
}

func Evaluate(rule string, charge map[string]interface{}) (bool, error) {
	if rule == "" {
		return false, fmt.Errorf("[ruler.evaluate] error: empty rule: [%s]", rule)
	}

	newRuleString := strings.ReplaceAll(rule, "\\", "")
	ev, err := parser.NewEvaluator(newRuleString)
	if err != nil {
		er := fmt.Errorf("[ruler.evaluate] error: Making evaluator from the rule [%v], [%v]", rule, err)

		return false, er
	}
	ans, err := ev.Process(charge)
	if err != nil {
		er := fmt.Errorf("[ruler.evaluate] error: Making process from the rule [%v], [%v]", rule, err)

		return false, er
	}
	if err := ev.LastDebugErr(); err != nil {
		er := fmt.Errorf("[ruler.evaluate] error: Last debug error [%v], [%v]", rule, err)

		return false, er
	}

	return ans, nil
}

func InsertRules() {
	configs := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(configs)
	ruleService := rules.NewRulesService(config.Config{}, nil, nil, nil, nil)
	now := time.Now().Truncate(time.Millisecond)

	companyID := "60ad5c44926c8400016cbfdc"
	familyCompanyID := "61e9b8490f32b5e40fa25176"

	rule := entities.Rule{
		CreatedAt:       now,
		CreatedBy:       "jesus.vega@conekta.com",
		IsTest:          false,
		Module:          "policy_compliance",
		Description:     "Regla para validar contratos con OXXO",
		CompanyID:       &companyID,
		FamilyCompanyID: &familyCompanyID,
		IsGlobal:        false,
		Rule:            "(monthly_installments > 8 and monthly_installments < 11)",
		Rules: []entities.RuleContent{
			{
				Field:     "monthly_installments",
				Operator:  ">",
				Value:     "8",
				Condition: "and",
			},
			{
				Field:     "monthly_installments",
				Operator:  "<",
				Value:     "11",
				Condition: "or",
			},
		},
	}
	rule.Rule = ruleService.BuildRule(rule.Rules)

	for i := 0; i < 50; i++ {
		rulesCollection := mongoDB.Collection(configs.MongoDB.Collections.Rules)
		result, err := rulesCollection.InsertOne(context.TODO(), rule, nil)
		if err != nil {
			return
		}
		fmt.Println("result: ", result)
	}
}

func InsertWhiteLists(size int) {
	configs := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(configs)
	emailTemplate := "loyal-g-%d@gmail.com"
	rule := entities.List{
		CreatedAt:   time.Now().Truncate(time.Millisecond),
		CreatedBy:   "jesus.vega@conekta.com",
		IsTest:      false,
		IsGlobal:    true,
		Description: "CLientes fieles",
		Type:        entities.White.String(),
		Field:       entities.EmailField,
		Decision:    entities.Accepted,
	}

	rulesCollection := mongoDB.Collection(configs.MongoDB.Collections.Lists)
	for i := 0; i < size; i++ {
		rule.Value = fmt.Sprintf(emailTemplate, i)
		result, err := rulesCollection.InsertOne(context.TODO(), rule, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("InsertWhiteLists ID: ", result.InsertedID)
	}
	fmt.Println("InsertWhiteLists End")
}

func InsertBlackLists(size int) {
	configs := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(configs)
	emailTemplate := "fraud-g-%d@gmail.com"
	rule := entities.List{
		CreatedAt:   time.Now().Truncate(time.Millisecond),
		CreatedBy:   "santiago.ceron@conekta.com",
		IsTest:      false,
		IsGlobal:    true,
		Description: "fraud email detected by owasp",
		Field:       entities.EmailField,
		Type:        entities.Black.String(),
		Decision:    entities.Declined,
	}

	rulesCollection := mongoDB.Collection(configs.MongoDB.Collections.Lists)

	for i := 0; i < size; i++ {
		rule.Value = fmt.Sprintf(emailTemplate, i)

		result, err := rulesCollection.InsertOne(context.TODO(), rule, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("InsertBlackLists ID: ", result.InsertedID)
	}
	fmt.Println("InsertBlackLists End")
}

func InsertFamilyCompanies(size int) {
	configs := config.NewConfig()
	mongoDB := mongodb.NewMongoDB(configs)
	familyCompanies := entities.FamilyCompanies{
		CompanyIDs: nil,
		CreatedAt:  time.Time{},
		CreatedBy:  "",
		UpdatedAt:  nil,
		UpdatedBy:  nil,
	}

	familyCompaniesCollection := mongoDB.Collection(configs.MongoDB.Collections.FamilyCompanies)

	for i := 0; i < size; i++ {
		var companyIDs []string
		companyIDsSlotsMax := 29
		companyIDsSlots := rand.Intn(companyIDsSlotsMax) + 1
		companyIDMaxValue := 1000
		for j := 0; j < companyIDsSlots; j++ {
			companyIDs = append(companyIDs, fmt.Sprintf("%d", rand.Intn(companyIDMaxValue)))
		}
		familyCompanies.Name = fmt.Sprintf("%d", i)
		familyCompanies.CompanyIDs = companyIDs

		result, err := familyCompaniesCollection.InsertOne(context.TODO(), familyCompanies, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("InsertFamilyCompanies ID: ", result.InsertedID)
	}
	fmt.Println("InsertFamilyCompanies End")
}

func main() {
	registerSize := 1000
	InsertRules()
	InsertBlackLists(registerSize)
	InsertWhiteLists(registerSize)
	InsertFamilyCompanies(registerSize)
}
