package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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
			}
		}
		count = resp.Payload.Count
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
		//Create a CSV filename for the group an see if the file exists.
		filename := fmt.Sprintf("%v.csv", r.Name)
		filename = strings.Replace(filename, "/", "-", -1)
		_, err := os.Stat(filename)
		if err == nil || os.IsExist(err) {
			log.Printf("ReadSupporters %v: %-32v skipped, file exists\n", id, r.Name)
		} else {
			log.Printf("ReadSupporters %v: %-32v start\n", id, r.Name)
			// Create a file using the ID and write to it.  We'll rename it to the group
			// when all of the supporters are gathered.
			temp := fmt.Sprintf("%v.csv", r.SegmentID)
			f, err := os.Create(temp)
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
				ok := false
				for !ok {
					err = n.Do()
					if err != nil {
						fmt.Printf("ReadSupporters %v: %-32v %v\n", id, r.Name, err)
						time.Sleep(time.Second)
					} else {
						ok = true
					}
				}
				for _, s := range resp.Payload.Supporters {
					email := goengage.FirstEmail(s)
					if email != nil {
						a := []string{*email}
						w.Write(a)
					}
				}
				w.Flush()
				count = resp.Payload.Count
				offset += int32(count)
			}
			err = os.Rename(temp, filename)
			if err != nil {
				return err
			}
			log.Printf("ReadSupporters %v: %-32v done\n", id, r.Name)
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

	log.Printf("main: napping...\n")
	time.Sleep(10 * time.Second)
	log.Printf("main: waiting...\n")
	wg.Wait()
	log.Printf("main: done")
}
