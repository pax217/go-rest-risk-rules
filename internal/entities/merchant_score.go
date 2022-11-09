package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type MerchantScore struct {
	ID        primitive.ObjectID `json:"id" csv:"id" bson:"_id,omitempty"`
	CompanyID string             `json:"company_id" csv:"company_id" bson:"company_id"`
	Score     float64            `json:"score" csv:"score_value" bson:"score"`
}
