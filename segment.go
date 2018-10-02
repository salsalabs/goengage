package goengage

//SegSearch is used to search for segments.
const SegSearch = "/api/integration/ext/v1/segments/search"

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
	Identifiers    []string `json:"identifiers,omitempty"`
	IdentifierType string   `json:"identifierType,omitempty"`
	Offset         int32    `json:"offset,omitempty"`
	Count          int32    `json:"count,omitempty"`
	MemberCounts   bool     `json:"includeMemberCounts,omitempty"`
}

//SegSearchResult is returned when supporters are found by a search.
type SegSearchResult struct {
	Payload struct {
		Count    int32     `json:"count,omitempty"`
		Offset   int32     `json:"offset,omitempty"`
		Total    int32     `json:"total,omitempty"`
		Segments []Segment `json:"segments,omitempty"`
	} `json:"payload,omitempty"`
}
