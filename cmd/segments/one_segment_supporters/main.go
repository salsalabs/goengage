// App to search for supporters who are members of a group.
// Displays contents of the "Recording" object.
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Number of input queue listeners.
	ListenerCount = 5
)

// Recording is the content that appears in the CSV file.
type Recording struct {
	SupporterID  string
	DateCreated  string
	Email        string
	AddressLine1 string
	PostalCode   string
}

// Runtime area for this app.
type Runtime struct {
	E         *goengage.Environment
	InChan    chan goengage.Supporter
	CsvChan   chan Recording
	DoneChan  chan bool
	SegmentID string
	Results   string
	Logger    *goengage.UtilLogger
}

// NewRuntime populates a new runtime.
func NewRuntime(env *goengage.Environment, segmentID string, results string, verbose bool) Runtime {
	r := Runtime{
		E:         env,
		InChan:    make(chan goengage.Supporter),
		CsvChan:   make(chan Recording),
		DoneChan:  make(chan bool),
		SegmentID: segmentID,
		Results:   results,
	}
	if verbose {
		logger, err := goengage.NewUtilLogger()
		if err != nil {
			log.Fatalf("unable to create logger, %v", err)
		}
		r.Logger = logger
	}
	return r
}

// Drive finds all of the supporters in the specified group and writes
// them to the input channel.
func Drive(rt Runtime) (err error) {
	log.Println("Drive: begin")
	count := rt.E.Metrics.MaxBatchSize
	offset := int32(0)
	for count == rt.E.Metrics.MaxBatchSize {
		payload := goengage.SegmentMembershipRequestPayload{
			SegmentID: rt.SegmentID,
			Offset:    offset,
			Count:     count,
		}
		rqt := goengage.SegmentMembershipRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}
		var resp goengage.SegmentMembershipResponse

		n := goengage.NetOp{
			Host:     rt.E.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.SegmentSearchMembers,
			Token:    rt.E.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err = n.Do()
		if err != nil {
			return err
		}
		if offset%100 == 0 {
			log.Printf("Drive: Read %4d of %4d\n", offset, resp.Payload.Total)
		}
		for _, s := range resp.Payload.Supporters {
			rt.InChan <- s
		}
		count = int32(len(resp.Payload.Supporters))
		offset += count
	}
	close(rt.InChan)
	log.Println("Drive: end")
	return nil
}

// Record reads Recording records from the CSV channel and writes
// them to the results file.
func Record(rt Runtime) (err error) {
	log.Println("Record: begin")
	f, err := os.Create(rt.Results)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	headers := []string{
		"SegmentID",
		"SupporterID",
		"DateCreated",
		"Email",
		"AddressLine1",
		"ZipCode",
	}
	err = writer.Write(headers)

	for {
		r, okay := <-rt.CsvChan
		if !okay {
			break
		}
		row := []string{
			rt.SegmentID,
			r.SupporterID,
			r.DateCreated,
			r.Email,
			r.AddressLine1,
			r.PostalCode,
		}
		writer.Write(row)
	}
	writer.Flush()
	log.Println("Record: end")
	return nil
}

// Update retrieves Supporters from the input channel. Each
// supporter is formatted as a SupporterKludgeFix record and
// submitted to Engage for repair. With that done, a Recording
// is dropped on the CSV channel showing the results.
func Update(rt Runtime, id int) (err error) {
	log.Printf("Update-{%d}: start\n", id)
	for {
		s, okay := <-rt.InChan
		if !okay {
			break
		}
		e := goengage.FirstEmail(s)
		email := ""
		if e != nil {
			email = *e
		}
		row := Recording{
			s.SupporterID,
			s.CreatedDate.Format("2006-01-02"),
			email,
			s.Address.AddressLine1,
			s.Address.PostalCode,
		}
		rt.CsvChan <- row
	}
	log.Printf("Update-{%d}: end", id)
	rt.DoneChan <- true
	return nil
}

// Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app       = kingpin.New("one-segment-supporters", "Selectively display info about supporters in a segment.")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		segmentID = app.Flag("segment-id", "segmentID for the group").Default("f4be4a19-b85f-4d69-baae-e027a86fd676").String()
		results   = app.Flag("results", "filename of CSV file to record results").Default("one_segment_supporters.csv").String()
		verbose   = app.Flag("verbose", "Log contents of all network actions. *Really* noisy").Default("false").Bool()
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

	rt := NewRuntime(e, *segmentID, *results, *verbose)
	var wg sync.WaitGroup

	//Start the recording listener.
	wg.Add(1)
	go (func(rt Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := Record(rt)
		if err != nil {
			panic(err)
		}
	})(rt, &wg)

	//Start supporter listeners.
	for i := 1; i <= ListenerCount; i++ {
		wg.Add(1)
		go (func(rt Runtime, wg *sync.WaitGroup, i int) {
			defer wg.Done()
			err := Update(rt, i)
			if err != nil {
				panic(err)
			}
		})(rt, &wg, i)
	}
	//Start done listener.
	wg.Add(1)
	go (func(rt Runtime, n int, wg *sync.WaitGroup) {
		defer wg.Done()
		goengage.DoneListener(rt.DoneChan, n)
		close(rt.CsvChan)
	})(rt, ListenerCount, &wg)

	//Start supporter reader.
	wg.Add(1)
	go (func(rt Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := Drive(rt)
		if err != nil {
			log.Fatalf("driver error: %v\n", err)
		}
	})(rt, &wg)

	d, err := time.ParseDuration("2s")
	if err != nil {
		panic(err)
	}
	log.Printf("main: sleeping for %s", d)
	time.Sleep(d)
	log.Printf("main:  waiting...")
	wg.Wait()
	log.Printf("main: done")
}
