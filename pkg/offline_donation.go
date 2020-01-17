package goengage

import (
	"encoding/json"
	"log"
	"time"
)

//Engage endpoint for offline donation upsert.
const (
	OfflineUpsert = "api/integration/ext/v1/offlineDonations"
)

//DonorAddress is the donor's address.  Shorter than a standard address...
type DonorAddress struct {
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	County       string `json:"county,omitempty"`
	Country      string `json:"country,omitempty"`
}

//Donation contains information about an individual donation.
type Donation struct {
	AccountType              string             `json:"accountType,omitempty"`
	AccountNumber            string             `json:"accountNumber,omitempty"`
	AccountExpiration        *time.Time         `json:"accountExpiration,omitempty"`
	AccountProvider          string             `json:"accountProvider,omitempty"`
	Fund                     string             `json:"fund,omitempty"`
	Campaign                 string             `json:"campaign,omitempty"`
	Appeal                   string             `json:"appeal,omitempty"`
	DedicationType           string             `json:"dedicationType,omitempty"`
	Dedication               string             `json:"dedication,omitempty"`
	Type                     string             `json:"type,omitempty"`
	Date                     *time.Time         `json:"date,omitempty"`
	Amount                   float64            `json:"amount,omitempty"`
	DeductibleAmount         float64            `json:"deductibleAmount,omitempty"`
	FeesPaid                 float64            `json:"feesPaid,omitempty"`
	GatewayTransactionID     string             `json:"gatewayTransactionId,omitempty"`
	GatewayAuthorizationCode string             `json:"gatewayAuthorizationCode,omitempty"`
	ActivityFormName         string             `json:"activityFormName,omitempty"`
	CustomFieldValues        []CustomFieldValue `json:"customFieldValues,omitempty"`
	Supporter                Supporter          `json:"supporter,omitempty"`
}

//DonationUpsertRequest tells Engage to add and modify offline donations.
//Note that Engage does not provide a way to process donations through a gateway
//via the API.  Note, too, that there are rules to follow.
//See https://help.salsalabs.com/hc/en-us/articles/360002275693-Engage-API-Offline-Donations#addingupdating-offline-donations
type DonationUpsertRequest struct {
	Payload struct {
		Donations []Donation `json:"donations,omitempty"`
	} `json:"payload,omitempty"`
}

//DonationUpsertResponse contains information about the add/update donations
//request.
type DonationUpsertResponse struct {
	Header  Header
	Payload struct {
		Donations []ResultDonation `json:"donations,omitempty"`
	} `json:"payload,omitempty"`
}

//ResultDonation contains a donation preceded by errors.  If "Errors" is
//not present, then the donation was successful.
type ResultDonation struct {
	Errors []struct {
		ID        string `json:"id,omitempty"`
		Code      int    `json:"code,omitempty"`
		Message   string `json:"message,omitempty"`
		FieldName string `json:"fieldName,omitempty"`
	} `json:"errors,omitempty"`
	ActivityID string `json:"activityId,omitempty"`
	Result     string `json:"result,omitempty"`
	Donation
}

//ToString converts Donation record to a JSON string.
func (r Donation) ToString() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatalf("InputSink: %v\n", err)
	}
	return string(b)
}
