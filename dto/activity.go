package dto

type Activity struct {
	ActivityType string  `json:"activityType,omitempty"`
	Probability  float32 `json:"probability"`
}
