package dto

type Records struct {
	Locations []RecordLocation `json:"locations,omitempty"`
}

type RecordLocation struct {
	LatitudeE7  int `json:"latitudeE7,omitempty"`
	LongitudeE7 int `json:"longitudeE7,omitempty"`
}
