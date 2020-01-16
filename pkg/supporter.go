package goengage

import (
	"time"
)

//Identifier types for supporter requests
const (
	SupporterIDType  = "SUPPORTER_ID"
	EmailAddressType = "EMAIL_ADDRESS"
	ExternalIDType   = "EXTERNAL_ID"
)

//Engage endpoints for supporters.
const (
	SearchSupporter = "/api/integration/ext/v1/supporters/search"
	UpsertSupporter = "/api/integration/ext/v1/supporters"
	DeleteSupporter = "/api/integration/ext/v1/supporters"
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

//Status types.
const (
	OptIn  = "OPT_IN"
	OptOut = "OPT_OUT"
)

//Address holds a street address and geolocation stuff for a supporter.
type Address struct {
	AddressLine1         string     `json:"addressLine1,omitempty"`
	AddressLine2         string     `json:"addressLine2,omitempty"`
	AddressLine3         string     `json:"addressLine3,omitempty"`
	City                 string     `json:"city,omitempty"`
	State                string     `json:"state,omitempty"`
	PostalCode           string     `json:"postalCode,omitempty"`
	County               string     `json:"county,omitempty"`
	Country              string     `json:"country,omitempty"`
	FederalHouseDistrict string     `json:"federalHouseDistrict,omitempty"`
	StateHouseDistrict   string     `json:"stateHouseDistrict,omitempty"`
	StateSenateDistrict  string     `json:"stateSenateDistrict,omitempty"`
	CountyDistrict       string     `json:"countyDistrict,omitempty"`
	MunicipalityDistrict string     `json:"municipalityDistrict,omitempty"`
	Lattitude            float64    `json:"lattitude,omitempty"`
	Longitude            float64    `json:"longitude,omitempty"`
	Status               string     `json:"status,omitempty"`
	OptInDate            *time.Time `json:"optInDate,omitempty"`
}

//CustomFieldValue contains information about a custom field.  Note that
//a supporter/activity will only have custom fields if the values have been
//set in the supporter/activity record.
type CustomFieldValue struct {
	FieldID    string     `json:"fieldId" gorm:"field_id,primarykey"`
	Name       string     `json:"name"`
	Value      string     `json:"value"`
	OptInDate  *time.Time `json:"optInDate,omitempty"`
	OptOutDate *time.Time `json:"optOutDate,omitempty"`
	//Foreign key for GORM.
	SupporterID string `json:"-" gorm:"supporter_id"`
}

//Contact holds a way to communicate with a supporter.  Typical contacts
//include email address and phone numbers.
type Contact struct {
	Type   string `json:"type,omitempty"`
	Value  string `json:"value,omitempty"`
	Status string `json:"status,omitempty,omitempty" gorm:"status,omitempty"`
	//Foreign key for GORM.``
	SupporterID string `json:"-" gorm:"supporter_id"`
	ContactID   string `json:"-" gorm:"contact_id,primarykey,autoincrement"`
}

//Supporter describes a single Engage supporter.
type Supporter struct {
	SupporterID       string             `json:"supporterId,omitempty" gorm:"primary_key"`
	Result            string             `json:"result,omitempty"`
	Title             string             `json:"title,omitempty"`
	FirstName         string             `json:"firstName,omitempty"`
	MiddleName        string             `json:"middleName,omitempty"`
	LastName          string             `json:"lastName,omitempty"`
	Suffix            string             `json:"suffix,omitempty"`
	DateOfBirth       *time.Time         `json:"dateOfBirth,omitempty"`
	Gender            string             `json:"gender,omitempty"`
	CreatedDate       *time.Time         `json:"createdDate,omitempty"`
	LastModified      *time.Time         `json:"lastModified,omitempty"`
	ExternalSystemID  string             `json:"externalSystemId,omitempty"`
	Address           Address            `json:"address,omitempty"`
	Contacts          []Contact          `json:"contacts,omitempty" gorm:"foreignkey:supporter_id"`
	CustomFieldValues []CustomFieldValue `json:"customFieldValues,omitempty" gorm:"foreignkey:supporter_id"`
}

//SupporterSearch provides the criteria to match when searching
//for supporters.  Providing no criterria will return all supporters.
//"modifiedTo" and/or "modifiedFrom" are mutually exclusive to searching
//by identifiers.
type SupporterSearch struct {
	Header  RequestHeader          `json:"header,omitempty"`
	Payload SupporterSearchPayload `json:"payload,omitempty"`
}

//SupporterSearchPayload holds the search criteria.  There are rules
//that you need to know about.  See those here
//https://help.salsalabs.com/hc/en-us/articles/224470107-Engage-API-Supporter-Data#searching-for-supporters
type SupporterSearchPayload struct {
	Identifiers    []string `json:"identifiers,omitempty"`
	IdentifierType string   `json:"identifierType,omitempty"`
	ModifiedFrom   string   `json:"modifiedFrom,omitempty"`
	ModifiedTo     string   `json:"modifiedTo,omitempty"`
	Offset         int32    `json:"offset,omitempty"`
	Count          int32    `json:"count,omitempty"`
}

//SupporterSearchResults lists the supporters that match the search criteria.
//Note that Supporter is common throughout Engage.
type SupporterSearchResults struct {
	ID        string     `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	Header    Header     `json:"header"`
	Payload   struct {
		Count      int32       `json:"count,omitempty"`
		Offset     int32       `json:"offset,omitempty"`
		Total      int32       `json:"total,omitempty"`
		Supporters []Supporter `json:"supporters,omitempty"`
	} `json:"payload,omitempty"`
}

//UpdateRequest provides a list of modified supporter records that
//the caller wants to be updated in the database.
type UpdateRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		Supporters []Supporter `json:"supporters,omitempty"`
	} `json:"payload,omitempty"`
}

//UpdateResponse provides results for the updated supporters.
type UpdateResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Supporters []Supporter `json:"supporters,omitempty"`
	} `json:"payload,omitempty"`
}

//DeleteRequest is used to delete supporter records.  By the way,
//deleted records are gone forever -- they are not coming back, Jim.
type DeleteRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		Supporters []Supporter `json:"supporters,omitempty"`
	} `json:"payload,omitempty"`
}

//DeletedResponse returns the results of deleting supporters.
type DeletedResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Supporters []struct {
			SupporterID string `json:"supporterId,omitempty"`
			Result      string `json:"result,omitempty"`
		} `json:"supporters,omitempty"`
	} `json:"payload,omitempty"`
}

//FetchSupporter retrieves a supporter record for Engage using the SupporterID
//in the provided record.
func FetchSupporter(e *Environment, k string) (*Supporter, error) {
	payload := SupporterSearchPayload{
		Identifiers:    []string{k},
		IdentifierType: SupporterIDType,
		Offset:         int32(0),
		Count:          e.Metrics.MaxBatchSize,
	}
	request := SupporterSearch{
		Header:  RequestHeader{},
		Payload: payload,
	}
	var response SupporterSearchResults
	n := NetOp{
		Host:     e.Host,
		Endpoint: SearchSupporter,
		Method:   SearchMethod,
		Token:    e.Token,
		Request:  &request,
		Response: &response,
	}
	err := n.Do()
	if err != nil {
		return nil, err
	}
	count := int32(len(response.Payload.Supporters))
	if count == 0 {
		return nil, nil
	}
	for _, s := range response.Payload.Supporters {
		// This should always be true, BTW
		if s.SupporterID == k {
			if s.Result == Found {
				return &s, nil
			}
		}
	}
	return nil, nil
}
