package goengage

import (
	"time"
)

//Constants for Engage endpoints.
const (
	SearchSegment          = "/api/integration/ext/v1/segments/search"
	SupporterSearchSegment = "/api/integration/ext/v1/segments/Supporters/search"
	UpsertSegment          = "/api/integration/ext/v1/supporters"
	DeleteSegment          = "/api/integration/ext/v1/supporters"
)

//Constants to drive counting, or not counting, supporters on a segment read.
//Counting is expensive, sometimes prohibitively so.
const (
	CountNo  = false
	CountYes = true
)

//Segment types.
const (
	TypeDefault = "DEFAULT"
	TypeCustom  = "CUSTOM"
)

//UpsertRequest is used to add or modify segments.
type UpsertRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		Segments []struct {
			ID               string `json:"segmentId,omitempty"`
			Name             string `json:"name"`
			Description      string `json:"description"`
			ExternalSystemID string `json:"externalSystemId,omitempty"`
		} `json:"segments"`
	} `json:"payload"`
}

//UpsertResponse returns the results from a UpsertRequest.
type UpsertResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Segments []Segment `json:"segments"`
	} `json:"payload"`
}

//Segment contains the results of an upsert.
type Segment struct {
	ID               string  `json:"segmentId,omitempty"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ExternalSystemID string  `json:"externalSystemId"`
	Result           string  `json:"result"`
	Errors           []Error `json:"errors,omitempty"`
}

//SegmentDeleteRequest is used to remove a group.
type SegmentDeleteRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		Segments []struct {
			SegmentID string `json:"segmentId"`
		} `json:"segments"`
	} `json:"payload"`
}

//SegmentDeleteResponse contains the results from deleting one or more segments.
type SegmentDeleteResponse struct {
	ID        string     `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	Header    Header     `json:"header"`
	Payload   struct {
		Segments []struct {
			SegmentID string `json:"segmentId"`
			Result    string `json:"result"`
		} `json:"segments"`
		Count int32 `json:"count"`
	} `json:"payload"`
}

//SegmentSearchRequest contains parameters for searching for segments.  Please
//see the documentation for details.  Note that true in "includeSupporterCounts"
//really, *really* slows this call down.  A bunch.
type SegmentSearchRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		Offset                 int32    `json:"offset"`
		Count                  int32    `json:"count"`
		Identifiers            []string `json:"identifiers"`
		IdentifierType         string   `json:"identifierType"`
		IncludeSupporterCounts bool     `json:"includeSupporterCounts"`
		JoinedSince            string   `json:"joinedSince"`
	} `json:"payload,omitempty"`
}

//SegmentSearchResponse contains the results returned by searching for segments.
type SegmentSearchResponse struct {
	Payload struct {
		Count    int32 `json:"count"`
		Offset   int32 `json:"offset"`
		Total    int32 `json:"total"`
		Segments []struct {
			SegmentID        string `json:"segmentId,omitempty"`
			Name             string `json:"name,omitempty"`
			Description      string `json:"description,omitempty"`
			Type             string `json:"type,omitempty"`
			TotalSupporters  int    `json:"totalSupporters,omitempty"`
			Result           string `json:"result"`
			ExternalSystemID string `json:"externalSystemId"`
		} `json:"segments"`
	} `json:"payload"`
}

//AssignSupportersRequest provides the segment and list of supporter IDs
//that need to be added.
type AssignSupportersRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		SegmentID    string   `json:"segmentId"`
		SupporterIds []string `json:"supporterIds"`
	} `json:"payload"`
}

//AssignSupportersResponse carries the results of adding supporters to a segment.
type AssignSupportersResponse struct {
	Payload struct {
		Supporters []struct {
			SupporterID string `json:"supporterId"`
			Result      string `json:"result"`
		} `json:"supporters"`
		Count int32 `json:"count"`
	} `json:"payload"`
}

//DeleteSupportersRequest provides the segment and list of supporter IDs
//that need to be added.
type DeleteSupportersRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		SegmentID    string   `json:"segmentId"`
		SupporterIds []string `json:"supporterIds"`
	} `json:"payload"`
}

//DeleteSupportersResponse carries the results of adding supporters to a segment.
type DeleteSupportersResponse struct {
	Payload struct {
		Supporters []struct {
			SupporterID string `json:"supporterId"`
			Result      string `json:"result"`
		} `json:"supporters"`
		Count int32 `json:"count"`
	} `json:"payload"`
}

//SupporterSearchRequest requests a list of supporters.  Supplying
//"supporterIds" constrains the results to just those supporters.
type SupporterSearchRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		SegmentID    string   `json:"segmentId"`
		Offset       int32    `json:"offset"`
		Count        int32    `json:"count"`
		SupporterIds []string `json:"supporterIds"`
	} `json:"payload"`
}

//SupporterSearchResponse contains a list of supporters that match
//the search criteria.
type SupporterSearchResponse struct {
	ID        string     `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	Header    Header     `json:"header"`
	Payload   struct {
		Total      int32       `json:"total"`
		Supporters []Supporter `json:"supporters"`
		Count      int32       `json:"count"`
	} `json:"payload"`
}
