package goengage

//Describes search for transactions. Note that transactions are listed
//as part of actions *and* can be searched without actions being involved.

// Engage endpoints for transactions.
const (
	SearchTransactionDetails = "/api/integration/ext/v1/transactionDetails/search"
)

// Transaction is an Engage donation template. Part of transactions in
// the API documentation.
type DonationTransaction struct {
	CreatedBy              string  `json:"createdBy,omitempty"`
	CreatedDate            string  `json:"createdDate,omitempty"`
	ModifiedBy             string  `json:"modifiedBy,omitempty"`
	LastModified           string  `json:"lastModified,omitempty"`
	TransactionID          string  `json:"transactionId,omitempty"`
	ActivityID             string  `json:"activityId,omitempty"`
	ActivityFormID         string  `json:"activityFormId,omitempty"`
	ActivityName           string  `json:"activityName,omitempty"`
	SupporterID            string  `json:"supporterId,omitempty"`
	AccountExpiration      string  `json:"accountExpiration,omitempty"`
	AccountType            string  `json:"accountType,omitempty"`
	Amount                 float32 `json:"amount,omitempty"`
	TemplateID             string  `json:"templateId,omitempty"`
	RelatedTransactionID   string  `json:"relatedTransactionId,omitempty"`
	TransactionDate        string  `json:"transactionDate,omitempty"`
	Transaction            string  `json:"transaction,omitempty"`
	TransactionType        string  `json:"transactionType,omitempty"`
	ConfirmationEmailID    string  `json:"confirmationEmailId,omitempty"`
	ConfirmationEmailSent  string  `json:"confirmationEmailSent,omitempty"`
	WasImported            bool    `json:"wasImported,omitempty"`
	WasOffline             bool    `json:"wasOffline,omitempty"`
	DeductibleAmount       float32 `json:"deductibleAmount,omitempty"`
	FeesPaid               float32 `json:"feesPaid,omitempty"`
	WasAPIImported         bool    `json:"wasApiImported,omitempty"`
	ConvertedValidityCheck bool    `json:"convertedValidityCheck,omitempty"`
	ValidityCheckDate      string  `json:"validityCheckDate,omitempty"`
	Result                 string  `json:"result,omitempty"`
}

// TransactionSearchRequest contains parameters for searching for segments.
// Please see the documentation for details.
type TransactionSearchRequest struct {
	Header  RequestHeader                   `json:"header,omitempty"`
	Payload TransactionSearchRequestPayload `json:"payload,omitempty"`
}

// TransactionSearchRequestPayload contains the payload for searching
// for transaction templates..
type TransactionSearchRequestPayload struct {
	Identifiers     []string `json:"identifiers,omitempty"`
	IdentifierType  string   `json:"identifierType,omitempty"`
	TransactionFrom string   `json:"transactionFrom,omitempty"`
	TransactionTo   string   `json:"transactionTo,omitempty"`
	Offset          int32    `json:"offset,omitempty"`
	Count           int32    `json:"count,omitempty"`
}

// TransactionSearchResponse contains the results returned by searching for segments.
type TransactionSearchResponse struct {
	ID        string                           `json:"id,omitempty"`
	Timestamp *string                          `json:"timestamp,omitempty"`
	Header    Header                           `json:"header,omitempty"`
	Payload   TransactionSearchResponsePayload `json:"payload,omitempty"`
	Errors    []Error                          `json:"errors,omitempty"`
}

// TransactionWrapper is a segment with errors and warnings.
type TransactionWrapper struct {
	Errors   []Error `json:"errors,omitempty"`
	Warnings []Error `json:"warnings,omitempty"`
	DonationTransaction
}

// TransactionSearchResponsePayload wraps the response payload for a
// segment search.
type TransactionSearchResponsePayload struct {
	Count        int32                `json:"count,omitempty"`
	Offset       int32                `json:"offset,omitempty"`
	Total        int32                `json:"total,omitempty"`
	Transactions []TransactionWrapper `json:"transactions,omitempty"`
}
