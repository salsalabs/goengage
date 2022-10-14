// App to extract phone numbers for a list of supporters. The extracted
// data is stored in a CSV.  The CSV has one row per supporter. Each row
// contains SupporterID, Home Phone, Cell Phone, and Work Phone.
//
// Unlike Classic, phone numbers in Engage are not fixed fields. They are
// elements in the "Contacts" part of the supporter record.
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	reportSupporter "github.com/salsalabs/goengage/pkg/report"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Runtime area for this app.
type Runtime struct {
	E         *goengage.Environment
	WriteChan chan goengage.Supporter
	DoneChan  chan bool
	IDFile    string
	IDs       []string
	IdOffset  int32
	CSVOut    *csv.Writer
}

// RequestedIDs returns the list of supporterIDs from the ID file.
// Each line of the file is a single Supporter ID.
func (rt *Runtime) RequestedIds() error {
	r, err := os.Open(rt.IDFile)
	if err != nil {
		return err
	}
	var a []string
	defer r.Close()
	fs := bufio.NewScanner(r)
	fs.Split(bufio.ScanLines)
	for fs.Scan() {
		id := fs.Text()
		id = strings.Trim(id, "'\" \t")
		if len(id) != 36 {
			err = fmt.Errorf("not a vaid supporterID, file %v, entry '%v'", rt.IDFile, id)
			return err
		}
		a = append(a, id)
	}
	rt.IDs = a
	return nil
}

// NewRuntime populates a new runtime.
func NewRuntime(env *goengage.Environment, idFile string, out *csv.Writer) Runtime {
	r := Runtime{
		E:         env,
		WriteChan: make(chan goengage.Supporter, 100),
		DoneChan:  make(chan bool),
		IDFile:    idFile,
		IdOffset:  0,
		CSVOut:    out,
	}
	return r
}

// Visit implements SupporterGuide.Visit and does something with
// a supporter record
func (r *Runtime) Visit(s goengage.Supporter) error {
	if s.Contacts == nil {
		return nil
	}
	var row []string
	for i := 0; i < 4; i++ {
		row = append(row, "")
	}
	row[0] = s.SupporterID
	for _, c := range s.Contacts {
		switch c.Type {
		case goengage.ContactHome:
			row[1] = c.Value
		case goengage.ContactCell:
			row[2] = c.Value
		case goengage.ContactWork:
			row[3] = c.Value
		}
	}
	err := r.CSVOut.Write(row)
	if err != nil {
		return err
	}
	r.CSVOut.Flush()
	log.Println(row)
	return nil
}

// Finalize implements SupporterGuide.Finalize and does nothing
// in this app.
func (r *Runtime) Finalize() error {
	return nil
}

// Payload implements SupporterGuide.Payload and provides a payload
// that will retrieve all supporters.
func (r *Runtime) Payload() goengage.SupporterSearchRequestPayload {
	low := float64(r.IdOffset)
	remaining := float64(len(r.IDs)) - float64(low)
	high := low + math.Min(remaining, float64(r.E.Metrics.MaxBatchSize))
	max := len(r.IDs)
	current := r.IDs[int32(low):int32(high):max]
	payload := goengage.SupporterSearchRequestPayload{
		IdentifierType: goengage.SupporterIDType,
		Identifiers:    current,
		Offset:         0,
		Count:          0,
	}
	r.IdOffset += r.E.Metrics.MaxBatchSize
	return payload
}

// Channel implements SupporterGuide.Channnel and provides the
// supporter channel.
func (r *Runtime) Channel() chan goengage.Supporter {
	return r.WriteChan
}

// DoneChannel implements SupporterGuide.DoneChannel to provide
// a channel that  receives a true when the listener is done.
func (r *Runtime) DoneChannel() chan bool {
	return r.DoneChan
}

// Offset returns the offset for the first read.
// Useful for restarts.
func (r *Runtime) Offset() int32 {
	return 0
}

//AdjustOffset allows us to change the read offset to
//zero for the payload that we're using.

func (r *Runtime) AdjustOffset(offset int32) int32 {
	if r.IdOffset > int32(len(r.IDs)) {
		// Offset that wll hopefully trigger end of file.
		// Yeah, this needs some work...
		return 1_000_000_000
	} else {
		return 0
	}
}

// Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app     = kingpin.New("phone_numbers", "Write a CSV of supporterIDs and phone numbers")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		idFile  = app.Flag("input", "Text with list of Engage supporterIDs to look up").Required().String()
		outFile = app.Flag("output", "CSV filename to store supporter-segment data").Default("phone_numbers.csv").String()
		//debug   = app.Flag("debug", "Write requests and responses to a log file in JSON").Bool()
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

	f, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	headers := []string{
		"SupporterID",
		"HomePhone",
		"CellPhone",
		"WorkPhone",
	}
	err = writer.Write(headers)
	if err != nil {
		panic(err)
	}

	r := NewRuntime(e, *idFile, writer)
	err = r.RequestedIds()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	//Start supporter listener. Only one of these because Visit is quick
	//in this app. More than one cases "concurrent map writes" errors.
	wg.Add(1)
	go (func(e *goengage.Environment, r *Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		reportSupporter.ProcessSupporters(r.E, r)
	})(e, &r, &wg)

	//Start done listener.
	wg.Add(1)
	go (func(r *Runtime, n int, wg *sync.WaitGroup) {
		defer wg.Done()
		goengage.DoneListener(r.DoneChan, n)
	})(&r, 1, &wg)

	//Start supporter reader.
	wg.Add(1)
	go (func(e *goengage.Environment, r *Runtime, wg *sync.WaitGroup) {
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
