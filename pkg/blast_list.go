package goengage

import "time"

// Web Developer list-of-blasts endpoint.
// Dee https://api.salsalabs.org/help/web-dev#operation/getBlastList

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
type BlastListResults []struct {
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
	Total   int                `json:"total,omitempty"`
	Offset  int                `json:"offset,omitempty"`
	Count   int                `json:"count,omitempty"`
	Results []BlastListResults `json:"results,omitempty"`
}

// BlastListResponse is returned from Engage for the list-of-blasts endpoint.
type BlastListResponse struct {
	ID        string                   `json:"id,omitempty"`
	Timestamp *time.Time               `json:"timestamp,omitempty"`
	Header    Header                   `json:"header,omitempty"`
	Payload   BlastListResponsePayload `json:"payload,,omitempty"`
	Errors    []Error                  `json:"errors,,omitempty"`
}
