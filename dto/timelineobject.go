package dto

type TimelineObject struct {
	ActivitySegment ActivitySegment `json:"activitySegment,omitempty"`
	PlaceVisit      PlaceVisit      `json:"placeVisit,omitempty"`
}
