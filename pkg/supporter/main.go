package goengage

import "time"

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

//Address holds a street address and geolocation stuff for a supporter.
type Address struct {
	AddressLine1 string    `json:"addressLine1,omitempty"`
	AddressLine2 string    `json:"addressLine2,omitempty"`
	City         string    `json:"city,omitempty"`
	State        string    `json:"state,omitempty"`
	PostalCode   string    `json:"postalCode,omitempty"`
	County       string    `json:"county,omitempty"`
	Country      string    `json:"country,omitempty"`
	Lattitude    float64   `json:"lattitude,omitempty"`
	Longitude    float64   `json:"longitude,omitempty"`
	Status       string    `json:"status,omitempty"`
	OptInDate    time.Time `json:"optInDate,omitempty"`
}

//Contact holds a way to communicate with a supporter.  Typical contacts
//include email address and phone numbers.
type Contact struct {
	Type   string `json:"type,omitempty"`
	Value  string `json:"value,omitempty"`
	Status string `json:"status,omitempty,omitempty"`
}

//Supporter describes a single Engage supporter.
type Supporter struct {
	SupporterID       string             `json:"supporterId,omitempty"`
	Result            string             `json:"result,omitempty"`
	Title             string             `json:"title,omitempty"`
	FirstName         string             `json:"firstName,omitempty"`
	MiddleName        string             `json:"middleName,omitempty"`
	LastName          string             `json:"lastName,omitempty"`
	Suffix            string             `json:"suffix,omitempty"`
	DateOfBirth       time.Time          `json:"dateOfBirth,omitempty"`
	Gender            string             `json:"gender,omitempty"`
	CreatedDate       time.Time          `json:"createdDate,omitempty"`
	LastModified      time.Time          `json:"lastModified,omitempty"`
	ExternalSystemID  string             `json:"externalSystemId,omitempty"`
	Address           Address            `json:"address,omitempty"`
	Contacts          []Contact          `json:"contacts,omitempty"`
	CustomFieldValues []CustomFieldValue `json:"customFieldValues,omitempty"`
}

//SearchResultsPayload wraps the supporters found by a
//supporter search request.
type SearchResultsPayload struct {
	Count      int         `json:"count,omitempty"`
	Offset     int         `json:"offset,omitempty"`
	Total      int         `json:"total,omitempty"`
	Supporters []Supporter `json:"supporters,omitempty"`
}

//UpdateRequest provides a list of modified supporter records that
//the caller wants to be updated in the database.
type UpdateRequest struct {
	Payload UpdateRequestPayload `json:"payload,omitempty"`
}

//UpdateRequestPayload carries the list of supporters to be modified.
type UpdateRequestPayload struct {
	Supporters []Supporter `json:"supporters,omitempty"`
}

//UpdateResponse provides results for the updated supporters.
type UpdateResponse struct {
	Payload UpdateResponsePayload `json:"payload,omitempty"`
}

//UpdateResponsePayload contains the results of modifying supporters.
type UpdateResponsePayload struct {
	Supporters []Supporter `json:"supporters,omitempty"`
}

//DeleteRequest is used to delete supporter records.  By the way,
//deleted records are gone forever -- they are not coming back, Jim.
type DeleteRequest struct {
	Payload DeleteRequestPayload `json:"payload,omitempty"`
}

//DeleteRequestPayload contains the list of supporters to be deleted.
type DeleteRequestPayload struct {
	Supporters []Supporter `json:"supporters,omitempty"`
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
