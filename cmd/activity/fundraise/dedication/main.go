package main

//Application scan for fundraising activities with dedications
//and write them to a CSV.
import (
	"fmt"
	"log"
	"os"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	//TimeFormat is used to parse text into Go time.
	TimeFormat = "2006-01-02"
)

func main() {
	// Default end date is the most recent Monday at midnight.
	// Default start date is the Monday before the end date.
	// Easy text format like Classic, "YYYY-MM-DD"
	now := time.Now()
	startDelta := 6 + int(now.Weekday())
	startTime := now.AddDate(0, 0, -startDelta)
	endTime := startTime.AddDate(0, 0, 6)
	start := startTime.Format(TimeFormat)
	end := endTime.Format(TimeFormat)

	var (
		app       = kingpin.New("dedications", "Write dedications to a CSV")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		startDate = app.Flag("startDate", "Start date, YYYY-MM-YY, default is Monday of last week at midnight").Default(start).String()
		endDate   = app.Flag("endDate", "End date, YYYY-MM-YY, default is the most recent Monday at midnight").Default(end).String()
		timeZone  = app.Flag("timezone", "Client's timezone, defaults to New York").Default("America/New_York").String()
	)
	app.Parse(os.Args[1:])

	loc, err := time.LoadLocation(*timeZone)
	if err != nil {
		log.Fatalf("%v, '%v'\n", err, *timeZone)
	}
	startTime, err = time.ParseInLocation(TimeFormat, *startDate, loc)
	d, _ := time.ParseDuration("0h0m0.0s")
	startTime = startTime.Add(d)
	if err != nil {
		log.Fatalf("%v, '%v'\n", err, *startDate)
	}
	d, _ = time.ParseDuration("23h59m59.999s")
	endTime, err = time.ParseInLocation(TimeFormat, *endDate, loc)
	endTime = endTime.Add(d)
	if err != nil {
		log.Fatalf("%v, '%v'\n", err, *endDate)
	}

	if endTime.Before(startTime) {
		start = endTime.Format(TimeFormat)
		end = startTime.Format(TimeFormat)
		log.Fatalf("End date (%v) is before start date (%v)", start, end)
	}

	//Engage expects Zulu time.
	_, offset := startTime.Zone()
	zt := fmt.Sprintf("%ds", -offset)
	d, err = time.ParseDuration(zt)
	startTime = startTime.Add(d)
	endTime = endTime.Add(d)

	//Convert to funky Engage format.
	engageStart := startTime.Format(goengage.EngageDateFormat)
	engageEnd := endTime.Format(goengage.EngageDateFormat)

	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	service := goengage.NewDedicationService()
	err = goengage.ReportFundraising(e, service, engageStart, engageEnd)
	if err != nil {
		panic(err)
	}
}
