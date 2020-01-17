package goengage

import (
	"net/http"
)

const (
	//UATHost is the hostname for Engage instances on the test server.
	UATHost = "hq.uat.igniteaction.net"
	//APIHost is the hostname for Engage instances on the production server.
	APIHost = "api.salsalabs.org"
	//ContentType is always Javascript.
	ContentType = "application/json"
	//SearchMethod is always "POST" in Engage.
	SearchMethod = http.MethodPost
)

//Segment constants
const (
	//Added indicates that the provided segment was added to the system
	Added = "ADDED"
	//Updated indicates that the provided segment was updated
	Updated = "UPDATED"
	//NotAllowed indicates that the segment represented by the provided id
	//is not allowed to be modified via the API.
	NotAllowed = "NOT_ALLOWED"
)

//Merge supporter records esult value constants.
const (
	//Found will be reported for the destination supporter if no updates were
	//specified to be performed.
	Found = "FOUND"
	//Update will be reported for the destination supporter if updates were
	//specified. It will also be reported on the main payload if the merge
	//operation was successful.
	Update = "UPDATE"
	//NotFound will be reported for the destination or source supporter if the
	//provided id(s) do not exist.
	NotFound = "NOT_FOUND"
	//Deleted will be reported for the source supporter on a successful merge.
	Deleted = "DELETED"
	//ValidationError will be reported on the main payload if either the source
	//or the destination supporter is not found, or a request to update the
	//destination was specified and validation errors occurred during that
	//update.
	ValidationError = "VALIDATION_ERROR"
	//SystemError if the merge could not be completed.
	SystemError = "SYSTEM_ERROR"
)

// Types for searching for email results.
const (
	//Email is used for searching for blasts.
	//Email = "Email"
	//CommSeries is used for searching email series.
	CommSeries = "CommSeries"
)

//Contact types.
const (
	Email     = "EMAIL"
	HomePhone = "HOME_PHONE"
	CellPhone = "CELL_PHONE"
	WorkPhone = "WORK_PHONE"
	Facebook  = "FACEBOOK_ID"
	Twitter   = "TWITTER_ID"
	Linkedin  = "LINKEDIN_ID"
)

//Error is used to report Engage errors.
type Error struct {
	ID        string `json:"id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   string `json:"details,omitempty"`
	FieldName string `json:"fieldName,omitempty"`
}

//RequestHeader provides a reference ID.
type RequestHeader struct {
	RefID string `json:"refId,omitEmpty"`
}

//Header returns server-side information for Engage API calls.
type Header struct {
	ProcessingTime int    `json:"processingTime"`
	ServerID       string `json:"serverId"`
}

//Engage endpoints for activities
const (
	SearchActivity = "/api/integration/ext/v1/activities/search"
)

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

//DonationType type
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
	EventDonation = "DONATION"
	EventTicket   = "EVENT_TICKET"
)

//Transaction type
const (
	Charge   = "CHARGE"
	Refund   = "REFUND"
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

//Ticket status.  Note that some of these definitions are for P2P only.
const (
	Valid                = "VALID"
	Refunded             = "REFUNDED"
	Cancelled            = "CANCELLED"
	CancelledAndRefunded = "CANCELLED_AND_REFUNDED"
	ValidAndRefunded     = "VALID_AND_REFUNDED"
)
