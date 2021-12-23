package goengage

//Describes search for transactions. Note that transactions are listed
//as part of actions *and* can be searched without actions being involved.

//Engage endpoints for transactions.
const (
	SearchTransactionDetails = "/api/integration/ext/v1/transactionDetails/search"
)

//Transaction is an Engage donation template. Part of transactions in
//the API documentation.
type DonationTransaction struct {
	CreatedBy              string `json:"createdBy,omitempty"`
	CreatedDate            string `json:"createdDate,omitempty"`
	ModifiedBy             string `json:"modifiedBy,omitempty"`
	LastModified           string `json:"lastModified,omitempty"`
	TransactionID          string `json:"templateId,omitempty"`
	ActivityDate           string `json:"activityDate,omitempty"`
	PersonID               string `json:"personId,omitempty"`
	ActivityID             string `json:"activityId,omitempty"`
	ActivityFormID         string `json:"activityFormId,omitempty"`
	ActivityName           string `json:"activityName,omitempty"`
	AccountExpiration      string `json:"accountExpiration,omitempty"`
	AccountType            string `json:"accountType,omitempty"`
	DonationType           string `json:"donationType,omitempty"`
	LastTransactionType    string `json:"lastTransactionType,omitempty"`
	OneTimeAmount          int    `json:"oneTimeAmount,omitempty"`
	RecurringAmount        int    `json:"recurringAmount,omitempty"`
	TotalReceivedAmount    int    `json:"totalReceivedAmount,omitempty"`
	RecurringEnd           string `json:"recurringEnd,omitempty"`
	RecurringInterval      string `json:"recurringInterval,omitempty"`
	RecurringStart         string `json:"recurringStart,omitempty"`
	RecurringTransactionID string `json:"recurringTransactionId,omitempty"`
	IsFirstDonation        bool   `json:"isFirstDonation,omitempty"`
	Dedication             string `json:"dedication,omitempty"`
	DedicationType         string `json:"dedicationType,omitempty"`
	Designation            string `json:"designation,omitempty"`
	Notify                 string `json:"notify,omitempty"`
	WasImported            bool   `json:"wasImported,omitempty"`
	ReceivedAmountDonation int    `json:"receivedAmountDonation,omitempty"`
	FeesPaid               int    `json:"feesPaid,omitempty"`
	ReceivedAmountTickets  int    `json:"receivedAmountTickets,omitempty"`
	Appeal                 string `json:"appeal,omitempty"`
	Campaign               string `json:"campaign,omitempty"`
	AppealName             string `json:"appealName,omitempty"`
	CampaignName           string `json:"campaignName,omitempty"`
	Fund                   string `json:"fund,omitempty"`
	FundName               string `json:"fundName,omitempty"`
	ReceivedAmountProducts int    `json:"receivedAmountProducts,omitempty"`
	IsAnonymous            bool   `json:"isAnonymous,omitempty"`
	DisplayName            string `json:"displayName,omitempty"`
	WasAPIImported         bool   `json:"wasApiImported,omitempty"`
	ExternalID             string `json:"externalId,omitempty"`
	DoNotSyncCrm           bool   `json:"doNotSyncCrm,omitempty"`
	HideAmount             bool   `json:"hideAmount,omitempty"`
	SmartAmount            bool   `json:"smartAmount,omitempty"`
	OpenEnded              bool   `json:"openEnded,omitempty"`
	Result                 string `json:"result,omitempty"`
}

//TransactionSearchRequest contains parameters for searching for segments.
// Please see the documentation for details.
type TransactionSearchRequest struct {
	Header  RequestHeader                   `json:"header,omitempty"`
	Payload TransactionSearchRequestPayload `json:"payload,omitempty"`
}

//TransactionSearchRequestPayload contains the payload for searching
//for transaction templates..
type TransactionSearchRequestPayload struct {
	Identifiers     []string `json:"identifiers,omitEmpty"`
	IdentifierType  string   `json:"identifierType,omitEmpty"`
	TransactionFrom string   `json:"transactionFrom,omitEmpty"`
	TransactionTo   string   `json:"transactionTo,omitEmpty"`
	Offset          int      `json:"offset,omitEmpty"`
	Count           int      `json:"count,omitEmpty"`
}

//TransactionSearchResponse contains the results returned by searching for segments.
type TransactionSearchResponse struct {
	ID        string                           `json:"id,omitempty"`
	Timestamp *string                          `json:"timestamp,omitempty"`
	Header    Header                           `json:"header,omitempty"`
	Payload   TransactionSearchResponsePayload `json:"payload,omitempty"`
	Errors    []Error                          `json:"errors,omitempty"`
}

//TransactionWrapper is a segment with errors and warnings.
type TransactionWrapper struct {
	Errors   []Error `json:"errors,omitempty"`
	Warnings []Error `json:"warnings,omitempty"`
	DonationTransaction
}

//TransactionSearchResponsePayload wraps the response payload for a
//segment search.
type TransactionSearchResponsePayload struct {
	Count        int32                `json:"count,omitempty"`
	Offset       int32                `json:"offset,omitempty"`
	Total        int32                `json:"total,omitempty"`
	Transactions []TransactionWrapper `json:"transactionTransactions,omitempty"`
}
