package testdata

import (
	"github.com/conekta/risk-rules/internal/entities"
	"time"
)

func GetDefaultCharge() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		Console:       SetDefaultConsole(),
		MerchantScore: -1,
		MarketSegment: "long tail SMB",
	}
}

func GetChargeConsoleIsEmptyOnlyRules() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		Console:       SetDefaultConsoleEmpty(),
		MerchantScore: -1,
		MarketSegment: "long tail SMB",
	}
}

func GetChargeConsoleCompanyRules() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		Console:       SetDefaultConsoleCompany(),
		MerchantScore: -1,
		MarketSegment: "long tail SMB",
	}
}

func GetChargeConsoleFamilyRules() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleFamily(),
		MarketSegment: "long tail SMB",
	}
}

func GetChargeConsoleFamilyMccRules() entities.ChargeRequest {
	companyMCC := "6171ad238730441d5ec9537b"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleFamilyMcc(),
		MarketSegment: "long tail SMB",
	}
}

func GetChargeConsoleGlobalRules() entities.ChargeRequest {
	return entities.ChargeRequest{
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "768345736444",
		CompanyMCC:          "1234",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "block@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "",
			Country:  "US",
		},
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleGlobal(),
		MarketSegment: "long tail SMB",
	}
}

func GetChargeConsoleRules() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "768345736444",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "block@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleOnlyRules(),
		MarketSegment: "long tail SMB",
	}
}

func GetChargeWithOutConsoleRules() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "768345736444",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "block@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
	}
}

func GetDefaultChargeWithoutConsole() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MarketSegment: "long tail SMB",
	}
}

func GetDefaultChargeFamily() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleFamily(),
		MarketSegment: "long tail SMB",
	}
}

func GetDefaultChargeFamilyMcc() entities.ChargeRequest {
	companyMCC := "6171ad238730441d5ec9537b"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleFamilyMcc(),
		MarketSegment: "long tail SMB",
	}
}

func GetDefaultChargeBlacklist() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleBlacklist(),
		MarketSegment: "long tail SMB",
	}
}

func GetDefaultChargeInGraylist() entities.ChargeRequest {

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    true,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleGrayList(),
		MarketSegment: "long tail SMB",
	}
}

func GetDefaultChargeInGraylistAndRule() entities.ChargeRequest {

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    true,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleGrayListAndGlobalRule(),
		MarketSegment: "long tail SMB",
	}
}

func GetDefaultAggregation() entities.Aggregation {
	return entities.Aggregation{
		PayerAggregation: entities.AggregationAttribute{
			Charge: entities.AggregationEvent{
				H1: entities.AggregationEventProperties{
					Sum:   100,
					Count: 200,
				},
				H2: entities.AggregationEventProperties{
					Sum:   300,
					Count: 400,
				},
				H12: entities.AggregationEventProperties{
					Sum:   500,
					Count: 600,
				},
			},
		},
		PayerCompany: entities.AggregationAttribute{
			Charge: entities.AggregationEvent{
				H1: entities.AggregationEventProperties{
					Sum:   34,
					Count: 19,
				},
				H2: entities.AggregationEventProperties{
					Sum:   45,
					Count: 21,
				},
				H12: entities.AggregationEventProperties{
					Sum:   70,
					Count: 45,
				},
			},
		},
		BinNumber: entities.AggregationAttribute{
			Charge: entities.AggregationEvent{
				H1: entities.AggregationEventProperties{
					Sum:   100,
					Count: 20,
				},
				H2: entities.AggregationEventProperties{
					Sum:   150,
					Count: 32,
				},
				H12: entities.AggregationEventProperties{
					Sum:   250,
					Count: 41,
				},
			},
		},
		CardHash: entities.AggregationAttribute{
			Charge: entities.AggregationEvent{
				H1: entities.AggregationEventProperties{
					Count: 100,
					Sum:   50,
				},
				H2: entities.AggregationEventProperties{
					Count: 200,
					Sum:   100,
				},
				H12: entities.AggregationEventProperties{
					Count: 300,
					Sum:   150,
				},
			},
		},
	}
}

func GetZeroDivisionAggregation() entities.Aggregation {
	return entities.Aggregation{
		PayerAggregation: entities.AggregationAttribute{
			Charge: entities.AggregationEvent{
				H1: entities.AggregationEventProperties{
					Sum:   0,
					Count: 0,
				},
				H2: entities.AggregationEventProperties{
					Sum:   300,
					Count: 400,
				},
				H12: entities.AggregationEventProperties{
					Sum:   500,
					Count: 600,
				},
			},
		},
	}
}

func GetChargeWithDeviceFingerprintBlocked() entities.ChargeRequest {
	return entities.ChargeRequest{
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "",
			Country:  "US",
		},
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsole(),
	}
}

func GetChargeWithEmailBlocked() entities.ChargeRequest {
	return entities.ChargeRequest{
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "768345736444",
		CompanyMCC:          "1234",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "block@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "",
			Country:  "US",
		},
		IsGraylist: false,
		Omniscore:  float64(-1),
		Console:    SetDefaultConsole(),
	}
}

func GetChargeWithEmailBlockedGlobal() entities.ChargeRequest {
	return entities.ChargeRequest{
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "768345736444",
		CompanyMCC:          "1234",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "block@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "",
			Country:  "US",
		},
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleGlobal(),
	}
}

func GetChargeWithEmailProximity() entities.ChargeRequest {
	now := time.Now().UTC().Truncate(time.Millisecond)
	return entities.ChargeRequest{
		Amount:              4540,
		DeviceFingerprint:   "fingerblockeed",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "768345736444",
		CompanyMCC:          "1234",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "block@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "",
			Country:  "US",
		},
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		EmailProximity: entities.EmailEvaluationResponse{
			IsValid: false,
			Stats: entities.EmailEvaluationStats{
				ObservationsCount: 8,
				VariationsCount:   2,
				LastSeen:          &now,
			},
		},
		Console: SetConsoleIdentityModule(),
	}
}

func GetChargeYellowFlag() entities.ChargeRequest {
	return entities.ChargeRequest{
		Amount:              801,
		DeviceFingerprint:   "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "6204339524bc5717b46b19de",
		CompanyMCC:          "1234",
		MonthlyInstallments: 2,
		IsYellowFlag:        true,
		Details: entities.DetailsRequest{
			Email:     "block@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "",
			Country:  "US",
		},
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetConsoleYellowFlag(),
	}
}

func GetChargeYellowFlagAndGlobal() entities.ChargeRequest {
	return entities.ChargeRequest{
		Amount:              801,
		DeviceFingerprint:   "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "6204339524bc5717b46b19de",
		CompanyMCC:          "1234",
		MonthlyInstallments: 2,
		IsYellowFlag:        true,
		Details: entities.DetailsRequest{
			Email:     "block@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:    "visa",
			CardType: "credit",
			CardHash: "",
			Country:  "US",
		},
		IsGraylist:    false,
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetConsoleYellowFlagAndGlobal(),
	}
}

func GetChargeMalformedRequest() string {
	request := `{
		"_id": "60902f89bb312295544b9d23",
		"amount": 30,
		"device_fingerprint": "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
		"order_id": "289225303828",
		"status": "pending_payment",
		"company_id": "60ad5c44926c8400016cbfdc",
		"monthly_installments": 100,
		"details": {
			"email": "eliosf2712121@gmail.com",
			"ip_address": "127.0.0.1",
			"phone": "52477266334212",
			"name": "de M12"
		},
		"payment_method": {
			"brand": "visa",
			"card_type": "credit",
			"card_hash": "cGfNEDJZjyj12121",
			"country": "US",
			"bin_number": "3233",
			"issuer": "AMEX",
		}
	}`
	return request
}

func GetChargeInvalidRequest() string {
	request := `{
		"_id": "60902f89bb312295544b9d23",
		"amount": 30,
		"device_fingerprint": "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
		"order_id": "289225303828",
		"status": "pending_payment",
		"company_id": "60ad5c44926c8400016cbfdc",
		"details": {
			"ip_address": "127.0.0.1",
			"phone": "52477266334212",
			"name": "de M12"
		},
		"payment_method": {
			"card_type": "credit",
			"card_hash": "cGfNEDJZjyj12121",
			"country": "US",
			"bin_number": "3233",
			"issuer": "test"
		}
	}`
	return request
}

func GetChargeInvalidRequestCompanyIDRequired() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation:   GetDefaultAggregation(),
		IsGraylist:    false,
		Omniscore:     float64(-1),
		Console:       SetDefaultConsole(),
		MerchantScore: -1,
	}
}

func GetChargeAggregationRequest() string {
	request := `{
	  "_id": "60902f89bb312295544b9d23",
	  "amount": 30,
	  "device_fingerprint": "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
	  "order_id": "289225303828",
	  "monthly_installments": 0,
	  "status": "pending_payment",
	  "company_id": "60ad5c44926c8400016cbfdc",
      "company_mcc": "1234",
	  "details": {
		"email": "eliosf2712121@gmail.com",
		"ip_address": "127.0.0.1",
		"phone": "52477266334212",
		"name": "de M12"
	  },
	  "payment_method": {
		"brand": "visa",
		"card_type": "credit",
		"card_hash": "cGfNEDJZjyj12121",
		"country": "US",
		"bin_number": "3233",
		"issuer": "test"
	  },
	  "aggregation": {
		"payer": {
		  "id": "maria_jose@conekta.com",
		  "charge": {
			"180": {
			  "sum": 1,
			  "count": 10500
			}
		  }
		},
		"payer_company": {},
		"card_hash": {
			"charge": {
				"12": {
					"count": 0,
					"sum": 0
				},
				"3": {
					"count": 0,
					"sum": 0
				},
				"9": {
					"count": 0,
					"sum": 0
				},
			}
	  	}
	}`
	return request
}

func GetChargeWithoutAggregationRequest() string {
	request := `{
	  "_id": "60902f89bb312295544b9d23",
	  "amount": 30,
	  "device_fingerprint": "7uG97A4VkKWdZW5gEBmYZe3m7VQM46oy",
	  "order_id": "289225303828",
	  "monthly_installments": 0,
	  "status": "pending_payment",
	  "company_id": "60ad5c44926c8400016cbfdc",
      "company_mcc": "1234",
	  "details": {
		"email": "eliosf2712121@gmail.com",
		"ip_address": "127.0.0.1",
		"phone": "52477266334212",
		"name": "de M12"
	  },
	  "payment_method": {
		"brand": "visa",
		"card_type": "credit",
		"card_hash": "cGfNEDJZjyj12121",
		"country": "US",
		"bin_number": "3233",
		"issuer": "test"
	  },
	  "aggregation": {}
	}`
	return request
}

func GetChargeWithoutCompanyMccInRequest() entities.ChargeRequest {
	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation: GetDefaultAggregation(),
	}
}

func GetDefaultChargeWithoutAggregation() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
	}
}

func GetDefaultFamilyCompaniesIDCharge() entities.ChargeRequest {
	companyMCC := "1234"

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          companyMCC,
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation: GetDefaultAggregation(),
		Omniscore:   float64(-1),
	}
}

func GetDefaultChargeWithChargebacks() entities.ChargeRequest {

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          "",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation: GetDefaultAggregation(),
		IsGraylist:  false,
		Payer: entities.PayerRequest{
			Chargebacks: 1,
		},
		Omniscore:     float64(-1),
		MerchantScore: -1,
		Console:       SetDefaultConsoleGlobal(),
	}
}

func GetChargeRequestForOmniscoreRule() entities.ChargeRequest {
	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          "",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation: GetDefaultAggregation(),
		IsGraylist:  false,
		Payer: entities.PayerRequest{
			Chargebacks: 0,
		},
		MerchantScore: -1,
		Console:       SetDefaultConsoleGlobal(),
	}
}

func GetChargeForOmniscoreRule() entities.ChargeRequest {

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          "",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation: GetDefaultAggregation(),
		IsGraylist:  false,
		Payer: entities.PayerRequest{
			Chargebacks: 0,
		},
		Omniscore:     0.4,
		MerchantScore: -1,
		Console:       SetDefaultConsoleGlobal(),
	}
}

func GetChargeForOmniscoreRuleNotApplied() entities.ChargeRequest {

	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          "",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation: GetDefaultAggregation(),
		IsGraylist:  false,
		Payer: entities.PayerRequest{
			Chargebacks: 0,
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Console:       SetDefaultConsoleGlobal(),
	}
}

func SetDefaultConsole() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.WhitelistType,
			Priority: []entities.Decision{entities.Accepted},
		},
		{
			Name:     entities.CompanyRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
	}
}

func SetDefaultConsoleOnlyRules() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.CompanyRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
		{
			Name:     entities.FamilyCompanyRulesType,
			Priority: []entities.Decision{entities.Declined, entities.Accepted, entities.Undecided},
		},
		{
			Name:     entities.FamilyMccRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
		{
			Name:     entities.GlobalRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
	}
}

func SetConsoleIdentityModule() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.CompanyRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
		{
			Name:     entities.FamilyCompanyRulesType,
			Priority: []entities.Decision{entities.Declined, entities.Accepted, entities.Undecided},
		},
		{
			Name:     entities.FamilyMccRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
		{
			Name:     entities.GlobalRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
		{
			Name:     entities.IdentityModuleType,
			Priority: []entities.Decision{entities.Undecided},
		},
	}
}

func SetConsoleYellowFlag() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.YellowFlagType,
			Priority: []entities.Decision{entities.Undecided},
		},
	}
}

func SetConsoleYellowFlagAndGlobal() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.YellowFlagType,
			Priority: []entities.Decision{entities.Undecided},
		},
		{
			Name:     entities.GlobalRulesType,
			Priority: []entities.Decision{entities.Declined, entities.Accepted, entities.Undecided},
		},
	}
}

func SetDefaultConsoleBlacklist() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.BlacklistType,
			Priority: []entities.Decision{entities.Declined},
		},
		{
			Name:     entities.CompanyRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
	}
}

func SetDefaultConsoleEmpty() []entities.Component {
	return []entities.Component{}
}

func SetDefaultConsoleCompany() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.CompanyRulesType,
			Priority: []entities.Decision{entities.Declined, entities.Accepted, entities.Undecided},
		},
	}
}

func SetDefaultConsoleFamily() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.FamilyCompanyRulesType,
			Priority: []entities.Decision{entities.Declined, entities.Accepted, entities.Undecided},
		},
	}
}

func SetDefaultConsoleFamilyMcc() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.FamilyMccRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
	}
}

func SetDefaultConsoleGlobal() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.GlobalRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
	}
}

func SetDefaultConsoleGrayList() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.GraylistType,
			Priority: []entities.Decision{entities.Undecided},
		},
	}
}

func SetDefaultConsoleGrayListAndGlobalRule() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.GraylistType,
			Priority: []entities.Decision{entities.Undecided},
		},
		{
			Name:     entities.GlobalRulesType,
			Priority: []entities.Decision{entities.Accepted, entities.Declined, entities.Undecided},
		},
	}
}

func SetDefaultConsoleOnlyBlacklist() []entities.Component {
	return []entities.Component{
		{
			Name:     entities.BlacklistType,
			Priority: []entities.Decision{entities.Declined},
		},
	}
}

func GetChargeRequestForBlacklist() entities.ChargeRequest {
	return entities.ChargeRequest{
		ID:                  "615324eb5bc1dea9ce66068f",
		Amount:              4540,
		DeviceFingerprint:   "cGfNEDJZjyj7W1N7DGbFtQi5RbkxhAvn",
		OrderID:             "289225303828",
		Status:              "pending_payment",
		CompanyID:           "7683457364",
		CompanyMCC:          "",
		MonthlyInstallments: 0,
		Details: entities.DetailsRequest{
			Email:     "mail@hotmail.com",
			IPAddress: "127.0.0.1",
			Phone:     "55-5555-5555",
			Name:      "Mario Moreno",
		},
		PaymentMethod: entities.PaymentMethodRequest{
			Brand:     "visa",
			CardType:  "credit",
			CardHash:  "hash-acd7s-34354",
			Country:   "US",
			BinNumber: "sfc23-5678-c50d",
			Issuer:    "AMEX",
		},
		Aggregation: GetDefaultAggregation(),
		IsGraylist:  false,
		Payer: entities.PayerRequest{
			Chargebacks: 0,
		},
		Omniscore:     -1,
		MerchantScore: -1,
		Console:       SetDefaultConsoleOnlyBlacklist(),
	}
}
