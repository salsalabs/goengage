package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	//SegmentListeners is the number of listeners for segments info records.
	SegmentListeners = 5
	//RowsPerCSV is the maximum number of rows in a CSV.  Keeps the individual
	//files to a reasonable number.
	RowsPerCSV = 100_000
)

//ReadSegments reads segments from Engage and writes them to
//a channel.  The channel is closed when all segments have been
//written.
func ReadSegments(e *goengage.Environment, offset int32, c chan goengage.Segment) (err error) {
	log.Println("ReadSegments: begin")
	count := e.Metrics.MaxBatchSize
	for count == e.Metrics.MaxBatchSize {
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
			return err
		}
		//Criteria for the first pass is to take groups containing "ALS" and
		//to ignore gruops that contain "test".  That gives us a chance to get
		//the useful data done more quickly.
		for _, s := range resp.Payload.Segments {
			name := strings.ToLower(s.Name)
			if strings.Contains(name, "als") && !strings.Contains(name, "test") {
				c <- s
				log.Printf("ReadSegments: pushed %-16v %v", s.Type, s.Name)
			}
		}
		count = resp.Payload.Count
		log.Printf("ReadSegments: %3d + %3d = %3d of %4d\n", offset, count, offset+int32(count), resp.Payload.Total)
		offset += int32(count)
	}
	close(c)
	log.Println("ReadSegments: end")
	return nil
}

//ReadSupporters reads from the segment channel and writes all chapter members
//to a CSV file names for the segment.  Note that an existing CSV  causes a segment
//to be ignored.
func ReadSupporters(e *goengage.Environment, c1 chan goengage.Segment, done chan bool, id int) (err error) {
	log.Printf("ReadSupporters %v: begin\n", id)
	for true {
		r, ok := <-c1
		if !ok {
			break
		}
		log.Printf("ReadSupporters %v: popped %v\n", id, r.Name)

		//Create a CSV filename for the group an see if the file exists.
		filename := fmt.Sprintf("%v.csv", r.Name)
		_, err := os.Stat(filename)
		if os.IsExist(err) {
			log.Printf("ReadSupporters: %v already exists, ignoring group\n", filename)
		} else {
			f, err := os.Create(filename)
			if err != nil {
				return err
			}
			w := csv.NewWriter(f)
			headers := []string{"Email"}
			w.Write(headers)

			// Read all supporters and write info to the group's CSV.
			count := e.Metrics.MaxBatchSize
			offset := int32(0)
			for count == e.Metrics.MaxBatchSize {
				payload := goengage.SupporterSearchRequestPayload{
					SegmentID: r.SegmentID,
					Offset:    offset,
					Count:     count,
				}
				rqt := goengage.SupporterSearchRequest{
					Header:  goengage.RequestHeader{},
					Payload: payload,
				}
				var resp goengage.SupporterSearchResponse

				n := goengage.NetOp{
					Host:     e.Host,
					Method:   goengage.SearchMethod,
					Endpoint: goengage.SupporterSearchSegment,
					Token:    e.Token,
					Request:  &rqt,
					Response: &resp,
				}
				err = n.Do()
				if err != nil {
					return err
				}
				for _, s := range resp.Payload.Supporters {
					email := goengage.FirstEmail(s)
					a := []string{*email}
					w.Write(a)
				}
				w.Flush()
				count = resp.Payload.Count
				log.Printf("ReadSupporters %v: %-32v %6d + %3d = %6d of %6d\n",
					id,
					r.Name,
					offset,
					count,
					offset+int32(count),
					resp.Payload.Total)
				offset += int32(count)
			}
		}
	}
	done <- true
	log.Printf("ReadSupporters %v: end\n", id)
	return nil
}

//WaitTerminations waits for "SegmentListeners" supporter readers to complete.
func WaitTerminations(done chan bool) {
	remaining := SegmentListeners
	for remaining > 0 {
		log.Printf("WaitTerminations: waiting for %d listeners\n", remaining)
		_ = <-done
		remaining--
	}
}

//Program entry point.
func main() {
	var (
		app     = kingpin.New("segments_and_supporters", "A command-line app to write Engage segments and email addresses to CSV files.")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("csv", "CSV filename to store segment info").Default("segment_and_supporters.csv").String()
		offset  = app.Flag("offset", "Starting offset").Default("0").Int32()
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

	segChan := make(chan goengage.Segment, 50)
	done := make(chan bool, SegmentListeners)
	var wg sync.WaitGroup

	//Start segment listeners(s)
	for id := 1; id <= SegmentListeners; id++ {
		go (func(e *goengage.Environment, c1 chan goengage.Segment, done chan bool, id int, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			err := ReadSupporters(e, c1, done, id)
			if err != nil {
				panic(err)
			}
		})(e, segChan, done, id, &wg)
	}
	log.Printf("main: %v segment listeners started\n", SegmentListeners)

	//Start "done" listener to keep track of segment listeners.
	go (func(done chan bool, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		WaitTerminations(done)
	})(done, &wg)
	log.Println("main: terminations listener started")

	//Start segment reader.
	go (func(e *goengage.Environment, c chan goengage.Segment, offset int32, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := ReadSegments(e, offset, c)
		if err != nil {
			panic(err)
		}
	})(e, segChan, *offset, &wg)
	log.Printf("main: segment reader started\n")

	log.Printf("main: waiting...\n")
	wg.Wait()
	log.Printf("main: done")
}
