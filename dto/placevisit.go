package dto

type PlaceVisit struct {
	Location                Location          `json:"location"`
	Duration                Duration          `json:"duration"`
	PlaceConfidence         string            `json:"placeConfidence"`
	CentreLatE7             int               `json:"centerLatE7"`
	CentreLngE7             int               `json:"centerLngE7"`
	VisitConfidence         int               `json:"visitConfidence"`
	OtherCandidateLocations []Location        `json:"otherCandidateLocations,omitempty"`
	EditConfirmationStatus  string            `json:"editConfirmationStatus"`
	SimplifiedRawPath       SimplifiedRawPath `json:"simplifiedRawPath"`
	LocationConfidence      int               `json:"locationConfidence"`
	PlaceVisitType          string            `json:"placeVisitType"`
	PlaceVisitImportance    string            `json:"placeVisitImportance"`
}
