package goengage

const (
	//UatHost is the hostname for Engage instances on the test server.
	UatHost = "hq.uat.igniteaction.net"
	//ProdHost is the hostname for Engage instances on the production server.
	ProdHost = "api.salsalabs.org"
	//ContentType is always Javascript.
	ContentType = "application/json"
	//SearchMethod is always "POST" in Engage.
	SearchMethod = "POST"
)

//EngEnv is the Engage environment.
type EngEnv struct {
	Host  string
	Token string
}

//Error is used to report validation and input errors.
type Error struct {
	ID        string
	Code      int
	Message   string
	Details   string
	FieldName string
}

//Header contains an optional refID.
type Header struct {
	RefID string `json:"refId,omitempty"`
}

//RequestBase is the common structure for a request.
type RequestBase struct {
	//Header  Header      `json:"header,omitempty"`
	Payload interface{} `json:"payload"`
}
