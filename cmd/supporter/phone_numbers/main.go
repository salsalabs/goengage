//App to extract phone numbers for a list of supporters. The extracted
//data is stored in a CSV.  The CSV has one row per supporter. Each row
//contains SupporterID, Home Phone, Cell Phone, and Work Phone.
//
//Unlike Classic, phone numbers are not fixed fields. They are elements
//in the "Contacts" part of the supporter record.  Note that Engage does
//not send output for fields that do not have values. You'll see a lot of
//"use this field if it exists or substitute an empty string" in this app.
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	reportSupporter "github.com/salsalabs/goengage/pkg/report"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Runtime area for this app.
type Runtime struct {
	E         *goengage.Environment
	ChunkChan chan []string
	WriteChan chan goengage.Supporter
	DoneChan  chan bool
	IDFile    string
	CSVOut    *csv.Writer
}

//RequestedIDs returns the list of supporterIDs from the ID file.
//Each line of the file is a single Supporter ID.
func (rt *Runtime) RequestedIds() (a []string, err error) {
	r, err := os.Open(rt.IDFile)
	if err != nil {
		return a, err
	}
	defer r.Close()
	fs := bufio.NewScanner(r)
	fs.Split(bufio.ScanLines)
	for fs.Scan() {
		id := fs.Text()
		id = strings.Trim(id, "'\" \t")
		if len(id) != 36 {
			log.Fatalf("RequestedIds: file %v, '%v' is not a valid id\n", rt.IDFile, id)
		}
		a = append(a, id)
	}
	return a, err
}

//NewRuntime populates a new runtime.
func NewRuntime(env *goengage.Environment, idFile string, out *csvWriter) Runtime {
	r := Runtime{
		E:         env,
		ChunkChan: make(chan []string, 100),
		WriteChan: make(chan goengage.Supporter, 100),
		DoneChan:  make(chan bool),
		IDFile:    idFile,
		CSVOut:    out,
	}
	return r
}

//Visit implements SupporterGuide.Visit and does something with
//a supporter record
func (r *Runtime) Visit(s goengage.Supporter) error {
	if s.Contacts == nil {
		return nil
	}
	row := make([]string, 4)
	row[0] = s.SupporterID
	for _, c := range s.Contacts {
		switch c.Type {
		case goengage.ContactHome:
			row[1] = c.Value
		case goengage.ContactCell:
			row[2] = c.Value
		case goEngage.ContactWork:
			row[3] = c.Value
		}
	}
	err := r.CSVOut.write(row)
	if err != nil {
		return err
	}
	log.Println(row)
	return nil
}

//Finalize implements SupporterGuide.Filnalize and does nothing
//in this app.
func (r *Runtime) Finalize() error {
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
		app     = kingpin.New("segments_for_supporters", "Write a CSV of supporters and segments for a list of supporter IDs")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		idFile  = app.Flag("input", "Text with list of Engage supporterIDs to look up").Required().String()
		outFile = app.Flag("output", "CSV filename to store supporter-segment data").Default("supporters_and_segments.csv").String()
		debug   = app.Flag("debug", "Write requests and responses to a log file in JSON").Bool()
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

	f, err := os.Create(csvFile)
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
