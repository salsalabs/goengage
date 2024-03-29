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
	//SeeAddressName is the supporter custom field name that contains
	//the dedication address.
	SeeAddressName = "Address of Recipient to Notify"

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

//SeeGuide is the Guide proxy.
type SeeGuide struct {
	Span         report.Span
	Timezone     *time.Location
	DonationType string
	ReadOffset   int32
}

//NewSeeGuide returns an initialized SeeGuide.
func NewSeeGuide(span report.Span, location *time.Location, donationType string, readOffset int32) SeeGuide {
	return SeeGuide{
		Span:         span,
		Timezone:     location,
		DonationType: donationType,
		ReadOffset:   readOffset,
	}
}

//TypeActivity returns the kind of activity being read.
//Implements goengage.report.Guide.
func (g SeeGuide) TypeActivity() string {
	return goengage.FundraiseType
}

//Filter returns true if the record should be used.
//Implements goengage.report.Guide.
func (g SeeGuide) Filter(f goengage.Fundraise) bool {
	switch g.DonationType {
	case "All":
		return true
	case goengage.OneTime:
		return f.DonationType == goengage.OneTime
	case goengage.Recurring:
		return f.DonationType == goengage.Recurring
	}
	return false
}

//Headers returns column headers for a CSV file.
//Implements goengage.report.Guide.
func (g SeeGuide) Headers() []string {
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
		"Fund",
		"Campaign",
		"Appeal",
		"Designation",
		"DedicationType",
		"Dedication",
		"TrackingCode",
	}
	return a
}

//Line returns a list of strings to go in to the CSV file.
//Implements goengage.report.Guide.
func (g SeeGuide) Line(f goengage.Fundraise) []string {
	activityDate := f.ActivityDate.In(g.Location())
	transactionDate := activityDate.Format(BriefFormat)
	h := g.Headers()
	var a []string
	for _, s := range h {
		var b string
		switch s {
		case "SupporterID":
			b = f.SupporterID
		case "FirstName":
			b = f.Supporter.FirstName
		case "LastName":
			b = f.Supporter.LastName
		case "PersonEmail":
			b = f.PersonEmail
		case "TransactionDate":
			b = transactionDate
		case "DonationType":
			b = f.DonationType
		case "DonationID":
			b = f.DonationID
		case "ActivityType":
			b = f.ActivityType
		case "ActivityID":
			b = f.ActivityID
		case "TransactionType":
			b = goengage.ToTitle(f.Transactions[0].Type)
		case "TransactionID":
			b = f.Transactions[0].TransactionID
		case "Amount":
			b = fmt.Sprintf("%.2f", f.TotalReceivedAmount)
		case "Fund":
			b = f.Fund
		case "Campaign":
			b = f.Campaign
		case "Designation":
			b = f.Designation
		case "Appeal":
			b = f.Appeal
		case "DedicationType":
			b = f.DedicationType
		case "Dedication":
			b = f.Dedication
		case "TrackingCode":
			b = f.TrackingCode
		}
		a = append(a, b)
	}
	return a
}

//Location returns the local location. Useful for date conversions.
func (g SeeGuide) Location() *time.Location {
	return g.Timezone
}

//Readers returns the number of readers to start.
func (g SeeGuide) Readers() int {
	return ReaderCount
}

//Filename returns the CSV filename.
func (g SeeGuide) Filename() string {
	s := g.Span.S.Format(BriefFormat)
	t := strings.ToLower(g.DonationType)
	return fmt.Sprintf("%s_see_%s.csv", s, t)
}

//Offset returns the offset for the first read.
//Useful for restarting after a service interruption.
func (g SeeGuide) Offset() int32 {
	return g.ReadOffset
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

//ValidateDonationType returns an error if the provided
//donation type is invalid.
func ValidateDonationType(d string) error {
	switch d {
	case "All":
		return nil
	case goengage.ToTitle(goengage.OneTime):
		return nil
	case goengage.ToTitle(goengage.Recurring):
		return nil
	}
	return fmt.Errorf("Not a valid donation type, '%s'", d)

}
func main() {
	start, end := DefaultDates()
	donationTypePrompt := fmt.Sprintf("Choose All, %s or %s", goengage.ToTitle(goengage.OneTime), goengage.ToTitle(goengage.Recurring))
	var (
		app          = kingpin.New("see", "Write all donations for a timeframe to a CSV")
		login        = app.Flag("login", "YAML file with API token").Required().String()
		startDate    = app.Flag("startDate", "Start date, YYYY-MM-YY, default is Monday of last week at midnight").Default(start).String()
		endDate      = app.Flag("endDate", "End date, YYYY-MM-YY, default is the most recent Monday at midnight").Default(end).String()
		timeZone     = app.Flag("timezone", "Client's timezone, defaults to EST/EDT").Default("America/New_York").String()
		donationType = app.Flag("donationType", donationTypePrompt).Default("All").String()
		readOffset   = app.Flag("readOffset", "Read reading here, useful for restarts").Default("0").Int32()
	)
	app.Parse(os.Args[1:])

	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = ValidateDonationType(*donationType)
	if err != nil {
		log.Fatalf("%v", err)
	}
	location, err := time.LoadLocation(*timeZone)
	if err != nil {
		log.Fatalf("%v", err)
	}
	span := report.ValidateSpan(*startDate, *endDate, location)
	guide := NewSeeGuide(span, location, *donationType, *readOffset)
	ts := report.NewTimeSpan(span.S, span.E)
	err = report.ReportFundraising(e, guide, ts)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
