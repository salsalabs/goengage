package goengage

import (
	"time"
)

// Constants for Engage endpoints.
const (
	SearchSegment        = "/api/integration/ext/v1/segments/search"
	SegmentSearchMembers = "/api/integration/ext/v1/segments/members/search"
	UpsertSegment        = "/api/integration/ext/v1/segments"
	DeleteSegment        = "/api/integration/ext/v1/segments"
)

// Constants to drive counting, or not counting, supporters on a segment read.
// Counting is expensive, sometimes prohibitively so.
const (
	CountNo  = false
	CountYes = true
)

// Segment types.
const (
	TypeDefault = "DEFAULT"
	TypeCustom  = "CUSTOM"
)

// Segment contains the information for a single Segment (group)>
type Segment struct {
	SegmentID         string `json:"segmentId,omitempty"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Type              string `json:"type,omitempty"`
	TotalMembers      int    `json:"totalMembers,omitempty"`
	ExternalSystemID  string `json:"externalSystemId,omitempty"`
	Result            string `json:"result,omitempty"`
	MailingList       bool   `json:"mailingList,omitempty"`
	PublicName        string `json:"publicName,omitempty"`
	PublicDescription string `json:"publicDescription,omitempty"`
	ParameterName     string `json:"parameterName,omitempty"`
}

// SegmentError describes the errors that can return when searching
// for segments.
type SegmentError struct {
	Error
	ContentType string `json:"contentType,omitempty"`
	ContentID   string `json:"contentId,omitempty"`
}

// UpsertRequest is used to add or modify segments.
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

// UpsertResponse returns the results from a UpsertRequest.
type UpsertResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Segments []Segment `json:"segments"`
	} `json:"payload"`
}

// SegmentDeleteRequest is used to remove a group.
type SegmentDeleteRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		Segments []struct {
			SegmentID string `json:"segmentId"`
		} `json:"segments"`
	} `json:"payload"`
}

// SegmentDeleteResponse contains the results from deleting one or more segments.
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

// SegmentSearchRequest contains parameters for searching for segments.  Please
// see the documentation for details.  Note that true in "includeMemberCounts"
// really, *really* slows this call down.  A bunch.
type SegmentSearchRequest struct {
	Header  RequestHeader               `json:"header,omitempty"`
	Payload SegmentSearchRequestPayload `json:"payload,omitempty"`
}

// SegmentSearchRequestPayload contains the payload for searching for segments.
type SegmentSearchRequestPayload struct {
	Offset              int32    `json:"offset,omitempty"`
	Count               int32    `json:"count,omitempty"`
	Identifiers         []string `json:"identifiers,omitempty"`
	IdentifierType      string   `json:"identifierType,omitempty"`
	IncludeMemberCounts bool     `json:"includeMemberCounts,omitempty"`
}

// SegmentSearchResponse contains the results returned by searching for segments.
type SegmentSearchResponse struct {
	ID        string                       `json:"id,omitempty"`
	Timestamp *time.Time                   `json:"timestamp,omitempty"`
	Header    Header                       `json:"header,omitempty"`
	Payload   SegmentSearchResponsePayload `json:"payload,omitempty"`
	Errors    []SegmentError               `json:"errors,omitempty"`
}

// SegmentWrapper is a segment with errors and warnings.
type SegmentWrapper struct {
	Errors   []Error `json:"errors,omitempty"`
	Warnings []Error `json:"warnings,omitempty"`
	Segment
}

// SegmentSearchResponsePayload wraps the response payload for a
// segment search.
type SegmentSearchResponsePayload struct {
	Count    int32            `json:"count,omitempty"`
	Offset   int32            `json:"offset,omitempty"`
	Total    int32            `json:"total,omitempty"`
	Segments []SegmentWrapper `json:"segments,omitempty"`
}

// SegmentMembershipRequest contains parameters for searching for segment
// (group) members.
type SegmentMembershipRequest struct {
	Header  RequestHeader                   `json:"header,omitempty"`
	Payload SegmentMembershipRequestPayload `json:"payload,omitempty"`
}

// SegmentMembershipRequestPayload contains the payload for searching for
// segment members.
type SegmentMembershipRequestPayload struct {
	SegmentID    string         `json:"segmentId,omitempty"`
	SupporterIds *[]interface{} `json:"supporterIds,omitempty"`
	JoinedSince  string         `json:"joinedSince,omitempty"`
	Offset       int32          `json:"offset,omitempty"`
	Count        int32          `json:"count,omitempty"`
	SortOrder    string         `json:"sortOrder,omitempty"`
}

// SegmentMembershipResponse contains the results returned by searching for
// segment members.
type SegmentMembershipResponse struct {
	Header  Header                           `json:"header,omitempty"`
	Payload SegmentMembershipResponsePayload `json:"payload,omitempty"`
}

// SegmentMembershipResponsePayload carries a batch of supporters for
// the provided segment.
type SegmentMembershipResponsePayload struct {
	Count      int32       `json:"count,omitempty"`
	Total      int32       `json:"total,omitempty"`
	Supporters []Supporter `json:"supporters,omitempty"`
}

// AssignSupportersRequest provides the segment and list of supporter IDs
// that need to be added.
type AssignSupportersRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		SegmentID    string   `json:"segmentId"`
		SupporterIds []string `json:"supporterIds"`
	} `json:"payload"`
}

// AssignSupportersResponse carries the results of adding supporters to a segment.
type AssignSupportersResponse struct {
	Payload struct {
		Supporters []struct {
			SupporterID string `json:"supporterId"`
			Result      string `json:"result"`
		} `json:"supporters"`
		Count int32 `json:"count"`
	} `json:"payload"`
}

// DeleteSupportersRequest provides the segment and list of supporter IDs
// that need to be added.
type DeleteSupportersRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		SegmentID    string   `json:"segmentId"`
		SupporterIds []string `json:"supporterIds"`
	} `json:"payload"`
}

// DeleteSupportersResponse carries the results of adding supporters to a segment.
type DeleteSupportersResponse struct {
	Payload struct {
		Supporters []struct {
			SupporterID string `json:"supporterId"`
			Result      string `json:"result"`
		} `json:"supporters"`
		Count int32 `json:"count"`
	} `json:"payload"`
}
