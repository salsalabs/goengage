package goengage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
)

//Service provides the basic tools to read and filter records then
//write them to a CSV file.
type Service interface {
	//WhichActivity returns the kind of activity being read.
	WhichActivity() string
	//Filter returns true if the record should be used.
	Filter() bool
	//Headers returns column headers for a CSV file.
	Headers() []string
	//Line returns a list of strings to go in to the CSV file for each
	//fundraising record.
	Line() []string
	//Readers returns the number of readers to start.
	Readers() int
	//Filename returns the CSV filename.
	Filename() string
}

//MaxRecords returns the maximum number of activity records
//of a particular type.
func MaxRecords(e *Environment, s Service, start string, end string) (int32, error) {
	resp, err := ReadBatch(e, s, int32(0), int32(0), start, end)
	if err != nil {
		return int32(0), err
	}
	return resp.Payload.Total, err
}

//ReadActivities retrieves activity records from Engage, filters them,
//then writes them to the Service channel. The offset channel tells
//us where to start reading.  When no items are available from the
//offset channel, we'll write a true to the done channel.
func ReadActivities(e *Environment,
	s Service,
	i int,
	oc chan int32,
	sc chan Service,
	dc chan bool,
	start string,
	end string) {

	n := fmt.Sprintf("ReadActivities-%d", i)
	log.Printf("%s: begin", n)
	for true {
		offset, ok := <-oc
		if !ok {
			break
		}
		resp, err := ReadBatch(e, s, offset, e.Metrics.MaxBatchSize, start, end)
		if err != nil {
			// panic(err)
			log.Printf("%s: offset %6d error %s\n", n, offset, err)
			break
		}
		if !ok {
			break
		}
		if resp.Payload.Count == 0 {
			break
		}
		pass := int32(0)
		for _, r := range resp.Payload.Activities {
			if r.Filter() {
				sc <- r
				pass++
			}
		}
		log.Printf("%s: offset %6d, %3d passed + %3d skipped = %3d\n",
			n,
			offset,
			pass,
			resp.Payload.Count-pass,
			resp.Payload.Count)
	}
	dc <- true
	log.Printf("%s: end", n)
}

//ReadBatch is a utility function to read activity records. Returns the
//response object and an error code.
func ReadBatch(e *Environment,
	s Service,
	offset int32,
	count int32,
	start string,
	end string) (resp *FundraiseResponse, err error) {

	// log.Printf("ReadBatch: offset %d, count %d, start %v, end %v\n", offset, count, start, end)
	payload := ActivityRequestPayload{
		Type:         s.WhichActivity(),
		Offset:       offset,
		Count:        count,
		ModifiedFrom: start,
		ModifiedTo:   end,
	}
	rqt := ActivityRequest{
		Header:  RequestHeader{},
		Payload: payload,
	}
	n := NetOp{
		Host:     e.Host,
		Method:   SearchMethod,
		Endpoint: SearchActivity,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err = n.Do()
	return resp, err
}

// ReportFundraising on a Service by reading all records, filtering, then
// writing survivors to a CSV file.
func ReportFundraising(e *Environment, s Service, start string, end string) (err error) {
	sc := make(chan Service, 100)
	dc := make(chan bool)
	oc := make(chan int32, 100)
	var wg sync.WaitGroup

	//Start the reader waiter.  It waits until all readers are done,
	//then closes the service channel.
	go (func(s Service, sc chan Service, done chan bool, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		WaitForReaders(s, sc, done)
	})(s, sc, dc, &wg)

	//Start the CSV writer. It receives Service records from readers and
	//writes them to a CSV.
	go (func(s Service, sc chan Service, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := WriteCSV(s, sc)
		if err != nil {
			panic(err)
		}
	})(s, sc, &wg)

	//Start the readers. Readers get offsets from the offset channel, read activities,
	//filter them, then put them onto the Service channel.
	for i := 0; i < s.Readers(); i++ {
		go (func(e *Environment,
			s Service,
			i int,
			oc chan int32,
			sc chan Service,
			done chan bool,
			start string,
			end string,
			wg *sync.WaitGroup) {

			wg.Add(1)
			defer wg.Done()
			ReadActivities(e, s, i, oc, sc, dc, start, end)
		})(e, s, i, oc, sc, dc, start, end, &wg)
	}

	// Push offsets onto the offset channel.
	maxRecords, err := MaxRecords(e, s, start, end)
	log.Printf("ReportFundraising: processing %d %s records\n", maxRecords, FundraiseType)
	maxRecords = maxRecords + int32(e.Metrics.MaxBatchSize-1)
	for offset := int32(0); offset <= maxRecords; offset += e.Metrics.MaxBatchSize {
		oc <- offset
	}
	close(oc)

	//Wait...
	log.Printf("ReportFundraising: waiting for terminations")
	wg.Wait()
	log.Printf("ReportFundraising done")
	return err
}

//WaitForReaders waits for readers to send to a done channel.
//The number of readers is specified in the provided Service.
//Closes the csv channel when all readers are done.
func WaitForReaders(s Service, sc chan Service, done chan bool) {
	count := s.Readers()
	for count > 0 {
		log.Printf("WaitForReaders: Waiting for %d readers\n", count)
		_, ok := <-done
		if !ok {
			break
		}
		count--
	}
	close(sc)
	log.Println("WaitForReaders: done")
}

//WriteCSV waits for Service records to appear on the queue, then
//Writes them to a CSV file.
func WriteCSV(s Service, sc chan Service) error {
	log.Println("WriteCSV: begin")
	f, err := os.Create(s.Filename())
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.Write(s.Headers())

	for true {
		r, ok := <-sc
		if !ok {
			break
		}
		w.Write(r.Line())
		w.Flush()
	}
	f.Close()
	log.Println("WriteCSV: done")
	return err
}
