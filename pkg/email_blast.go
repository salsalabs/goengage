package goengage

//Search for email blasts and return activity.
// Search: see https://api.salsalabs.org/help/integration#operation/emailsSearch
// Activity: see https://help.salsalabs.com/hc/en-us/articles/360019505914-Engage-API-Email-Results

const (
	//EmailType indicates a search for an email blast
	EmailType = "EMAIL"

	//CommSeriesType indicates a search for a communications series.
	CommSeriesType = "CommSeries"

	//EmailBlastSearch is used to find blasts.
	EmailBlastSearch = "/api/integration/ext/v1/emails/search"

	//EmailIndividualBlast is used to retrieve a single blast
	//as well as as all of the recipient data.
	EmailIndividualBlast = "/api/integration/ext/v1/emails/individualResults"
)

//EmailComponent are used in comm series.
type EmailComponent struct {
	ContentID     string `json:"contentId,omitempty"`
	MessageNumber string `json:"messageNumber,omitempty"`
}

//EmailActivity describes the contents of the email.
type EmailActivity struct {
	ID          string            `json:"id,omitempty"`
	Topic       string            `json:"topic,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	PublishDate string            `json:"publishDate,omitempty"`
	Components  *[]EmailComponent `json:"components,omitempty"`
	EmailErrors *[]EmailError     `json:"errors,omitempty"`
}

//Conversion hold information about any donations made as a result of an
//email blast.
type Conversion struct {
	ConversionDate string `json:"conversionDate,omitempty"`
	ActivityType   string `json:"activityType,omitempty"`
	ActivityName   string `json:"activityName,omitempty"`
	ActivityID     string `json:"activityId,omitempty"`
	ActivityFormID string `json:"activityFormId,omitempty"`
	Amount         string `json:"amount,omitempty"`
	DonationType   string `json:"donationType,omitempty"`
}

//SingleBlastRecipient holds identity, blast statistics
//and conversion info.
type SingleBlastRecipient struct {
	SalesforceID         string       `json:"salesforceId,omitempty"`
	SupporterID          string       `json:"supporterId,omitempty"`
	ExternalID           string       `json:"externalId,omitempty"`
	SupporterEmail       string       `json:"supporterEmail,omitempty"`
	FirstName            string       `json:"firstName,omitempty"`
	LastName             string       `json:"lastName,omitempty"`
	City                 string       `json:"city,omitempty"`
	State                string       `json:"state,omitempty"`
	Country              string       `json:"country,omitempty"`
	TimeSent             string       `json:"timeSent,omitempty"`
	SplitName            string       `json:"splitName,omitempty"`
	EmailSeriesName      string       `json:"emailSeriesName,omitempty"`
	Status               string       `json:"status,omitempty"`
	Opened               bool         `json:"opened,omitempty"`
	Clicked              bool         `json:"clicked,omitempty"`
	Converted            bool         `json:"converted,omitempty"`
	Unsubscribed         bool         `json:"unsubscribed,omitempty"`
	FirstOpenDate        string       `json:"firstOpenDate,omitempty"`
	NumberOfLinksClicked string       `json:"numberOfLinksClicked,omitempty"`
	BounceCategory       string       `json:"bounceCategory,omitempty"`
	BounceCode           string       `json:"bounceCode,omitempty"`
	ConversionData       []Conversion `json:"conversionData,omitempty"`
}

//EmailError describes issues found in the email blast search call.
type EmailError struct {
	ID          string `json:"id,omitempty"`
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Details     string `json:"details,omitempty"`
	FieldName   string `json:"fieldName,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	ContentID   string `json:"contentId,omitempty"`
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
	Total           int32           `json:"total,omitempty"`
	Offset          int32           `json:"offset,omitempty"`
	Count           int32           `json:"count,omitempty"`
	EmailActivities []EmailActivity `json:"emailActivities,omitempty"`
}

//EmailBlastSearchResponse wraps a response payload.
type EmailBlastSearchResponse struct {
	ID        string                          `json:"id,omitempty"`
	TimeStamp string                          `json:"timestamp,omitempty"`
	Header    Header                          `json:"header,omitempty"`
	Payload   EmailBlastSearchResponsePayload `json:"payload,omitempty"`
}

//IndivualBlastRequestPayload sets the criteria for
//the blasts to read.
type IndivualBlastRequestPayload struct {
	ID            string `json:"id,omitempty"`
	ContentID     string `json:"contentId,omitempty"`
	Cursor        string `json:"cursor,omitempty"`
	PublishedFrom string `json:"publishedFrom,omitempty"`
	PublishedTo   string `json:"publishedTo,omitempty"`
	Type          string `json:"type,omitempty"`
	Offset        int    `json:"offset,omitempty"`
	Count         int    `json:"count,omitempty"`
}

//IndivualBlastRequest wraps the request payload.
type IndivualBlastRequest struct {
	ID      string                      `json:"id"`
	Header  RequestHeader               `json:"header,omitempty"`
	Payload IndivualBlastRequestPayload `json:"payload,omitempty"`
}

//IndivualBlastResponsePayload holds the response content.
type IndivualBlastResponsePayload struct {
	Total                      int32                     `json:"total,omitempty"`
	Offset                     int32                     `json:"offset,omitempty"`
	IndividualEmalActivityData []IndividualEmailActivity `json:"indivualEmailActivityData,omitempty"`
	EmailErrors                []EmailError              `json:"EmailErrors,omitempty"`
}

//IndividualEmailActivity contains the email activity for one blast.
type IndividualEmailActivity []struct {
	ID             string                    `json:"id,omitempty"`
	Cursor         string                    `json:"cursor,omitempty"`
	Name           string                    `json:"name,omitempty"`
	RecipientsData SingleBlastRecipientsData `json:"recipientsData,omitempty"`
}

//SingleBlastRecipientsData is a wrapper around the recipients
//for a blast.
type SingleBlastRecipientsData struct {
	Recipients []SingleBlastRecipient `json:"recipients,omitempty"`
	Total      int32                  `json:"total,omitempty"`
}
