package goengage

//Segment is a named group of supporters.
type Segment struct {
	SegmentID        string
	Name             string
	Description      string
	Type             string
	TotalMembers     int32
	Result           string
	ExternalSystemID string
}

//SegSearchRequest is used to ask for supporters.
type SegSearchRequest struct {
	Identifiers    []string `json:"identifiers"`
	IdentifierType string   `json:"identifierType"`
	Offset         int32    `json:"offset"`
	Count          int32    `json:"count"`
	MemberCounts   bool     `json:"includeMemberCounts"`
}

//SegSearchResult is returned when supporters are found by a search.
type SegSearchResult struct {
	Payload struct {
		Count    int32     `json:"count"`
		Offset   int32     `json:"offset"`
		Total    int32     `json:"total"`
		Segments []Segment `json:"segments"`
	} `json:"payload"`
}
