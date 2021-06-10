//App to search for supporters whose addressLine1 and city fields
//contain the Zip code.  Each matching record is modified to erase
//addressLine1 and City fields.
package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	reporter "github.com/salsalabs/goengage/pkg/report"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Runtime area for this app.
type Runtime struct {
	E          *goengage.Environment
	InChan     chan goengage.Supporter
	DoneChan   chan bool
	Keys       []string
	FieldName  string
	ReadOffset int32
}

//NewRuntime populates a new runtime.
func NewRuntime(env *goengage.Environment, offset int32) Runtime {
	r := Runtime{
		E:          env,
		InChan:     make(chan goengage.Supporter),
		DoneChan:   make(chan bool),
		ReadOffset: offset,
	}
	return r
}

//Visit implements SupporterGuide.Visit and does something with
//a supporter record. In this case, the kludged addressLine1
//and city are cleared away.
func (r *Runtime) Visit(s goengage.Supporter) error {
	if s.Address != nil {
		a := s.Address
		if a.AddressLine1 == a.City && a.City == a.PostalCode {
			a.AddressLine1 = ""
			a.City = ""
			log.Printf("%+v", s.Address)
		}
	}
	return nil
}

//Finalize implements SupporterGuide.Filnalize and outputs the
//distribution results.
func (r *Runtime) Finalize() error {
	return nil
}

//Payload implements SupporterGuide.Payload and provides a payload
//that will retrieve all supporters.
func (r *Runtime) Payload() goengage.SupporterSearchRequestPayload {
	payload := goengage.SupporterSearchRequestPayload{
		IdentifierType: goengage.SupporterIDType,
		ModifiedFrom:   "2000-01-01T00:00:00.00000Z",
		ModifiedTo:     "2050-01-01T00:00:00.00000Z",
		Offset:         0,
		Count:          0,
	}
	return payload
}

//Channel implements SupporterGuide.Channnel and provides the
//supporter channel.
func (r *Runtime) Channel() chan goengage.Supporter {
	return r.InChan
}

//DoneChannel implements SupporterGuide.DoneChannel to provide
// a channel that  receives a true when the listener is done.
func (r *Runtime) DoneChannel() chan bool {
	return r.DoneChan
}

//Offset returns the offset for the first read.
//Useful for restarts.
func (r *Runtime) Offset() int32 {
	return r.ReadOffset
}

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app    = kingpin.New("custom_field-distribution", "Find and fix supporter records with malformed addressLine1 and City")
		login  = app.Flag("login", "YAML file with API token").Required().String()
		offset = app.Flag("offset", "Read offset. Useful when network goes away").Default("0").Int32()
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

	r := NewRuntime(e, *offset)
	var wg sync.WaitGroup

	//Start supporter listener. Only one of these because Visit is quick
	//in this app. More than one cases "concurrent map writes" errors.
	go (func(e *goengage.Environment, r *Runtime, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		reporter.ProcessSupporters(r.E, r)
	})(e, &r, &wg)

	//Start done listener.
	go (func(r *Runtime, n int, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		goengage.DoneListener(r.DoneChan, n)
	})(&r, 1, &wg)

	//Start supporter reader.
	go (func(e *goengage.Environment, r *Runtime, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		reporter.ReadSupporters(r.E, r)
	})(e, &r, &wg)

	d, err := time.ParseDuration("10s")
	if err != nil {
		panic(err)
	}
	log.Printf("main: sleeping for %s", d)
	time.Sleep(d)
	log.Printf("main:  waiting...")
	wg.Wait()
	log.Printf("main: done")
}
