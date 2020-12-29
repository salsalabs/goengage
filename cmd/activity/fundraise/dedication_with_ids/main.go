package main

//Application scan for fundraising activities with dedications
//and write them to a CSV.
import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	//DedicationAddressName is the supporter custom field name that contains
	//the dedication address.
	DedicationAddressName = "Address of Recipient to Notify"
)

//DedicationGuide is the Guide proxy for a Fundraise record.
type DedicationGuide struct{}

//NewDedicationGuide returns an record.
func NewDedicationGuide() DedicationGuide {
	e := DedicationGuide{}
	return e
}

//WhichActivity returns the kind of activity being read.
func (g DedicationGuide) WhichActivity() string {
	return goengage.FundraiseType
}

//Filter returns true if the record should be used.
func (g DedicationGuide) Filter(f goengage.Fundraise) bool {
	return len(f.Dedication) > 0
}

//Headers returns column headers for a CSV file.
func (g DedicationGuide) Headers() []string {
	return []string{
		"ActivityID",
		"DonationID",
		"SupporterID",
		"PersonName",
		"PersonEmail",
		"AddressLine1",
		"AddressLine2",
		"City",
		"State",
		"Zip",
		"TransactionDate",
		"Amount",
		"DedicationType",
		"Dedication",
		"Notify",
		"DedicationAddress",
	}
}

//Line returns a list of strings to go in to the CSV file.
func (g DedicationGuide) Line(f goengage.Fundraise) []string {
	// log.Printf("Line: %+v", f)
	addressLine1 := ""
	addressLine2 := ""
	city := ""
	state := ""
	postalCode := ""
	dedicationAddress := ""
	dedication := strings.Replace(f.Dedication, "\n", " ", -1)
	dedication = strings.Replace(dedication, "\r", " ", -1)
	dedication = strings.Replace(dedication, "\t", " ", -1)

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
					break
				}
			}
		}
	}
	return []string{
		f.ActivityID,
		f.DonationID,
		f.SupporterID,
		f.PersonName,
		f.PersonEmail,
		addressLine1,
		addressLine2,
		city,
		state,
		postalCode,
		fmt.Sprintf("%s", f.ActivityDate),
		fmt.Sprintf("%.2f", f.TotalReceivedAmount),
		f.DedicationType,
		f.Dedication,
		f.Notify,
		dedicationAddress,
	}
}

//Readers returns the number of readers to start.
func (g DedicationGuide) Readers() int {
	return 5
}

//Filename returns the CSV filename.
func (g DedicationGuide) Filename() string {
	return "dedications.csv"
}

const (
	//TimeFormat is used to parse text into Go time.
	TimeFormat = "2006-01-02"

	//StartDuration is text to initialize a duration for start times.
	//Used in converting Go time strings to Engage times.
	StartDuration = "0h0m0.0s"

	//EndDuration is text to initialize a duration for end times.
	//Used in converting Go time strings to Engage times.
	EndDuration = "23h59m59.999s"
)

//DefaultDates computes the default start and end dates.
//Default end date is just before the most recent Monday at midnight.
//Default start date is the Monday before the end date at 00:00.
//Formatted like Classic, "YYYY-MM-DD".
func DefaultDates() (start, end string) {
	now := time.Now()
	startDelta := 6 + int(now.Weekday())
	startTime := now.AddDate(0, 0, -startDelta)
	endTime := startTime.AddDate(0, 0, 6)
	start = startTime.Format(TimeFormat)
	end = endTime.Format(TimeFormat)
	return start, end
}

//Parse accepts a date in TimeFormat and returns a Go time. The
// duration text defines the time-of-day to add to the date
//before converting to Engage format. Errors are internal and fatal.
func Parse(s string, loc *time.Location, durationText string) time.Time {
	t, err := time.ParseInLocation(TimeFormat, s, loc)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	d, err := time.ParseDuration(StartDuration)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	t = t.Add(d)
	// Engage wants Zulu time.
	_, offset := t.Zone()
	zt := fmt.Sprintf("%ds", -offset)
	d, err = time.ParseDuration(zt)
	t = t.Add(d)
	return t
}

//Validate validates the provided start and end dates.  Errors are fatal.
//Converts the dates from the provided timezone to Zulu, then formats them
//suitably for submission to Engage.  Return the validated and formatted dates.
func Validate(startDate, endDate, timeZone string) (string, string) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	st := Parse(startDate, loc, StartDuration)
	et := Parse(endDate, loc, EndDuration)

	if et.Before(st) {
		s := st.Format(TimeFormat)
		e := et.Format(TimeFormat)
		log.Fatalf("end date '%v' is before start date '%v'", s, e)
	}
	s := st.Format(goengage.EngageDateFormat)
	e := et.Format(goengage.EngageDateFormat)
	return s, e
}

func main() {
	start, end := DefaultDates()
	var (
		app       = kingpin.New("dedications", "Write dedications to a CSV")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		startDate = app.Flag("startDate", "Start date, YYYY-MM-YY, default is Monday of last week at midnight").Default(start).String()
		endDate   = app.Flag("endDate", "End date, YYYY-MM-YY, default is the most recent Monday at midnight").Default(end).String()
		timeZone  = app.Flag("timezone", "Client's timezone, defaults to EST/EDT").Default("America/New_York").String()
	)
	app.Parse(os.Args[1:])

	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("%v", err)
	}
	engageStart, engageEnd := Validate(*startDate, *endDate, *timeZone)

	guide := NewDedicationGuide()
	err = goengage.ReportFundraising(e, guide, engageStart, engageEnd)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
