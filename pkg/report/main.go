package goengage

import (
	"fmt"
	"log"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

const (
	//BriefFormat is used to parse text into Classic-looking time.
	BriefFormat = "2006-01-02"

	//StartDuration is text to initialize a duration for start times.
	//Used in converting Go time strings to Engage times.
	StartDuration = "0h0m0.0s"

	//EndDuration is text to initialize a duration for end times.å
	//Used in converting Go time strings to Engage times.
	EndDuration = "23h59m59.999s"

	//DayDuration is used to scan a Span for new months.
	DayDuration = "24h"

	//BackupDuration is used to back a date up to the last
	//millisecond in the previous day.
	BackupDuration = "-1ms"
)

// Guide provides the basic tools to read and filter records then
// write them to a CSV file.
type Guide interface {
	//TypeActivity returns the kind of activity being read.
	TypeActivity() string

	//Filter returns true if the record should be used.
	Filter(goengage.Fundraise) bool

	//Headers returns column headers for a CSV file.
	Headers() []string

	//Line returns a list of strings to go in to the CSV file for each
	//fundraising record.
	Line(goengage.Fundraise) []string

	//Readers returns the number of readers to start.
	Readers() int

	//Filename returns the CSV filename.
	Filename() string

	//Location returns the location used to adjust transactions.
	//Transactions are Zulu.  Timezone is used to covert them to local.
	Location() *time.Location

	//Offset() returns the offset to start reading.  Useful for
	//restarting after a service interruption.
	Offset() int32
}

// Span is a pair of Time objects for the start and end of a time span.
type Span struct {
	S time.Time
	E time.Time
}

// TimeSpan contains a start and end time in Engage time format.
type TimeSpan struct {
	Start string
	End   string
}

// NewTimeSpan creates a Timespan using two Time objects.
func NewTimeSpan(s, e time.Time) TimeSpan {
	return TimeSpan{
		Start: s.Format(goengage.EngageDateFormat),
		End:   e.Format(goengage.EngageDateFormat),
	}
}

// Parse accepts a date in BriefFormat and returns a Go time. Engage
// needs a date and time.  Parameter "todText" defines the time to add.
// Errors are internal and fatal.
func Parse(s string, loc *time.Location, todText string) time.Time {
	t, err := time.ParseInLocation(BriefFormat, s, loc)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	d, err := time.ParseDuration(todText)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	t = t.Add(d)

	// Engage wants Zulu time.
	// TODO: handle positive offsets correctly.
	_, offset := t.Zone()
	zt := fmt.Sprintf("%ds", -offset)
	d, _ = time.ParseDuration(zt)
	t = t.Add(d)
	return t
}

// ValidateSpan validates the provided start and end dates.
// Errors are internal and fatal.
func ValidateSpan(startDate string, endDate string, loc *time.Location) Span {
	st := Parse(startDate, loc, StartDuration)
	et := Parse(endDate, loc, EndDuration)

	if et.Before(st) {
		log.Fatalf("end date '%v' is before start date '%v'", startDate, endDate)
	}
	span := Span{st, et}
	return span
}
