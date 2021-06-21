package main

// An application to read all email blasts and write CSVs for
// the blasts and for components.  Components are said to be
// valid only for comm series.  We'll test that.  If no CSV
// appears for components, then we'll know that assertion to
// be the truth.

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
	report "github.com/salsalabs/goengage/pkg/report"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	//SettleDuration is the app's settle time in seconds before it
	//starts waiting for things to terminate.
	SettleDuration = "5"
)

//Runtime contains the configuration parts that this app needs.
type Runtime struct {
	Env              *goengage.Environment
	BlastChan        chan goengage.EmailActivity
	BlastCSVChan     chan goengage.EmailActivity
	ComponentChan    chan goengage.EmailActivity
	DoneChan         chan bool
	BlastOffset      int32
	BlastCursor      *string
	BlastCSVFile     string
	ComponentCSVFile string
}

//Visit does something with the blast. Errors terminate.
//Implements goengage.EmailBlastGuide.
func (rt *Runtime) Visit(s goengage.EmailActivity) error {
	rt.BlastCSVChan <- s
	rt.ComponentChan <- s
	return nil
}

//Finalize is called after all blasts have been processed.
//Implements goengage.EmailBlastGuide.
func (rt *Runtime) Finalize() error {
	close(rt.BlastCSVChan)
	close(rt.ComponentChan)
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

//WriteBlasts accepts a blast from the channel and writes it to a CSV
//file.
func (rt *Runtime) WriteBlasts() error {
	log.Printf("WriteBlasts: begin")
	f, err := os.Create(rt.BlastCSVFile)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	headers := []string{
		"ID",
		"Topic",
		"Name",
		"Description",
		"PublishDate",
	}
	err = writer.Write(headers)
	if err != nil {
		return err
	}
	for {
		r, ok := <-rt.BlastCSVChan
		if !ok {
			break
		}
		row := []string{
			r.ID,
			r.Topic,
			r.Name,
			r.Description,
			r.PublishDate,
		}
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}
	log.Printf("WriteBlasts: end")
	return nil
}

//WriteComponents accepts a blast from the channel and writes
//any Components to a CSVfile.
func (rt *Runtime) WriteComponents() error {
	log.Printf("WriteComponents: begin")
	f, err := os.Create(rt.BlastCSVFile)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	headers := []string{
		"EmailActivityID",
		"ContentID",
		"Message",
	}
	err = writer.Write(headers)
	if err != nil {
		return err
	}
	for {
		r, ok := <-rt.ComponentChan
		if !ok {
			break
		}
		if r.Components != nil && len(*r.Components) > 0 {
			for _, c := range *r.Components {
				row := []string{
					r.ID,
					c.ContentID,
					c.MessageNumber,
				}
				err = writer.Write(row)
				if err != nil {
					return err
				}
			}
		}
	}
	log.Printf("WriteComponents: end")
	return nil
}

//Program entry point.
func main() {
	var (
		app              = kingpin.New("blasts_and_components", "Write all email activity (blast) and component info to CSV files")
		login            = app.Flag("login", "YAML file with API token").Required().String()
		blastCSVFile     = app.Flag("blast-csv", "CSV filename to store blast info").Default("email_activity.csv").Required().String()
		componentCSVFile = app.Flag("component-csv", "CSV filename to store component info").Default("email_component.csv").Required().String()
		offset           = app.Flag("blast-offset", "Start here if you lose network connectivity").Default("0").Int32()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		log.Fatalf("Error --login is required.")
		os.Exit(1)
	}
	if blastCSVFile == nil || len(*blastCSVFile) == 0 {
		log.Fatalf("Error --blast-csv is required.")
		os.Exit(1)
	}
	if componentCSVFile == nil || len(*componentCSVFile) == 0 {
		log.Fatalf("Error --csv is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("Error %v\n", err)
		os.Exit(1)
	}

	rtx := Runtime{
		Env:              e,
		BlastChan:        make(chan goengage.EmailActivity, 100),
		BlastCSVChan:     make(chan goengage.EmailActivity, 100),
		ComponentChan:    make(chan goengage.EmailActivity, 100),
		DoneChan:         make(chan bool),
		BlastOffset:      *offset,
		BlastCursor:      nil,
		BlastCSVFile:     *blastCSVFile,
		ComponentCSVFile: *componentCSVFile,
	}
	rt := &rtx
	var wg sync.WaitGroup

	//Start the blast writer.
	wg.Add(1)
	go (func(rt *Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := rt.WriteBlasts()
		if err != nil {
			panic(err)
		}
	})(rt, &wg)
	log.Printf("main: started blast writer")

	//Start the component writer.
	wg.Add(1)
	go (func(rt *Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := rt.WriteComponents()
		if err != nil {
			panic(err)
		}
	})(rt, &wg)
	log.Printf("main: started component writer")

	//Start the blast processor. It calls the functions found in the
	//BlastGuide interface.
	wg.Add(1)
	go (func(rt *Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := report.ProcessEmailBlasts(rt.Env, rt)
		if err != nil {
			panic(err)
		}
	})(rt, &wg)
	log.Printf("main: started blast processor")

	//Start the blast reader.
	wg.Add(1)
	go (func(rt *Runtime, wg *sync.WaitGroup) {
		defer wg.Done()
		err := report.ReadEmailBlasts(rt.Env, rt)
		if err != nil {
			panic(err)
		}
	})(rt, &wg)
	log.Printf("main: started blast reader")

	//Settle time.
	d, _ := time.ParseDuration(SettleDuration)
	log.Printf("main: waiting %d seconds to let things settle", d)
	time.Sleep(d)
}
