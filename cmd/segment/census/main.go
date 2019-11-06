package main

// Program to do a census of all groups and supporters.  Output is a CSV file with
// segment name and supporter email.

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//store accepts Census objects from a channel and writes them to a CSV writer.
//The writer is managed internally.  You provide a filename.
func store(c chan goengage.Census, fn string, ids bool) error {
	//CSV output goes here.
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	r := []string{"GroupName", "Email"}
	if ids {
		r = []string{"GroupID", "GroupName", "SupporterID", "Email"}
	}
	err = w.Write(r)
	if err != nil {
		return err
	}
	for true {
		s, ok := <-c
		if !ok {
			break
		}
		sName := s.Segment.Name
		for _, u := range s.Supporters {
			email := goengage.FirstEmail(u)
			if email != nil {
				r := []string{sName, *email}
				if ids {
					r = []string{sName, *email, s.SegmentID, u.SupporterID}
				}
				err := w.Write(r)
				if err != nil {
					return err
				}
			}
		}
		fmt.Printf("%-32s %5d\n", sName, len(s.Supporters))
		w.Flush()
	}
	w.Flush()
	f.Close()
	return nil
}

func main() {
	var (
		app   = kingpin.New("see-segment-census", "A command-line app to display segment names and supporters.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		ids   = app.Flag("ids", "Segment and supporterIDs in output").Bool()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	//store receives Census records and writes them to a CSV file.
	c := make(chan goengage.Census)
	fn := "group_census.csv"
	var wg sync.WaitGroup

	go (func(c chan goengage.Census, fn string, wg *sync.WaitGroup, ids bool) {
		wg.Add(1)
		defer wg.Done()
		err := store(c, fn, ids)
		if err != nil {
			panic(err)
		}
	})(c, fn, &wg, *ids)
	fmt.Println("Started store...")

	go (func(e *goengage.Environment, c chan goengage.Census, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		goengage.AllSegmentCensus(e, c)
	})(e, c, &wg)
	fmt.Println("Started AllSegmentCensus...")

	fmt.Println("Waiting...")
	wg.Wait()
	fmt.Println("Done!")
	fmt.Printf("Output can be found in %v\n", fn)
}