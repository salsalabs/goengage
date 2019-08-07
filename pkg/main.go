package goengage

import (
	"net/http"
)

const (
	//UatHost is the hostname for Engage instances on the test server.
	UatHost = "hq.uat.igniteaction.net"
	//ProdHost is the hostname for Engage instances on the production server.
	ProdHost = "api.salsalabs.org"
	//ContentType is always Javascript.
	ContentType = "application/json"
	//SearchMethod is always "POST" in Engage.
	SearchMethod = http.MethodPost
)

//Environment is the Engage environment.
type Environment struct {
	Host    string
	Token   string
	Metrics MetricData
}

//Error is used to report validation and input errors.
type Error struct {
	ID        string
	Code      int
	Message   string
	Details   string
	FieldName string
}

//RequestHeader contains an optional refID.
type RequestHeader struct {
	RefID string `json:"refId,omitempty"`
}

//RequestBase is the common structure for a request.  YOur request object
//gets stored in Payload automatically by net.Do().
type RequestBase struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload interface{}   `json:"payload"`
}

//ResponseHeader is the common object returned by calls to Engage.
//Payloads are defined by the objects receiving the data, since they
//need the payload.
type ResponseHeader struct {
	ProcessingTime string `json:"processingTime,omitempty"`
	ServerID       string `json:"serverID,omitempty"`
}

//ResponseBase is the common structure for a request.  YOur request object
//gets stored in Payload automatically by net.Do().
type ResponseBase struct {
	Header  ResponseHeader `json:"header,omitempty"`
	Payload interface{}    `json:"payload"`
}

//NewEnvironment creates a new Environment and initializes the metrics.
//Panics if updating the metrics returns an error.
func NewEnvironment(h string, t string) Environment {
	e := Environment{
		Host:  h,
		Token: t,
	}
	err := (&e).UpdateMetrics()
	if err != nil {
		panic(err)
	}
	return e
}

//UpdateMetrics reads metrics and returns them.
func (e *Environment) UpdateMetrics() error {
	var resp MetResponse
	n := NetOp{
		Host:     e.Host,
		Endpoint: FragMetrics,
		Method:   http.MethodGet,
		Token:    e.Token,
		Request:  nil,
		Response: &resp,
	}
	err := n.Do()
	if err != nil {
		return err
	}
	e.Metrics = resp.Payload
	return nil
}
