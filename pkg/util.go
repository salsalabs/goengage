package goengage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
func Date(s string) (t *time.Time) {
	if len(p) == 0 {
		return t
	}
	p := strings.Replace(time.RFC3339Nano, "9999999Z07:00", "Z", -1)
	t, err := time.Parse(p, s)
	if err != nil {
		panic(err)
	}
	return &t
}

//UtilLogger is an environment to support a file logger.  It contains
//a log.Logger attached to a file.
type UtilLogger struct {
	File   *os.File
	Logger *log.Logger
}

//NewUtilLogger creates a file and attaches a logger to it.  The file is generic
//looking with the date-time that this object was created.
func NewUtilLogger() (*UtilLogger, error) {
	u := UtilLogger{}
	now := time.Now()
	t := now.Format("2006-01-02T15:04:05")
	t = fmt.Sprintf("%v_log.txt", t)
	f, err := os.OpenFile(t, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return &u, err
	}
	u.File = f
	u.Logger = log.New(u.File, "", 0)
	return &u, err
}

//LogJSON is used to write the contents of a byte slice to the log
//as formatted JSON.  Note that no writes are performed if the Logger
//object hasn't been initialized.
func (u *UtilLogger) LogJSON(b []byte) {
	if u.Logger != nil {
		var x interface{}
		_ = json.Unmarshal(b, &x)
		t, _ := json.MarshalIndent(x, "", "\t")
		u.Logger.Println(string(t))
	}
}
