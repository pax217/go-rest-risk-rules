package testdata

import (
	"time"

	"github.com/conekta/risk-rules/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDefaultPayer() entities.ChargebackRequest {
	return entities.ChargebackRequest{
		ChargebackID: "chbk_2rUKtAda8Ljz16Y7J",
		BankReason:   "bank_reason",
		CreatedAt:    time.Now().UTC(),
		Reason:       "lost",
		Status:       "status",
		UpdatedAt:    time.Now().UTC(),
		ChargeID:     "1924045e-ec99-4bb0-aad3-e0a2ab88806a",
		Currency:     "USD",
		Amount:       500,
		CompanyID:    "2",
		Email:        "me@gmail.com",
	}
}

func GetPayerWithDistinctChargebackID() entities.ChargebackRequest {
	return entities.ChargebackRequest{
		ChargebackID: "chbk_2rUKtAda8Ljz16FR4",
		BankReason:   "bank_reason",
		CreatedAt:    time.Now().UTC(),
		Reason:       "lost",
		Status:       "status",
		UpdatedAt:    time.Now().UTC(),
		ChargeID:     "1924045e-ec99-4bb0-aad3-e0a2ab88806a",
		Currency:     "USD",
		Amount:       500,
		CompanyID:    "2",
		Email:        "me@gmail.com",
	}
}

func GetPayerDefeult() entities.Payer {
	return entities.Payer{
		ID:    primitive.NewObjectID(),
		Email: "mail@hotmail.com",
		Chargebacks: []entities.Chargebacks{
			{
				ChargebackID: "chbk_2rUKtAda8Ljz16Y7J",
				ChargeID:     "1924045e-ec99-4bb0-aad3-e0a2ab88806a",
				CompanyID:    "7683457364",
				Status:       "status",
				Reason:       "lost",
				CreatedAt:    time.Now().UTC(),
				Currency:     "USD",
				Amount:       500,
			},
		},
	}
}
