package main

// An application to read all email blasts and write detailed info
// into a CSV for each blast.  Uses the getBlastList endpoint from
// Engage;s Web Developer API.
//
// See: https://api.salsalabs.org/help/web-dev#operation/getBlastList

import (
	"encoding/csv"
	"log"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	report "github.com/salsalabs/goengage/pkg/report"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Runtime contains the configuration parts that this app needs.
type Runtime struct {
	Env            *goengage.Environment
	Headers        []string
	RequestPayload goengage.BlastListRequest
	ResultChan     chan goengage.BlastListResult
	DoneChan       chan bool
	BlastOffset    int32
	CSVFile        *os.File
	CSVWriter      *csv.Writer
}

// VisitContent formats and writes a BlastContent record.  Errors terminate.
// Implements goengage.BlastListGuide.
func (rt *Runtime) VisitContent(s goengage.BlastListResult, t goengage.BlastContent) error {
	record := []string{
		s.ID,
		s.Name,
		s.Status,
		s.PublishDate.Format(goengage.EngageDateFormat),
		t.Subject,
		t.PageURL,
	}
	rt.CSVWriter.Write((record))
	rt.CSVWriter.Flush()
	return nil
}

// Finalize is called after all blasts have been processed.
// Implements goengage.BlastListGuide.
func (rt *Runtime) Finalize() error {
	rt.CSVWriter.Flush()
	err := rt.CSVFile.Close()
	return err
}

// Payload is a convenience method to define which blasts to return.
// Each item is turned into a URL query at execution time.
// Implements goengage.BlastListGuide.
func (rt *Runtime) Payload() *goengage.BlastListRequest {
	return &rt.RequestPayload
}

// ResultChannel is the listener channel for blast info.
// Implements goengage.BlastListGuide.
func (rt *Runtime) ResultChannel() chan goengage.BlastListResult {
	return rt.ResultChan
}

// DoneChannel() receives a true when the listener is done.
// Implements goengage.BlastListGuide.
func (rt *Runtime) DoneChannel() chan bool {
	return rt.DoneChan
}

//Offset() returns the offset to start reading.  Useful for
//restarting after a service interruption.
// Implements goengage.BlastListGuide.

func (rt *Runtime) Offset() int32 {
	return rt.BlastOffset
}

//Writer() returns the CSV writer for the output file.

func (rt *Runtime) Writer() *csv.Writer {
	return rt.CSVWriter
}

// Program entry point.
func main() {
	var (
		app          = kingpin.New("blast_urls", "Write email blast info (including URLs) to a CSV")
		login        = app.Flag("login", "YAML file with API token").Required().String()
		blastCSVFile = app.Flag("blast-csv", "CSV filename to store blast info").Default("email_blast_info.csv").String()
		offset       = app.Flag("blast-offset", "Start here if you lose network connectivity").Default("0").Int32()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		log.Fatalf("Error --login is required.")
	}
	if blastCSVFile == nil || len(*blastCSVFile) == 0 {
		log.Fatalf("Error --blast-csv is required.")
	}

	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("Error %v\n", err)
	}
	f, err := os.Create(*blastCSVFile)
	if err != nil {
		log.Fatalf("Error %v on %v\n", err, *blastCSVFile)
	}
	defer f.Close()
	writer := csv.NewWriter(f)

	requestPayload := goengage.BlastListRequest{
		StartDate: "2010-01-01T00:00:00.0000Z",
		EndDate:   "2030-01-01T00:00:00.0000Z",
		Criteria:  "",
		SortField: "name",
		SortOrder: "ASCENDING",
		Count:     20,
		Offset:    0,
	}

	rtx := Runtime{
		Env: e,
		Headers: []string{
			"ID",
			"Name",
			"Status",
			"PublishDate",
			"Subject",
			"PageURL",
		},
		ResultChan:     make(chan goengage.BlastListResult, 1000),
		RequestPayload: requestPayload,
		DoneChan:       make(chan bool),
		BlastOffset:    *offset,
		CSVFile:        f,
		CSVWriter:      writer,
	}
	rtx.CSVWriter.Write(rtx.Headers)

	//Start running.  The Guide does everything for this app.
	err = report.ReportBlastLists(e, &rtx)
	if err != nil {
		log.Fatalf("Error: %v running ReportBlastLists", err)
	}
}
