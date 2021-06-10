package goengage

import (
	"fmt"
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
	SearchSupporter       = "/api/integration/ext/v1/supporters/search"
	UpsertSupporter       = "/api/integration/ext/v1/supporters"
	DeleteSupporter       = "/api/integration/ext/v1/supporters"
	SupporterSearchGroups = "/api/integration/ext/v1/supporters/groups"
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

// CustomFieldError is returned when an attempt to change a custom field fails.
type CustomFieldError struct {
	ID        string `json:"id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   string `json:"details,omitempty"`
	FieldName string `json:"fieldName,omitempty"`
}

//CustomFieldValue contains information about a custom field.  Note that
//a supporter/activity will only have custom fields if the values have been
//set in the supporter/activity record.
type CustomFieldValue struct {
	FieldID    string             `json:"fieldId,omitempty" gorm:"field_id,primarykey,omitempty"`
	Name       string             `json:"name"`
	Value      string             `json:"value"`
	OptInDate  *time.Time         `json:"optInDate,omitempty"`
	OptOutDate *time.Time         `json:"optOutDate,omitempty"`
	Errors     []CustomFieldError `json:"errors,omitempty"`
	//Foreign key for GORM.
	SupporterID string `json:"-" gorm:"supporter_id"`
}

//Contact holds a way to communicate with a supporter.  Typical contacts
//include email address and phone numbers.
type Contact struct {
	Type   string `json:"type,omitempty"`
	Value  string `json:"value,omitempty"`
	Status string `json:"status,omitempty" gorm:"status,omitempty"`
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
	Address           *Address           `json:"address,omitempty"`
	Contacts          []Contact          `json:"contacts,omitempty" gorm:"foreignkey:supporter_id"`
	CustomFieldValues []CustomFieldValue `json:"customFieldValues,omitempty" gorm:"foreignkey:supporter_id"`
	ReadOnly          bool               `json:"readOnly,omitempty" gorm:"readOnly,omitempty"`
}

//SupporterSegment is returned when searching for segments that a
//supporter belongs to.
type SupporterSegment struct {
	SupporterID string    `json:"supporterId,omitempty"`
	Segments    []Segment `json:"segments,omitempty"`
	Result      string    `json:"result,omitempty"`
}

//SupporterSearchRequest provides the criteria to match when searching
//for supporters.  Providing no criterria will return all supporters.
//"modifiedTo" and/or "modifiedFrom" are mutually exclusive to searching
//by identifiers.
type SupporterSearchRequest struct {
	Header  RequestHeader                 `json:"header,omitempty"`
	Payload SupporterSearchRequestPayload `json:"payload,omitempty"`
}

//SupporterSearchRequestPayload holds the search criteria.  There are rules
//that you need to know about.  See those here
//https://help.salsalabs.com/hc/en-us/articles/224470107-Engage-API-Supporter-Data#searching-for-supporters
type SupporterSearchRequestPayload struct {
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
	ID        string                         `json:"id"`
	Timestamp *time.Time                     `json:"timestamp"`
	Header    Header                         `json:"header"`
	Payload   SupporterSearchResponsePayload `json:"payload,omitempty"`
}

//SupporterSearchResponsePayload holds the payload for a single supporter search
//operation.
type SupporterSearchResponsePayload struct {
	Count      int32       `json:"count,omitempty"`
	Offset     int32       `json:"offset,omitempty"`
	Total      int32       `json:"total,omitempty"`
	Supporters []Supporter `json:"supporters,omitempty"`
}

//SupporterUpdatePayload holds the list of supporter records to be updated.
type SupporterUpdatePayload struct {
	Supporters []Supporter `json:"supporters,omitempty"`
}

//SupporterUpdateRequest provides a list of modified supporter records that
//the caller wants to be updated in the database.
type SupporterUpdateRequest struct {
	Header  RequestHeader          `json:"header,omitempty"`
	Payload SupporterUpdatePayload `json:"payload,omitempty"`
}

//SupporterUpdateResponse provides results for the updated supporters.
type SupporterUpdateResponse struct {
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

//SupporterUpsert upserts the provided supporter into Engage.
func SupporterUpsert(e *Environment, s *Supporter, logger *UtilLogger) (*Supporter, error) {
	payload := SupporterUpdatePayload{
		Supporters: []Supporter{*s},
	}
	request := SupporterUpdateRequest{
		Header:  RequestHeader{},
		Payload: payload,
	}
	var response SupporterUpdateResponse
	n := NetOp{
		Host:     e.Host,
		Endpoint: UpsertSupporter,
		Method:   UpdateMethod,
		Token:    e.Token,
		Request:  &request,
		Response: &response,
		Logger:   logger,
	}
	err := n.Do()
	if err != nil {
		return s, err
	}
	count := int32(len(response.Payload.Supporters))
	if count != 0 {
		s = &response.Payload.Supporters[0]
		switch s.Result {
		case Added:
			err = nil
		case Updated:
			err = nil
		case ValidationError:
			err = fmt.Errorf("engage returned %s for ID %s", s.Result, s.SupporterID)
		case SystemError:
			err = fmt.Errorf("engage returned %s for ID %s", s.Result, s.SupporterID)
		case NotFound:
			err = fmt.Errorf("engage returned %s for ID %s", s.Result, s.SupporterID)
		}
	} else {
		err = fmt.Errorf("engage return zero responses for ID %s", s.SupporterID)

	}
	return s, err
}

//SupporterGroupRequest requests the groups (segments) that a supporter
//belongs to.
type SupporterGroupRequest struct {
	Header  RequestHeader                `json:"header,omitempty"`
	Payload SupporterGroupRequestPayload `json:"payload,omitempty"`
}

//SupporterGroupRequestPayload holds the search criteria.
//https://api.salsalabs.org/help/integration#operation/getGroupsForSupporters
type SupporterGroupRequestPayload struct {
	Identifiers     []string `json:"identifiers,omitempty"`
	IdentifierType  string   `json:"identifierType,omitempty"`
	SearchString    string   `json:"searchString,omitempty"`
	ModifiedFrom    string   `json:"modifiedFrom,omitempty"`
	ModifiedTo      string   `json:"modifiedTo,omitempty"`
	Offset          int32    `json:"offset,omitempty"`
	Count           int32    `json:"count,omitempty"`
	IncludeCellOnly bool     `json:"includeCellOnly,omitempty"`
	IncludeNormal   bool     `json:"includeNormal,omitempty"`
}

//SupporterGroupResponse provides results for the updated supporters.
type SupporterGroupResponse struct {
	ID        string                        `json:"id,omitempty"`
	Timestamp string                        `json:"timestamp,omitempty"`
	Header    Header                        `json:"header,omitempty"`
	Payload   SupporterGroupResponsePayload `json:"payload,omitempty"`
	Errors    []Error                       `json:"errors,omitempty"`
}

//SupporterGroupResponsePayload lists the supporters that match the search criteria.
//Note that Supporter is common throughout Engage.
type SupporterGroupResponsePayload struct {
	Total   int                `json:"total,omitempty"`
	Offset  int32              `json:"offset,omitempty"`
	Count   int32              `json:"count,omitempty"`
	Results []SupporterSegment `json:"results,omitempty"`
}

//SupporterByID retrieves a supporter record for Engage using the SupporterID
//in the provided record.
func SupporterByID(e *Environment, k string) (*Supporter, error) {
	payload := SupporterSearchRequestPayload{
		Identifiers:    []string{k},
		IdentifierType: SupporterIDType,
		Offset:         int32(0),
		Count:          e.Metrics.MaxBatchSize,
	}
	request := SupporterSearchRequest{
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

// SupporterByEmail returns the first supporter whose email
// matches the provided email.  Duplicates are gleefully ignored.
func SupporterByEmail(e *Environment, email string) (s *Supporter, err error) {
	offset := int32(0)
	payload := SupporterSearchRequestPayload{
		Identifiers:    []string{email},
		IdentifierType: EmailAddressType,
		Offset:         offset,
		Count:          e.Metrics.MaxBatchSize,
	}
	rqt := SupporterSearchRequest{
		Header:  RequestHeader{},
		Payload: payload,
	}
	var resp SupporterSearchResults
	n := NetOp{
		Host:     e.Host,
		Method:   SearchMethod,
		Endpoint: SearchSupporter,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err = n.Do()
	if err != nil {
		return s, err
	}
	count := resp.Payload.Count
	if count != 0 {
		for _, s := range resp.Payload.Supporters {
			// This should always be true, BTW
			x := FirstEmail(s)
			if x != nil && *x == email && s.Result == Found {
				return &s, nil
			}
		}
	}
	err = fmt.Errorf("error: %s is not a valid email", email)
	return s, err
}
