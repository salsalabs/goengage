package goengage

import "time"

//Constants for Engage endpoints.
const (
	SegmentSearch          = "/api/integration/ext/v1/segments/search"
	SegmentSupporterSearch = "/api/integration/ext/v1/segments/Supporters/search"
	SemgentUpsert          = "/api/integration/ext/v1/supporters"
	SegmentDelete          = "/api/integration/ext/v1/supporters"
)

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

//SegmentUpsertRequest is used to add or modify segments.
type SegmentUpsertRequest struct {
	Payload struct {
		Segments []struct {
			SegmentID        string `json:"segmentId,omitempty"`
			Name             string `json:"name"`
			Description      string `json:"description"`
			ExternalSystemID string `json:"externalSystemId,omitempty"`
		} `json:"segments"`
	} `json:"payload"`
}

//SegmentUpsertResponse returns the results from a SegmentUpsertRequest.
type SegmentUpsertResponse struct {
	Payload SegmentUpsertPayload `json:"payload"`
}

//Segment contains the results of an upsert.
type Segment struct {
	SegmentID        string  `json:"segmentId,omitempty"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ExternalSystemID string  `json:"externalSystemId"`
	Result           string  `json:"result"`
	Errors           []Error `json:"errors,omitempty"`
}

//SegmentUpsertPayload wraps the response for a segment upsert.
type SegmentUpsertPayload struct {
	Segments []Segment `json:"segments"`
}

//SegmentDeleteRequest is used to remove a group.
type SegmentDeleteRequest struct {
	Payload struct {
		Segments []struct {
			SegmentID string `json:"segmentId"`
		} `json:"segments"`
	} `json:"payload"`
}

//SegmentDeleteResponse contains the results from deleting one or more segments.
type SegmentDeleteResponse struct {
	ID        string               `json:"id"`
	Timestamp time.Time            `json:"timestamp"`
	Header    Header               `json:"header"`
	Payload   SegmentDeletePayload `json:"payload"`
}

//SegmentDeleteResult describes the result from a single segment delete.
type SegmentDeleteResult struct {
	SegmentID string `json:"segmentId"`
	Result    string `json:"result"`
}

//SegmentDeletePayload is a wrapper about for details about deleting segments.
type SegmentDeletePayload struct {
	Segments []SegmentDeleteResult `json:"segments"`
	Count    int                   `json:"count"`
}

//SegmentSearchRequest contains parameters for searching for segments.  Please
//see the documentation for details.  Note that true in "includeSupporterCounts"
//really, *really* slows this call down.  A bunch.
type SegmentSearchRequest struct {
	Header struct {
		RefID string `json:"refId"`
	} `json:"header"`
	Payload struct {
		Offset                 int      `json:"offset"`
		Count                  int      `json:"count"`
		Identifiers            []string `json:"identifiers"`
		IdentifierType         string   `json:"identifierType"`
		IncludeSupporterCounts bool     `json:"includeSupporterCounts"`
		JoinedSince            string   `json:"joinedSince"`
	} `json:"payload"`
}

//SegmentSearchResult contains the results of a search for a segment.
//Different from a Segment by the fact that it contains extra fields.
type SegmentSearchResult struct {
	SegmentID        string `json:"segmentId,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Type             string `json:"type,omitempty"`
	TotalSupporters  int    `json:"totalSupporters,omitempty"`
	Result           string `json:"result"`
	ExternalSystemID string `json:"externalSystemId"`
}

//SegmentSearchResponse contains the results returned by searching for segments.
type SegmentSearchResponse struct {
	Payload struct {
		Count    int                   `json:"count"`
		Offset   int                   `json:"offset"`
		Total    int                   `json:"total"`
		Segments []SegmentSearchResult `json:"segments"`
	} `json:"payload"`
}

//AssignSupportersRequest provides the segment and list of supporter IDs
//that need to be added.
type AssignSupportersRequest struct {
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
	Count      int                      `json:"count"`
}

//DeleteSupportersRequest provides the segment and list of supporter IDs
//that need to be added.
type DeleteSupportersRequest struct {
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
	Count      int                      `json:"count"`
}

//SegmentSupporterSearchRequest requests a list of supporters.  Supplying
//"supporterIds" constrains the results to just those supporters.
type SegmentSupporterSearchRequest struct {
	Header  SegmentSupporterSearchHeader  `json:"header"`
	Payload SegmentSupporterSearchPayload `json:"payload"`
}

//SegmentSupporterSearchHeader contains a reference ID provided by the caller.
type SegmentSupporterSearchHeader struct {
	RefID string `json:"refId"`
}

//SegmentSupporterSearchPayload provides the reqest body.
type SegmentSupporterSearchPayload struct {
	SegmentID    string   `json:"segmentId"`
	Offset       int      `json:"offset"`
	Count        int      `json:"count"`
	SupporterIds []string `json:"supporterIds"`
}

//SegmentSupporterSearchResponse contains a list of supporters that match
//the search criteria.
type SegmentSupporterSearchResponse struct {
	ID        string                                `json:"id"`
	Timestamp time.Time                             `json:"timestamp"`
	Header    Header                                `json:"header"`
	Payload   SegmentSupporterSearchResponsePayload `json:"payload"`
}

//SegmentSupporterSearchResponsePayload (whew) carries information about the found
//supporters.  Note that Supporter is common for all of Engage.
type SegmentSupporterSearchResponsePayload struct {
	Total      int         `json:"total"`
	Supporters []Supporter `json:"supporters"`
	Count      int         `json:"count"`
}
