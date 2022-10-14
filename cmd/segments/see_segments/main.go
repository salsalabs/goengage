package main

//Application to create a CSV of segments for a client. Output includes
//UUID, SegmentName and potentially the .  Can include more --
//find the output formatter and kludge away!

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Runtime holds the stuff that this app needs.
type Runtime struct {
	Env           *goengage.Environment
	IncludeCounts bool
	CSVFilename   string
	Logger        *goengage.UtilLogger
}

func NewRuntime(e *goengage.Environment, c bool, f string, v bool) (*Runtime, error) {
	rt := Runtime{
		Env:           e,
		IncludeCounts: c,
		CSVFilename:   f,
	}
	if v {
		logger, err := goengage.NewUtilLogger()
		if err != nil {
			return nil, err
		}
		rt.Logger = logger
	}
	return &rt, nil
}

// Run finds and displays all segments.
func Run(rt *Runtime) error {
	log.Println("Run: begin")
	f, err := os.Create(rt.CSVFilename)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	headers := []string{
		"SegmentID",
		"SegmentName",
		"MemberCount",
	}
	err = writer.Write(headers)
	if err != nil {
		return err
	}

	count := rt.Env.Metrics.MaxBatchSize
	offset := int32(0)

	for count == rt.Env.Metrics.MaxBatchSize {
		payload := goengage.SegmentSearchRequestPayload{
			Count:               count,
			Offset:              offset,
			IncludeMemberCounts: rt.IncludeCounts,
		}

		rqt := goengage.SegmentSearchRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}

		var resp goengage.SegmentSearchResponse

		n := goengage.NetOp{
			Host:     rt.Env.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.SearchSegment,
			Token:    rt.Env.Token,
			Request:  &rqt,
			Response: &resp,
			Logger:   rt.Logger,
		}
		err := n.Do()
		if err != nil {
			return err
		}
		if offset%100 == 0 {
			log.Printf("Run: %6d: %2d of %6d\n",
				offset,
				len(resp.Payload.Segments),
				resp.Payload.Total)
		}

		var cache [][]string
		for _, s := range resp.Payload.Segments {
			record := []string{
				s.SegmentID,
				s.Name,
				fmt.Sprintf("%v", s.TotalMembers),
			}
			cache = append(cache, record)
		}
		err = writer.WriteAll(cache)
		if err != nil {
			return err
		}
		count = resp.Payload.Count
		offset += int32(count)
	}
	log.Printf("Run: end")
	return nil
}

// Program entry point.
func main() {
	var (
		app           = kingpin.New("one_segment_xref", "Creates a CSV of segments for a client")
		login         = app.Flag("login", "YAML file with API token").Required().String()
		csvFile       = app.Flag("csv", "CSV filename to create").Default("segments.csv").String()
		includeCounts = app.Flag("include-counts", "Output will contain the number of members, too").Bool()
		verbose       = app.Flag("verbose", "Log all requests and responses to a file.  Verrrry noisy...").Bool()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		log.Fatalf("Error --login is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("main: %+v\n", e)
	}
	rt, err := NewRuntime(e, *includeCounts, *csvFile, *verbose)
	if err != nil {
		log.Fatalf("main: %v\n", err)
	}
	err = Run(rt)
	if err != nil {
		log.Fatalf("main: %v\n", err)
	}
}
