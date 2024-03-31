package dto

import "time"

type Point struct {
	LatE7          int       `json:"latE7,omitempty"`
	LngE7          int       `json:"lngE7,omitempty"`
	AccuracyMetres int       `json:"accuracyMeters,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}
