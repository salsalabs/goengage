package goengage

//Custom is a custom field in Engage.
type Custom struct {
	FieldID    string
	Name       string
	Value      string
	Type       string
	OptInDate  string
	OptOutDate string
}

//Contact describes a way to contact a supporter.
type Contact struct {
	Type   string
	Value  string
	Status string
}

//Address is a geographic locaiton for a supporter.
type Address struct {
	AddressLine1         string
	AddressLine2         string
	City                 string
	State                string
	PostalCode           string
	County               string
	Country              string
	FederalDistrict      string
	StateHouseDistrict   string
	StateSenateDistrict  string
	CountyDistrict       string
	MunicipalityDistrict string
	Lattitude            float32
	Longitude            float32
	Status               string
}

//Supporter is a supporter from the database or being saved to the database.
type Supporter struct {
	SupporterID       string
	Result            string
	Title             string
	FirstName         string
	MiddleName        string
	LastName          string
	Suffix            string
	DateOfBirth       string
	Gender            string
	CreatedDate       string
	LastModified      string
	ExternalSystemID  string
	Address           Address
	Contacts          []Contact
	CustomFieldValues []Custom
}

//SupporterPayload is returned when supporters are found by a search.
type SupporterPayload struct {
	Payload struct {
		Count      string
		Offset     string
		Total      string
		supporters []Supporter
	}
}
