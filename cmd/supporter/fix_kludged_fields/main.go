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

const (
	// Number of input queue listeners.
	ListenerCount = 5
)

//Runtime area for this app.
type Runtime struct {
	E          *goengage.Environment
	InChan     chan goengage.Supporter
	DoneChan   chan bool
	ReadOffset int32
	Logger     *goengage.UtilLogger
}

//NewRuntime populates a new runtime.
func NewRuntime(env *goengage.Environment, offset int32, verbose bool) Runtime {
	r := Runtime{
		E:          env,
		InChan:     make(chan goengage.Supporter),
		DoneChan:   make(chan bool),
		ReadOffset: offset,
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

//Visit implements SupporterGuide.Visit and does something with
//a supporter record. In this case, the kludged addressLine1
//and city are cleared away.
func (r *Runtime) Visit(s goengage.Supporter) error {
	if s.Address != nil {
		a := s.Address
		if a.AddressLine1 == a.City && a.City == a.PostalCode && (len(a.PostalCode) > 0) {
			a.AddressLine1 = ""
			a.City = ""
			sdk := goengage.NewSupporterKludgeFix(s)
			updated, err := goengage.SupporterKludgeFixUpsert(r.E, &sdk, r.Logger)
			if err != nil {
				return err
			}
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
		app     = kingpin.New("custom_field-distribution", "Find and fix supporter records with malformed addressLine1 and City")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		offset  = app.Flag("offset", "Read offset. Useful when network goes away").Default("0").Int32()
		verbose = app.Flag("verbose", "See contents of all network actions.  *Really* noisy").Default("false").Bool()
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

	r := NewRuntime(e, *offset, *verbose)
	var wg sync.WaitGroup

	//Start supporter listeners.
	for i := 1; i <= ListenerCount; i++ {
		wg.Add(1)
		go (func(e *goengage.Environment, r *Runtime, wg *sync.WaitGroup) {
			defer wg.Done()
			reporter.ProcessSupporters(r.E, r)
		})(e, &r, &wg)
	}
	//Start done listener.
	wg.Add(1)
	go (func(r *Runtime, n int, wg *sync.WaitGroup) {
		defer wg.Done()
		goengage.DoneListener(r.DoneChan, n)
	})(&r, ListenerCount, &wg)

	//Start supporter reader.
	wg.Add(1)
	go (func(e *goengage.Environment, r *Runtime, wg *sync.WaitGroup) {
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
