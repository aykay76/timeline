package dto

type WaypointPath struct {
	Waypoints      []Point `json:"waypoints,omitempty"`
	Source         string  `json:"source,omitempty"`
	RoadSegment    []Place `json:"roadSegment,omitempty"`
	DistanceMetres float32 `json:"distanceMeters,omitempty"`
	TravelMode     string  `json:"travelMode,omitempty"`
	Confidence     float32 `json:"confidence,omitempty"`
}
