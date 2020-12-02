package main

import (
	"encoding/csv"
	"log"
	"os"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	//ListenerCount is the number of listeners.
	ListenerCount = 5
)

//SegmentInfo is used to identify segments that need processing.
type SegmentInfo struct {
	ID   string
	Name string
}

//OutputRecord contains segment and supporter info.
type OutputRecord struct {
	Segment   SegmentInfo
	Supporter goengage.Supporter
}

//ReadAll reads segments from Engage and writes them to
//a channel.  The channel is closed when all segments have been
//written.
func ReadAll(e *goengage.Environment, c chan SegmentInfo) (err error) {
	log.Println("ReadAll: begin")
	count := e.Metrics.MaxBatchSize
	offset := int32(0)
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
			a := SegmentInfo{
				ID:   s.SegmentID,
				Name: s.Name,
			}
			c <- a
		}
		count = resp.Payload.Count
		log.Printf("ReadAll: %3d + %3d = %3d of %4d\n", offset, count, offset+int32(count), resp.Payload.Total)
		offset += int32(count)
	}
	close(c)
	log.Println("ReadAll: end")
	return nil
}

//EchoOutput reads from the segment channel and writes the contents to
//the log.  Used to test the segment driver.
func EchoOutput(e *goengage.Environment, c chan SegmentInfo, id int) (err error) {
	log.Printf("EchoOut %v: begin\n", id)
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		log.Printf("EchoOut %v: %+v\n", id, r)
	}
	log.Printf("EchoOut %v: end\n", id)
	return nil
}

//Program entry point.  Summarize segments.  No details.
func main() {
	var (
		app     = kingpin.New("segments_details", "A command-line app to write Engage segments and email addresses to a CSV file..")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("csv", "CSV filename to store segment info").Default("segment_and_supporters.csv").String()
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
	headers := []string{"GroupID",
		"GroupName",
		"SupporterID",
		"Email",
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

	segChan := make(chan SegmentInfo)
	// outChan := make(chan OutputRecord)
	var wg sync.WaitGroup

	//Start listener(s)
	for id := 1; id <= ListenerCount; id++ {
		go (func(e *goengage.Environment, c chan SegmentInfo, id int, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			err := EchoOutput(e, c, id)
			if err != nil {
				panic(err)
			}
		})(e, segChan, id, &wg)
	}
	log.Printf("main: %v listeners started\n", ListenerCount)

	//Start reader.
	go (func(e *goengage.Environment, c chan SegmentInfo, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := ReadAll(e, c)
		if err != nil {
			panic(err)
		}
	})(e, segChan, &wg)

	log.Printf("main: reader started\n")
	log.Printf("main: waiting for goroutines\n")
	wg.Wait()
	w.Flush()
	log.Printf("Done.  Output is in %v\n", *csvFile)
}
