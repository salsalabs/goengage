package goengage

import (
	"time"
)

const (
	//EmailBlastSearchRequest is used to find blasts.
	SearchEmailBlast = "/api/integration/ext/v1/emails/search"

	//EmailBlastGetSingleBlast is used to retrieve a single blast
	//as well as as all of the recipient data.
	EmailBlastGetSingleBlast = "/api/integration/ext/v1/emails/individualResults"
)

//Conversion hold information about any donatoins made as a result of an
//email blast.
type Conversion struct {
	ConversionDate *time.Time `json:"conversionDate"`
	ActivityType   string     `json:"activityType"`
	ActivityName   string     `json:"activityName"`
	ActivityID     string     `json:"activityId"`
	Amount         string     `json:"amount"`
	DonationType   string     `json:"donationType"`
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
	TimeSent             *time.Time   `json:"timeSent"`
	SplitName            string       `json:"splitName"`
	Status               string       `json:"status"`
	Opened               bool         `json:"opened"`
	Clicked              bool         `json:"clicked"`
	Converted            bool         `json:"converted"`
	Unsubscribed         bool         `json:"unsubscribed"`
	FirstOpenDate        *time.Time   `json:"firstOpenDate,omitempty"`
	NumberOfLinksClicked string       `json:"numberOfLinksClicked"`
	Conversions          []Conversion `json:"conversionData,omitempty"`
}

//EmailBlastSearchRequestPayload contains the criteria used to retrieve blasts.
type EmailBlastSearchRequestPayload struct {
	ID            string `json:"id,omitempty"`
	ContentID     string `json:"contentId,omitempty"`
	Cursor        string `json:"cursor,omitempty"`
	PublishedFrom string `json:"publishedFrom,omitempty"`
	PublishedTo   string `json:"publishedTo,omitempty"`
	Type          string `json:"type,omitempty"`
	Offset        int32  `json:"offset,omitempty"`
	Count         int32  `json:"count,omitempty"`
}

//EmailBlastSearchRequest wraps the request payload.
type EmailBlastSearchRequest struct {
	ID      string                         `json:"id,omitempty"`
	Header  RequestHeader                  `json:"header,omitempty"`
	Payload EmailBlastSearchRequestPayload `json:"payload,omitempty"`
}

//EmailBlastSearchResponsePayload contains the results of a
//search.
type EmailBlastSearchResponsePayload struct {
	Total           int32        `json:"total,omitempty"`
	Offset          int32        `json:"offset,omitempty"`
	Count           int32        `json:"count,omitempty"`
	EmailActivities []EmailBlast `json:"emailActivities,omitempty"`
}

//EmailBlastSearchResponse wraps a response payload.
type EmailBlastSearchResponse struct {
	ID        string                          `json:"id,omitempty"`
	TimeStamp string                          `json:"timestamp,omitempty"`
	Header    Header                          `json:"header,omitempty"`
	Payload   EmailBlastSearchResponsePayload `json:"payload,omitempty"`
}

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
	ID        string     `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	Header    Header     `json:"header"`
	Payload   struct {
		Total           int32        `json:"total"`
		Offset          int32        `json:"offset"`
		EmailActivities []EmailBlast `json:"emailActivities"`
		Count           int32        `json:"count"`
	} `json:"payload"`
}

//EmailActivity describes the contents of the email.
type EmailBlast struct {
	ID          string     `json:"id"`
	Topic       string     `json:"topic"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	PublishDate *time.Time `json:"publishDate"`
	Components  []struct {
		ContentID     string `json:"contentId"`
		MessageNumber string `json:"messageNumber"`
	} `json:"components"`
}

//SeriesResponse response is returned when the request type is "CommSeries".
type SeriesResponse struct {
	ID        string     `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	Header    Header     `json:"header"`
	Payload   struct {
		IndividualEmailActivityData []struct {
			ID             string `json:"id"`
			Cursor         string `json:"cursor"`
			Name           string `json:"name"`
			RecipientsData struct {
				Recipients []Recipient `json:"recipients"`
				Total      int32       `json:"total"`
			} `json:"recipientsData"`
		} `json:"individualEmailActivityData"`
	} `json:"payload"`
}
