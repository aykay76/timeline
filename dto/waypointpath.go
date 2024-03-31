package dto

type WaypointPath struct {
	Waypoints      []Point `json:"waypoints,omitempty"`
	Source         string  `json:"source,omitempty"`
	RoadSegment    []Place `json:"roadSegment,omitempty"`
	DistanceMetres float64 `json:"distanceMeters,omitempty"`
	TravelMode     string  `json:"travelMode,omitempty"`
	Confidence     float64 `json:"confidence,omitempty"`
}
