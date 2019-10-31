package goengage

import "time"

//RecurringResponse is returned when the request type is "RECURRING".
type RecurringResponse struct {
	Header  Header                   `json:"header,omitempty"`
	Payload RecurringResponsePayload `json:"payload,omitempty"`
}

//RecurringActivity contains the details about a recurring transaction.
//Note that the base of this struct is the same as BaseActivity.
type RecurringActivity struct {
	ActivityID             string        `json:"activityId,omitempty"`
	ActivityFormName       string        `json:"activityFormName,omitempty"`
	ActivityFormID         string        `json:"activityFormId,omitempty"`
	SupporterID            string        `json:"supporterId,omitempty"`
	ActivityDate           time.Time     `json:"activityDate,omitempty"`
	ActivityType           string        `json:"activityType,omitempty"`
	LastModified           time.Time     `json:"lastModified,omitempty"`
	DonationID             string        `json:"donationId,omitempty"`
	TotalReceivedAmount    float64       `json:"totalReceivedAmount,omitempty"`
	RecurringAmount        float64       `json:"recurringAmount,omitempty"`
	DonationType           string        `json:"donationType,omitempty"`
	RecurringInterval      string        `json:"recurringInterval,omitempty"`
	RecurringCount         int           `json:"recurringCount,omitempty"`
	RecurringTransactionID string        `json:"recurringTransactionId,omitempty"`
	RecurringStart         time.Time     `json:"recurringStart,omitempty"`
	RecurringEnd           time.Time     `json:"recurringEnd,omitempty"`
	AccountType            string        `json:"accountType,omitempty"`
	AccountNumber          string        `json:"accountNumber,omitempty"`
	AccountExpiration      time.Time     `json:"accountExpiration,omitempty"`
	AccountProvider        string        `json:"accountProvider,omitempty"`
	PaymentProcessorName   string        `json:"paymentProcessorName,omitempty"`
	FundName               string        `json:"fundName,omitempty"`
	FundGLCode             string        `json:"fundGLCode,omitempty"`
	Designation            string        `json:"designation,omitempty"`
	DedicationType         string        `json:"dedicationType,omitempty"`
	Dedication             string        `json:"dedication,omitempty"`
	Notify                 string        `json:"notify,omitempty"`
	Transactions           []Transaction `json:"transactions,omitempty"`
}

//RecurringResponsePayload contains the data returned for a RECURRING
//search.
type RecurringResponsePayload struct {
	Total      int                 `json:"total,omitempty"`
	Offset     int                 `json:"offset,omitempty"`
	Count      int                 `json:"count,omitempty"`
	Activities []RecurringActivity `json:"activities,omitempty"`
}
