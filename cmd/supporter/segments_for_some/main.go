package main

// An application to accept a list of supporterIDs, find the groups
// that they belong to, and write a CSV file. Each row of the CSV
// file will contain the supporterID, supporter's email and a comma-
// separated list of groups.
//
// Input is provided by a YAML file that contains a "supporterIDs:"
// field.  The field will be a list of supporterIDs of interest.
// For convenience's sake, the list of IDs can be added to the
// login file, then the login file gets supplied twice in the
// calling arguments.
import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	//SettleDuration is the app's settle time in seconds before it
	//starts waiting for things to terminate.
	SettleDuration = "5s"

	//ReaderCount is the number of SupporterID readers.
	ReaderCount = 5
)

//OutRecord is the data that we want written to a CSV file.
type OutRecord struct {
	SupporterID string
	Email       string
	Segments    []goengage.Segment
}

//Runtime contains the configuration parts that this app needs.
type Runtime struct {
	Env         *goengage.Environment
	ReaderCount int
	IDChan      chan string
	OutChan     chan OutRecord
	DoneChan    chan bool
	IDFile      string
	OutFile     string
	Logger      *goengage.UtilLogger
}

//BuildOut accepts a supporter key from a channel and
// writes an OutRecord to the out channel.
func (rt *Runtime) BuildOut(id int) error {
	for {
		supporterId, ok := <-rt.IDChan
		if !ok {
			break
		}
		s, err := goengage.SupporterByID(rt.Env, supporterId)
		if err != nil {
			return err
		}
		if s == nil {
			//log.Printf("BuildOut-%d: %v does not locate a supporter\n", id, supporterId)
		} else {
			email := ""
			e := goengage.FirstEmail(*s)
			if err != nil {
				email = *e
			}
			segments, err := goengage.SupporterSegments(rt.Env, supporterId)
			if err != nil {
				return err
			}
			if len(segments) > 0 {
				r := OutRecord{
					SupporterID: supporterId,
					Email:       email,
					Segments:    segments,
				}
				rt.OutChan <- r
				log.Printf("BuildOut-%d: %v %d segments\n", id, supporterId, len(segments))
			} else {
				log.Printf("BuildOut-%d: %v does not belong to any segments\n", id, s)
			}
		}
	}
	rt.DoneChan <- true
	log.Printf("Buildout-%d: end", id)
	return nil
}

//RequestedIDs returns the list of supporterIDs from a text file. Each line
//is an id.
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
		if len(id) == 36 {
			a = append(a, id)
		} else {
			//log.Printf("RequestedIds: file %v, '%v' is not a valid id\n", rt.IDFile, id)
		}
	}
	sort.Strings(a)
	return a, err
}

//WriteOut accepts an OutRecord from a channel and writes
//it to a CSV file.
func (rt *Runtime) WriteOut() error {
	f, err := os.Create(rt.OutFile)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	headers := []string{
		"SupporterID",
		"Email",
		"Groups",
	}
	err = writer.Write(headers)
	if err != nil {
		return err
	}
	for {
		r, ok := <-rt.OutChan
		if !ok {
			break
		}
		var segments []string

		for _, s := range r.Segments {
			segments = append(segments, s.Name)
		}
		groups := strings.Join(segments, ",")
		row := []string{
			r.SupporterID,
			r.Email,
			groups,
		}
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	log.Printf("WriteOut: end\n")
	return nil
}

//WaitForReaders waits for readers to send to the done channel.
//Closes the out channel when all readers are done.
func (rt *Runtime) WaitForReaders() {
	count := rt.ReaderCount
	for count > 0 {
		log.Printf("WaitForReaders: Waiting for %d readers\n", count)
		_, ok := <-rt.DoneChan
		if !ok {
			break
		}
		count--
	}
	close(rt.OutChan)
	log.Println("WaitForReaders: done")
}

//Program entry point.
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
		log.Fatalf("Error --login is required.")
		os.Exit(1)
	}
	if idFile == nil || len(*idFile) == 0 {
		idFile = login
	}
	if outFile == nil || len(*outFile) == 0 {
		log.Fatalf("Error --output is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("Error %v\n", err)
		os.Exit(1)
	}

	var logger *goengage.UtilLogger
	if *debug {
		logger, err = goengage.NewUtilLogger()
		if err != nil {
			log.Fatalf("Error %v\n", err)
			os.Exit(1)
		}
	}

	rtx := Runtime{
		Env:         e,
		ReaderCount: ReaderCount,
		IDChan:      make(chan string, 100),
		OutChan:     make(chan OutRecord, 100),
		DoneChan:    make(chan bool),
		IDFile:      *idFile,
		OutFile:     *outFile,
		Logger:      logger,
	}
	rt := &rtx

	requestedIds, err := rt.RequestedIds()
	if err != nil {
		log.Fatalf("Error %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go (func(rt *Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		rt.WaitForReaders()
	})(rt, &wg)
	log.Println("main: started reader waiter")

	wg.Add(1)
	go (func(rt *Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := rt.WriteOut()
		if err != nil {
			log.Fatalf("WriteOut error %v\n", err)
			os.Exit(1)
		}
	})(rt, &wg)
	log.Println("main: started output writer")

	for i := 1; i <= rt.ReaderCount; i++ {
		wg.Add(1)
		go (func(rt *Runtime, wg *sync.WaitGroup, i int) {
			defer wg.Done()
			err := rt.BuildOut(i)
			if err != nil {
				log.Fatalf("BuildOut-%d: error %v\n", i, err)
				os.Exit(1)
			}
		})(rt, &wg, i)
	}
	log.Println("main: started output builders")

	// Load the input queue.
	for _, id := range requestedIds {
		rt.IDChan <- id
	}
	close(rt.IDChan)
	log.Printf("main: main queue loaded with %d supporterIds\n", len(requestedIds))

	//Settle time.
	d, _ := time.ParseDuration(SettleDuration)
	log.Printf("main: waiting %v seconds to let things settle\n", d.Seconds())
	time.Sleep(d)
	log.Println("main: running...")
	wg.Wait()
	log.Println("main: done")
}
