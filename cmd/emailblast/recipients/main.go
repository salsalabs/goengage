package main

// An application to read all email blasts and write blast and
//write blast and recipient activity to CSVs.  This should give
//an Engage client a way to create history in a data warehouse.
import (
	goengage "github.com/salsalabs/goengage/pkg"
)

//RecipientRecord contains both blast and recipient info.
type RecipientRecord struct {
	Recipient goengage.SingleBlastRecipient
	Blast     goengage.EmailActivity
}

//Runtime contains the configuration parts that this app needs.
type Runtime struct {
	BlastChan       chan goengage.EmailActivity
	BlastCSVChan    chan goengage.EmailActivity
	ConversionChan  chan goengage.Conversion
	RecipientChan   chan RecipientRecord
	DoneChan        chan bool
	BlastOffset     int32
	RecipientOffset int32
	BlastCursor     *string
}

//Bongos are drums.
type Bongos struct {
	Count int32
}

//These declarations are used by the goengage blast report
//handler (goengage.ProcessEmailBlasts) to process blasts
//and history.

//Visit some bongos.
func (b *Bongos) Visit() {

}

//Visit does something with the blast. Errors terminate.
//Implements goengage.EmailBlastGuide.
func (rt *Runtime) Visit(s goengage.EmailActivity) error {
	rt.BlastCSVChan <- s
	return nil
}

//Finalize is called after all blasts have been processed.
//Implements goengage.EmailBlastGuide.
func (rt *Runtime) Finalize() error {
	return nil
}

//Payload is the request payload defining which supporters to retrieve.
//Implements goengage.EmailBlastGuide.
func (rt *Runtime) Payload() goengage.EmailBlastSearchRequestPayload {
	payload := goengage.EmailBlastSearchRequestPayload{}
	return payload
}

//Channel is the listener channel to use.
func (rt *Runtime) Channel() chan goengage.EmailActivity {
	return rt.BlastChan
}

//DoneChannel receives a true when the listener is done.
//Implements goengage.EmailBlastGuide.
func (rt *Runtime) DoneChannel() chan bool {
	return rt.DoneChan
}

//Offset returns the offset to start reading.
//Implements goengage.EmailBlastGuide.
func (rt *Runtime) Offset() int32 {
	return rt.BlastOffset
}
