package goengage

import (
	"time"
)

//Constants for Engage endpoints.
const (
	SearchSegment        = "/api/integration/ext/v1/segments/search"
	SegmentSearchMembers = "/api/integration/ext/v1/segments/members/search"
	UpsertSegment        = "/api/integration/ext/v1/segments"
	DeleteSegment        = "/api/integration/ext/v1/segments"
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

//Segment contains the information for a single Segment (group)>
type Segment struct {
	SegmentID        string  `json:"segmentId,omitempty"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Type             string  `json:"type,omitempty"`
	TotalSupporters  int     `json:"totalSupporters,omitempty"`
	ExternalSystemID string  `json:"externalSystemId,omitempty"`
	Result           string  `json:"result,omitempty"`
	Errors           []Error `json:"errors,omitempty"`
}

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
	Header  RequestHeader               `json:"header,omitempty"`
	Payload SegmentSearchRequestPayload `json:"payload,omitempty"`
}

//SegmentSearchRequestPayload contains the payload for searching for segments.
type SegmentSearchRequestPayload struct {
	Offset                 int32    `json:"offset"`
	Count                  int32    `json:"count"`
	Identifiers            []string `json:"identifiers"`
	IdentifierType         string   `json:"identifierType"`
	IncludeSupporterCounts bool     `json:"includeSupporterCounts"`
	JoinedSince            string   `json:"joinedSince"`
}

//SegmentSearchResponse contains the results returned by searching for segments.
type SegmentSearchResponse struct {
	Payload struct {
		Count    int32     `json:"count"`
		Offset   int32     `json:"offset"`
		Total    int32     `json:"total"`
		Segments []Segment `json:"segments"`
	} `json:"payload"`
}

//SegmentMembershipRequest contains parameters for searching for segment
//(group) members.
type SegmentMembershipRequest struct {
	Header  RequestHeader                   `json:"header,omitempty"`
	Payload SegmentMembershipRequestPayload `json:"payload,omitempty"`
}

//SegmentMembershipRequestPayload contains the payload for searching for
// segment members.
type SegmentMembershipRequestPayload struct {
	SegmentId    string         `json:"segmentId,omitempty"`
	SupporterIds *[]interface{} `json:"supporterIds,omitempty"`
	JoinedSince  string         `json:"joinedSince,omitempty"`
	Offset       int32          `json:"offset,omitempty"`
	Count        int32          `json:"count,omitempty"`
	SortOrder    string         `json:"sortOrder,omitempty"`
}

//SegmentMembershipResponse contains the results returned by searching for
//segment members.
type SegmentMembershipResponse struct {
	Payload struct {
		Count int32 `json:"count"`
		// Potential bug.  These fields do not appear in the response payload.
		// Offset     int32       `json:"offset"`
		// Total      int32       `json:"total"`
		Supporters []Supporter `json:"supporters"`
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
