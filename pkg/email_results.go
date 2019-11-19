package goengage

import (
	"time"
)

//EmailResultsRequest is used to request email blast activity
//records for a blast.
//
// See https://help.salsalabs.com/hc/en-us/articles/360019505914-Engage-API-Email-Results
type EmailResultsRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		Cursor    string `json:"cursor,omitempty"`
		Type      string `json:"type,omitempty"`
		ID        string `json:"id,omitempty"`
		ContentID string `json:"contentId,omitempty"`
	} `json:"payload,omitempty"`
}

//EmailResponse is returned when the request type is "Email".
type EmailResponse struct {
	ID        string               `json:"id"`
	Timestamp *time.Time            `json:"timestamp"`
	Header    Header               `json:"header"`
	Payload   EmailResponsePayload `json:"payload"`
}

//Component is something...
type Component struct {
	ContentID     string `json:"contentId"`
	MessageNumber string `json:"messageNumber"`
}

//EmailActivity describes the contents of the email.
type EmailActivity struct {
	ID          string      `json:"id"`
	Topic       string      `json:"topic"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	PublishDate *time.Time   `json:"publishDate"`
	Components  []Component `json:"components"`
}

//EmailResponsePayload holds the content for an email blast search.
type EmailResponsePayload struct {
	Total           int32           `json:"total"`
	Offset          int32           `json:"offset"`
	EmailActivities []EmailActivity `json:"emailActivities"`
	Count           int32           `json:"count"`
}

//SeriesResponse response is returned when the request type is "CommSeries".
type SeriesResponse struct {
	ID        string                `json:"id"`
	Timestamp *time.Time             `json:"timestamp"`
	Header    Header                `json:"header"`
	Payload   SeriesResponsePayload `json:"payload"`
}

//Conversion hold information about any donatoins made as a result of an
//email blast.
type Conversion struct {
	ConversionDate *time.Time `json:"conversionDate"`
	ActivityType   string    `json:"activityType"`
	ActivityName   string    `json:"activityName"`
	ActivityID     string    `json:"activityId"`
	Amount         string    `json:"amount"`
	DonationType   string    `json:"donationType"`
}

//Recipient contains information about an email blast message sent to a
//supporter.
type Recipient struct {
	SupporterID          string       `json:"supporterId"`
	SupporterEmail       string       `json:"supporterEmail"`
	FirstName            string       `json:"firstName"`
	LastName             string       `json:"lastName"`
	Country              string       `json:"country"`
	State                string       `json:"state"`
	City                 string       `json:"city"`
	TimeSent             *time.Time    `json:"timeSent"`
	SplitName            string       `json:"splitName"`
	Status               string       `json:"status"`
	Opened               bool         `json:"opened"`
	Clicked              bool         `json:"clicked"`
	Converted            bool         `json:"converted"`
	Unsubscribed         bool         `json:"unsubscribed"`
	FirstOpenDate        *time.Time    `json:"firstOpenDate,omitempty"`
	NumberOfLinksClicked string       `json:"numberOfLinksClicked"`
	Conversions          []Conversion `json:"conversionData,omitempty"`
}

//RecipientsData contains the list of recpients.
type RecipientsData struct {
	Recipients []Recipient `json:"recipients"`
	Total      int32       `json:"total"`
}

//IndividualEmailActivityData contains information about something.  Not
//sure, exactly...
type IndividualEmailActivityData struct {
	ID             string         `json:"id"`
	Cursor         string         `json:"cursor"`
	Name           string         `json:"name"`
	RecipientsData RecipientsData `json:"recipientsData"`
}

//SeriesResponsePayload holds the data returned by a series search.
type SeriesResponsePayload struct {
	IndividualEmailActivityData []IndividualEmailActivityData `json:"individualEmailActivityData"`
}
