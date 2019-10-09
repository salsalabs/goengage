package goengage

import "time"

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

//Header returns server-side information for Engage API calls.
type Header struct {
	ProcessingTime int    `json:"processingTime"`
	ServerID       string `json:"serverId"`
}

//CustomFieldValues contains information about a custom field.  Note that
//a supporter will only have custom fields if the values have been
//set for the supporter.
type CustomFieldValue struct {
	FieldID    string    `json:"fieldId"`
	Name       string    `json:"name"`
	Value      string    `json:"value"`
	Type       string    `json:"type"`
	OptInDate  time.Time `json:"optInDate,omitempty"`
	OptOutDate time.Time `json:"optOutDate,omitempty"`
}
