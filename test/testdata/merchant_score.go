package testdata

import "github.com/conekta/risk-rules/internal/entities"

func GetRawDataMerchantsScore() [][]string {
	return [][]string{
		{"id", "company_id", "score"},
		{"1", "1", "1.0"},
		{"2", "2", "0.56"},
	}
}

func GetDefaultMerchantScoreData() []entities.MerchantScore {
	return []entities.MerchantScore{
		{
			CompanyID: "1",
			Score:     1.0,
		},
		{
			CompanyID: "2",
			Score:     0.56,
		}}
}
