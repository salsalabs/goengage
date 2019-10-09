package goengage

import "time"

const (
	//Email is used for searching for blasts.
	Email = "Email"
	//CommSeries is used for searching email series.
	CommSeries = "CommSeries"
)

//EmailResultsRequest is used to request email blast activity
//records for a blast.
//
// See https://help.salsalabs.com/hc/en-us/articles/360019505914-Engage-API-Email-Results
type EmailResultsRequest struct {
	Payload struct {
		Cursor    string `json:"cursor,omitempty,omitempty"`
		Type      string `json:"type,omitempty"`
		ID        string `json:"id,omitempty"`
		ContentID string `json:"contentId,omitempty,omitempty"`
	} `json:"payload,omitempty,omitempty"`
}

//EmailResponse is returned when the request type is "Email".
type EmailResponse struct {
	ID        string    `json:"id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Header    struct {
		ProcessingTime int    `json:"processingTime,omitempty"`
		ServerID       string `json:"serverId,omitempty"`
	} `json:"header,omitempty"`
	Payload struct {
		Total           int `json:"total,omitempty"`
		Offset          int `json:"offset,omitempty"`
		EmailActivities []struct {
			ID          string    `json:"id,omitempty"`
			Topic       string    `json:"topic,omitempty"`
			Name        string    `json:"name,omitempty"`
			Description string    `json:"description,omitempty"`
			PublishDate time.Time `json:"publishDate,omitempty"`
			Components  []struct {
				ContentID     string `json:"contentId,omitempty"`
				MessageNumber string `json:"messageNumber,omitempty"`
			} `json:"components,omitempty"`
		} `json:"emailActivities,omitempty"`
		Count int `json:"count,omitempty"`
	} `json:"payload,omitempty"`
}

//SeriesResponse response is returned when the request type is "CommSeries".
type SeriesResponse struct {
	ID        string    `json:"id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Header    struct {
		ProcessingTime int    `json:"processingTime,omitempty"`
		ServerID       string `json:"serverId,omitempty"`
	} `json:"header,omitempty"`
	Payload struct {
		IndividualEmailActivityData []struct {
			ID             string `json:"id,omitempty"`
			Cursor         string `json:"cursor,omitempty"`
			Name           string `json:"name,omitempty"`
			RecipientsData struct {
				Recipients []struct {
					SupporterID          string    `json:"supporterId,omitempty"`
					SupporterEmail       string    `json:"supporterEmail,omitempty"`
					FirstName            string    `json:"firstName,omitempty"`
					LastName             string    `json:"lastName,omitempty"`
					Country              string    `json:"country,omitempty"`
					State                string    `json:"state,omitempty"`
					City                 string    `json:"city,omitempty"`
					TimeSent             time.Time `json:"timeSent,omitempty"`
					SplitName            string    `json:"splitName,omitempty"`
					Status               string    `json:"status,omitempty"`
					Opened               bool      `json:"opened,omitempty"`
					Clicked              bool      `json:"clicked,omitempty"`
					Converted            bool      `json:"converted,omitempty"`
					Unsubscribed         bool      `json:"unsubscribed,omitempty"`
					FirstOpenDate        time.Time `json:"firstOpenDate,omitempty,omitempty"`
					NumberOfLinksClicked string    `json:"numberOfLinksClicked,omitempty"`
					ConversionData       []struct {
						ConversionDate time.Time `json:"conversionDate,omitempty"`
						ActivityType   string    `json:"activityType,omitempty"`
						ActivityName   string    `json:"activityName,omitempty"`
						ActivityID     string    `json:"activityId,omitempty"`
						Amount         string    `json:"amount,omitempty"`
						DonationType   string    `json:"donationType,omitempty"`
					} `json:"conversionData,omitempty,omitempty"`
				} `json:"recipients,omitempty"`
				Total int `json:"total,omitempty"`
			} `json:"recipientsData,omitempty"`
		} `json:"individualEmailActivityData,omitempty"`
	} `json:"payload,omitempty"`
}
