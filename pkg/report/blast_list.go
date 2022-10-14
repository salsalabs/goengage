package goengage

//The report.blast_list package reads email blast information from the list-of-blasts
// endpoint in the Web Developer API.  Once a developer completes the interface, then
// the resulting program is fairly straightforward and independent. The output is
// a CSV of blast information, including timestamps and URLs.

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
)

// ChannelCount is the number of "done" notifications to wait for before terminate.
const ChannelCount = 3

// BlastListGuide is the interface to use when scanning all email blasts
// and doing something.
type BlastListGuide interface {

	//VisitResult does something with the blast. Errors terminate.
	VisitResult(s goengage.BlastListResult) error

	//VisitContent does something with the blast. Errors terminate.
	VisitContent(s goengage.BlastListResult, t goengage.BlastContent) error

	//Finalize is called after all blasts have been processed.
	Finalize() error

	//Payload is a convenience method to define which blasts to return.
	// Each item is turned into a URL query at execution time.
	Payload() goengage.BlastListRequest

	//ResultChannel is the listener channel for blast info.
	ResultChannel() chan goengage.BlastListResult

	//DoneChannel() receives a true when the listener is done.
	DoneChannel() chan bool

	//Offset() returns the offset to start reading.  Useful for
	//restarting after a service interruption.

	Offset() int32
}

// readBlastLists reads all blasts and pushes them onto a channel.
// Probably a good idea to start this as a go routine after the Listener
// is started...
func readBlastLists(e *goengage.Environment, g BlastListGuide) error {
	log.Println("ReadBlastLists: start")
	count := int32(e.Metrics.MaxBatchSize)
	offset := int32(g.Offset())
	for count == int32(e.Metrics.MaxBatchSize) {

		// --------------------------------------
		// TODO: Turn payload into URL + queries.
		// --------------------------------------

		payload := g.Payload()
		payload.Offset = offset
		payload.Count = count
		v := url.Values{}
		if len(payload.StartDate) != 0 {
			v.Set("startDate", payload.StartDate)
		}
		if len(payload.EndDate) != 0 {
			v.Set("endDate", payload.EndDate)
		}
		if len(payload.Criteria) != 0 {
			v.Set("criteria", payload.Criteria)
		}
		if len(payload.SortField) != 0 {
			v.Set("sortField", payload.SortField)
		}
		if len(payload.SortOrder) != 0 {
			v.Set("sortOrder", payload.SortOrder)
		}
		v.Set("count", fmt.Sprintf("%v", payload.Count))
		v.Set("offset", fmt.Sprintf("%v", payload.Offset))
		queries := v.Encode()

		endpoint := goengage.EmailBlastList + "&" + queries

		var resp goengage.BlastListResponse
		n := goengage.NetOp{
			Host:     e.Host,
			Method:   goengage.EnquireMethod,
			Endpoint: endpoint,
			Token:    e.Token,
			Response: &resp,
		}
		err := n.Do()
		if err != nil {
			return err
		}
		count = resp.Payload.Count
		log.Printf("ReadBlastLists: offset %5d, read %2d\n", offset, count)
		for _, s := range resp.Payload.Results {
			g.ResultChannel() <- s
		}
		offset += resp.Payload.Count
	}
	log.Println("ReadBlastLists: done")
	close(g.ResultChannel())
	g.DoneChannel() <- true
	return nil
}

// handleResults reads from the result channel and calls the result
// visitor.  When that returns, then the content visitor is called for
// each content item in the result.
func handleResults(e *goengage.Environment, g BlastListGuide) error {
	log.Println("ProcessBlastLists: start")
	for {
		s, ok := <-g.ResultChannel()
		if !ok {
			break
		}
		g.VisitResult(s)
		for _, c := range s.Content {
			g.VisitContent(s, c)
		}
	}
	close(g.ResultChannel())
	g.DoneChannel() <- true
	return nil
}

func ReportBlastLists(e *goengage.Environment, g BlastListGuide) error {
	var wg sync.WaitGroup
	log.Println("ReportBlastLists: start")

	// Start the results listener.
	go (func(e *goengage.Environment, g BlastListGuide, wg *sync.WaitGroup) {
		defer wg.Done()
		handleResults(e, g)
	})(e, g, &wg)

	// Start the reader.
	go (func(e *goengage.Environment, g BlastListGuide, wg *sync.WaitGroup) {
		defer wg.Done()
		readBlastLists(e, g)
	})(e, g, &wg)

	// Wait for things to finish.
	WaitForWorkers(ChannelCount, g.DoneChannel())
	g.Finalize()
	g.DoneChannel() <- true
	log.Println("ReportBlastLists: end")
	return nil
}

// WaitForWorkers waits for readers to send to a done channel.
func WaitForWorkers(count int32, done chan bool) {
	for count > 0 {
		log.Printf("WaitForWorkers: Waiting for %d readers\n", count)
		_, ok := <-done
		if !ok {
			break
		}
		count--
	}
	log.Println("WaitForWorkers: done")
}
