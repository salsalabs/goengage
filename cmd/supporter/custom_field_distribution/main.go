//App to search for supporters that have a custom field then
//to report on the distribution of custom field values.  The
//app also shows supporters who do not have the custom field.
//Unlike Classic, a supporter without an assigned custom field
//value is not equipped with that custom field.
package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	reportSupporter "github.com/salsalabs/goengage/pkg/report"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Runtime area for this app.
type Runtime struct {
	E          *goengage.Environment
	InChan     chan goengage.Supporter
	DoneChan   chan bool
	Cache      Cache
	Keys       []string
	FieldName  string
	ReadOffset int32
}

//Cache is used to store values and counts.
type Cache map[string]int32

//NewRuntime populates a new runtime.
func NewRuntime(env *goengage.Environment, f string) Runtime {
	c := Cache{
		Null:        0,
		NotEquipped: 0,
	}
	r := Runtime{
		E:          env,
		InChan:     make(chan goengage.Supporter),
		DoneChan:   make(chan bool),
		Cache:      c,
		Keys:       []string{Null, NotEquipped},
		FieldName:  f,
		ReadOffset: int32(0),
	}
	return r
}

//Constants for cache updates.
const (
	Null        = "Null"
	NotEquipped = "NotEquipped"
)

//Visit implements SupporterGuide.Visit and does something with
//a supporter record
func (r *Runtime) Visit(s goengage.Supporter) error {
	for _, f := range s.CustomFieldValues {
		if f.Name == r.FieldName {
			if len(f.Value) == 0 {
				r.Cache[Null] = r.Cache[Null] + 1
				// log.Printf("Visit: %s: %d\n", Null, r.Cache[Null])
				return nil
			}
			p, ok := r.Cache[f.Value]
			if ok {
				p = r.Cache[f.Value]
			} else {
				r.Keys = append(r.Keys, f.Value)
				p = 0
			}
			r.Cache[f.Value] = p + 1
			// log.Printf("Visit: %s: %d\n", f.Value, r.Cache[f.Value])
			return nil
		}
	}
	r.Cache[NotEquipped] = r.Cache[NotEquipped] + 1
	// log.Printf("Visit: %s: %d\n", NotEquipped, r.Cache[NotEquipped])
	return nil
}

//Finalize implements SupporterGuide.Finalize and outputs the
//distribution results.
func (r *Runtime) Finalize() error {
	fmt.Println("Value,Count")
	for _, k := range r.Keys {
		fmt.Printf("%s,%d\n", k, r.Cache[k])
	}
	log.Printf("%v\n", r.Cache)
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

//Adjust offset changes the proposed offset as needed.
//Useful for chunked ID reads.  Does nothing in this app.
func (r *Runtime) AdjustOffset(offset int32) int32 {
	return offset
}

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app       = kingpin.New("custom_field-distribution", "Search for a custom field and report on value distribution")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		fieldName = app.Flag("fieldName", "Custom field name to count").Required().String()
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

	r := NewRuntime(e, *fieldName)
	var wg sync.WaitGroup

	//Start supporter listener. Only one of these because Visit is quick
	//in this app. More than one cases "concurrent map writes" errors.
	go (func(e *goengage.Environment, r *Runtime, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		reportSupporter.ProcessSupporters(r.E, r)
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
		reportSupporter.ReadSupporters(r.E, r)
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
