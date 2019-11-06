package goengage

import (
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//FundraiseResponse is returned for requests of type "FUNDRAISE".
type FundraiseResponse struct {
	Header  goengage.Header          `json:"header,omitempty"`
	Payload FundraiseResponsePayload `json:"payload,omitempty"`
}

//Transaction holds a single monetary transaction.  Transactions are
//generally contained in "donations" as FundRaiseActivity.
type Transaction struct {
	TransactionID            string    `json:"transactionId,omitempty"`
	Type                     string    `json:"type,omitempty"`
	Reason                   string    `json:"reason,omitempty"`
	Date                     time.Time `json:"date,omitempty"`
	Amount                   float64   `json:"amount,omitempty"`
	DeductibleAmount         float64   `json:"deductibleAmount,omitempty"`
	FeesPaid                 float64   `json:"feesPaid,omitempty"`
	GatewayTransactionID     string    `json:"gatewayTransactionId,omitempty"`
	GatewayAuthorizationCode string    `json:"gatewayAuthorizationCode,omitempty"`
}

//Fundraise holds a single fundraising activity.  A fundraising
//activity is actually a base activity with fundraising-specific fields.
//Note:  Fundraise also contains recurring fields.  Those will be
//automatically populated when the ActivityType is "Recurring".
type Fundraise struct {
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
	OneTimeAmount          float64       `json:"oneTimeAmount,omitempty"`
	DonationType           string        `json:"donationType,omitempty"`
	RecurringInterval      string        `json:"recurringInterval,omitempty"`
	RecurringCount         int32           `json:"recurringCount,omitempty"`
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

//FundraiseResponsePayload holds the activities for a ONE_TIME search.
type FundraiseResponsePayload struct {
	Total      int         `json:"total,omitempty"`
	Offset     int         `json:"offset,omitempty"`
	Count      int32         `json:"count,omitempty"`
	Activities []Fundraise `json:"activities,omitempty"`
}
