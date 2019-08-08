package goengage

import (
	"fmt"
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
	ID        string `json:"id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   string `json:"details,omitempty"`
	FieldName string `json:"fieldName,omitempty"`
}

//RequestHeader contains an optional refID.
type RequestHeader struct {
	RefID string `json:"refId,omitempty"`
}

//Request is the common structure for a request.  YOur request object
//gets stored in Payload automatically by net.Do().
type Request struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload interface{}   `json:"payload"`
}

//ResponseHeader is the common object returned by calls to Engage.
//Payloads are defined by the objects receiving the data, since they
//need the payload.
type ResponseHeader struct {
	ProcessingTime int    `json:"processingTime,omitempty"`
	ServerID       string `json:"serverID,omitempty"`
}

//Response is the common structure for a request.  YOur request object
//gets stored in Payload automatically by net.Do().
type Response struct {
	ID        string         `json:"id,omitempty"`
	Timestamp string         `json:"timestamp"`
	Header    ResponseHeader `json:"header,omitempty"`
	Errors    []Error        `json:"errors,omitempty"`
	Payload   interface{}    `json:"payload,omitempty"`
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
	var resp MetricData
	n := NetOp{
		Host:     e.Host,
		Endpoint: MetricsCommand,
		Method:   http.MethodGet,
		Token:    e.Token,
		Request:  nil,
		Response: &resp,
	}
	fmt.Printf("Request: %+v\n", n)
	err := n.Do()
	if err != nil {
		return err
	}
	e.Metrics = resp
	return nil
}
