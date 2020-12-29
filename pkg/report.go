package goengage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
)

//Guide provides the basic tools to read and filter records then
//write them to a CSV file.
type Guide interface {
	//WhichActivity returns the kind of activity being read.
	WhichActivity() string
	//Filter returns true if the record should be used.
	Filter(Fundraise) bool
	//Headers returns column headers for a CSV file.
	Headers() []string
	//Line returns a list of strings to go in to the CSV file for each
	//fundraising record.
	Line(Fundraise) []string
	//Readers returns the number of readers to start.
	Readers() int
	//Filename returns the CSV filename.
	Filename() string
}

//MaxRecords returns the maximum number of activity records
//of a particular type.
func MaxRecords(e *Environment, guide Guide, start string, end string) (int32, error) {
	resp, err := ReadBatch(e, guide, int32(0), int32(0), start, end)
	if err != nil {
		return int32(0), err
	}
	return resp.Payload.Total, err
}

//ReadActivities retrieves activity records from Engage, filters them,
//then writes them to the Guide channel. The offset channel tells
//us where to start reading.  When no items are available from the
//offset channel, we'll write a true to the done channel.
func ReadActivities(e *Environment,
	guide Guide,
	i int,
	oc chan int32,
	gc chan Fundraise,
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
		resp, err := ReadBatch(e, guide, offset, e.Metrics.MaxBatchSize, start, end)
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
			if guide.Filter(r) {
				log.Printf("ReadActivities: r %+v\n", r)
				s, err := ReadSupporter(e, r)
				if err != nil {
					panic(err)
				}
				r.Supporter = *s
				gc <- r
				pass++
			}
		}
		log.Printf("%s: offset %6d, %3d of %3d passed\n",
			n,
			offset,
			pass,
			resp.Payload.Count)
	}
	dc <- true
	log.Printf("%s: end", n)
}

//ReadBatch is a utility function to read activity records. Returns the
//response object and an error code.
func ReadBatch(e *Environment,
	guide Guide,
	offset int32,
	count int32,
	start string,
	end string) (resp *FundraiseResponse, err error) {

	// log.Printf("ReadBatch: offset %d, count %d, start %v, end %v\n", offset, count, start, end)
	payload := ActivityRequestPayload{
		Type:         guide.WhichActivity(),
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

//ReadSupporter reads a supporter record for the specified ID.
func ReadSupporter(e *Environment, f Fundraise) (supporter *Supporter, err error) {
	payload := SupporterSearchPayload{
		Identifiers:    []string{f.Supporter.SupporterID},
		IdentifierType: SupporterIDType,
	}

	rqt := SupporterSearch{
		Header:  RequestHeader{},
		Payload: payload,
	}
	log.Printf("ReadSupporter: request %+v\n", rqt)
	var resp SupporterSearchResults
	n := NetOp{
		Host:     e.Host,
		Method:   SearchMethod,
		Endpoint: SearchSupporter,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err = n.Do()
	if err != nil {
		return supporter, err
	}
	log.Printf("ReadSupporter: response %+v\n", resp)
	s := resp.Payload.Supporters[0]
	log.Printf("supporter: %+v\n", s)
	return &s, err
}

// ReportFundraising on a Guide by reading all records, filtering, then
// writing survivors to a CSV file.
func ReportFundraising(e *Environment, guide Guide, start string, end string) (err error) {
	gc := make(chan Fundraise, 100)
	dc := make(chan bool)
	oc := make(chan int32, 100)
	var wg sync.WaitGroup

	log.Println("")
	log.Println("-------------------------------------------------------------------------------------------------")
	log.Println("ReportFundraising:  This version reads Supporter because Fundraise does not have the right data.")
	log.Println("-------------------------------------------------------------------------------------------------")
	log.Println("")

	//Start the reader waiter.  It waits until all readers are done,
	//then closes the service channel.
	go (func(guide Guide, gc chan Fundraise, done chan bool, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		WaitForReaders(guide, gc, done)
	})(guide, gc, dc, &wg)

	//Start the CSV writer. It receiveguide Guide records from readers and
	//writes them to a CSV.
	go (func(guide Guide, gc chan Fundraise, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := WriteCSV(guide, gc)
		if err != nil {
			panic(err)
		}
	})(guide, gc, &wg)

	//Start the readers. Readers get offsets from the offset channel, read activities,
	//filter them, then put them onto the Guide channel.
	for i := 0; i < guide.Readers(); i++ {
		go (func(e *Environment,
			guide Guide,
			i int,
			oc chan int32,
			gc chan Fundraise,
			done chan bool,
			start string,
			end string,
			wg *sync.WaitGroup) {

			wg.Add(1)
			defer wg.Done()
			ReadActivities(e, guide, i, oc, gc, dc, start, end)
		})(e, guide, i, oc, gc, dc, start, end, &wg)
	}

	// Push offsets onto the offset channel.
	maxRecords, err := MaxRecords(e, guide, start, end)
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
//The number of readers is specified in the provided Guide.
//Closes the csv channel when all readers are done.
func WaitForReaders(guide Guide, gc chan Fundraise, done chan bool) {
	count := guide.Readers()
	for count > 0 {
		log.Printf("WaitForReaders: Waiting for %d readers\n", count)
		_, ok := <-done
		if !ok {
			break
		}
		count--
	}
	close(gc)
	log.Println("WaitForReaders: done")
}

//WriteCSV waits for Guide records to appear on the queue, then
//Writes them to a CSV file.
func WriteCSV(guide Guide, gc chan Fundraise) error {
	log.Println("WriteCSV: begin")
	f, err := os.Create(guide.Filename())
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.Write(guide.Headers())

	for true {
		r, ok := <-gc
		if !ok {
			break
		}
		w.Write(guide.Line(r))
		w.Flush()
	}
	f.Close()
	log.Println("WriteCSV: done")
	return err
}
