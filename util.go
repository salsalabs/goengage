package goengage

import (
	"strings"
	"time"
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
