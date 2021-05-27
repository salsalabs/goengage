package main

//Application to accept a segmentId and output the supporters that belong
//to the segment.  Output includes a list of the other segments that a
//supporter belongs to.  Produces a CSV of supporter_KEY, Email, Groups.

import (
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	//SupporterListeners is the number of listeners for segments info records.
	SupporterListeners = 5
	//XrefListeners is the channel volume for the xref listeners.
	XrefListeners = 500
	//StartDate is used for the starting of Joined and LastModified ranges.
	StartDate = "2001-01-01T01:01:01.001Z"
)

//XrefRecord is the container for the information that goes to the output.
type XrefRecord struct {
	SupporterID string
	Email       string
	Segments    []string
}

//NewXrefRecord creates an XrefRecord and returns a reference.
func NewXrefRecord(s string, e string) *XrefRecord {
	x := XrefRecord{
		SupporterID: s,
		Email:       e,
		Segments:    make([]string, 0),
	}
	return &x
}

//Runtime holds the common data used by the tasks in this app.
type Runtime struct {
	E         *goengage.Environment
	SegmentID string
	C1        chan *XrefRecord
	C2        chan *XrefRecord
	D         chan bool
	F         string
	L         *goengage.UtilLogger
	N         *regexp.Regexp
}

//Members accepts a segmentId and writes the segment members to the
//provided channel.
func Members(rt Runtime) (err error) {
	log.Println("Members: begin")

	count := rt.E.Metrics.MaxBatchSize
	offset := int32(0)
	for count == rt.E.Metrics.MaxBatchSize {
		payload := goengage.SegmentMembershipRequestPayload{
			SegmentId:   rt.SegmentID,
			Offset:      offset,
			Count:       count,
			JoinedSince: StartDate,
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
			Logger:   rt.L,
		}
		err = n.Do()
		if err != nil {
			return err
		}
		if offset%500 == 0 {
			log.Printf("Members: %6d: %2d of %6d\n",
				offset,
				len(resp.Payload.Supporters),
				resp.Payload.Total)
		}

		for _, s := range resp.Payload.Supporters {
			p := goengage.FirstEmail(s)
			email := ""
			if p != nil {
				email = *p
			}
			x := NewXrefRecord(s.SupporterID, email)
			rt.C1 <- x
		}
		count = resp.Payload.Count
		offset += int32(count)
	}
	close(rt.C1)
	log.Println("Members: end")
	return nil
}

//Segments accepts an xref record from the channel, populates the Groups field, then
//pushes the completed record into the write channel. Notifies done with the input
//channel is empty.
func Segments(rt Runtime, id int) (err error) {
	log.Printf("Segments %+v: begin\n", id)
	for {
		x, ok := <-rt.C1
		if !ok {
			break
		}

		// Read groups, sort, then pass them to the writer's channel.
		count := rt.E.Metrics.MaxBatchSize
		offset := int32(0)

		for count == rt.E.Metrics.MaxBatchSize {
			payload := goengage.SupporterGroupRequestPayload{
				Identifiers:    []string{x.SupporterID},
				IdentifierType: goengage.SupporterIDType,
				ModifiedFrom:   StartDate,
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
				Logger:   rt.L,
			}
			err = n.Do()
			if err != nil {
				return err
			}
			respPayload := resp.Payload
			results := respPayload.Results
			for _, s := range results {
				for _, t := range s.Segments {
					if t.SegmentID != rt.SegmentID {
						if !(rt.N != nil && rt.N.MatchString(t.Name)) {
							x.Segments = append(x.Segments, t.Name)
						}
					}
				}
			}
			count = resp.Payload.Count
			offset += int32(count)
		}
		rt.C2 <- x
	}

	rt.D <- true
	return nil
}

//OutputCSV accepts Xref records from a channeland and writes them to
//a CSV file.
func OutputCSV(rt Runtime) error {
	log.Printf("OutputCSV: begin")
	f, err := os.Create(rt.F)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	headers := []string{"SupporterId", "Email", "Groups"}
	w.Write(headers)
	for {
		x, ok := <-rt.C2
		if !ok {
			break
		}
		s := strings.Join(x.Segments, ",")
		row := []string{
			x.SupporterID,
			x.Email,
			s,
		}
		w.Write(row)
	}
	w.Flush()
	f.Close()
	log.Printf("OutputCSV: end")
	return nil
}

//WaitTerminations waits for "SupporterListeners" supporter readers to
//complete.  That triggers a close for the CSV writer channel.
func WaitTerminations(rt Runtime) {
	log.Printf("WaitTerminations: begin")
	remaining := SupporterListeners
	for remaining > 0 {
		log.Printf("WaitTerminations: waiting for %d listeners\n", remaining)
		<-rt.D
		remaining--
	}
	close(rt.C2)
	log.Printf("WaitTerminations: end")
}

//Program entry point.
func main() {
	var (
		app       = kingpin.New("one_segment_xref", "Find supporters for a segment. Display supporters and lists of groups they belong to.")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		segmentID = app.Flag("segmentId", "primary key for the segment of interest").Default("0d2b6078-6a5c-42c0-b62d-e01208b468cd").String()
		notThis   = app.Flag("ignore-segments-like", "Regex for segments not to consider.  Defaults to nil").String()
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
	if segmentID == nil || len(*segmentID) != 36 {
		log.Fatalf("Error --segmentId is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("Error: %+v\n", e)
	}

	logger, err := goengage.NewUtilLogger()
	if err != nil {
		panic(err)
	}

	var regex *regexp.Regexp
	if notThis != nil {
		regex = regexp.MustCompile(*notThis)
	}
	rt := Runtime{
		E:         e,
		SegmentID: *segmentID,
		C1:        make(chan *XrefRecord, XrefListeners),
		C2:        make(chan *XrefRecord, XrefListeners),
		D:         make(chan bool, SupporterListeners),
		F:         *csvFile,
		L:         logger,
		N:         regex,
	}
	var wg sync.WaitGroup

	//Start the CSV output listener.  Note wg.Add before. Should
	//reduce race conditions.
	wg.Add(1)
	go (func(rt Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := OutputCSV(rt)
		if err != nil {
			panic(err)
		}
	})(rt, &wg)
	log.Printf("main: CSV writer started\n")

	//Start segment listeners
	for id := 1; id <= SupporterListeners; id++ {
		wg.Add(1)
		go (func(rt Runtime, id int, wg *sync.WaitGroup) {
			defer wg.Done()
			err := Segments(rt, id)
			if err != nil {
				panic(err)
			}
		})(rt, id, &wg)
	}
	log.Printf("main: %+v segment listeners started\n", SupporterListeners)

	//Start "done" listener to keep track of segment listeners.
	wg.Add(1)
	go (func(rt Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		WaitTerminations(rt)
	})(rt, &wg)
	log.Println("main: terminations listener started")

	//Start segment reader.
	wg.Add(1)
	go (func(rt Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := Members(rt)
		if err != nil {
			panic(err)
		}
	})(rt, &wg)
	log.Printf("main: segment reader started\n")

	log.Printf("main: napping...\n")
	time.Sleep(10 * time.Second)
	log.Printf("main: waiting...\n")
	wg.Wait()
	log.Printf("main: done")
}
