package goengage

import (
	"sync"
)

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

//Census describes a segment and the supproters that are members.  This is
//an aggregate structure used by SegmentCensus.
type Census struct {
	Segment
	Supporters []Supporter
}

//AllSegments retrieves all segments (groups) from an Engage instance.  Each
//segment is pushed onto the provided channel.  The channel is closed when the
//last segment is pushed.
//
//The boolean argument is true (CountYes) if the segment records should contain
//the number of supporters in the group.  Note tht counting supporters is very
//expensive in terms of clock time.  *Very* expensive.  Use CountNo unless you
//must have the number of supporters.
func AllSegments(e *Environment, c bool, x chan Segment) error {
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
	for rqt.Count > 0 {
		err := n.Do()
		if err != nil {
			return err
		}
		for _, s := range resp.Segments {
			x <- s
		}
		count := len(resp.Segments)
		rqt.Count = int32(count)
		rqt.Offset = rqt.Offset + int32(count)
	}
	close(x)
	return nil
}

//SegmentCensus returns the supporters in a group.
func SegmentCensus(e *Environment, s Segment) ([]Supporter, error) {
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

//AllSegmentCensus creates a Census object for every segment, then pushes
//it onto the census channel.  The census channel is closed when all segments
//have been procssed.
//
//Internally, there are goroutines that retrieves segments and pushes census records.
//
//Panicking on errors until we find an elegant way to handle them in goroutines.
//
//Note!  This function currently discriminates against the DEFAULT segments. TODO: Add filter function.
func AllSegmentCensus(e *Environment, d chan Census) {
	c := make(chan Segment)
	var wg sync.WaitGroup

	//Receiver creates a Census object and pushes it to the output channel.
	//Census channel is closed after the last segment is processed.
	//
	go (func(c chan Segment, d chan Census, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		for true {
			s, ok := <-c
			if !ok {
				break
			}
			if s.Type != SegmentTypeDefault {
				supporters, err := SegmentCensus(e, s)
				if err != nil {
					panic(err)
				}
				census := Census{
					Segment:    s,
					Supporters: supporters,
				}
				d <- census
			}
		}
		close(d)
	})(c, d, &wg)

	//Sender sends all segments to the Census maker.
	//
	//Panicking on errors until we find an elegant way to handle them
	//in goroutines.
	go (func(c chan Segment, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := AllSegments(e, CountNo, c)
		if err != nil {
			panic(err)
		}
	})(c, &wg)
	wg.Wait()
}
