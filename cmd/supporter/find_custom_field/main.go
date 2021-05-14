package main

//Application to accept a custom field name and write supporters
//who have that custom field to a CSV file.

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
	//OffsetListeners is the number of listeners for supporters info records.
	OffsetListeners = 5
	//StartDate is used for the starting of Joined and LastModified ranges.
	StartDate = "2001-01-01T01:01:01.001Z"
)

//CFRecord is the container for the information that goes to the output.
type CFRecord struct {
	SupporterID      string
	Email            string
	CustomFieldID    string
	CustomFieldName  string
	CustomFieldValue string
}

//NewCFRecord creates an CFRecord and returns a reference.
func NewCFRecord(s goengage.Supporter, c goengage.CustomFieldValue) *CFRecord {
	f := goengage.FirstEmail(s)
	e := "(None)"
	if f != nil {
		e = *f
	}
	x := CFRecord{
		SupporterID:      s.SupporterID,
		Email:            e,
		CustomFieldID:    c.FieldID,
		CustomFieldName:  c.Name,
		CustomFieldValue: c.Value,
	}
	return &x
}

//Runtime holds the common data used by the tasks in this app.
type Runtime struct {
	E          *goengage.Environment
	Name       string
	MaxRecords int32
	C0         chan int32
	C1         chan goengage.Supporter
	C2         chan *CFRecord
	D          chan bool
	F          string
	L          *goengage.UtilLogger
}

//Offsets accepts a number of records and writes offsets
//to the offsets queue in chunks of MaxBatchSize.
func Offsets(rt Runtime, total int32) {
	for i := int32(0); i < total; i += rt.E.Metrics.MaxBatchSize {
		rt.C0 <- i
	}
	close(rt.C0)
}

//One accepts an offset and passes along the ones that match the filter requirements.
func Supporters(rt Runtime, i int) (err error) {
	log.Printf("Supporters %d: begin\n", i)
	for {
		d, ok := <-rt.C0
		if !ok {
			break
		}
		log.Printf("Supporters %d: %d\n", i, d)
		payload := goengage.SupporterSearchRequestPayload{
			ModifiedFrom: StartDate,
			Offset:       d,
			Count:        rt.E.Metrics.MaxBatchSize,
		}
		rqt := goengage.SupporterSearchRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}
		var resp goengage.SupporterSearchResults
		n := goengage.NetOp{
			Host:     rt.E.Host,
			Endpoint: goengage.SearchSupporter,
			Method:   goengage.SearchMethod,
			Token:    rt.E.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err = n.Do()
		if err != nil {
			return err
		}
		for _, s := range resp.Payload.Supporters {
			found := false
			for _, c := range s.CustomFieldValues {
				if !found && c.Name == rt.Name {
					x := NewCFRecord(s, c)
					rt.C2 <- x
				}
			}
		}
	}
	rt.D <- true
	log.Printf("Supporters %d: begin\n", i)
	return nil
}

//TotalRecords reads the first page of supporters and returns
//the number of records.
func TotalRecords(rt Runtime) (int32, error) {
	payload := goengage.SupporterSearchRequestPayload{
		ModifiedFrom: StartDate,
		Offset:       int32(0),
		Count:        rt.E.Metrics.MaxBatchSize,
	}
	rqt := goengage.SupporterSearchRequest{
		Header:  goengage.RequestHeader{},
		Payload: payload,
	}
	var resp goengage.SupporterSearchResults
	n := goengage.NetOp{
		Host:     rt.E.Host,
		Endpoint: goengage.SearchSupporter,
		Method:   goengage.SearchMethod,
		Token:    rt.E.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err := n.Do()
	if err != nil {
		return 0, err
	}
	return resp.Payload.Total, nil
}

//OutputCSV accepts CF records from a channel and and writes them to
//a CSV file.
func OutputCSV(rt Runtime) error {
	log.Printf("OutputCSV: begin")
	f, err := os.Create(rt.F)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	headers := []string{
		"SupporterID",
		"Email",
		"CustomFieldID",
		"CustomFieldName",
		"CustomFieldValue",
	}
	w.Write(headers)
	for {
		x, ok := <-rt.C2
		if !ok {
			break
		}
		a := []string{
			x.SupporterID,
			x.Email,
			x.CustomFieldID,
			x.CustomFieldName,
			x.CustomFieldValue,
		}
		w.Write(a)
		w.Flush()
	}
	w.Flush()
	log.Printf("OutputCSV: end")
	return nil
}

//WaitTerminations waits for "OffsetListeners" supporter readers to
//complete.  That triggers a close for the CSV writer channel.
func WaitTerminations(rt Runtime) {
	log.Printf("WaitTerminations: begin")
	remaining := OffsetListeners
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
		app         = kingpin.New("find_custom_field", "Find supporters that have values for the provided cusotm field.")
		login       = app.Flag("login", "YAML file with API token").Required().String()
		customField = app.Flag("custom_field", "Search for this custom field").Required().String()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		log.Fatalf("Error --login is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("Error: %+v\n", e)
	}

	// logger, err := goengage.NewUtilLogger()
	// if err != nil {
	// 	panic(err)
	// }

	rt := Runtime{
		E:    e,
		Name: *customField,
		C0:   make(chan int32, 1000),
		C1:   make(chan goengage.Supporter, 20_000),
		C2:   make(chan *CFRecord, 100),
		D:    make(chan bool, OffsetListeners),
		F:    "find_custom_field.csv",
		// L:         logger,
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

	//Start offset listeners
	for id := 1; id <= OffsetListeners; id++ {
		wg.Add(1)
		go (func(rt Runtime, id int, wg *sync.WaitGroup) {
			defer wg.Done()
			err := Supporters(rt, id)
			if err != nil {
				panic(err)
			}
		})(rt, id, &wg)
	}
	log.Printf("main: %+v offset listeners started\n", OffsetListeners)

	//Start "done" listener to keep track of offsets listeners.
	wg.Add(1)
	go (func(rt Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		WaitTerminations(rt)
	})(rt, &wg)
	log.Println("main: terminations listener started")

	total, err := TotalRecords(rt)
	if err != nil {
		panic(err)
	}

	//Start offset seeder
	log.Printf("main: processing %d offsets\n", total)
	wg.Add(1)
	go (func(rt Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		Offsets(rt, total)
	})(rt, &wg)
	log.Printf("main: offset seeder started\n")

	log.Printf("main: napping...\n")
	time.Sleep(10 * time.Second)
	log.Printf("main: waiting...\n")
	wg.Wait()
	log.Printf("main: done")
}
