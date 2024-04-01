package dto

type ActivitySegment struct {
	Vao               uint32            `json:"-"`
	StartLocation     Location          `json:"startLocation"`
	EndLocation       Location          `json:"endLocation"`
	Duration          Duration          `json:"duration"`
	Distance          int               `json:"distance"`
	ActivityType      string            `json:"activityType"`
	Confidence        string            `json:"confidence"`
	Activities        []Activity        `json:"activities"`
	WaypointPath      WaypointPath      `json:"waypointPath"`
	SimplifiedRawPath SimplifiedRawPath `json:"simplifiedRawPath"`
	ParkingEvent      ParkingEvent      `json:"parkingEvent"`
}
