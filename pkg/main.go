package goengage

import (
	"net/http"
	"time"
)

const (
	//UATHost is the hostname for Engage instances on the test server.
	UATHost = "hq.uat.igniteaction.net"
	//APIHost is the hostname for Engage instances on the production server.
	APIHost = "api.salsalabs.org"
	//ContentType is always Javascript.
	ContentType = "application/json"
	//SearchMethod is always "POST" in Engage.
	SearchMethod = http.MethodPost
)

//Segment constants
const (
	//Added indicates that the provided segment was added to the system
	Added = "ADDED"
	//Updated indicates that the provided segment was updated
	Updated = "UPDATED"
	//NotAllowed indicates that the segment represented by the provided id
	//is not allowed to be modified via the API.
	NotAllowed = "NOT_ALLOWED"
)

//Merge supporter records esult value constants.
const (
	//Found will be reported for the destination supporter if no updates were
	//specified to be performed.
	Found = "FOUND"
	//Update will be reported for the destination supporter if updates were
	//specified. It will also be reported on the main payload if the merge
	//operation was successful.
	Update = "UPDATE"
	//NotFound will be reported for the destination or source supporter if the
	//provided id(s) do not exist.
	NotFound = "NOT_FOUND"
	//Deleted will be reported for the source supporter on a successful merge.
	Deleted = "DELETED"
	//ValidationError will be reported on the main payload if either the source
	//or the destination supporter is not found, or a request to update the
	//destination was specified and validation errors occurred during that
	//update.
	ValidationError = "VALIDATION_ERROR"
	//SystemError if the merge could not be completed.
	SystemError = "SYSTEM_ERROR"
)

// Types for searching for email results.
const (
	//Email is used for searching for blasts.
	//Email = "Email"
	//CommSeries is used for searching email series.
	CommSeries = "CommSeries"
)

//Contact types.
const (
	Email     = "EMAIL"
	HomePhone = "HOME_PHONE"
	CellPhone = "CELL_PHONE"
	WorkPhone = "WORK_PHONE"
	Facebook  = "FACEBOOK_ID"
	Twitter   = "TWITTER_ID"
	Linkedin  = "LINKEDIN_ID"
)

//Error is used to report Engage errors.
type Error struct {
	ID        string `json:"id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   string `json:"details,omitempty"`
	FieldName string `json:"fieldName,omitempty"`
}

//RequestHeader provides a reference ID.
type RequestHeader struct {
	RefID string `json:"refId,omitEmpty"`
}

//Header returns server-side information for Engage API calls.
type Header struct {
	ProcessingTime int    `json:"processingTime"`
	ServerID       string `json:"serverId"`
}

//CustomFieldValue contains information about a custom field.  Note that
//a supporter/activity will only have custom fields if the values have been
//set in the supporter/activity record.
type CustomFieldValue struct {
	FieldID    string     `json:"fieldId" gorm:"primary_key"`
	Name       string     `json:"name"`
	Value      string     `json:"value"`
	Type       string     `json:"type"`
	OptInDate  *time.Time `json:"optInDate,omitempty" gorm:"optInDate"`
	OptOutDate *time.Time `json:"optOutDate,omitempty" gorm:"optOutDate"`
}

//Address holds a street address and geolocation stuff for a supporter.
type Address struct {
	AddressLine1 string     `json:"addressLine1,omitempty" gorm:"addressLine1"`
	AddressLine2 string     `json:"addressLine2,omitempty" gorm:"addressLine2"`
	City         string     `json:"city,omitempty" gorm:"city"`
	State        string     `json:"state,omitempty" gorm:"state"`
	PostalCode   string     `json:"postalCode,omitempty" gorm:"postalCode"`
	County       string     `json:"county,omitempty" gorm:"county"`
	Country      string     `json:"country,omitempty" gorm:"country"`
	Lattitude    float64    `json:"lattitude,omitempty" gorm:"lattitude"`
	Longitude    float64    `json:"longitude,omitempty" gorm:"longitude"`
	Status       string     `json:"status,omitempty" gorm:"status"`
	OptInDate    *time.Time `json:"optInDate,omitempty" gorm:"optInDate"`
}

//Contact holds a way to communicate with a supporter.  Typical contacts
//include email address and phone numbers.
type Contact struct {
	Type   string `json:"type,omitempty" gorm:"type"`
	Value  string `json:"value,omitempty" gorm:"value"`
	Status string `json:"status,omitempty,omitempty" gorm:"status,omitempty"`
}

//Supporter describes a single Engage supporter.
type Supporter struct {
	SupporterID       string             `json:"supporterId,omitempty" gorm:"supporterId"`
	Result            string             `json:"result,omitempty" gorm:"result"`
	Title             string             `json:"title,omitempty" gorm:"title"`
	FirstName         string             `json:"firstName,omitempty" gorm:"firstName"`
	MiddleName        string             `json:"middleName,omitempty" gorm:"middleName"`
	LastName          string             `json:"lastName,omitempty" gorm:"lastName"`
	Suffix            string             `json:"suffix,omitempty" gorm:"suffix"`
	DateOfBirth       *time.Time         `json:"dateOfBirth,omitempty" gorm:"dateOfBirth"`
	Gender            string             `json:"gender,omitempty" gorm:"gender"`
	CreatedDate       *time.Time         `json:"createdDate,omitempty" gorm:"createdDate"`
	LastModified      *time.Time         `json:"lastModified,omitempty" gorm:"lastModified"`
	ExternalSystemID  string             `json:"externalSystemId,omitempty" gorm:"externalSystemId"`
	Address           Address            `json:"address,omitempty" gorm:"address"`
	Contacts          []Contact          `json:"contacts,omitempty" gorm:"-"`
	CustomFieldValues []CustomFieldValue `json:"customFieldValues,omitempty" gorm:"-"`
}
