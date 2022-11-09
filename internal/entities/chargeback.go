package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payer struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	Chargebacks []Chargebacks      `json:"chargebacks" bson:"chargebacks"`
}

type Chargebacks struct {
	ChargebackID string    `json:"chargeback_id" bson:"chargeback_id"`
	ChargeID     string    `json:"charge_id" bson:"charge_id"`
	CompanyID    string    `json:"company_id" bson:"company_id"`
	Status       string    `json:"status" bson:"status"`
	Reason       string    `json:"reason" bson:"reason"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	Currency     string    `json:"currency" bson:"currency"`
	Amount       float64   `json:"amount" bson:"amount"`
}

type ChargebackRequest struct {
	ChargebackID string    `json:"_id"`
	BankReason   string    `json:"bank_reason"`
	CreatedAt    time.Time `json:"created_at"`
	Reason       string    `json:"reason"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
	ChargeID     string    `json:"charge_id"`
	Currency     string    `json:"currency"`
	Amount       float64   `json:"amount"`
	CompanyID    string    `json:"company_id"`
	Email        string    `json:"email"`
}

func (c ChargebackRequest) NewPayerFromPostRequest() Payer {
	nowStr := time.Now().UTC().Truncate(time.Millisecond)
	return Payer{
		ID:    primitive.NewObjectID(),
		Email: c.Email,
		Chargebacks: []Chargebacks{{
			ChargebackID: c.ChargebackID,
			ChargeID:     c.ChargeID,
			CompanyID:    c.CompanyID,
			Status:       c.Status,
			Reason:       c.Reason,
			CreatedAt:    nowStr,
			Currency:     c.Currency,
			Amount:       c.Amount,
		}},
	}
}

func (p *Payer) ExistChargeback(chargebacks []Chargebacks) bool {
	for _, chargeback := range chargebacks {
		if chargeback.ChargebackID == p.Chargebacks[0].ChargebackID {
			return true
		}
	}

	return false
}
