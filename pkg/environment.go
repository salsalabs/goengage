package goengage

import (
	"errors"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Environment is the Engage environment.
type Environment struct {
	Host    string
	Token   string
	Metrics Metrics
}

// NewEnvironment creates a new Environment and initializes the metrics.
// Panics if updating the metrics returns an error.
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

// Credentials reads a YAML file containing an Engage API host
// and an Engage API token.  These are then stored into an
// environment object.
func Credentials(fn string) (*Environment, error) {
	if len(fn) == 0 {
		return nil, errors.New(" configuration file is *required*")
	}
	var c struct {
		Token string `json:"token"`
		Host  string `json:"host"`
	}
	raw, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	if len(c.Host) == 0 {
		c.Host = APIHost
	}
	e := NewEnvironment(c.Host, c.Token)
	return &e, nil
}

// UpdateMetrics reads metrics and returns them.
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
	e.Metrics = resp.Payload
	return nil
}
