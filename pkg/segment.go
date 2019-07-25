package goengage

import (
	"net/http"
)

//SegSearch is used to search for segments.
const SegSearch = "/api/integration/ext/v1/segments/search"

//SegSupporterSearch is used to search segments for supporters.
const SegSupporterSearch = "/api/integration/ext/v1/segments/members/search"

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

//SegSupporterSearchRequest is used to find supporters in a segment.
//Be sure to pass the correct IdentifierType.
//SegSearchRequest is used to ask for supporters.
type SegSupporterSearchRequest struct {
	SegmentID    string   `json:"segmentId"`
	SupporterIDs []string `json:"supporterIds"`
	Offset       int32    `json:"offset,omitempty"`
	Count        int32    `json:"count,omitempty"`
}

//SegSupporterSearchResult is returned when supporters are found by a search.
type SegSupporterSearchResult struct {
	Payload struct {
		Count      int32       `json:"count,omitempty"`
		Offset     int32       `json:"offset,omitempty"`
		Total      int32       `json:"total,omitempty"`
		Supporters []Supporter `json:"supporters,omitempty"`
	} `json:"payload,omitempty"`
}

//AllSegments returns all groups.
func AllSegments(e *EngEnv, m *MetricData, c bool) ([]Segment, error) {
	rqt := SegSearchRequest{
		Offset:       0,
		Count:        m.MaxBatchSize,
		MemberCounts: c,
	}
	var resp SegSearchResult
	n := NetOp{
		Host:     e.Host,
		Fragment: SegSearch,
		Method:   http.MethodPost,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	var a []Segment
	for rqt.Count > 0 {
		err := n.Do()
		if err != nil {
			panic(err)
		}
		for _, s := range resp.Payload.Segments {
			a = append(a, s)
		}
		count := len(resp.Payload.Segments)
		rqt.Count = int32(count)
		rqt.Offset = rqt.Offset + int32(count)
	}
	return a, nil
}
