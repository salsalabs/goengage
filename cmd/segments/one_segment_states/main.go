//App to retrieve segment (group) members and write their State fields
//to a CSV.  Raw data for segment-based demographic analysis.
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//ListenerCount is the number of go tasks listening for
//supporters.  Trivial in this version.
const ListenerCount = 1

//Cache is the content that appears in the CSV file.
type Cache map[string]int32

//Runtime area for this app.
type Runtime struct {
	E         *goengage.Environment
	InChan    chan goengage.Supporter
	Cache     Cache
	DoneChan  chan bool
	SegmentID string
	Results   string
	Logger    *goengage.UtilLogger
}

//NewRuntime populates a new runtime.
func NewRuntime(env *goengage.Environment, segmentID string, results string, verbose bool) Runtime {
	r := Runtime{
		E:         env,
		InChan:    make(chan goengage.Supporter),
		Cache:     make(Cache),
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

//Drive finds all of the supporters in the specified group and writes
//them to the input channel.
func (rt Runtime) Drive() (err error) {
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

//SaveCache writes an alphabetized CSV of states and counts to
//the output CSV file.
func (rt Runtime) SaveCache() (err error) {
	f, err := os.Create(rt.Results)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	headers := []string{
		"State",
		"Members",
	}
	err = writer.Write(headers)
	if err != nil {
		return err
	}

	keys := make([]string, len(rt.Cache))
	for k := range rt.Cache {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := rt.Cache[k]
		record := []string{
			k,
			fmt.Sprintf("%v", v),
		}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

//Update retrieves Supporters from the input channel and updates
//the cache in the runtime.  Writes the CSV at end of data.
func (rt Runtime) Update() (err error) {
	log.Println("Update: start")
	for {
		s, okay := <-rt.InChan
		if !okay {
			break
		}
		if s.Address != nil {
			v, ok := rt.Cache[s.Address.State]
			if !ok {
				v = 0
				rt.Cache[s.Address.State] = v
			}
			rt.Cache[s.Address.State] = v + 1
		}
	}
	rt.SaveCache()
	log.Println("Update: end")
	rt.DoneChan <- true
	return nil
}

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app       = kingpin.New("one_segment_states", "Tabulate segment member counts by state")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		segmentID = app.Flag("segment-id", "segmentID for the group").Default("f4be4a19-b85f-4d69-baae-e027a86fd676").String()
		results   = app.Flag("results", "filename of CSV file to record results").Default("one_segment_states.csv").String()
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

	//Start Update task. More than one leads to multiple CSV files.
	for i := 1; i <= ListenerCount; i++ {
		wg.Add(1)
		go (func(rt Runtime, wg *sync.WaitGroup, i int) {
			defer wg.Done()
			err := rt.Update()
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
	})(rt, ListenerCount, &wg)

	//Start supporter reader.
	wg.Add(1)
	go (func(rt Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := rt.Drive()
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
