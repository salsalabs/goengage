package goengage

import "time"

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

//ActivityRequest is used to retrieve activities from Engage.
//Note that ActivityRequest can be used to retrieve activities based
//on three types of criteria: activity IDs, activity form IDs, modified
//date range.  Choose one and provide the necessary data.  The remainder
//will be ignored when the request is sent to Engage.
type ActivityRequest struct {
	Type            string   `json:"type,omitempty"`
	Offset          int32    `json:"offset,omitempty"`
	Count           int32    `json:"count,omitempty"`
	ActivityIDs     []string `json:"activityIds,omitempty,omitempty"`
	ActivityFormIDs []string `json:"activityFormIds,omitempty,omitempty"`
	ModifiedFrom    string   `json:"modifiedFrom,omitempty,omitempty"`
	ModifiedTo      string   `json:"modifiedTo,omitempty,omitempty"`
}

//BaseResponse is the set of common fields returned for all activities.
//Some activities (like SUBSCRIBE or SUBSCRIPTION_MANAGEMENT) only return
//ActivityBase.  Other activities, like donations, events and P2P, return
//data appended to the base.
type BaseResponse struct {
	Header struct {
		ProcessingTime int    `json:"processingTime,omitempty"`
		ServerID       string `json:"serverId,omitempty"`
	} `json:"header,omitempty"`
	Payload struct {
		Total      int `json:"total,omitempty"`
		Offset     int `json:"offset,omitempty"`
		Count      int `json:"count,omitempty"`
		Activities []struct {
			ActivityType     string    `json:"activityType,omitempty"`
			ActivityID       string    `json:"activityId,omitempty"`
			ActivityFormName string    `json:"activityFormName,omitempty"`
			ActivityFormID   string    `json:"activityFormId,omitempty"`
			SupporterID      string    `json:"supporterId,omitempty"`
			ActivityDate     time.Time `json:"activityDate,omitempty"`
			LastModified     time.Time `json:"lastModified,omitempty"`
		} `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}

//FundraiseResponse is returned for requests of type "FUNDRAISE".
type FundraiseResponse struct {
	Header struct {
		ProcessingTime int    `json:"processingTime,omitempty"`
		ServerID       string `json:"serverId,omitempty"`
	} `json:"header,omitempty"`
	Payload struct {
		Total  int `json:"total,omitempty"`
		Offset int `json:"offset,omitempty"`
		Count  int `json:"count,omitempty"`
	} `json:"payload,omitempty"`
	Activities []struct {
		ActivityID           string    `json:"activityId,omitempty"`
		ActivityFormName     string    `json:"activityFormName,omitempty"`
		ActivityFormID       string    `json:"activityFormId,omitempty"`
		SupporterID          string    `json:"supporterId,omitempty"`
		ActivityDate         time.Time `json:"activityDate,omitempty"`
		ActivityType         string    `json:"activityType,omitempty"`
		LastModified         time.Time `json:"lastModified,omitempty"`
		DonationID           string    `json:"donationId,omitempty"`
		TotalReceivedAmount  float64   `json:"totalReceivedAmount,omitempty"`
		OneTimeAmount        float64   `json:"oneTimeAmount,omitempty"`
		DonationType         string    `json:"donationType,omitempty"`
		AccountType          string    `json:"accountType,omitempty"`
		AccountNumber        string    `json:"accountNumber,omitempty"`
		AccountExpiration    time.Time `json:"accountExpiration,omitempty"`
		AccountProvider      string    `json:"accountProvider,omitempty"`
		PaymentProcessorName string    `json:"paymentProcessorName,omitempty"`
		FundName             string    `json:"fundName,omitempty"`
		FundGLCode           string    `json:"fundGLCode,omitempty"`
		Designation          string    `json:"designation,omitempty"`
		DedicationType       string    `json:"dedicationType,omitempty"`
		Dedication           string    `json:"dedication,omitempty"`
		Notify               string    `json:"notify,omitempty"`
		Transactions         []struct {
			TransactionID            string    `json:"transactionId,omitempty"`
			Type                     string    `json:"type,omitempty"`
			Reason                   string    `json:"reason,omitempty"`
			Date                     time.Time `json:"date,omitempty"`
			Amount                   float64   `json:"amount,omitempty"`
			DeductibleAmount         float64   `json:"deductibleAmount,omitempty"`
			FeesPaid                 float64   `json:"feesPaid,omitempty"`
			GatewayTransactionID     string    `json:"gatewayTransactionId,omitempty"`
			GatewayAuthorizationCode string    `json:"gatewayAuthorizationCode,omitempty"`
		} `json:"transactions,omitempty"`
	} `json:"activities,omitempty"`
}

//RecurringResponse is returned when the request type is "RECURRING".
type RecurringResponse struct {
	Header struct {
		ProcessingTime int    `json:"processingTime,omitempty"`
		ServerID       string `json:"serverId,omitempty"`
	} `json:"header,omitempty"`
	Payload struct {
		Total      int `json:"total,omitempty"`
		Offset     int `json:"offset,omitempty"`
		Count      int `json:"count,omitempty"`
		Activities []struct {
			ActivityID           string    `json:"activityId,omitempty"`
			ActivityFormName     string    `json:"activityFormName,omitempty"`
			ActivityFormID       string    `json:"activityFormId,omitempty"`
			SupporterID          string    `json:"supporterId,omitempty"`
			ActivityDate         time.Time `json:"activityDate,omitempty"`
			ActivityType         string    `json:"activityType,omitempty"`
			LastModified         time.Time `json:"lastModified,omitempty"`
			DonationID           string    `json:"donationId,omitempty"`
			TotalReceivedAmount  float64   `json:"totalReceivedAmount,omitempty"`
			OneTimeAmount        float64   `json:"oneTimeAmount,omitempty"`
			DonationType         string    `json:"donationType,omitempty"`
			AccountType          string    `json:"accountType,omitempty"`
			AccountNumber        string    `json:"accountNumber,omitempty"`
			AccountExpiration    time.Time `json:"accountExpiration,omitempty"`
			AccountProvider      string    `json:"accountProvider,omitempty"`
			PaymentProcessorName string    `json:"paymentProcessorName,omitempty"`
			FundName             string    `json:"fundName,omitempty"`
			FundGLCode           string    `json:"fundGLCode,omitempty"`
			Designation          string    `json:"designation,omitempty"`
			DedicationType       string    `json:"dedicationType,omitempty"`
			Dedication           string    `json:"dedication,omitempty"`
			Notify               string    `json:"notify,omitempty"`
			Transactions         []struct {
				TransactionID            string    `json:"transactionId,omitempty"`
				Type                     string    `json:"type,omitempty"`
				Reason                   string    `json:"reason,omitempty"`
				Date                     time.Time `json:"date,omitempty"`
				Amount                   float64   `json:"amount,omitempty"`
				DeductibleAmount         float64   `json:"deductibleAmount,omitempty"`
				FeesPaid                 float64   `json:"feesPaid,omitempty"`
				GatewayTransactionID     string    `json:"gatewayTransactionId,omitempty"`
				GatewayAuthorizationCode string    `json:"gatewayAuthorizationCode,omitempty"`
			} `json:"transactions,omitempty"`
		} `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}

//PetitionResponse is returned when the request type is "PETITION".
type PetitionResponse struct {
	Header struct {
		ProcessingTime int    `json:"processingTime,omitempty"`
		ServerID       string `json:"serverId,omitempty"`
	} `json:"header,omitempty"`
	Payload struct {
		Total      int `json:"total,omitempty"`
		Offset     int `json:"offset,omitempty"`
		Count      int `json:"count,omitempty"`
		Activities []struct {
			ActivityID               string    `json:"activityId,omitempty"`
			ActivityFormName         string    `json:"activityFormName,omitempty"`
			ActivityFormID           string    `json:"activityFormId,omitempty"`
			SupporterID              string    `json:"supporterId,omitempty"`
			ActivityDate             time.Time `json:"activityDate,omitempty"`
			ActivityType             string    `json:"activityType,omitempty"`
			LastModified             time.Time `json:"lastModified,omitempty"`
			Comment                  string    `json:"comment,omitempty"`
			ModerationState          string    `json:"moderationState,omitempty"`
			DisplaySignaturePublicly bool      `json:"displaySignaturePublicly,omitempty"`
			DisplayCommentPublicly   bool      `json:"displayCommentPublicly,omitempty"`
		} `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}

//TargetedLetterResponse is returned when the request is "TARGETED_LETTERS".
type TargetedLetterResponse struct {
	Header struct {
		ProcessingTime int    `json:"processingTime,omitempty"`
		ServerID       string `json:"serverId,omitempty"`
	} `json:"header,omitempty"`
	Payload struct {
		Total      int `json:"total,omitempty"`
		Offset     int `json:"offset,omitempty"`
		Count      int `json:"count,omitempty"`
		Activities []struct {
			ActivityID       string    `json:"activityId,omitempty"`
			ActivityFormName string    `json:"activityFormName,omitempty"`
			ActivityFormID   string    `json:"activityFormId,omitempty"`
			SupporterID      string    `json:"supporterId,omitempty"`
			ActivityDate     time.Time `json:"activityDate,omitempty"`
			ActivityType     string    `json:"activityType,omitempty"`
			LastModified     time.Time `json:"lastModified,omitempty"`
			Letters          []struct {
				Name               string `json:"name,omitempty"`
				Subject            string `json:"subject,omitempty"`
				Message            string `json:"message,omitempty"`
				AdditionalComment  string `json:"additionalComment,omitempty"`
				SubjectWasModified bool   `json:"subjectWasModified,omitempty"`
				MessageWasModified bool   `json:"messageWasModified,omitempty"`
				Targets            []struct {
					TargetID            string `json:"targetId,omitempty"`
					TargetName          string `json:"targetName,omitempty"`
					TargetTitle         string `json:"targetTitle,omitempty"`
					PoliticalParty      string `json:"politicalParty,omitempty"`
					TargetType          string `json:"targetType,omitempty"`
					State               string `json:"state,omitempty"`
					DistrictID          string `json:"districtId,omitempty"`
					DistrictName        string `json:"districtName,omitempty"`
					Role                string `json:"role,omitempty"`
					SentEmail           string `json:"sentEmail,omitempty"`
					SentFacebook        string `json:"sentFacebook,omitempty"`
					SentTwitter         string `json:"sentTwitter,omitempty"`
					MadeCall            string `json:"madeCall,omitempty"`
					CallDurationSeconds string `json:"callDurationSeconds,omitempty"`
					CallResult          string `json:"callResult,omitempty"`
				} `json:"targets,omitempty"`
			} `json:"letters,omitempty"`
		} `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}

//TicketedEventResponse is returned when the request type is "TICKETED_EVENT".
type TicketedEventResponse struct {
	Header struct {
		ProcessingTime int    `json:"processingTime,omitempty"`
		ServerID       string `json:"serverId,omitempty"`
	} `json:"header,omitempty"`
	Payload struct {
		Total      int `json:"total,omitempty"`
		Offset     int `json:"offset,omitempty"`
		Count      int `json:"count,omitempty"`
		Activities []struct {
			ActivityID           string    `json:"activityId,omitempty"`
			ActivityFormName     string    `json:"activityFormName,omitempty"`
			ActivityFormID       string    `json:"activityFormId,omitempty"`
			SupporterID          string    `json:"supporterId,omitempty"`
			ActivityDate         time.Time `json:"activityDate,omitempty"`
			ActivityType         string    `json:"activityType,omitempty"`
			LastModified         time.Time `json:"lastModified,omitempty"`
			DonationID           string    `json:"donationId,omitempty"`
			TotalReceivedAmount  float64   `json:"totalReceivedAmount,omitempty"`
			OneTimeAmount        float64   `json:"oneTimeAmount,omitempty"`
			DonationType         string    `json:"donationType,omitempty"`
			AccountType          string    `json:"accountType,omitempty"`
			AccountNumber        string    `json:"accountNumber,omitempty"`
			AccountExpiration    time.Time `json:"accountExpiration,omitempty"`
			AccountProvider      string    `json:"accountProvider,omitempty"`
			PaymentProcessorName string    `json:"paymentProcessorName,omitempty"`
			ActivityResult       string    `json:"activityResult,omitempty"`
			Transactions         []struct {
				TransactionID            string    `json:"transactionId,omitempty"`
				Type                     string    `json:"type,omitempty"`
				Reason                   string    `json:"reason,omitempty"`
				Date                     time.Time `json:"date,omitempty"`
				Amount                   float64   `json:"amount,omitempty"`
				DeductibleAmount         float64   `json:"deductibleAmount,omitempty"`
				FeesPaid                 float64   `json:"feesPaid,omitempty"`
				GatewayTransactionID     string    `json:"gatewayTransactionId,omitempty"`
				GatewayAuthorizationCode string    `json:"gatewayAuthorizationCode,omitempty"`
				ReasonID                 string    `json:"reasonId,omitempty,omitempty"`
			} `json:"transactions,omitempty"`
			Tickets []struct {
				TicketID         string    `json:"ticketId,omitempty"`
				TicketName       string    `json:"ticketName,omitempty"`
				TransactionID    string    `json:"transactionId,omitempty"`
				LastModified     time.Time `json:"lastModified,omitempty"`
				TicketStatus     string    `json:"ticketStatus,omitempty"`
				TicketCost       float64   `json:"ticketCost,omitempty"`
				DeductibleAmount float64   `json:"deductibleAmount,omitempty"`
			} `json:"tickets,omitempty"`
			Questions []struct {
				ID       string `json:"id,omitempty"`
				Question string `json:"question,omitempty"`
				Answer   string `json:"answer,omitempty"`
			} `json:"questions,omitempty"`
			Attendees []struct {
				AttendeeID         string    `json:"attendeeId,omitempty"`
				FirstName          string    `json:"firstName,omitempty"`
				Type               string    `json:"type,omitempty"`
				Status             string    `json:"status,omitempty"`
				LastName           string    `json:"lastName,omitempty"`
				Email              string    `json:"email,omitempty"`
				AdressLine1        string    `json:"adressLine1,omitempty"`
				AdressLine2        string    `json:"adressLine2,omitempty"`
				City               string    `json:"city,omitempty"`
				State              string    `json:"state,omitempty"`
				Phone              string    `json:"phone,omitempty"`
				IsCurrentSupporter bool      `json:"isCurrentSupporter,omitempty"`
				LastModified       time.Time `json:"lastModified,omitempty"`
				Questions          []struct {
					ID       string `json:"id,omitempty"`
					Question string `json:"question,omitempty"`
					Answer   string `json:"answer,omitempty"`
				} `json:"questions,omitempty"`
			} `json:"attendees,omitempty"`
			Purchases []struct {
				PurchaseID string  `json:"purchaseId,omitempty"`
				TicketID   string  `json:"ticketId,omitempty"`
				AttendeeID string  `json:"attendeeId,omitempty"`
				Name       string  `json:"name,omitempty"`
				Cost       float64 `json:"cost,omitempty"`
				Quantity   string  `json:"quantity,omitempty"`
				Status     string  `json:"status,omitempty"`
				Details    []struct {
					ID       string `json:"id,omitempty"`
					Question string `json:"question,omitempty"`
					Answer   string `json:"answer,omitempty"`
				} `json:"details,omitempty"`
				PurchasedDate string `json:"purchasedDate,omitempty"`
			} `json:"purchases,omitempty"`
		} `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}
