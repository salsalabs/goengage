package goengage

import (
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//Engage endpoints for supporters.
const (
	Search = "/api/integration/ext/v1/supporters/search"
	Upsert = "/api/integration/ext/v1/supporters"
	Delete = "/api/integration/ext/v1/supporters"
)

//Contact types.
const (
	ContactEmail    = "EMAIL"
	ContactHome     = "HOME_PHONE"
	ContactCell     = "CELL_PHONE"
	ContactWork     = "WORK_PHONE"
	ContactFacebook = "FACEBOOK_ID"
	ContactTwitter  = "TWITTER_ID"
	ContactLinkedin = "LINKEDIN_ID"
)

//SearchRequest provides the criteria to match when searching
//for supporters.  Providing no criterria will return all supporters.
//"modifiedTo" and/or "modifiedFrom" are mutually exclusive to searching
//by identifiers.
type SearchRequest struct {
	Payload SearchRequestPayload `json:"payload,omitempty"`
}

//SearchRequestPayload holds the search criteria.  There are rules
//that you need to know about.  See those here
//https://help.salsalabs.com/hc/en-us/articles/224470107-Engage-API-Supporter-Data#searching-for-supporters
type SearchRequestPayload struct {
	Identifiers    []string  `json:"identifiers,omitempty"`
	IdentifierType string    `json:"identifierType,omitempty"`
	ModifiedFrom   time.Time `json:"modifiedFrom,omitempty"`
	ModifiedTo     time.Time `json:"modifiedTo,omitempty"`
	Offset         int       `json:"offset,omitempty"`
	Count          int       `json:"count,omitempty"`
}

//SearchResults lists the supporters that match the search criteria.
//Note that Supporter is common throughout Engage.
type SearchResults struct {
	Payload SearchResultsPayload `json:"payload,omitempty"`
}

//SearchResultsPayload wraps the supporters found by a
//supporter search request.
type SearchResultsPayload struct {
	Count      int                  `json:"count,omitempty"`
	Offset     int                  `json:"offset,omitempty"`
	Total      int                  `json:"total,omitempty"`
	Supporters []goengage.Supporter `json:"supporters,omitempty"`
}

//UpdateRequest provides a list of modified supporter records that
//the caller wants to be updated in the database.
type UpdateRequest struct {
	Payload UpdateRequestPayload `json:"payload,omitempty"`
}

//UpdateRequestPayload carries the list of supporters to be modified.
type UpdateRequestPayload struct {
	Supporters []goengage.Supporter `json:"supporters,omitempty"`
}

//UpdateResponse provides results for the updated supporters.
type UpdateResponse struct {
	Payload UpdateResponsePayload `json:"payload,omitempty"`
}

//UpdateResponsePayload contains the results of modifying supporters.
type UpdateResponsePayload struct {
	Supporters []goengage.Supporter `json:"supporters,omitempty"`
}

//DeleteRequest is used to delete supporter records.  By the way,
//deleted records are gone forever -- they are not coming back, Jim.
type DeleteRequest struct {
	Payload DeleteRequestPayload `json:"payload,omitempty"`
}

//DeleteRequestPayload contains the list of supporters to be deleted.
type DeleteRequestPayload struct {
	Supporters []goengage.Supporter `json:"supporters,omitempty"`
}

//DeletedResponse returns the results of deleting supporters.
type DeletedResponse struct {
	Payload DeletedResponsePayload `json:"payload,omitempty"`
}

//DeleteResult describes the results of deleting a single supporter.
type DeleteResult struct {
	SupporterID string `json:"supporterId,omitempty"`
	Result      string `json:"result,omitempty"`
}

//DeletedResponsePayload contains the delete results.
type DeletedResponsePayload struct {
	Supporters []DeleteResult `json:"supporters,omitempty"`
}