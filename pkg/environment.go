package goengage

import (
	"fmt"
	"net/http"
)

const (
	//UATHost is the hostname for Engage instances on the test server.
	UATHost = "hq.uat.igniteaction.net"
	//APIHost is the hostname for Engage instances on the production server.
	APIHost = "api.salsalabs.org"
	//ContentType is always Javascript.
	ContentType = "application/json"
	//SearchMethod is always "POST" in Engage.
	SearchMethod = http.MethodPost
)

//Environment is the Engage environment.
type Environment struct {
	Host    string
	Token   string
	Metrics Metrics
}

//Error is used to report Engage errors.
type Error struct {
	ID        string `json:"id,omitempty"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   string `json:"details,omitempty"`
	FieldName string `json:"fieldName,omitempty"`
}

//Request is the common structure for a request.  YOur request object
//gets stored in Payload automatically by net.Do().
type Request struct {
	Header struct {
		RefID string `json:"refId,omitempty"`
	} `json:"header,omitempty"`
	Payload interface{} `json:"payload"`
}

//Response is the common structure for a request.  YOur request object
//gets stored in Payload automatically by net.Do().
type Response struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Header    struct {
		ProcessingTime int    `json:"processingTime"`
		ServerID       string `json:"serverId"`
	} `json:"header,omitempty"`
	Errors  []Error     `json:"errors,omitempty"`
	Payload interface{} `json:"payload"`
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
	var resp MetricsResponse
	n := NetOp{
		Host:     e.Host,
		Endpoint: MetricsCommand,
		Method:   http.MethodGet,
		Token:    e.Token,
		Request:  nil,
		Response: &resp,
	}
	err := n.Do()
	if err != nil {
		return err
	}
	fmt.Printf("Updated metrics: %+v\n", resp)
	e.Metrics = resp.Payload
	return nil
}
