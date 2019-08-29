package goengage

import "fmt"

//SegSearch is used to search for segments.
const SegSearch = "/api/integration/ext/v1/segments/search"

//SegSupporterSearch is used to search segments for supporters.
const SegSupporterSearch = "/api/integration/ext/v1/segments/members/search"

//Constants to drive counting, or not counting, supporters on a segment read.
//Counting is expensive, sometimes prohibitively so.
const (
	CountNo  = false
	CountYes = true
)

//Segment types.
const (
	SegmentTypeDefault = "DEFAULT"
	SegmentTypeCustom  = "CUSTOM"
)

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
	Count    int32     `json:"count,omitempty"`
	Offset   int32     `json:"offset,omitempty"`
	Total    int32     `json:"total,omitempty"`
	Segments []Segment `json:"segments,omitempty"`
}

//SegSupporterSearchRequest is used to find supporters in a segment.
//Be sure to pass the correct IdentifierType.
type SegSupporterSearchRequest struct {
	SegmentID    string   `json:"segmentId"`
	SupporterIDs []string `json:"supporterIds"`
	Offset       int32    `json:"offset,omitempty"`
	Count        int32    `json:"count,omitempty"`
}

//SegSupporterSearchResult is returned when supporters are found by a search.
type SegSupporterSearchResult struct {
	Count      int32       `json:"count,omitempty"`
	Offset     int32       `json:"offset,omitempty"`
	Total      int32       `json:"total,omitempty"`
	Supporters []Supporter `json:"supporters,omitempty"`
}

//Census describes a segment and the supproters that are members.
type Census struct {
	Segment
	Supporters []Supporter
}

//AllSegments returns all groups.
func AllSegments(e *Environment, c bool) ([]Segment, error) {
	rqt := SegSearchRequest{
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
		MemberCounts: c,
	}
	var resp SegSearchResult
	n := NetOp{
		Host:     e.Host,
		Endpoint: SegSearch,
		Method:   SearchMethod,
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
		for _, s := range resp.Segments {
			a = append(a, s)
		}
		count := len(resp.Segments)
		rqt.Count = int32(count)
		rqt.Offset = rqt.Offset + int32(count)
	}
	return a, nil
}

//SegmentCensus returns the supporters in a group.
func SegmentCensus(e *Environment, s Segment) ([]Supporter, error) {
	fmt.Printf("SegmentCensus: retrieving %d members for %v\n", s.TotalMembers, s.Name)
	var a []Supporter
	rqt := SegSupporterSearchRequest{
		SegmentID: s.SegmentID,
		Offset:    0,
		Count:     e.Metrics.MaxBatchSize,
	}

	var resp SegSupporterSearchResult
	n := NetOp{
		Host:     e.Host,
		Endpoint: SegSupporterSearch,
		Method:   SearchMethod,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}

	for rqt.Count > 0 {
		err := n.Do()
		if err != nil {
			return a, err
		}
		for _, s := range resp.Supporters {
			a = append(a, s)
		}
		count := len(resp.Supporters)
		rqt.Count = int32(count)
		rqt.Offset = rqt.Offset + int32(count)
	}
	return a, nil
}

//AllSegmentCensus returns a Census object for all segments.  The Census
//object describes a segment and all of its supporters.
func AllSegmentCensus(e *Environment) ([]Census, error) {
	var a []Census
	segments, err := AllSegments(e, CountNo)
	if err != nil {
		return a, err
	}
	fmt.Printf("AllSegmentCensus found %d segments\n", len(segments))
	for _, s := range segments {
		if s.Type != SegmentTypeDefault {
			supporters, err := SegmentCensus(e, s)
			if err != nil {
				return a, err
			}
			spop := Census{
				Segment:    s,
				Supporters: supporters,
			}
			a = append(a, spop)
		}
	}
	return a, nil
}
