package goengage

import (
	"time"
)

//Attendee holds information about event attendees.  Note that an attendee
//may not necessarily be a supporter.
type Attendee struct {
	AttendeeID         string     `json:"attendeeId,omitempty"`
	FirstName          string     `json:"firstName,omitempty"`
	Type               string     `json:"type,omitempty"`
	Status             string     `json:"status,omitempty"`
	LastName           string     `json:"lastName,omitempty"`
	Email              string     `json:"email,omitempty"`
	AdressLine1        string     `json:"adressLine1,omitempty"`
	AdressLine2        string     `json:"adressLine2,omitempty"`
	City               string     `json:"city,omitempty"`
	State              string     `json:"state,omitempty"`
	Phone              string     `json:"phone,omitempty"`
	IsCurrentSupporter bool       `json:"isCurrentSupporter,omitempty"`
	LastModified       *time.Time `json:"lastModified,omitempty"`
	Questions          []Question `json:"questions,omitempty"`
}

//Purchase contains the information about a single purchase for a P2P
//event attendee.
type Purchase struct {
	PurchaseID string  `json:"purchaseId,omitempty"`
	TicketID   string  `json:"ticketId,omitempty"`
	AttendeeID string  `json:"attendeeId,omitempty"`
	Name       string  `json:"name,omitempty"`
	Cost       float64 `json:"cost,omitempty"`
	Quantity   int64   `json:"quantity,omitempty"`
	Status     string  `json:"status,omitempty"`
}

//Question holds a question and a response from the event signup process.
type Question struct {
	ID       string `json:"id,omitempty"`
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

//Ticket hold information about event tickets.
type Ticket struct {
	TicketID         string     `json:"ticketId,omitempty"`
	TicketName       string     `json:"ticketName,omitempty"`
	TransactionID    string     `json:"transactionId,omitempty"`
	LastModified     *time.Time `json:"lastModified,omitempty"`
	TicketStatus     string     `json:"ticketStatus,omitempty"`
	TicketCost       float64    `json:"ticketCost,omitempty"`
	DeductibleAmount float64    `json:"deductibleAmount,omitempty"`
	Questions        []Question `json:"questions,omitempty"`
	Attendees        []Attendee `json:"attendees,omitempty"`
}

//TicketedEvent holds information about signup event for a
//ticketed event. Note that the Purchases field is only filled in when
//a P2P event attendee attends an event.
type TicketedEvent struct {
	ActivityID           string        `json:"activityId,omitempty"`
	ActivityFormName     string        `json:"activityFormName,omitempty"`
	ActivityFormID       string        `json:"activityFormId,omitempty"`
	SupporterID          string        `json:"supporterId,omitempty"`
	ActivityDate         *time.Time    `json:"activityDate,omitempty"`
	ActivityType         string        `json:"activityType,omitempty"`
	LastModified         *time.Time    `json:"lastModified,omitempty"`
	DonationID           string        `json:"donationId,omitempty"`
	TotalReceivedAmount  float64       `json:"totalReceivedAmount,omitempty"`
	OneTimeAmount        float64       `json:"oneTimeAmount,omitempty"`
	DonationType         string        `json:"donationType,omitempty"`
	AccountType          string        `json:"accountType,omitempty"`
	AccountNumber        string        `json:"accountNumber,omitempty"`
	AccountExpiration    *time.Time    `json:"accountExpiration,omitempty"`
	AccountProvider      string        `json:"accountProvider,omitempty"`
	PaymentProcessorName string        `json:"paymentProcessorName,omitempty"`
	ActivityResult       string        `json:"activityResult,omitempty"`
	Transactions         []Transaction `json:"transactions,omitempty"`
	Tickets              []Ticket      `json:"tickets,omitempty"`
	Purchases            []Purchase    `json:"purchases,omitempty"`
}

//TicketedEventResponse is returned when the request type is "TICKETED_EVENT"
//or "P2P_EVENT".  P2P events differ from ticketed events by containing purchase
//information.  A ticketed event will not have a 'Purchases" field.
type TicketedEventResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Total      int32           `json:"total,omitempty"`
		Offset     int32           `json:"offset,omitempty"`
		Count      int32           `json:"count,omitempty"`
		Activities []TicketedEvent `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}
