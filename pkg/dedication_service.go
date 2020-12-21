package goengage

import (
	"fmt"
)

//DedicationService is the Service proxy for a Fundraise record.
type DedicationService = Fundraise

//NewDedicationService returns an record.
func NewDedicationService() DedicationService {
	f := Fundraise{}
	return f
}

//WhichActivity returns the kind of activity being read.
func (f DedicationService) WhichActivity() string {
	return FundraiseType
}

//Filter returns true if the record should be used.
func (f DedicationService) Filter() bool {
	return len(f.Dedication) > 0
}

//Headers returns column headers for a CSV file.
func (f DedicationService) Headers() []string {
	return []string{
		"PersonName",
		"PersonEmail",
		"AddressLine1",
		"AddressLine2",
		"City",
		"State",
		"Zip",
		"TransactionDate",
		"Amount",
		"DedicationType",
		"Dedication",
	}
}

//Line returns a list of strings to go in to the CSV file.
func (f DedicationService) Line() []string {
	// log.Printf("Line: %+v", f)
	addressLine1 := ""
	addressLine2 := ""
	city := ""
	state := ""
	postalCode := ""
	s := &f.Supporter
	if s == nil {
		addressLine1 = f.Supporter.Address.AddressLine1
		addressLine2 = f.Supporter.Address.AddressLine2
		city = f.Supporter.Address.City
		state = f.Supporter.Address.State
		postalCode = f.Supporter.Address.PostalCode
	}
	return []string{
		f.PersonName,
		f.PersonEmail,
		addressLine1,
		addressLine2,
		city,
		state,
		postalCode,
		fmt.Sprintf("%s", f.ActivityDate),
		fmt.Sprintf("%.2f", f.TotalReceivedAmount),
		f.DedicationType,
		f.Dedication,
	}
}

//Readers returns the number of readers to start.
func (f DedicationService) Readers() int {
	return 5
}

//Filename returns the CSV filename.
func (f DedicationService) Filename() string {
	return "dedications.csv"
}
