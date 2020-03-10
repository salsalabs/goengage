package main

import (
	"encoding/csv"
	"log"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Program entry point.  Summarize segments.  No details.
func main() {
	var (
		app     = kingpin.New("see_distrincts", "A command-line app to to show districs for segment for email address(es).")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("csv", "CSV filename to store segment info").Default("segments.csv").String()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		log.Fatalf("Error --login is required.")
		os.Exit(1)
	}
	if csvFile == nil || len(*csvFile) == 0 {
		log.Fatalf("Error --csv is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	headers := []string{"ID",
		"GroupName",
	}
	f, err := os.Create(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.Write(headers)
	if err != nil {
		log.Fatal(err)
	}

	//Read segments and save them.
	count := e.Metrics.MaxBatchSize
	offset := int32(0)
	for count > 0 {
		payload := goengage.SegmentSearchRequestPayload{
			Offset: offset,
			Count:  count,
		}
		rqt := goengage.SegmentSearchRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}
		var resp goengage.SegmentSearchResponse

		n := goengage.NetOp{
			Host:     e.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.SearchSegment,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err = n.Do()
		if err != nil {
			panic(err)
		}
		for _, s := range resp.Payload.Segments {
			var a []string
			a = append(a, s.SegmentID)
			a = append(a, s.Name)
			err = w.Write(a)
		}
		count = resp.Payload.Count
		log.Printf("Main: read %3d from offset %4d\n", count, offset)
		offset += int32(count)
	}
	w.Flush()
	log.Printf("Done.  Output is in %v\n", *csvFile)
}
