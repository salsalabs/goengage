package main

//Application to scan for fundraising activities and write them to a CSV.
import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	report "github.com/salsalabs/goengage/pkg/report"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	//SupporterAddressName is the supporter custom field name that contains
	//the dedication address.
	SupporterAddressName = "Address of Recipient to Notify"

	//BriefFormat is used to parse text into Classic-looking time.
	BriefFormat = "2006-01-02"

	//StartDuration is text to initialize a duration for start times.
	//Used in converting Go time strings to Engage times.
	StartDuration = "0h0m0.0s"

	//EndDuration is text to initialize a duration for end times.
	//Used in converting Go time strings to Engage times.
	EndDuration = "23h59m59.999s"

	//DayDuration is used to scan a Span for new months.
	DayDuration = "24h"

	//BackupDuration is used to back a date up to the last
	//millisecond in the previous day.
	BackupDuration = "-1ms"

	//ReaderCount is the number of Engage reaaders to start.
	ReaderCount = 3
)

//SupporterGuide is the Guide proxy.
type SupporterGuide struct {
	Span        Span
	Timezone    *time.Location
	SupporterID string
}

//NewSupporterGuide returns an initialized SupporterGuide.
func NewSupporterGuide(span Span, location *time.Location, supporterID string) SupporterGuide {
	return SupporterGuide{
		Span:        span,
		Timezone:    location,
		SupporterID: supporterID,
	}
}

//Span is a pair of Time objects for the start and end of a time span.
type Span struct {
	S time.Time
	E time.Time
}

//TypeActivity returns the kind of activity being read.
func (g SupporterGuide) TypeActivity() string {
	return goengage.FundraiseType
}

//Filter returns true if the record should be used.
func (g SupporterGuide) Filter(f goengage.Fundraise) bool {
	return f.SupporterID == g.SupporterID
}

//Headers returns column headers for a CSV file.
func (g SupporterGuide) Headers() []string {
	a := []string{
		"SupporterID",
		"FirstName",
		"LastName",
		"PersonEmail",
		"TransactionDate",
		"DonationType",
		"DonationID",
		"ActivityType",
		"ActivityID",
		"TransactionType",
		"TransactionID",
		"Amount",
	}
	return a
}

//Line returns a list of strings to go in to the CSV file.
func (g SupporterGuide) Line(f goengage.Fundraise) []string {
	activityDate := f.ActivityDate.In(g.Location())
	transactionDate := activityDate.Format(BriefFormat)

	a := []string{
		f.SupporterID,
		f.Supporter.FirstName,
		f.Supporter.LastName,
		f.PersonEmail,
		transactionDate,
		f.DonationType,
		f.DonationID,
		ToTitle(f.ActivityType),
		f.ActivityID,
		ToTitle(f.Transactions[0].Type),
		f.Transactions[0].TransactionID,
		fmt.Sprintf("%.2f", f.TotalReceivedAmount),
	}
	return a
}

//Location returns the local location. Useful for date conversions.
func (g SupporterGuide) Location() *time.Location {
	return g.Timezone
}

//Readers returns the number of readers to start.
func (g SupporterGuide) Readers() int {
	return ReaderCount
}

//Filename returns the CSV filename.
func (g SupporterGuide) Filename() string {
	s := g.Span.S.Format(BriefFormat)
	return fmt.Sprintf("%s_supporter.csv", s)
}

//DefaultDates computes the default start and end dates.
//Default end date is just before the most recent Monday at midnight.
//Default start date is the Monday before the end date at 00:00.
//Formatted like Classic, "YYYY-MM-DD".
func DefaultDates() (start, end string) {
	now := time.Now()
	startDelta := 6 + int(now.Weekday())
	startTime := now.AddDate(0, 0, -startDelta)
	endTime := startTime.AddDate(0, 0, 6)
	start = startTime.Format(BriefFormat)
	end = endTime.Format(BriefFormat)
	return start, end
}

//Parse accepts a date in BriefFormat and returns a Go time. Engage
//needs a date and time.  Parameter "todText" defines the time to add.
//Errors are internal and fatal.
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
	d, err = time.ParseDuration(zt)
	t = t.Add(d)
	return t
}

//ToTitle converts engage constants to title-case.  Underbars
//are treated as word separators.
func ToTitle(s string) string {
	parts := strings.Split(s, "_")
	var a []string
	for _, x := range parts {
		a = append(a, strings.Title(strings.ToLower(x)))
	}
	return strings.Join(a, "_")
}

// Validate validates the provided start and end dates. Errors are
// noisily fatal.  Returns a span.
func Validate(startDate string, endDate string, loc *time.Location) Span {
	st := Parse(startDate, loc, StartDuration)
	et := Parse(endDate, loc, EndDuration)

	if et.Before(st) {
		log.Fatalf("end date '%v' is before start date '%v'", startDate, endDate)
	}
	span := Span{st, et}
	return span
}

func main() {
	start, end := DefaultDates()
	var (
		app         = kingpin.New("see2", "Write donations to a CSV")
		login       = app.Flag("login", "YAML file with API token").Required().String()
		startDate   = app.Flag("startDate", "Start date, YYYY-MM-YY, default is Monday of last week at midnight").Default(start).String()
		endDate     = app.Flag("endDate", "End date, YYYY-MM-YY, default is the most recent Monday at midnight").Default(end).String()
		timeZone    = app.Flag("timezone", "Client's timezone, defaults to EST/EDT").Default("America/New_York").String()
		supporterID = app.Flag("supporterID", "Show donations for this supporter").Required().String()
	)
	app.Parse(os.Args[1:])

	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("%v", err)
	}
	location, err := time.LoadLocation(*timeZone)
	if err != nil {
		log.Fatalf("%v", err)
	}
	span := Validate(*startDate, *endDate, location)
	guide := NewSupporterGuide(span, location, *supporterID)
	ts := report.NewTimeSpan(span.S, span.E)
	err = report.ReportFundraising(e, guide, ts)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
