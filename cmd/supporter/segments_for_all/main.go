// App to write a CSV of supporters and segments.  Each row is a single
// supporter-segment relationship.  A row contains
// * supporterId
// * Email
// * segmentId
// * segmentName
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	report "github.com/salsalabs/goengage/pkg/report"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	//SupporterListenerCount is the number of channels waiting for
	//supporter records.
	SupporterListenerCount = 5
)

// OutRec holds a supporter-segment relationship.
type OutRec struct {
	Supporter goengage.Supporter
	Segment   goengage.Segment
}

// Runtime area for this app.
type Runtime struct {
	E             *goengage.Environment
	SupporterChan chan goengage.Supporter
	DoneChan      chan bool
	OutChan       chan OutRec
	CSVOut        *csv.Writer
}

// NewRuntime populates a new runtime.
func NewRuntime(env *goengage.Environment, out *csv.Writer) Runtime {
	r := Runtime{
		E:             env,
		SupporterChan: make(chan goengage.Supporter, 100),
		DoneChan:      make(chan bool),
		OutChan:       make(chan OutRec, 100),
		CSVOut:        out,
	}
	return r
}

// Adjust offset changes the proposed offset as needed.
// Implements SupporterGuide.AdjustOffset.
// Useful for chunked ID reads.  Does nothing in this app.
func (r *Runtime) AdjustOffset(offset int32) int32 {
	return offset
}

// Visit implements SupporterGuide.Visit and does something with
// a supporter record.  In this case, Visit retrieves segments
// for a supporter and writes them to OutChan.
func (r *Runtime) Visit(s goengage.Supporter) error {
	segments, err := goengage.SupporterSegments(r.E, s.SupporterID)
	if err != nil {
		return err
	}
	for _, g := range segments {
		outRec := OutRec{s, g}
		r.OutChan <- outRec
	}
	return nil
}

// Finalize implements SupporterGuide.Finalize and does nothing
// in this app.
func (r *Runtime) Finalize() error {
	return nil
}

// Payload implements SupporterGuide.Payload and provides a payload
// that will retrieve all supporters.
func (r *Runtime) Payload() goengage.SupporterSearchRequestPayload {
	payload := goengage.SupporterSearchRequestPayload{
		IdentifierType: goengage.SupporterIDType,
		ModifiedFrom:   "2000-01-01T00:00:00.00000Z",
		ModifiedTo:     "2050-01-01T00:00:00.00000Z",
		Offset:         0,
		Count:          0,
	}
	return payload
}

// Channel implements SupporterGuide.Channnel and provides the
// supporter channel.
func (r *Runtime) Channel() chan goengage.Supporter {
	return r.SupporterChan
}

// DoneChannel implements SupporterGuide.DoneChannel to provide
// a channel that receives a true when the listener(s) are done.
func (r *Runtime) DoneChannel() chan bool {
	return r.DoneChan
}

// Offset returns the offset for the first read.
// Useful for restarts.
func (r *Runtime) Offset() int32 {
	return 0
}

// Writer accepts items from OutChan and writes them to the CSV.
func (r *Runtime) Writer() error {
	count := int32(0)
	log.Printf("Writer: begin")
	for {
		s, okay := <-r.OutChan
		if !okay {
			break
		}
		email := ""
		e := goengage.FirstEmail(s.Supporter)
		if e != nil {
			email = *e
		}
		row := []string{
			s.Supporter.SupporterID,
			email,
			s.Segment.SegmentID,
			s.Segment.Name,
		}
		err := r.CSVOut.Write(row)
		if err != nil {
			return err
		}
		r.CSVOut.Flush()
		if count%1000 == 0 {
			log.Printf("Writer: %d\n", count)
		}
		count++
	}
	log.Printf("Writer: end")
	return nil
}

// Program entry point. Scan through supporters.  Write supporter-group
// data to a CSV file.
func main() {
	var (
		app     = kingpin.New("supporter_segments", "Write a CSV of supporters and segments")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		outFile = app.Flag("output", "CSV filename to store supporter-segment data").Default("supporter_segments.csv").String()
		//debug   = app.Flag("debug", "Write requests and responses to a log file in JSON").Bool()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		fmt.Println("Error --login is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	headers := []string{
		"SupporterID",
		"Email",
		"SegmentID",
		"SegmentName",
	}
	err = writer.Write(headers)
	if err != nil {
		log.Fatalf("%s on %s during header write\n", err, *outFile)
	}

	r := NewRuntime(e, writer)
	var wg sync.WaitGroup

	//start CSV writer. It waits for (supporter, segment) records
	//to arrive on OutChan.
	wg.Add(1)
	go (func(r *Runtime) {
		err := r.Writer()
		if err != nil {
			log.Fatalf("%s on CSV writer", err)
		}
	})(&r)

	//Start supporter processor.  This drives all of the "Guide"-based
	//calls in this  source file.
	for i := 0; i < SupporterListenerCount; i++ {
		wg.Add(1)
		go (func(r *Runtime, wg *sync.WaitGroup) {
			defer wg.Done()
			report.ProcessSupporters(r.E, r)
		})(&r, &wg)
	}
	log.Printf("main: started %d supporter listeners\n", SupporterListenerCount)

	//Start done listener. It waits for all of the supporter
	//readers to complete.
	wg.Add(1)
	go (func(r *Runtime, n int, wg *sync.WaitGroup) {
		defer wg.Done()
		goengage.DoneListener(r.DoneChan, n)
		close(r.OutChan)
	})(&r, SupporterListenerCount, &wg)

	//Start supporter reader. Reads all supporters and puts them
	//onto a supporter channel.
	wg.Add(1)
	go (func(e *goengage.Environment, r *Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		report.ReadSupporters(r.E, r)
	})(e, &r, &wg)

	d, err := time.ParseDuration("10s")
	if err != nil {
		panic(err)
	}
	log.Printf("main: sleeping for %s seconds", d)
	time.Sleep(d)
	log.Printf("main:  waiting...")
	wg.Wait()
	log.Printf("main: done")
}
