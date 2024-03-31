package dto

type Location struct {
	LatitudeE7            int        `json:"latitudeE7,omitempty"`
	LongitudeE7           int        `json:"longitudeE7,omitempty"`
	SourceInfo            SourceInfo `json:"sourceInfo,omitempty"`
	PlaceID               string     `json:"placeId,omitempty"`
	Address               string     `json:"address,omitempty"`
	Name                  string     `json:"name,omitempty"`
	LocationConfidence    int        `json:"placeConfidence,omitempty"`
	CalibratedProbability float64    `json:"calibratedProbability,omitempty"`
	AccuracyMetres        int        `json:"accuracyMeters,omitempty"`
}
