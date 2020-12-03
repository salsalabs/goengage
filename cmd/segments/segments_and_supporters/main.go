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
	SegmentListeners = 2
	//RowsPerCSV is the maximum number of rows in a CSV.  Keeps the individual
	//files to a reasonable number.
	RowsPerCSV = 100_000
)

//OutputRecord contains segment and supporter info.
type OutputRecord struct {
	Segment   goengage.Segment
	Supporter goengage.Supporter
}

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
		for _, s := range resp.Payload.Segments {
			c <- s
			log.Printf("ReadSegments: pushed %v", s.Name)
		}
		count = resp.Payload.Count
		log.Printf("ReadSegments: %3d + %3d = %3d of %4d\n", offset, count, offset+int32(count), resp.Payload.Total)
		offset += int32(count)
	}
	close(c)
	log.Println("ReadSegments: end")
	return nil
}

//ReadSupporters reads from the segment channel and writes the contents to
//the log.  Used to test the segment driver.
func ReadSupporters(e *goengage.Environment, c1 chan goengage.Segment, c2 chan OutputRecord, done chan bool, id int) (err error) {
	log.Printf("ReadSupporters %v: begin\n", id)
	for true {
		r, ok := <-c1
		if !ok {
			break
		}
		log.Printf("ReadSupporters %v: popped %v\n", id, r.Name)
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
				a := OutputRecord{
					Segment:   r,
					Supporter: s,
				}
				c2 <- a
			}
			count = resp.Payload.Count
			log.Printf("ReadSupporters: %32v %3d + %3d = %3d of %4d\n",
				r.Name,
				offset,
				count,
				offset+int32(count),
				resp.Payload.Total)
			offset += int32(count)
		}
		log.Println("ReadSegments: end")
		return nil
	}
	done <- true
	log.Printf("ReadSupporters %v: end\n", id)
	return nil
}

//WaitTerminations waits for "SegmentListeners" supporter readers to complete.
func WaitTerminations(c2 chan OutputRecord, done chan bool) {
	remaining := SegmentListeners
	for remaining > 0 {
		log.Printf("WaitTerminations: waiting for %d listeners\n", remaining)
		_ = <-done
		remaining--
	}
	close(c2)
}

//WriteOutput reads from the output queue and writes to the CSV file.
func WriteOutput(e *goengage.Environment, c chan OutputRecord, csvFile string) (err error) {
	headers := []string{"GroupID",
		"GroupName",
		"SupporterID",
		"Email",
	}
	rows := RowsPerCSV
	current := 1
	var f *os.File
	var w *csv.Writer

	log.Printf("WriteOutput: begin\n")
	for true {
		r, ok := <-c
		if !ok {
			break
		}

		// Open an output file as needed.
		if rows >= RowsPerCSV {
			if f != nil {
				err := f.Close()
				if err != nil {
					return err
				}
				f = nil
			}
			parts := strings.Split(csvFile, ".")
			s := fmt.Sprintf("%s-%03d.%s", parts[0], current, parts[1])
			current++
			f, err := os.Create(s)
			if err != nil {
				return err
			}
			w = csv.NewWriter(f)
			err = w.Write(headers)
			if err != nil {
				log.Fatal(err)
			}
			rows = 0
			log.Printf("WriteOutput: opened %s\n", s)
		}

		var a []string
		a = append(a, r.Segment.SegmentID)
		a = append(a, r.Segment.Name)
		a = append(a, r.Supporter.SupporterID)
		email := goengage.FirstEmail(r.Supporter)
		a = append(a, *email)
		err = w.Write(a)
		if err != nil {
			return (err)
		}
		w.Flush()
	}
	if f != nil {
		err := f.Close()
		if err != nil {
			return err
		}
		f = nil
	}

	log.Printf("WriteOutput: end\n")
	return nil
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
	outChan := make(chan OutputRecord, 1000)
	done := make(chan bool, SegmentListeners)
	var wg sync.WaitGroup

	//Start output writer.
	go (func(e *goengage.Environment, c chan OutputRecord, csvFile string, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := WriteOutput(e, c, csvFile)
		if err != nil {
			panic(err)
		}
	})(e, outChan, *csvFile, &wg)
	log.Println("main: output writer started")

	//Start segment listeners(s)
	for id := 1; id <= SegmentListeners; id++ {
		go (func(e *goengage.Environment, c1 chan goengage.Segment, c2 chan OutputRecord, done chan bool, id int, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			err := ReadSupporters(e, c1, c2, done, id)
			if err != nil {
				panic(err)
			}
		})(e, segChan, outChan, done, id, &wg)
	}
	log.Printf("main: %v segment listeners started\n", SegmentListeners)

	//Start "done" listener to keep track of segment listeners.
	go (func(c2 chan OutputRecord, done chan bool, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		WaitTerminations(c2, done)
	})(outChan, done, &wg)
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
