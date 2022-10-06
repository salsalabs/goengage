package goengage

//Provides a way to search for transaction templates.
//If you're want to search for transactions, use pkg/transactions.go.

// Engage endpoints for transactions.
const (
	SearchTransactionTransactionTemplates = "/api/integration/ext/v1/transactionTransactionTemplates/search"
)

// TransactionTemplate is an Engage donation template. Part of transactions in
// the API documentation.
type TransactionTemplate struct {
	CreatedBy              string `json:"createdBy"`
	CreatedDate            string `json:"createdDate"`
	ModifiedBy             string `json:"modifiedBy"`
	LastModified           string `json:"lastModified"`
	TransactionTemplateID  string `json:"templateId"`
	ActivityDate           string `json:"activityDate"`
	PersonID               string `json:"personId"`
	ActivityID             string `json:"activityId"`
	ActivityFormID         string `json:"activityFormId"`
	ActivityName           string `json:"activityName"`
	AccountExpiration      string `json:"accountExpiration"`
	AccountType            string `json:"accountType"`
	DonationType           string `json:"donationType"`
	LastTransactionType    string `json:"lastTransactionType"`
	OneTimeAmount          int    `json:"oneTimeAmount"`
	RecurringAmount        int    `json:"recurringAmount"`
	TotalReceivedAmount    int    `json:"totalReceivedAmount"`
	RecurringEnd           string `json:"recurringEnd"`
	RecurringInterval      string `json:"recurringInterval"`
	RecurringStart         string `json:"recurringStart"`
	RecurringTransactionID string `json:"recurringTransactionId"`
	IsFirstDonation        bool   `json:"isFirstDonation"`
	Dedication             string `json:"dedication"`
	DedicationType         string `json:"dedicationType"`
	Designation            string `json:"designation"`
	Notify                 string `json:"notify"`
	WasImported            bool   `json:"wasImported"`
	ReceivedAmountDonation int    `json:"receivedAmountDonation"`
	FeesPaid               int    `json:"feesPaid"`
	ReceivedAmountTickets  int    `json:"receivedAmountTickets"`
	Appeal                 string `json:"appeal"`
	Campaign               string `json:"campaign"`
	AppealName             string `json:"appealName"`
	CampaignName           string `json:"campaignName"`
	Fund                   string `json:"fund"`
	FundName               string `json:"fundName"`
	ReceivedAmountProducts int    `json:"receivedAmountProducts"`
	IsAnonymous            bool   `json:"isAnonymous"`
	DisplayName            string `json:"displayName"`
	WasAPIImported         bool   `json:"wasApiImported"`
	ExternalID             string `json:"externalId"`
	DoNotSyncCrm           bool   `json:"doNotSyncCrm"`
	HideAmount             bool   `json:"hideAmount"`
	SmartAmount            bool   `json:"smartAmount"`
	OpenEnded              bool   `json:"openEnded"`
	Result                 string `json:"result"`
}

// TransactionTemplateSearchRequest contains parameters for searching for segments.
// Please see the documentation for details.
type TransactionTemplateSearchRequest struct {
	Header  RequestHeader                           `json:"header,omitempty"`
	Payload TransactionTemplateSearchRequestPayload `json:"payload,omitempty"`
}

// TransactionTemplateSearchRequestPayload contains the payload for searching
// for transaction templates..
type TransactionTemplateSearchRequestPayload struct {
	Identifiers    []string `json:"identifiers,omitempty"`
	IdentifierType string   `json:"identifierType,omitempty"`
	CreatedFrom    string   `json:"createdFrom,omitempty"`
	CreatedTo      string   `json:"createdTo,omitempty"`
	ModifiedFrom   string   `json:"modifiedFrom,omitempty"`
	ModifiedTo     string   `json:"modifiedTo,omitempty"`
	Offset         int      `json:"offset,omitempty"`
	Count          int      `json:"count,omitempty"`
}

// TransactionTemplateSearchResponse contains the results returned by searching for segments.
type TransactionTemplateSearchResponse struct {
	ID        string                                   `json:"id,omitempty"`
	Timestamp *string                                  `json:"timestamp,omitempty"`
	Header    Header                                   `json:"header,omitempty"`
	Payload   TransactionTemplateSearchResponsePayload `json:"payload,omitempty"`
	Errors    []Error                                  `json:"errors,omitempty"`
}

// TransactionTemplateWrapper is a segment with errors and warnings.
type TransactionTemplateWrapper struct {
	Errors   []Error `json:"errors,omitempty"`
	Warnings []Error `json:"warnings,omitempty"`
	TransactionTemplate
}

// TransactionTemplateSearchResponsePayload wraps the response payload for a
// segment search.
type TransactionTemplateSearchResponsePayload struct {
	Count                int32                        `json:"count,omitempty"`
	Offset               int32                        `json:"offset,omitempty"`
	Total                int32                        `json:"total,omitempty"`
	TransactionTemplates []TransactionTemplateWrapper `json:"transactionTransactionTemplates,omitempty"`
}
