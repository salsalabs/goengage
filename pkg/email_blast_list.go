package goengage

import "time"

// Web Developer list-of-blasts endpoint.
// Dee https://api.salsalabs.org/help/web-dev#operation/getBlastList

// BlastListRequest is a convenience struct used to record the
// search criteria for the list-of-blasts endpoint. Note that the
// contents are used to append queries to the endpoint URL.
type BlastListRequest struct {
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	Criteria  string `json:"criteria,omitempty"`
	SortField string `json:"sortField,omitempty"`
	SortOrder string `json:"sortOrder,omitempty"`
	Count     int32  `json:"count,omitempty"`
	Offset    int32  `json:"offset,omitempty"`
}

// BlastContent describes a single blast.
type BlastContent struct {
	Subject                string    `json:"subject,omitempty"`
	PageTitle              string    `json:"pageTitle,omitempty"`
	PageURL                string    `json:"pageUrl,omitempty"`
	WebVersionEnabled      bool      `json:"webVersionEnabled,omitempty"`
	WebVersionRedirectURL  string    `json:"webVersionRedirectUrl,omitempty"`
	WebVersionRedirectDate time.Time `json:"webVersionRedirectDate,omitempty"`
}

// BlastListResponsePayload carries the tthe blast information that's
// not considered content.
type BlastListResult struct {
	ID           string         `json:"id,omitempty"`
	Name         string         `json:"name,omitempty"`
	Description  string         `json:"description,omitempty"`
	Status       string         `json:"status,omitempty"`
	Topic        string         `json:"topic,omitempty"`
	PublishDate  time.Time      `json:"publishDate,omitempty"`
	ScheduleDate time.Time      `json:"scheduleDate,omitempty"`
	Content      []BlastContent `json:"content,omitempty"`
}

// BlastListResponse is the payload from Engage for the list-of-blasts endpoint.
type BlastListResponsePayload struct {
	Total   int32             `json:"total,omitempty"`
	Offset  int32             `json:"offset,omitempty"`
	Count   int32             `json:"count,omitempty"`
	Results []BlastListResult `json:"results,omitempty"`
}

// BlastListResponse is returned from Engage for the list-of-blasts endpoint.
type BlastListResponse struct {
	ID        string                   `json:"id,omitempty"`
	Timestamp *time.Time               `json:"timestamp,omitempty"`
	Header    Header                   `json:"header,omitempty"`
	Payload   BlastListResponsePayload `json:"payload,,omitempty"`
	Errors    []Error                  `json:"errors,,omitempty"`
}
