package main

//Application to accept a segmentId and output the supporters that belong
//to the segment.  Output includes a list of the other segments that a
//supporter belongs to.  Produces a CSV of supporter_KEY, Email, Groups.

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	//SupporterListeners is the number of listeners for segments info records.
	SupporterListeners = 5
)

//XrefRecord is the container for the information that goes to the output.
type XrefRecord struct {
	SupporterID string
	Email       string
	Segments    []string
}

//Runtime holds the common data used by the tasks in this app.
type Runtime struct {
	E         goengage.Environment
	SegmentId string
	C1        chan XrefRecord
	C2        chan XrefRecord
	D         chan bool
	W         *csv.Writer
}

//Members accepts a segmentId and writes the segment members to the
//provided channel.
func Members(rt Runtime) (err error) {
	log.Println("Members: begin")
	count := rt.E.Metrics.MaxBatchSize
	offset := int32(0)
	for count == rt.E.Metrics.MaxBatchSize {
		payload := goengage.SegmentMembershipRequestPayload{
			SegmentId: rt.SegmentId,
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
		for _, s := range resp.Payload.Supporters {
			p := goengage.FirstEmail(s)
			email := ""
			if p != nil {
				email = *p
			}
			x := XrefRecord{
				SupporterID: s.SupporterID,
				Email:       email,
				Segments:    make([]string, 0),
			}
			rt.C1 <- x
		}
		count = resp.Payload.Count
		offset += int32(count)
	}
	close(rt.C1)
	log.Println("Members: end")
	return nil
}

//Supporters accepts an xref record from the channel, populates the Groups field, then
//pushes the completed record into the write channel. Notifies done with the input
//channel is empty.
func Segments(rt Runtime, id int) (err error) {
	log.Printf("Supporters %v: begin\n", id)
	for true {
		x, ok := <-rt.C1
		if !ok {
			break
		}
		// Read groups, sort, then pass them to the writer's channel.
		count := rt.E.Metrics.MaxBatchSize
		offset := int32(0)
		var identifiers = make([]string, 0)
		identifiers = append(identifiers, x.SupporterID)

		for count == rt.E.Metrics.MaxBatchSize {
			payload := goengage.SupporterGroupRequestPayload{
				Identifiers:    []string{x.SupporterID},
				IdentifierType: goengage.SupporterIDType,
				Offset:         offset,
				Count:          count,
			}
			rqt := goengage.SupporterGroupRequest{
				Header:  goengage.RequestHeader{},
				Payload: payload,
			}
			var resp goengage.SupporterGroupResponse

			n := goengage.NetOp{
				Host:     rt.E.Host,
				Method:   goengage.SearchMethod,
				Endpoint: goengage.SupporterSearchGroups,
				Token:    rt.E.Token,
				Request:  &rqt,
				Response: &resp,
			}
			err = n.Do()
			if err != nil {
				return err
			}
			respPayload := resp.Payload
			results := respPayload.Results
			for _, s := range results {
				for _, t := range s.Segments {
					x.Segments = append(x.Segments, t)
				}
			}
			count = resp.Payload.Count
			offset += int32(count)
		}
		rt.C2 <- x
	}

	rt.D <- true
	log.Printf("Supporters %v: end\n", id)
	return nil
}

//WaitTerminations waits for "SupporterListeners" supporter readers to complete.
func WaitTerminations(done chan bool) {
	remaining := SupporterListeners
	for remaining > 0 {
		log.Printf("WaitTerminations: waiting for %d listeners\n", remaining)
		_ = <-done
		remaining--
	}
}

//Program entry point.
func main() {
	var (
		app       = kingpin.New("segments_and_supporters", "A command-line app to write Engage segments and email addresses to CSV files.")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		segmentId = app.Flag("segmentId", "primary key for the segment of interest").Required().String()
		csvFile   = app.Flag("csv", "CSV filename to store supporter-segment info").Default("supporter_segment.csv").String()
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

	c1 := make(chan XrefRecord, 50)
	c2 := make(chan XrefRecord, 100)
	doneChan := make(chan bool, SupporterListeners)
	f, err := os.Create(*csvFile)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	headers := []string{"SupporterId", "Email", "Groups"}
	w.Write(headers)

	rt := Runtime{
		Env:       e,
		SegmentId: *segmentId,
		C1:        c1,
		C2:        c2,
		D:         doneChan,
		Writer:    w,
	}
	var wg sync.WaitGroup

	//Start segment listeners(s)
	for id := 1; id <= SupporterListeners; id++ {
		go (func(e *goengage.Environment, c1 chan goengage.Segment, done chan bool, id int, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			err := Supporters(e, c1, done, id)
			if err != nil {
				panic(err)
			}
		})(e, segChan, done, id, &wg)
	}
	log.Printf("main: %v segment listeners started\n", SupporterListeners)

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
		err := Memberss(e, offset, c)
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
