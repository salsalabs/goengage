package goengage

//Custom is a custom field in Engage.
type Custom struct {
	FieldID    string `json:"fieldID:omitempty"`
	Name       string `json:"name:omitempty"`
	Value      string `json:"value:omitempty"`
	Type       string `json:"type:omitempty"`
	OptInDate  string `json:"optInDate:omitempty"`
	OptOutDate string `json:"optOutDate:omitempty"`
}

//Contact describes a way to contact a supporter.
type Contact struct {
	Type   string  `json:"type,omitempty"`
	Value  string  `json:"value,omitempty"`
	Status string  `json:"status,omitempty"`
	Errors []Error `json:"errors,omitempty"`
}

//Address is a geographic locaiton for a supporter.
type Address struct {
	AddressLine1         string  `json:"addressLine1,omitempty"`
	AddressLine2         string  `json:"addressLine2,omitempty"`
	City                 string  `json:"city,omitempty"`
	State                string  `json:"state,omitempty"`
	PostalCode           string  `json:"postalCode,omitempty"`
	County               string  `json:"county,omitempty"`
	Country              string  `json:"country,omitempty"`
	FederalDistrict      string  `json:"federalDistrict,omitempty"`
	StateHouseDistrict   string  `json:"stateHouseDistrict,omitempty"`
	StateSenateDistrict  string  `json:"stateSenateDistrict,omitempty"`
	CountyDistrict       string  `json:"countyDistrict,omitempty"`
	MunicipalityDistrict string  `json:"municipalityDistrict,omitempty"`
	Lattitude            float32 `json:"lattitude,omitempty"`
	Longitude            float32 `json:"longitude,omitempty"`
	Status               string  `json:"status,omitempty"`
	Errors               []Error `json:"errors,omitempty"`
}

//Supporter is a supporter from the database or being saved to the database.
type Supporter struct {
	SupporterID       string    `json:"supporterId,omitempty"`
	Result            string    `json:"result,omitempty"`
	Title             string    `json:"title,omitempty"`
	FirstName         string    `json:"firstName,omitempty"`
	MiddleName        string    `json:"middleName,omitempty"`
	LastName          string    `json:"lastName,omitempty"`
	Suffix            string    `json:"suffix,omitempty"`
	DateOfBirth       string    `json:"dateOfBirth,omitempty"`
	Gender            string    `json:"gender,omitempty"`
	CreatedDate       string    `json:"createdDate,omitempty"`
	LastModified      string    `json:"lastModified,omitempty"`
	ExternalSystemID  string    `json:"externalSystemId,omitempty"`
	Status            string    `json:"status,omitempty"`
	Address           *Address  `json:"address,omitempty"`
	Contacts          []Contact `json:"contacts,omitempty"`
	CustomFieldValues []Custom  `json:"customFieldValues,omitempty"`
	//Timezone          string    `json:"timezone,omitempty"`
	//LanguageCode      string    `json:"languageCode,omitempty"`
}

//SupSearchRequest is used to ask for supporters.
type SupSearchRequest struct {
	ModifiedFrom   string   `json:"modifiedFrom,omitempty"`
	ModifiedTo     string   `json:"modifiedTo,omitempty"`
	Offset         int32    `json:"offset,omitempty"`
	Count          int32    `json:"count,omitempty"`
	Identifiers    []string `json:"identifiers,omitempty"`
	IdentifierType string   `json:"identifierType,omitempty"`
}

//SupUpsertRequest is a request to change/insert a supporter.
type SupUpsertRequest struct {
	Supporters []Supporter `json:"supporters"`
}

//SupUpsertResult is returned after an upsert (add/modify)
type SupUpsertResult struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Header    struct {
		ProcessingTime int32  `json:"processingTime"`
		ServerID       string `json:"serverId"`
	} `json:"header"`
	Payload struct {
		Supporters []Supporter `json:"supporters"`
	} `json:"payload"`
}

//DeletingSupporters contains the list of supporters IDs that will be
//deleted
type DeletingSupporters struct {
	SupporterID string `json:"supporterId"`
}

//DeleteResults contains the list of supporters and reasons returned
//after supporters are deleted.
type DeleteResults struct {
	ReadOnly    bool   `json:"readOnly"`
	SupporterID string `json:"supporterId"`
	Result      string `json:"result"`
}

//SupDeleteRequest is a request to delete supporters.  Deleting
//supporters uses a (supporterId, result) duple.  Note that this
//works because empty fields are ignored.
type SupDeleteRequest struct {
	Supporters []DeletingSupporters `json:"supporters"`
}

//SupDeleteResult is the result of deleting a supporter.  Deleting
//supporters returns a (supporterId, result) duple.  Note that this
//works because empty fields are ignored.
type SupDeleteResult struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Header    struct {
		ProcessingTime int32  `json:"processingTime"`
		ServerID       string `json:"serverId"`
	} `json:"header"`
	Payload struct {
		Supporters []DeleteResults `json:"supporters"`
	} `json:"payload"`
}

//SupSearchResult is returned when supporters are found by a search.
type SupSearchResult struct {
	Payload struct {
		Count      int32       `json:"count"`
		Offset     int32       `json:"offset"`
		Total      int32       `json:"total"`
		Supporters []Supporter `json:"supporters"`
	} `json:"payload"`
}
