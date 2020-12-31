package goengage

import (
	"time"
)

//Transaction holds a single monetary transaction.  Transactions are
//generally contained in "donations" as FundRaiseActivity.
type Transaction struct {
	TransactionID            string     `json:"transactionId,omitempty"`
	Type                     string     `json:"type,omitempty"`
	Reason                   string     `json:"reason,omitempty"`
	Date                     *time.Time `json:"date,omitempty"`
	Amount                   float64    `json:"amount,omitempty"`
	DeductibleAmount         float64    `json:"deductibleAmount,omitempty"`
	FeesPaid                 float64    `json:"feesPaid,omitempty"`
	GatewayTransactionID     string     `json:"gatewayTransactionId,omitempty"`
	GatewayAuthorizationCode string     `json:"gatewayAuthorizationCode,omitempty"`
	ActivityID               string     `gorm:"activity_id"`
}

//Fundraise holds a single fundraising activity.  A fundraising
//activity is actually a base activity with fundraising-specific fields.
//Note:  Fundraise also contains recurring fields.  Those will be
//automatically populated when the ActivityType is "Recurring".
type Fundraise struct {
	BaseActivity
	DonationID             string        `json:"donationId,omitempty"`
	TotalReceivedAmount    float64       `json:"totalReceivedAmount,omitempty"`
	RecurringAmount        float64       `json:"recurringAmount,omitempty"`
	OneTimeAmount          float64       `json:"oneTimeAmount,omitempty"`
	DonationType           string        `json:"donationType,omitempty"`
	RecurringInterval      string        `json:"recurringInterval,omitempty"`
	RecurringCount         int32         `json:"recurringCount,omitempty"`
	RecurringTransactionID string        `json:"recurringTransactionId,omitempty"`
	RecurringStart         *time.Time    `json:"recurringStart,omitempty"`
	RecurringEnd           *time.Time    `json:"recurringEnd,omitempty"`
	AccountType            string        `json:"accountType,omitempty"`
	AccountNumber          string        `json:"accountNumber,omitempty"`
	AccountExpiration      *time.Time    `json:"accountExpiration,omitempty"`
	AccountProvider        string        `json:"accountProvider,omitempty"`
	PaymentProcessorName   string        `json:"paymentProcessorName,omitempty"`
	FundName               string        `json:"fundName,omitempty"`
	FundGLCode             string        `json:"fundGLCode,omitempty"`
	Designation            string        `json:"designation,omitempty"`
	DedicationType         string        `json:"dedicationType,omitempty"`
	Dedication             string        `json:"dedication,omitempty"`
	Notify                 string        `json:"notify,omitempty"`
	WasImported            bool          `json:"wasimported,omitEmpty"`
	Transactions           []Transaction `json:"transactions,omitempty" gorm:"foreignkey:activity_id"`
	Supporter              Supporter     `gorm:"foreignkey:supporter_id"`
	Month                  int
	Day                    int
	Year                   int
}

//FundraiseResponse is returned for requests of type "FUNDRAISE".
type FundraiseResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Total      int32       `json:"total,omitempty"`
		Offset     int32       `json:"offset,omitempty"`
		Count      int32       `json:"count,omitempty"`
		Activities []Fundraise `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}
