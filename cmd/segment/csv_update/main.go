package main

// Program to do a census of all groups and supporters.  Output is a CSV file with
// segment name and supporter email.

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//input is a record that contains at least a group name and an email.
type input struct {
	GroupName   string
	Email       string
	GroupID     string
	SupporterID string
}

//segmentMap is a map of segment names to segment IDs.
type segmentMap map[string]string

//findHead accepts a string slice and a value.  if the value is in the slice, then findHead
//returns the offset.  If the value is not in the slice, then findHead returns -1.
func findHead(s []string, v string) int {
	for i, t := range s {
		if t == v {
			return i
		}
	}
	return -1
}

//mapHeads accepts a list of headers and the first line of a CSV file.  It returns a map
//of headers to field offsets.  Field offsets are -1 if the header doesn't appear in the
//first CSV line.  Returns an error if any of the required fields (GroupName or Email) is
//missing.
func mapHeads(heads []string, first []string, fn string) (map[string]int, error) {
	e := 0
	m := make(map[string]int)
	for _, x := range heads {
		j := findHead(first, x)
		if j == -1 {
			switch x {
			case "GroupName":
				fmt.Printf("Error: %v missing in %v\n", x, fn)
				e++
			case "Email":
				fmt.Printf("Error: %v missing in %v\n", x, fn)
				e++
			default:
				fmt.Printf("Info: %v missing in %v\n", x, fn)
			}
		}
		m[x] = j
	}
	if e > 0 {
		x := fmt.Sprintf("Error: missing headers in %v\n", fn)
		return nil, errors.New(x)
	}
	return m, nil
}

//mapSegments retrieves all segments from Engage and creates a map of segment
//names to segment IDs.
func mapSegments(e *goengage.Environment) (segmentMap, error) {
	c := make(chan goengage.Segment)
	m := make(segmentMap)
	var wg sync.WaitGroup

	//listener
	go (func(c chan goengage.Segment, m segmentMap, wg *sync.WaitGroup) {
		wg.Add(1)
		for r := range c {
			fmt.Printf("mapSegments:listener %v\n", r)
			m[r.Name] = r.SegmentID
		}
		wg.Done()
	})(c, m, &wg)
	//talker
	go (func(e *goengage.Environment, c chan goengage.Segment, wg *sync.WaitGroup) {
		wg.Add(1)
		err := goengage.AllSegments(e, goengage.CountNo, c)
		if err != nil {
			panic(err)
		}
		wg.Done()
	})(e, c, &wg)
	wg.Wait()
	return m, nil
}

//process accepts input records, processes them, then sends them downstream.
func process(c chan input) error {
	for true {
		r, ok := <-c
		if !ok {
			break
		}
		fmt.Printf("process: %+v\n", r.GroupName)
	}
	return nil
}

//read retrieves the contents of the group CSV file, validates the contents, sorts
//the contents, then passes them down to the aggregator.
func read(c chan input, fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	r := csv.NewReader(f)
	a, err := r.ReadAll()
	if err != nil {
		return err
	}
	heads := []string{
		"GroupName",
		"Email",
		"GroupID",
		"SupporterID",
	}
	h := a[0]
	m, err := mapHeads(heads, h, fn)
	if err != nil {
		return err
	}
	for _, t := range a[1:] {
		rec := input{}
		for _, x := range heads {
			j := m[x]
			if j != -1 {
				switch x {
				case "GroupName":
					rec.GroupName = t[j]
				case "GroupID":
					rec.GroupID = t[j]
				case "Email":
					rec.Email = t[j]
				case "SupporterID":
					rec.SupporterID = t[j]
				}
			}
		}
		c <- rec
	}
	close(c)
	return nil
}

//store accepts Census objects from a channel and writes them to a CSV writer.
//The writer is managed internally.  You provide a filename.
func store(c chan goengage.Census, fn string, ids bool) error {
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
		app     = kingpin.New("csv-segment-update", "App to accept group-supporter file and update an org")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("csv", "group-supporter CSV file").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	c := make(chan input)
	var wg sync.WaitGroup

	segMap, err := mapSegments(e)
	fmt.Printf("main: segMap\n%v\n", segMap)
	go (func(c chan input, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := process(c)
		if err != nil {
			panic(err)
		}
	})(c, &wg)

	go (func(c chan input, fn string, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err = read(c, *csvFile)
		if err != nil {
			panic(err)
		}
	})(c, *csvFile, &wg)

	fmt.Println("Waiting...")
	wg.Wait()
	fmt.Println("Done")
	/*
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
	*/
}
