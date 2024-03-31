package dto

type Activity struct {
	ActivityType string  `json:"activityType,omitempty"`
	Probability  float64 `json:"probability"`
}
