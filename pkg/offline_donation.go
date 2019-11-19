package goengage

import "time"

//DonationUpsertRequest tells Engage to add and modify offline donations.
//Note that Engage does not provide a way to process donations through a gateway
//via the API.  Note, too, that there are rules to follow.
//See https://help.salsalabs.com/hc/en-us/articles/360002275693-Engage-API-Offline-Donations#addingupdating-offline-donations
type DonationUpsertRequest struct {
	Payload DonationUpsertRequestPayload `json:"payload,omitempty"`
}

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
	AccountExpiration        *time.Time          `json:"accountExpiration,omitempty"`
	AccountProvider          string             `json:"accountProvider,omitempty"`
	Fund                     string             `json:"fund,omitempty"`
	Campaign                 string             `json:"campaign,omitempty"`
	Appeal                   string             `json:"appeal,omitempty"`
	DedicationType           string             `json:"dedicationType,omitempty"`
	Dedication               string             `json:"dedication,omitempty"`
	Type                     string             `json:"type,omitempty"`
	Date                     *time.Time          `json:"date,omitempty"`
	Amount                   float64            `json:"amount,omitempty"`
	DeductibleAmount         float64            `json:"deductibleAmount,omitempty"`
	FeesPaid                 float64            `json:"feesPaid,omitempty"`
	GatewayTransactionID     string             `json:"gatewayTransactionId,omitempty"`
	GatewayAuthorizationCode string             `json:"gatewayAuthorizationCode,omitempty"`
	ActivityFormName         string             `json:"activityFormName,omitempty"`
	CustomFieldValues        []CustomFieldValue `json:"customFieldValues,omitempty"`
	Supporter                Supporter          `json:"supporter,omitempty"`
}

//DonationUpsertRequestPayload (argh) contains the list of added or modified
//donations.
type DonationUpsertRequestPayload struct {
	Donations []Donation `json:"donations,omitempty"`
}

//DonationUpsertResponse contains information about the add/update donations
//request.
type DonationUpsertResponse struct {
	Payload DonationUpsertResponsePayload `json:"payload,omitempty"`
}

//DonationUpsertError shows the details about upsert errors.
type DonationUpsertError struct {
	ID        string `json:"id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	FieldName string `json:"fieldName,omitempty"`
}

//ResultDonation contains a donation precended by errors.  If "Errors" is
//not present, then the donation was successful.
type ResultDonation struct {
	Errors                   []DonationUpsertError `json:"errors,omitempty"`
	ActivityID               string                `json:"activityId,omitempty"`
	Result                   string                `json:"result,omitempty"`
	AccountType              string                `json:"accountType,omitempty"`
	AccountNumber            string                `json:"accountNumber,omitempty"`
	AccountExpiration        *time.Time             `json:"accountExpiration,omitempty"`
	AccountProvider          string                `json:"accountProvider,omitempty"`
	Fund                     string                `json:"fund,omitempty"`
	Campaign                 string                `json:"campaign,omitempty"`
	Appeal                   string                `json:"appeal,omitempty"`
	DedicationType           string                `json:"dedicationType,omitempty"`
	Dedication               string                `json:"dedication,omitempty"`
	Type                     string                `json:"type,omitempty"`
	Date                     *time.Time             `json:"date,omitempty"`
	Amount                   float64               `json:"amount,omitempty"`
	DeductibleAmount         float64               `json:"deductibleAmount,omitempty"`
	FeesPaid                 float64               `json:"feesPaid,omitempty"`
	GatewayTransactionID     string                `json:"gatewayTransactionId,omitempty"`
	GatewayAuthorizationCode string                `json:"gatewayAuthorizationCode,omitempty"`
	ActivityFormName         string                `json:"activityFormName,omitempty"`
	CustomFieldValues        []CustomFieldValue    `json:"customFieldValues,omitempty"`
	Supporter                Supporter             `json:"supporter,omitempty"`
}

//DonationUpsertResponsePayload (argh) contains the upsert results.
type DonationUpsertResponsePayload struct {
	Donations []ResultDonation `json:"donations,omitempty"`
}
