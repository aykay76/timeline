package dto

import "time"

type Duration struct {
	StartTimestamp time.Time `json:"startTimestamp,omitempty"`
	EndTimestamp   time.Time `json:"endTimestamp,omitempty"`
}
