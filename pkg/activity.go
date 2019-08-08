package goengage

//ActSearch is used to search for activities.
const ActSearch = "/api/integration/ext/v1/activities/search"

//Activity types
const (
	SubscriptionManagementType = "SUBSCRIPTION_MANAGEMENT"
	SubscribeType              = "SUBSCRIBE"
	FundraiseType              = "FUNDRAISE"
	PetitionType               = "PETITION"
	TargetedLetterType         = "TARGETED_LETTER"
	TicketedEventType          = "TICKETED_EVENT"
	P2PEventType               = "P2P_EVENT"
)

//Donation type
const (
	OneTime   = "ONE_TIME"
	Recurring = "RECURRING"
)

//RecurringInterval
const (
	Monthly = "MONTHLY"
	Yearly  = "YEARLY"
)

//Account type
const (
	CreditCard = "CREDIT_CARD"
	ECheck     = "E_CHECK"
)

//Dedication type
const (
	None       = "NONE"
	InHonorOf  = "IN_HONOR_OF"
	InMemoryOf = "IN_MEMORY_OF"
)

//Event Reason
const (
	Donation    = "DONATION"
	EventTicket = "EVENT_TICKET"
)

//Transaction type
const (
	Charge   = "CHARGE"
	Refunc   = "REFUND"
	Cancel   = "CANCEL"
	Complete = "COMPLETE"
)

//Moderatiion state
const (
	Display     = "DISPLAY"
	DontDisplay = "DONT_DISPLAY"
	Pending     = "PENDING"
)

//Advocacy action target type
const (
	FederalExecutive = "Federal Executive"
	FederalSenate    = "Federal Senate"
	FederalHouse     = "Federal House"
	StateExecutive   = "State Executive"
	StateSenate      = "State Senate"
	StateHouse       = "State House"
	USCounty         = "US County"
	USMunicipality   = "US Municipality"
	CustomTarget     = "Custom Target"
)

//Result from action telephone calls
const (
	CallBusy       = "BUSY"
	CallFailed     = "FAILED"
	CallCompleted  = "COMPLETED"
	CallCancelled  = "CANCELLED"
	CallNoAnswer   = "NO_ANSWER"
	CallMachine    = "MACHINE"
	CallOverBudget = "OVER_BUDGET"
	CallSkipped    = "SKIPPED"
	CallNoCall     = "NO_CALL"
)

//Event activity result
const (
	DonationOnly       = "DONATION_ONLY"
	DonationAndTickets = "DONATION_AND_TICKETS"
	TicketsOnly        = "TICKETS_ONLY"
)

//Event ticket status
const (
	ValidTicket     = "VALID"
	RefundedTicket  = "REFUNDED"
	CancelledTicket = "CANCELLED"
)

// Attendee types
const (
	PurchaserAttendee = "PURCHASER"
	GuestAttendee     = "GUEST"
)

//Attendee status
const (
	RegisteredType   = "REGISTERED"
	UnregisteredType = "UNREGISTERED"
)

//P2P registration status
const (
	Valid                = "VALID"
	Refunded             = "REFUNDED"
	Cancelled            = "CANCELLED"
	CancelledAndRefunded = "CANCELLED_AND_REFUNDED"
	ValidAndRefunded     = "VALID_AND_REFUNDED"
)

//ActivityRequest is used to retrieve activities from Engage.
//Note that ActivityRequest can be used to retrieve activities based
//on three types of criteria: activity IDs, activity form IDs, modified
//date range.  Choose one and provide the necessary data.  The remainder
//will be ignored when the request is sent to Engage.
type ActivityRequest struct {
	Type            string   `json:"type,omitempty"`
	Offset          int32    `json:"offset,omitempty"`
	Count           int32    `json:"count,omitempty"`
	ActivityIDs     []string `json:"activityIds,omitempty"`
	ActivityFormIDs []string `json:"activityFormIds,omitempty"`
	ModifiedFrom    string   `json:"modifiedFrom,omitempty"`
	ModifiedTo      string   `json:"modifiedTo,omitempty"`
}

//ActivityResponse is returned in an activity response.  It contains
//a list of selected activities as well as the current position in the
//database.
type ActivityResponse struct {
	Payload struct {
		Total      int32      `json:"total,omitempty"`
		Offset     int32      `json:"offset,omitempty"`
		Count      int32      `json:"count,omitempty"`
		Activities []Activity `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}

//ActivityBase is the set of common fields returned for all activities.
//Some activities (like SUBSCRIBE or SUBSCRIPTION_MANAGEMENT) only return
//ActivityBase.  Other activities, like donations, events and P2P, return
//data appended to the base.
type ActivityBase struct {
	ActivityID       string `json:"activityID,omitempty"`
	ActivityFormName string `json:"activityFormName,omitempty"`
	ActivityFormID   string `json:"activityFormID,omitempty"`
	SupporterID      string `json:"supporterID,omitempty"`
	ActivityDate     string `json:"activityDate,omitempty"`
	ActivityType     string `json:"activityType,omitempty"`
	LastModified     string `json:"lastModified,omitempty"`
	//CustomFieldValues []something
}

//Fundraising is the additional information returned in an Activity for
//fundraising activities.
type Fundraising struct {
	DonationID           string        `json:"donationId,omitempty"`
	TotalReceivedAmount  float64       `json:"totalReceivedAmount,omitempty"`
	OneTimeAmount        float64       `json:"oneTimeAmount,omitempty"`
	DonationType         string        `json:"donationType,omitempty"`
	AccountType          string        `json:"accountType,omitempty"`
	AccountNumber        string        `json:"accountNumber,omitempty"`
	AccountExpiration    string        `json:"accountExpiration,omitempty"`
	AccountProvider      string        `json:"accountProvider,omitempty"`
	PaymentProcessorName string        `json:"paymentProcessorName,omitempty"`
	FundName             string        `json:"fundName,omitempty"`
	FundGLCode           string        `json:"fundGLCode,omitempty"`
	Designation          string        `json:"designation,omitempty"`
	DedicationType       string        `json:"dedicationType,omitempty"`
	Dedication           string        `json:"dedication,omitempty"`
	Notify               string        `json:"notify,omitempty"`
	Transactions         []Transaction `json:"transactions,omitempty"`
}

//Transaction is a single operation involving money.
type Transaction struct {
	TransactionID            string  `json:"transactionId"`
	Type                     string  `json:"type"`
	Reason                   string  `json:"reason"`
	Date                     string  `json:"date"`
	Amount                   float64 `json:"amount"`
	DeductibleAmount         float64 `json:"deductibleAmount"`
	FeesPaid                 float64 `json:"feesPaid"`
	GatewayTransactionID     string  `json:"gatewayTransactionId"`
	GatewayAuthorizationCode string  `json:"gatewayAuthorizationCode"`
}

//Activity is the wrapper for all retrieved activities.  This techinique
//works because the JSON decclarations are "omitempty".  Go simply ignores
//any empty fields during JSON Unmarshaling.
type Activity struct {
	ActivityBase
	Fundraising
}

//ActSearchResult is returned when supporters are found by a search.
type ActSearchResult struct {
	Payload struct {
		Count         int32      `json:"count,omitempty"`
		Offset        int32      `json:"offset,omitempty"`
		Total         int32      `json:"total,omitempty"`
		SupActivities []Activity `json:"Activities,omitempty"`
	} `json:"payload,omitempty"`
}
