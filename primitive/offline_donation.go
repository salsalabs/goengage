package goengage

import "time"

//DonationUpsertRequest tells Engage to add and modify offline donations.
//Note that Engage does not provide a way to process donations through a gateway
//via the API.  Note, too, that there are rules to follow.
//See https://help.salsalabs.com/hc/en-us/articles/360002275693-Engage-API-Offline-Donations#addingupdating-offline-donations
type DonationUpsertRequest struct {
	Payload DonationUpsertRequestPayload `json:"payload"`
}

//DonorAddress is the donor's address.  Shorter than a standard address...
type DonorAddress struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
	County       string `json:"county"`
	Country      string `json:"country"`
}

//Donation contains information about an individual donation.
type Donation struct {
	AccountType              string             `json:"accountType"`
	AccountNumber            string             `json:"accountNumber"`
	AccountExpiration        time.Time          `json:"accountExpiration"`
	AccountProvider          string             `json:"accountProvider"`
	Fund                     string             `json:"fund"`
	Campaign                 string             `json:"campaign"`
	Appeal                   string             `json:"appeal"`
	DedicationType           string             `json:"dedicationType"`
	Dedication               string             `json:"dedication"`
	Type                     string             `json:"type"`
	Date                     time.Time          `json:"date"`
	Amount                   float64            `json:"amount"`
	DeductibleAmount         float64            `json:"deductibleAmount"`
	FeesPaid                 float64            `json:"feesPaid"`
	GatewayTransactionID     string             `json:"gatewayTransactionId"`
	GatewayAuthorizationCode string             `json:"gatewayAuthorizationCode"`
	ActivityFormName         string             `json:"activityFormName"`
	CustomFieldValues        []CustomFieldValue `json:"customFieldValues"`
	Supporter                Supporter          `json:"supporter"`
}

//DonationUpsertRequestPayload (argh) contains the list of added or modified
//donations.
type DonationUpsertRequestPayload struct {
	Donations []Donation `json:"donations"`
}

//DonationUpsertResponse contains information about the add/update donations
//request.
type DonationUpsertResponse struct {
	Payload DonationUpsertResponsePayload `json:"payload"`
}

//DonationUpsertError shows the details about upsert errors.
type DonationUpsertError struct {
	ID        string `json:"id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	FieldName string `json:"fieldName"`
}

//ResultDonation contains a donation precended by errors.  If "Errors" is
//not present, then the donation was successful.
type ResultDonation struct {
	Errors                   []DonationUpsertError `json:"errors"`
	ActivityID               string                `json:"activityId"`
	Result                   string                `json:"result"`
	AccountType              string                `json:"accountType"`
	AccountNumber            string                `json:"accountNumber"`
	AccountExpiration        time.Time             `json:"accountExpiration"`
	AccountProvider          string                `json:"accountProvider"`
	Fund                     string                `json:"fund"`
	Campaign                 string                `json:"campaign"`
	Appeal                   string                `json:"appeal"`
	DedicationType           string                `json:"dedicationType"`
	Dedication               string                `json:"dedication"`
	Type                     string                `json:"type"`
	Date                     time.Time             `json:"date"`
	Amount                   float64               `json:"amount"`
	DeductibleAmount         float64               `json:"deductibleAmount"`
	FeesPaid                 float64               `json:"feesPaid"`
	GatewayTransactionID     string                `json:"gatewayTransactionId"`
	GatewayAuthorizationCode string                `json:"gatewayAuthorizationCode"`
	ActivityFormName         string                `json:"activityFormName"`
	CustomFieldValues        []CustomFieldValue    `json:"customFieldValues"`
	Supporter                Supporter             `json:"supporter"`
}

//DonationUpsertResponsePayload (argh) contains the upsert results.
type DonationUpsertResponsePayload struct {
	Donations []ResultDonation `json:"donations"`
}
