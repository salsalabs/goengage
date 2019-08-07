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
	TicketsOnly        = "TICKETS_ONLY	"
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

//ActivityIDs (one of the request options)
type ActivityIDs []string

//ActivityRequest is used to retrieve activities from Engage.
//Note that ActivityRequest can be used to retrieve activities based
//on three types of criteria: activity IDs, activity form IDs, modified
//date range.  Choose one and provide the necessary data.  The remainder
//will be ignored when the request is sent to Engage.
type ActivityRequest struct {
	ActivityIDs     []string `json:"activityIDs,omitEmpty"`
	ActivityFormIDs []string `json:"activityFormIDs,omitEmpty"`
	ModifiedFrom    string   `json:"modifiedFrom,omitempty"`
	ModifiedTo      string   `json:"modifiedTo,omitempty"`
	Offset          int32    `json:"offset,omitempty"`
	Count           int32    `json:"count,omitempty"`
	Type            string   `json:"type,omitempty"`
}

//ActivityResponse is the data returned by every activity request.  Typically,
//the base data is returned with activity specific data after it.  That means
//that you can use ActivityResponse to accept data from any activity call.
//JSON"s marshalling will not equip any object element that is not in The
//reply from Engage.
type ActivityResponse struct {
	ResponseHeader
	ActivityPayload
}

//ActivityPayload is returned in an activity response.  It contains
//a list of selected activities as well as the current position in the
//database.
type ActivityPayload struct {
	Payload struct {
		Total      int32    `json:"total,omitempty"`
		Offset     int32    `json:"offset,omitempty"`
		Count      int32    `json:"count,omitempty"`
		Activities []Activity `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}

//Activity is one activity of the list returned in an activity response.
type Activity struct {
	ActivityID       string `json:"activityID,omitempty"`
	ActivityFormName string `json:"activityFormName,omitempty"`
	ActivityFormID   string `json:"activityFormID,omitempty"`
	SupporterID      string `json:"supporterID,omitempty"`
	ActivityDate     string `json:"activityDate,omitempty"`
	ActivityType     string `json:"activityType,omitempty"`
	LastModified     string `json:"lastModified,omitempty"`
	//CustomFieldValues []something
}

//ActSearchResult is returned when supporters are found by a search.
type ActSearchResult struct {
	Payload struct {
		Count         int32         `json:"count,omitempty"`
		Offset        int32         `json:"offset,omitempty"`
		Total         int32         `json:"total,omitempty"`
		SupActivities []Activity `json:"Activities,omitempty"`
	} `json:"payload,omitempty"`
}
