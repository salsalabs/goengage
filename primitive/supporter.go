package goengage

import "time"

//SupporterSearchRequest provides the criteria to match when searching
//for supporters.  Providing no criterria will return all supporters.
//"modifiedTo" and/or "modifiedFrom" are mutually exclusive to searching
//by identifiers.
type SupporterSearchRequest struct {
	Payload SupporterSearchRequestPayload `json:"payload"`
}

//SupporterSearchRequestPayload holds the search criteria.  There are rules
//that you need to know about.  See those here
//https://help.salsalabs.com/hc/en-us/articles/224470107-Engage-API-Supporter-Data#searching-for-supporters
type SupporterSearchRequestPayload struct {
	Identifiers    []string  `json:"identifiers"`
	IdentifierType string    `json:"identifierType"`
	ModifiedFrom   time.Time `json:"modifiedFrom"`
	ModifiedTo     time.Time `json:"modifiedTo"`
	Offset         int       `json:"offset"`
	Count          int       `json:"count"`
}

//SupporterSearchResults lists the supporters that match the search criteria.
//Note that Supporter is common throughout Engage.
type SupporterSearchResults struct {
	Payload SupporterSearchResultsPayload `json:"payload"`
}

//Address holds a street address and geolocation stuff for a supporter.
type Address struct {
	AddressLine1 string    `json:"addressLine1"`
	AddressLine2 string    `json:"addressLine2"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	PostalCode   string    `json:"postalCode"`
	County       string    `json:"county"`
	Country      string    `json:"country"`
	Lattitude    float64   `json:"lattitude"`
	Longitude    float64   `json:"longitude"`
	Status       string    `json:"status"`
	OptInDate    time.Time `json:"optInDate"`
}

//Contact holds a way to communicate with a supporter.  Typical contacts
//include email address and phone numbers.
type Contact struct {
	Type   string `json:"type"`
	Value  string `json:"value"`
	Status string `json:"status,omitempty"`
}

//Supporter describes a single Engage supporter.
type Supporter struct {
	SupporterID       string             `json:"supporterId"`
	Result            string             `json:"result"`
	Title             string             `json:"title"`
	FirstName         string             `json:"firstName"`
	MiddleName        string             `json:"middleName"`
	LastName          string             `json:"lastName"`
	Suffix            string             `json:"suffix"`
	DateOfBirth       time.Time          `json:"dateOfBirth"`
	Gender            string             `json:"gender"`
	CreatedDate       time.Time          `json:"createdDate"`
	LastModified      time.Time          `json:"lastModified"`
	ExternalSystemID  string             `json:"externalSystemId"`
	Address           Address            `json:"address"`
	Contacts          []Contact          `json:"contacts"`
	CustomFieldValues []CustomFieldValue `json:"customFieldValues"`
}

//SupporterSearchResultsPayload wraps the supporters found by a
//supporter search request.
type SupporterSearchResultsPayload struct {
	Count      int         `json:"count"`
	Offset     int         `json:"offset"`
	Total      int         `json:"total"`
	Supporters []Supporter `json:"supporters"`
}

//SupporterUpdateRequest provides a list of modified supporter records that
//the caller wants to be updated in the database.
type SupporterUpdateRequest struct {
	Payload SupporterUpdateRequestPayload `json:"payload"`
}

//SupporterUpdateRequestPayload carries the list of supporters to be modified.
type SupporterUpdateRequestPayload struct {
	Supporters []Supporter `json:"supporters"`
}

//SupporterUpdateResponse provides results for the updated supporters.
type SupporterUpdateResponse struct {
	Payload SupporterUpdateResponsePayload `json:"payload"`
}

//SupporterUpdateResponsePayload contains the results of modifying supporters.
type SupporterUpdateResponsePayload struct {
	Supporters []Supporter `json:"supporters"`
}

//SupporterDeleteRequest is used to delete supporter records.  By the way,
//deleted records are gone forever -- they are not coming back, Jim.
type SupporterDeleteRequest struct {
	Payload SupporterDeleteRequestPayload `json:"payload"`
}

//Supporters is the way to define a single supporter ID.
type Supporters struct {
	SupporterID string `json:"supporterId"`
}

//SupporterDeleteRequestPayload contains the list of supporters to be deleted.
type SupporterDeleteRequestPayload struct {
	Supporters []Supporter `json:"supporters"`
}

//SupporterDeletedResponse returns the results of deleting supporters.
type SupporterDeletedResponse struct {
	Payload SupporterDeletedResponsePayload `json:"payload"`
}

//SupporterDeleteResult describes the results of deleting a single supporter.
type SupporterDeleteResult struct {
	SupporterID string `json:"supporterId"`
	Result      string `json:"result"`
}

//SupporterDeletedResponsePayload contains the delete results.
type SupporterDeletedResponsePayload struct {
	Supporters []SupporterDeleteResult `json:"supporters"`
}
