package goengage

import (
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

//FirstEmail returns the first email address for the provided supporter.
//Returns nil if the supporter does not have an email.  (As if...)
func FirstEmail(s Supporter) *string {
	c := s.Contacts
	if c == nil || len(c) == 0 {
		return nil
	}
	for _, x := range c {
		if x.Type == "EMAIL" {
			email := x.Value
			return &email
		}
	}
	return nil
}

//Date parses an Engage date and returns a Go time.
func Date(s string) time.Time {
	p := strings.Replace(time.RFC3339Nano, "9999999Z07:00", "Z", -1)
	t, err := time.Parse(p, s)
	if err != nil {
		panic(err)
	}
	return t
}

//Credentials reads a YAML file containing an Engage API host
//and an Engage API token.  These are then stored into an
//environment object.
func Credentials(fn string) (*Environment, error) {
	if len(fn) == 0 {
		return nil, errors.New("A configuration file is *required*.")
	}
	var c struct {
		Token string `json:"token"`
		Host  string `json:"host"`
	}
	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	e := NewEnvironment(c.Host, c.Token)
	return &e, nil
}
