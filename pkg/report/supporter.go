package goengage

//The report.supporter package reads supporter records using an interface-specified payload.
//Those records are written to a supporter channel provided in an interface-specied
//function.  A separate listener (or many listeners, for that matter) retrieves
//supporters from the channel and calls the Visit function from the interface.
//This continues until the channel is closed.  The listener calls Finalize before
//terminating and puts a true onto DoneChannel.

import (
	"log"

	goengage "github.com/salsalabs/goengage/pkg"
)

// SupporterGuide is the interface to use when scanning all supporters
// and doing something.
type SupporterGuide interface {

	//Visit does something with the supporter. Errors terminate.
	Visit(s goengage.Supporter) error

	//Finalize is called after all supporters have been processed.
	Finalize() error

	//Payload is the request payload defining which supporters to retrieve.
	Payload() goengage.SupporterSearchRequestPayload

	//Channel is the listener channel to use.
	Channel() chan goengage.Supporter

	//DoneChannel() receives a true when the listener is done.
	DoneChannel() chan bool

	//Offset() returns the offset to start reading.  Useful for
	//restarting after a service interruption.
	Offset() int32

	//AdjustOffset() allows you to change the offset just before the
	//read. This is most helpful when the payload contains a list of
	//ID's that need to be read in chunks. This method allows you to
	//set the offset to zero, thereby avoiding confusion.
	AdjustOffset(offset int32) int32
}

// ReadSupporters reads all supporters and pushes them onto a channel.
// Probably a good idea to start this as a go routine after the Listener
// is started...
func ReadSupporters(e *goengage.Environment, g SupporterGuide) error {
	log.Println("ReadSupporters: start")
	count := int32(e.Metrics.MaxBatchSize)
	offset := int32(g.Offset())
	for count == int32(e.Metrics.MaxBatchSize) {
		payload := g.Payload()
		payload.Offset = g.AdjustOffset(offset)
		payload.Count = count
		rqt := goengage.SupporterSearchRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}
		var resp goengage.SupporterSearchResults
		n := goengage.NetOp{
			Host:     e.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.SearchSupporter,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err := n.Do()
		if err != nil {
			return err
		}
		count = resp.Payload.Count
		log.Printf("ReadSupporters: offset %d\n", offset)
		for _, s := range resp.Payload.Supporters {
			g.Channel() <- s
		}
		offset += count
	}
	log.Println("ReadSupporters: done")
	close(g.Channel())
	return nil
}

// ProcessSupporters reads supporters from an interface-provided channel, then
// calls Visit in the interface.  At end of data, the app calls Finalize() then
// sends true to the DoneChannel.
func ProcessSupporters(e *goengage.Environment, g SupporterGuide) error {
	log.Println("ProcessSupporters: start")
	for {
		s, ok := <-g.Channel()
		if !ok {
			break
		}
		g.Visit(s)
	}
	g.Finalize()
	g.DoneChannel() <- true
	log.Println("ProcessSupporters: end")
	return nil
}
