package dto

type ParkingEvent struct {
	Location       Location `json:"location,omitempty"`
	Method         string   `json:"method,omitempty"`
	LocationSource string   `json:"locationSource,omitempty"`
	Timestamp      string   `json:"timestamp,omitempty"`
}
