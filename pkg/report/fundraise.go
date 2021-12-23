package goengage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
)

//MaxRecords returns the maximum number of activity records
//of a particular type.
func MaxRecords(e *goengage.Environment, guide Guide, ts TimeSpan) (int32, error) {
	resp, err := ReadBatch(e, guide, int32(0), ts)
	if err != nil {
		return int32(0), err
	}
	return resp.Payload.Total, err
}

//ReadActivities retrieves activity records from Engage, filters them,
//then writes them to the Guide channel. The offset channel tells
//us where to start reading.  When no items are available from the
//offset channel, we'll write a true to the done channel.
func ReadActivities(e *goengage.Environment,
	guide Guide,
	i int,
	oc chan int32,
	gc chan goengage.Fundraise,
	dc chan bool,
	ts TimeSpan) {

	n := fmt.Sprintf("ReadActivities-%d", i)
	log.Printf("%s: begin", n)
	for {
		offset, ok := <-oc
		if !ok {
			break
		}
		resp, err := ReadBatch(e, guide, offset, ts)
		if err != nil {
			// panic(err)
			log.Printf("%s: offset %6d error %s\n", n, offset, err)
			break
		}
		if resp.Payload.Count == 0 {
			break
		}
		pass := int32(0)
		total := resp.Payload.Total
		for _, r := range resp.Payload.Activities {
			if guide.Filter(r) {
				s, err := goengage.SupporterByID(e, r.SupporterID)
				if err != nil {
					panic(err)
				}
				r.Supporter = *s
				gc <- r
				pass++
			}
		}
		log.Printf("%s: offset %6d of %6d, %3d adds\n", n, offset, total, pass)
	}
	dc <- true
	log.Printf("%s: end", n)
}

//ReadBatch is a utility function to read activity records. Returns the
//response object and an error code.
func ReadBatch(e *goengage.Environment,
	guide Guide,
	offset int32,
	ts TimeSpan) (resp *goengage.FundraiseResponse, err error) {

	payload := goengage.ActivityRequestPayload{
		Type:         guide.TypeActivity(),
		Offset:       offset,
		Count:        e.Metrics.MaxBatchSize,
		ModifiedFrom: ts.Start,
		ModifiedTo:   ts.End,
	}
	rqt := goengage.ActivityRequest{
		Header:  goengage.RequestHeader{},
		Payload: payload,
	}
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   goengage.SearchMethod,
		Endpoint: goengage.SearchActivity,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err = n.Do()
	return resp, err
}

// ReportFundraising on a Guide by reading all records, filtering, then
// writing survivors to a CSV file.
func ReportFundraising(e *goengage.Environment, guide Guide, ts TimeSpan) (err error) {
	gc := make(chan goengage.Fundraise, 100)
	dc := make(chan bool)
	oc := make(chan int32, 100)
	var wg sync.WaitGroup

	//Start the reader waiter.  It waits until all readers are done,
	//then closes the fundraise channel.
	wg.Add(1)
	go (func(guide Guide, gc chan goengage.Fundraise, done chan bool, wg *sync.WaitGroup) {
		defer wg.Done()
		WaitForReaders(guide, gc, done)
	})(guide, gc, dc, &wg)

	//Start the CSV writer. It receives fundraise records from readers and
	//writes them to a CSV.
	wg.Add(1)
	go (func(guide Guide, gc chan goengage.Fundraise, wg *sync.WaitGroup) {
		defer wg.Done()
		err := Store(guide, gc)
		if err != nil {
			panic(err)
		}
	})(guide, gc, &wg)

	//Start the readers. Readers get offsets from the offset channel, read activities,
	//filter them, then put them onto the Guide channel.
	for i := 0; i < guide.Readers(); i++ {
		wg.Add(1)
		go (func(e *goengage.Environment,
			guide Guide,
			i int,
			oc chan int32,
			gc chan goengage.Fundraise,
			done chan bool,
			ts TimeSpan,
			wg *sync.WaitGroup) {
			defer wg.Done()
			ReadActivities(e, guide, i, oc, gc, dc, ts)
		})(e, guide, i, oc, gc, dc, ts, &wg)
	}

	// Push offsets onto the offset channel.
	maxRecords, err := MaxRecords(e, guide, ts)
	log.Printf("ReportFundraising: reporting on start time %s\n", ts.Start)
	log.Printf("ReportFundraising:              end   time %s\n", ts.End)
	log.Printf("ReportFundraising: %d donations\n", maxRecords)
	maxRecords = maxRecords + int32(e.Metrics.MaxBatchSize-1)
	for offset := int32(guide.Offset()); offset <= maxRecords; offset += e.Metrics.MaxBatchSize {
		oc <- offset
	}
	close(oc)

	//Wait...
	log.Printf("ReportFundraising: waiting for terminations")
	wg.Wait()
	log.Printf("ReportFundraising: done")
	return err
}

//WaitForReaders waits for readers to send to a done channel.
//The number of readers is specified in the provided Guide.
//Closes the csv channel when all readers are done.
func WaitForReaders(guide Guide, gc chan goengage.Fundraise, done chan bool) {
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

//Store waits for Guide records to appear on the queue, then
//Writes them to a CSV file.
func Store(guide Guide, gc chan goengage.Fundraise) error {
	log.Println("Store: begin")
	f, err := os.Create(guide.Filename())
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.Write(guide.Headers())

	for {
		r, ok := <-gc
		if !ok {
			break
		}
		w.Write(guide.Line(r))
		w.Flush()
	}
	f.Close()
	log.Println("Store: done")
	return err
}
