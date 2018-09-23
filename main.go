package goengage

//EngEnv is the Engage environment.
type EngEnv struct {
	Host  string
	Token string
}

const (
	//FragMetrics is used to retrieve runtime metrics.
	FragMetrics = "/api/integration/ext/v1/metrics"
	//SupSearch is used to search for supporters.
	SupSearch = "/api/integration/ext/v1/supporters/search"
	//UatHost is the hostname for Engage instances on the test server.
	UatHost = "hq.uat.igniteaction.net"
	//ProdHost is the hostname for Engage instances on the production server.
	ProdHost = "api.salsalabs.org/"
)

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
	Header struct {
		RefID string `json:"refId"`
	} `json:"header"`
}

//NetOp is the wrapper for calls to Engage.  Here to keep
//call complexity down.
type NetOp struct {
	Host     string
	Fragment string
	Token    string
	Request  interface{}
	Response interface{}
}

//RequestBase is the common structure for a request.
type RequestBase struct {
	Header  Header      `json:"header"`
	Payload interface{} `json:"payload"`
}
