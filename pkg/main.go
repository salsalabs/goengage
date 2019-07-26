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
	SearchMethod = "POST"
)

//EngEnv is the Engage environment.
type EngEnv struct {
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

//Header contains an optional refID.
type Header struct {
	RefID string `json:"refId,omitempty"`
}

//RequestBase is the common structure for a request.
type RequestBase struct {
	//Header  Header      `json:"header,omitempty"`
	Payload interface{} `json:"payload"`
}

//NewEngEnv creates a new EngEnv and initializes the metrics.
//Panics if updating the metrics returns an error.
func NewEngEnv(h string, t string) EngEnv {
	e := EngEnv{
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
func (e *EngEnv) UpdateMetrics() error {
	var resp MetResponse
	n := NetOp{
		Host:     e.Host,
		Fragment: FragMetrics,
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
