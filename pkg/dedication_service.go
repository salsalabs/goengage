package goengage

import (
	"fmt"
)

//DedicationService is the Service proxy for a Fundraise record.
type DedicationService = Fundraise

//ActivityType returns the kind of activity being read.
func (f DedicationService) ActivityType() string {
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
		"Dedication",
	}
}

//Line returns a list of strings to go in to the CSV file.
func (f DedicationService) Line() []string {
	return []string{
		f.PersonName,
		f.PersonEmail,
		f.Supporter.Address.AddressLine1,
		f.Supporter.Address.AddressLine2,
		f.Supporter.Address.City,
		f.Supporter.Address.State,
		f.Supporter.Address.PostalCode,
		fmt.Sprintf("%s", f.ActivityDate),
		fmt.Sprintf("%.2f", f.OneTimeAmount),
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
