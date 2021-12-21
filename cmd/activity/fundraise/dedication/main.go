package main

//Application scan for fundraising activities with dedications
//and write them to a CSV.
import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	report "github.com/salsalabs/goengage/pkg/report"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	//DedicationAddressName is the supporter custom field name that contains
	//the dedication address.
	DedicationAddressName = "Address of Recipient to Notify"

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
)

//DedicationGuide is the Guide proxy.
type DedicationGuide struct {
	Span     report.Span
	AddKeys  bool
	Timezone *time.Location
}

//NewDedicationGuide returns an initialized DedicationGuide.
func NewDedicationGuide(span report.Span, addKeys bool, location *time.Location) DedicationGuide {
	return DedicationGuide{
		Span:     span,
		AddKeys:  addKeys,
		Timezone: location,
	}
}

//TypeActivity returns the kind of activity being read.
//Implements goengage.report.Guide.
func (g DedicationGuide) TypeActivity() string {
	return goengage.FundraiseType
}

//Filter returns true if the record should be used.
//Implements goengage.report.Guide.
func (g DedicationGuide) Filter(f goengage.Fundraise) bool {
	return f.DedicationType != goengage.None && !f.ActivityDate.Before(g.Span.S) && !f.ActivityDate.After(g.Span.E)
}

//Headers returns column headers for a CSV file.
//Implements goengage.report.Guide.
func (g DedicationGuide) Headers() []string {
	a := []string{
		"FirstName",
		"LastName",
		"PersonEmail",
		"AddressLine1",
		"AddressLine2",
		"City",
		"State",
		"Zip",
		"TransactionDate",
		"DonationType",
		"ActivityType",
		"TransactionType",
		"Amount",
		"DedicationType",
		"Dedication",
		"Notify",
		"DedicationAddress",
	}
	if g.AddKeys {
		a = append(a, "ActivityID")
		a = append(a, "DonationID")
		a = append(a, "TransactionID")
		a = append(a, "SupporterID")
	}
	return a
}

//Line returns a list of strings to go in to the CSV file.
//Implements goengage.report.Guide.
func (g DedicationGuide) Line(f goengage.Fundraise) []string {
	addressLine1 := ""
	addressLine2 := ""
	city := ""
	state := ""
	postalCode := ""
	dedicationAddress := ""
	dedication := strings.Replace(f.Dedication, "\n", " ", -1)
	dedication = strings.Replace(dedication, "\r", " ", -1)
	dedication = strings.Replace(dedication, "\t", " ", -1)
	activityDate := f.ActivityDate.In(g.Location())
	transactionDate := activityDate.Format(BriefFormat)

	s := &f.Supporter
	if s != nil {
		if f.Supporter.Address != nil {
			addressLine1 = f.Supporter.Address.AddressLine1
			addressLine2 = f.Supporter.Address.AddressLine2
			city = f.Supporter.Address.City
			state = f.Supporter.Address.State
			postalCode = f.Supporter.Address.PostalCode
		}
		if f.Supporter.CustomFieldValues != nil {
			for _, c := range f.Supporter.CustomFieldValues {
				if c.Name == DedicationAddressName {
					dedicationAddress = strings.Replace(c.Value, "\n", " ", -1)
					dedicationAddress = strings.Replace(dedicationAddress, "\r", " ", -1)
					dedicationAddress = strings.Replace(dedicationAddress, "\t", " ", -1)
					re := regexp.MustCompile("[\r\n\t ]+")
					dedicationAddress = re.ReplaceAllString(dedicationAddress, " ")
					break
				}
			}
		}
	}
	a := []string{
		f.Supporter.FirstName,
		f.Supporter.LastName,
		f.PersonEmail,
		addressLine1,
		addressLine2,
		city,
		state,
		postalCode,
		transactionDate,
		goengage.ToTitle(f.DonationType),
		goengage.ToTitle(f.ActivityType),
		goengage.ToTitle(f.Transactions[0].Type),
		fmt.Sprintf("%.2f", f.TotalReceivedAmount),
		goengage.ToTitle(f.DedicationType),
		f.Dedication,
		f.Notify,
		dedicationAddress,
	}
	if g.AddKeys {
		a = append(a, f.ActivityID)
		a = append(a, f.DonationID)
		a = append(a, f.Transactions[0].TransactionID)
		a = append(a, f.SupporterID)
	}
	return a
}

//Location returns the local location. Useful for date conversions.
func (g DedicationGuide) Location() *time.Location {
	return g.Timezone
}

//Readers returns the number of readers to start.
func (g DedicationGuide) Readers() int {
	return 3
}

//Filename returns the CSV filename.
func (g DedicationGuide) Filename() string {
	s := g.Span.S.Format(BriefFormat)
	return fmt.Sprintf("%s_dedications.csv", s)
}

//Offset returns the starting offset for the first read.
func (g DedicationGuide) Offset() int32 {
	return int32(0)
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

// Validate validates the provided start and end dates.
// Converts the dates from the provided location to Zulu, checks for start
// time before end time, then returns a slice of Span objects.  Typically,
// the Slice is 1 entry.  It becomes multiple entries when interval between
// startDate and endDate crosses month boundaries.
// Errors are internal and fatal.
func Validate(startDate string, endDate string, loc *time.Location) []report.Span {
	st := report.Parse(startDate, loc, StartDuration)
	et := report.Parse(endDate, loc, EndDuration)

	if et.Before(st) {
		log.Fatalf("end date '%v' is before start date '%v'", endDate, startDate)
	}
	var a []report.Span
	day, _ := time.ParseDuration(DayDuration)
	yesterday, err := time.ParseDuration(BackupDuration)
	if err != nil {
		panic(err)
	}
	for ct := st; ct.Before(et); ct = ct.Add(day) {
		if ct.Month() != st.Month() {
			span := report.Span{st, ct.Add(yesterday)}
			a = append(a, span)
			st = ct
		}
	}
	span := report.Span{st, et}
	a = append(a, span)
	return a
}

func main() {
	start, end := DefaultDates()
	var (
		app       = kingpin.New("dedications", "Write dedications to a CSV")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		startDate = app.Flag("startDate", "Start date, YYYY-MM-YY, default is Monday of last week at midnight").Default(start).String()
		endDate   = app.Flag("endDate", "End date, YYYY-MM-YY, default is the most recent Monday at midnight").Default(end).String()
		timeZone  = app.Flag("timezone", "Client's timezone, defaults to EST/EDT").Default("America/New_York").String()
		addKeys   = app.Flag("keys", "Export activity, donation, transaction and supporter IDs").Bool()
	)
	app.Parse(os.Args[1:])

	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("%v", err)
	}
	location, err := time.LoadLocation(*timeZone)
	spans := Validate(*startDate, *endDate, location)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, span := range spans {
		guide := NewDedicationGuide(span, *addKeys, location)
		ts := report.NewTimeSpan(span.S, span.E)
		err = report.ReportFundraising(e, guide, ts)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
