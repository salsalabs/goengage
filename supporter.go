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

//SupporterHeader contains an optional refID.
type SupporterHeader struct {
	Header struct {
		RefID string `json:"refId"`
	} `json:"header"`
}

//SupSearchRequest is used to ask for supporters.
type SupSearchRequest struct {
	Payload struct {
		ModifiedFrom   string `json:"modifiedFrom"`
		ModifiedTo     string `json:"modifiedTo"`
		Offset         int32
		Count          int32
		Identifiers    []string `json:"identifiers"`
		IdentifierType string   `json:"identifierType"`
	} `json:"payload"`
}

//SupUpdateRequest is a request to change/insert a supporter.
type SupUpdateRequest struct {
	Payload struct {
		Supporters []Supporter `json:"supporters"`
	} `json:"payload"`
}

//SupSearchResult is returned when supporters are found by a search.
type SupSearchResult struct {
	Payload struct {
		Count      string      `json:"count"`
		Offset     string      `json:"offset"`
		Total      string      `json:"total"`
		Supporters []Supporter `json:"suporters"`
	}
}
