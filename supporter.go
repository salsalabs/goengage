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
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
	County       string
	Country      string
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
