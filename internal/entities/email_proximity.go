package entities

import "time"

type EmailEvaluationResponse struct {
	IsValid bool                 `json:"is_valid" mapstructure:"is_valid" bson:"is_valid"`
	Stats   EmailEvaluationStats `json:"stats" mapstructure:"stats" bson:"stats"`
}

type EmailEvaluationStats struct {
	ObservationsCount int        `json:"observations_count" mapstructure:"observations_count" bson:"observations_count"`
	VariationsCount   int        `json:"variations_count" mapstructure:"variations_count" bson:"variations_count"`
	LastSeen          *time.Time `json:"last_seen" mapstructure:"last_seen" bson:"last_seen"`
}
