package goengage

import (
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Constants for Engage endpoints.
const (
	Search          = "/api/integration/ext/v1/segments/search"
	SupporterSearch = "/api/integration/ext/v1/segments/Supporters/search"
	SemgentUpsert   = "/api/integration/ext/v1/supporters"
	Delete          = "/api/integration/ext/v1/supporters"
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
	Header  goengage.RequestHeader `json:"header,omitempty"`
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
	Header  goengage.Header `json:"header,omitempty"`
	Payload UpsertPayload   `json:"payload"`
}

//Segment contains the results of an upsert.
type Segment struct {
	ID               string           `json:"segmentId,omitempty"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	ExternalSystemID string           `json:"externalSystemId"`
	Result           string           `json:"result"`
	Errors           []goengage.Error `json:"errors,omitempty"`
}

//UpsertPayload wraps the response for a segment upsert.
type UpsertPayload struct {
	Segments []Segment `json:"segments"`
}

//DeleteRequest is used to remove a group.
type DeleteRequest struct {
	Header  goengage.RequestHeader `json:"header,omitempty"`
	Payload struct {
		Segments []struct {
			SegmentID string `json:"segmentId"`
		} `json:"segments"`
	} `json:"payload"`
}

//DeleteResponse contains the results from deleting one or more segments.
type DeleteResponse struct {
	ID        string          `json:"id"`
	Timestamp *time.Time       `json:"timestamp"`
	Header    goengage.Header `json:"header"`
	Payload   DeletePayload   `json:"payload"`
}

//DeleteResult describes the result from a single segment delete.
type DeleteResult struct {
	SegmentID string `json:"segmentId"`
	Result    string `json:"result"`
}

//DeletePayload is a wrapper about for details about deleting segments.
type DeletePayload struct {
	Segments []DeleteResult `json:"segments"`
	Count    int32          `json:"count"`
}

//SearchRequest contains parameters for searching for segments.  Please
//see the documentation for details.  Note that true in "includeSupporterCounts"
//really, *really* slows this call down.  A bunch.
type SearchRequest struct {
	Header  goengage.RequestHeader `json:"header,omitempty"`
	Payload SearchRequestPayload   `json:"payload,omitempty"`
}

//SearchRequestPayload contains the request criteria.
type SearchRequestPayload struct {
	Offset                 int32    `json:"offset"`
	Count                  int32    `json:"count"`
	Identifiers            []string `json:"identifiers"`
	IdentifierType         string   `json:"identifierType"`
	IncludeSupporterCounts bool     `json:"includeSupporterCounts"`
	JoinedSince            string   `json:"joinedSince"`
}

//SearchResult contains the results of a search for a segment.
//Different from a Segment by the fact that it contains extra fields.
type SearchResult struct {
	SegmentID        string `json:"segmentId,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Type             string `json:"type,omitempty"`
	TotalSupporters  int    `json:"totalSupporters,omitempty"`
	Result           string `json:"result"`
	ExternalSystemID string `json:"externalSystemId"`
}

//SearchResponse contains the results returned by searching for segments.
type SearchResponse struct {
	Payload struct {
		Count    int32          `json:"count"`
		Offset   int32          `json:"offset"`
		Total    int32          `json:"total"`
		Segments []SearchResult `json:"segments"`
	} `json:"payload"`
}

//AssignSupportersRequest provides the segment and list of supporter IDs
//that need to be added.
type AssignSupportersRequest struct {
	Header  goengage.RequestHeader  `json:"header,omitempty"`
	Payload AssignSupportersPayload `json:"payload"`
}

//AssignSupportersPayload carries the request details.
type AssignSupportersPayload struct {
	SegmentID    string   `json:"segmentId"`
	SupporterIds []string `json:"supporterIds"`
}

//AssignSupportersResponse carries the results of adding supporters to a segment.
type AssignSupportersResponse struct {
	Payload AssignSupportersResultPayload `json:"payload"`
}

//AssignSupportersResult contains the results of adding supporters to a segment.
type AssignSupportersResult struct {
	SupporterID string `json:"supporterId"`
	Result      string `json:"result"`
}

//AssignSupportersResultPayload (argh) wraps the supporters from
//an assigment request.
type AssignSupportersResultPayload struct {
	Supporters []AssignSupportersResult `json:"supporters"`
	Count      int32                    `json:"count"`
}

//DeleteSupportersRequest provides the segment and list of supporter IDs
//that need to be added.
type DeleteSupportersRequest struct {
	Header  goengage.RequestHeader  `json:"header,omitempty"`
	Payload DeleteSupportersPayload `json:"payload"`
}

//DeleteSupportersPayload carries the request details.
type DeleteSupportersPayload struct {
	SegmentID    string   `json:"segmentId"`
	SupporterIds []string `json:"supporterIds"`
}

//DeleteSupportersResponse carries the results of adding supporters to a segment.
type DeleteSupportersResponse struct {
	Payload DeleteSupportersResultPayload `json:"payload"`
}

//DeleteSupportersResult contains the results of adding supporters to a segment.
type DeleteSupportersResult struct {
	SupporterID string `json:"supporterId"`
	Result      string `json:"result"`
}

//DeleteSupportersResultPayload (argh) wraps the supporters from
//an assigment request.
type DeleteSupportersResultPayload struct {
	Supporters []DeleteSupportersResult `json:"supporters"`
	Count      int32                    `json:"count"`
}

//SupporterSearchRequest requests a list of supporters.  Supplying
//"supporterIds" constrains the results to just those supporters.
type SupporterSearchRequest struct {
	Header  goengage.RequestHeader `json:"header,omitempty"`
	Payload SupporterSearchPayload `json:"payload"`
}

//SupporterSearchPayload provides the reqest body.
type SupporterSearchPayload struct {
	SegmentID    string   `json:"segmentId"`
	Offset       int32    `json:"offset"`
	Count        int32    `json:"count"`
	SupporterIds []string `json:"supporterIds"`
}

//SupporterSearchResponse contains a list of supporters that match
//the search criteria.
type SupporterSearchResponse struct {
	ID        string                         `json:"id"`
	Timestamp *time.Time                      `json:"timestamp"`
	Header    goengage.Header                `json:"header"`
	Payload   SupporterSearchResponsePayload `json:"payload"`
}

//SupporterSearchResponsePayload (whew) carries information about the found
//supporters.  Note that Supporter is common for all of Engage.
type SupporterSearchResponsePayload struct {
	Total      int32                `json:"total"`
	Supporters []goengage.Supporter `json:"supporters"`
	Count      int32                `json:"count"`
}
