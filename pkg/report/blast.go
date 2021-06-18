package goengage

//The report.blast package reads email blasts using an interface-provided payload.
//Next, this app collects delivery history for each recipient of the blast.
//Those records are written to a channel provided in an interface-specied
//function.  A separate listener (or many listeners, for that matter) retrieves
//supporters from the channel and calls the Visit function from the interface.
//This continues until the channel is closed.  The listener calls Finalize before
//terminating and puts a true onto DoneChannel.

import (
	"log"

	goengage "github.com/salsalabs/goengage/pkg"
)

//EmailBlastGuide is the interface to use when scanning all supporters
//and doing something.
type EmailBlastGuide interface {

	//Visit does something with the blast. Errors terminate.
	Visit(s goengage.EmailActivity) error

	//Finalize is called after all blasts have been processed.
	Finalize() error

	//Payload is the request payload defining which supporters to retrieve.
	Payload() goengage.EmailBlastSearchRequestPayload

	//Channel is the listener channel to use.
	Channel() chan goengage.EmailActivity

	//DoneChannel() receives a true when the listener is done.
	DoneChannel() chan bool

	//Offset() returns the offset to start reading.  Useful for
	//restarting after a service interruption.

	Offset() int32
}

//ReadEmailBlasts reads all blasts and pushes them onto a channel.
//Probably a good idea to start this as a go routine after the Listener
//is started...
func ReadEmailBlasts(e *goengage.Environment, g EmailBlastGuide) error {
	log.Println("ReadEmailBlasts: start")
	count := int32(e.Metrics.MaxBatchSize)
	offset := int32(g.Offset())
	for count == int32(e.Metrics.MaxBatchSize) {
		payload := g.Payload()
		payload.Offset = offset
		payload.Count = count
		rqt := goengage.EmailBlastSearchRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}
		var resp goengage.EmailBlastSearchResponse
		n := goengage.NetOp{
			Host:     e.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.EmailBlastSearch,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err := n.Do()
		if err != nil {
			return err
		}
		count = resp.Payload.Count
		log.Printf("ReadEmailBlasts: offset %d", offset)
		for _, s := range resp.Payload.EmailActivities {
			g.Channel() <- s
		}
		offset += count
	}
	log.Println("ReadEmailBlasts: done")
	close(g.Channel())
	return nil
}

//ProcessEmailBlasts reads blasts from an interface-provided channel, then
//calls Visit in the interface.  At end of data, the app calls Finalize() then
//sends true to the DoneChannel.
func ProcessEmailBlasts(e *goengage.Environment, g EmailBlastGuide) error {
	log.Println("ProcessEmailBlasts: start")
	for true {
		s, ok := <-g.Channel()
		if !ok {
			break
		}
		g.Visit(s)
	}
	g.Finalize()
	g.DoneChannel() <- true
	log.Println("ProcessEmailBlasts: end")
	return nil
}
