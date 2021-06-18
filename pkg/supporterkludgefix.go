package goengage

import (
	"fmt"
	"time"
)

//Address holds a street address and geolocation stuff for a supporter.
//The kludge fix is to requre certain fields in the JSON by removing
//'omitempty'. That allows us to overwrite those fields with emptiness.
type AddressKludgeFix struct {
	AddressLine1         string     `json:"addressLine1"`
	AddressLine2         string     `json:"addressLine2"`
	City                 string     `json:"city"`
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

//NewAddressKludgeFix accepts and address and returned a kludged address.
func NewAddressKludgeFix(a *Address) *AddressKludgeFix {
	adf := AddressKludgeFix{
		AddressLine1:         a.AddressLine1,
		AddressLine2:         a.AddressLine2,
		City:                 a.City,
		State:                a.State,
		PostalCode:           a.PostalCode,
		County:               a.County,
		Country:              a.Country,
		FederalHouseDistrict: a.FederalHouseDistrict,
		StateHouseDistrict:   a.StateHouseDistrict,
		StateSenateDistrict:  a.StateSenateDistrict,
		CountyDistrict:       a.CountyDistrict,
		MunicipalityDistrict: a.MunicipalityDistrict,
		Lattitude:            a.Lattitude,
		Longitude:            a.Longitude,
		Status:               a.Status,
		OptInDate:            a.OptInDate,
	}
	return &adf
}

//Supporter describes a single Engage supporter.  The
//KludgeFix is implemented in AddressKludgeFix.
type SupporterKludgeFix struct {
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
	Address           *AddressKludgeFix  `json:"address,omitempty"`
	Contacts          []Contact          `json:"contacts,omitempty" gorm:"foreignkey:supporter_id"`
	CustomFieldValues []CustomFieldValue `json:"customFieldValues,omitempty" gorm:"foreignkey:supporter_id"`
	ReadOnly          bool               `json:"readOnly,omitempty" gorm:"readOnly,omitempty"`
}

//Accept a supporter, return n kludged supporter.
func NewSupporterKludgeFix(s Supporter) SupporterKludgeFix {
	skf := SupporterKludgeFix{
		SupporterID:       s.SupporterID,
		Result:            s.Result,
		Title:             s.Title,
		FirstName:         s.FirstName,
		MiddleName:        s.MiddleName,
		LastName:          s.LastName,
		Suffix:            s.Suffix,
		DateOfBirth:       s.DateOfBirth,
		Gender:            s.Gender,
		CreatedDate:       s.CreatedDate,
		LastModified:      s.LastModified,
		ExternalSystemID:  s.ExternalSystemID,
		Address:           NewAddressKludgeFix(s.Address),
		Contacts:          s.Contacts,
		CustomFieldValues: s.CustomFieldValues,
		ReadOnly:          s.ReadOnly,
	}
	return skf
}

//SupporterKludgeFixSearchRequestPayload holds the search criteria.  There are rules
//that you need to know about.  See those here
//https://help.salsalabs.com/hc/en-us/articles/224470107-Engage-API-Supporter-Data#searching-for-supporters
type SupporterKludgeFixSearchRequestPayload struct {
	Identifiers    []string `json:"identifiers,omitempty"`
	IdentifierType string   `json:"identifierType,omitempty"`
	ModifiedFrom   string   `json:"modifiedFrom,omitempty"`
	ModifiedTo     string   `json:"modifiedTo,omitempty"`
	Offset         int32    `json:"offset,omitempty"`
	Count          int32    `json:"count,omitempty"`
}

//SupporterKludgeFixSearchResults lists the supporters that match the search criteria.
//Note that SupporterKludgeFix is common throughout Engage.
type SupporterKludgeFixSearchResults struct {
	ID        string                                  `json:"id"`
	Timestamp *time.Time                              `json:"timestamp"`
	Header    Header                                  `json:"header"`
	Payload   SupporterKludgeFixSearchResponsePayload `json:"payload,omitempty"`
}

//SupporterKludgeFixSearchResponsePayload holds the payload for a single supporter search
//operation.
type SupporterKludgeFixSearchResponsePayload struct {
	Count               int32                `json:"count,omitempty"`
	Offset              int32                `json:"offset,omitempty"`
	Total               int32                `json:"total,omitempty"`
	SupporterKludgeFixs []SupporterKludgeFix `json:"supporters,omitempty"`
}

//SupporterKludgeFixUpdatePayload holds the list of supporter records to be updated.
type SupporterKludgeFixUpdatePayload struct {
	SupporterKludgeFixs []SupporterKludgeFix `json:"supporters,omitempty"`
}

//SupporterKludgeFixUpdateRequest provides a list of modified supporter records that
//the caller wants to be updated in the database.
type SupporterKludgeFixUpdateRequest struct {
	Header  RequestHeader                   `json:"header,omitempty"`
	Payload SupporterKludgeFixUpdatePayload `json:"payload,omitempty"`
}

//SupporterKludgeFixUpdateResponse provides results for the updated supporters.
type SupporterKludgeFixUpdateResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		SupporterKludgeFixs []SupporterKludgeFix `json:"supporters,omitempty"`
	} `json:"payload,omitempty"`
}

//SupporterKludgeFixUpsert upserts the provided supporter into Engage.
func SupporterKludgeFixUpsert(e *Environment, s *SupporterKludgeFix, logger *UtilLogger) (*SupporterKludgeFix, error) {
	payload := SupporterKludgeFixUpdatePayload{
		SupporterKludgeFixs: []SupporterKludgeFix{*s},
	}
	request := SupporterKludgeFixUpdateRequest{
		Header:  RequestHeader{},
		Payload: payload,
	}
	var response SupporterKludgeFixUpdateResponse
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
	count := int32(len(response.Payload.SupporterKludgeFixs))
	if count != 0 {
		s = &response.Payload.SupporterKludgeFixs[0]
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
