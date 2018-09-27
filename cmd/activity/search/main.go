package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/salsalabs/goengage"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Merged is a supporter and the acivity used to subscribe the supporter.
type Merged struct {
	Activity  goengage.SupActivity
	Supporter goengage.Supporter
}

//Out is the record that we're writing to the output.
type Out struct {
	SupporterID      string
	CreatedDate      string
	Email            string
	ActivityFormID   string
	ActivityFormName string
}

//OutHeads are the headers for Out.  Sets the order so that the CSV
//output is consistent.
const OutHeads = "SupporterID,CreatedDate,Email,ActivityFormID,ActivityFormName"

//Line accepts a merge record and returns an Output Record.
//Note that the fields are in the same order as OutHeads.
func Line(m Merged) []string {
	email := "None"
	e := goengage.FirstEmail(m.Supporter)
	if e != nil {
		email = *e
	}
	var a []string
	a = append(a, m.Supporter.SupporterID)
	a = append(a, email)
	a = append(a, m.Supporter.CreatedDate)
	a = append(a, m.Activity.ActivityFormID)
	a = append(a, m.Activity.ActivityFormName)

	return a
}

//Lookup accepts a slice of SupActivity from a channel, reads the associated supporter
//records then puts a slice of merged record onto the output channel.
func Lookup(e goengage.EngEnv, in chan []goengage.SupActivity, out chan []Merged) {
	log.Println("Lookup: start")
	for {
		sa, ok := <-in
		if !ok {
			log.Println("Lookup done!")
			close(out)
			return
		}
		//Make a map of supporter ID and supActivities.
		m := make(map[string]goengage.SupActivity)
		var a []string
		for _, r := range sa {
			m[r.SupporterID] = r
			a = append(a, r.SupporterID)
		}

		//log.Printf("Lookkup: received %+v\n", sa)
		rqt := goengage.SupSearchRequest{
			Identifiers:    a,
			IdentifierType: "SUPPORTER_ID",
			Offset:         0,
			Count:          int32(len(a)),
		}
		var resp goengage.SupSearchResult
		n := goengage.NetOp{
			Host:     e.Host,
			Fragment: goengage.SupSearch,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err := n.Search()
		if err != nil {
			panic(err)
		}
		var x []Merged
		for _, s := range resp.Payload.Supporters {
			email := goengage.FirstEmail(s)
			if s.Result == "FOUND" {
				mg := Merged{
					Activity:  m[s.SupporterID],
					Supporter: s,
				}
				x = append(x, mg)
			} else {
				log.Printf("Lookup: %v status %v\n", email, s.Result)
			}
		}
		if len(x) > 0 {
			out <- x
		}
	}
}

//View accepts a slice of merge records and writes them to disk
//in CSV format.
func View(e goengage.EngEnv, in chan []Merged) {
	log.Println("Merge: start")
	f, err := os.Create("supporter_page.csv")
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	h := strings.Split(OutHeads, ",")
	w.Write(h)
	for {
		m, ok := <-in
		if !ok {
			log.Println("View done!")
			return
		}
		var a [][]string
		for _, r := range m {
			a = append(a, Line(r))
		}
		w.WriteAll(a)
		w.Flush()
	}
}

//Drive finds all subscribe activities and pushes them onto a queue.
func Drive(e goengage.EngEnv, out chan []goengage.SupActivity) {
	// Use the metrics to determine buffer size.
	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Drive: max size is %d, we're using %d\n", m.MaxBatchSize, 20)

	// Search for all subscribe activities.  Retiurns a supporter ID
	// and activity information.
	rqt := goengage.ActSearchRequest{
		Offset:       0,
		Count:        20,
		Type:         "SUBSCRIBE",
		ModifiedFrom: "2010-01-01T00:00:00.000Z",
	}
	var resp goengage.ActSearchResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.ActSearch,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err = n.Search()
	if err != nil {
		panic(err)
	}

	// Do for all items in the results.  Send the SupActivity
	count := int32(rqt.Count)
	for count > 0 {
		err := n.Search()
		if err != nil {
			panic(err)
		}
		count = int32(len(resp.Payload.SupActivities))
		if count > 0 {
			log.Printf("Drive: read %d activities from offset %d\n", count, rqt.Offset)
			rqt.Offset = rqt.Offset + count
			out <- resp.Payload.SupActivities
		}
	}
	close(out)
}

//main starts all of the go routines and then waits for them to finish.
func main() {
	var (
		app   = kingpin.New("activity-search", "A command-line app to search for supporters added by activities.")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	c1 := make(chan []goengage.SupActivity)
	c2 := make(chan []Merged)

	go (func(wg *sync.WaitGroup) {
		wg.Add(1)
		Lookup(*e, c1, c2)
		wg.Done()
	})(&wg)

	go (func(wg *sync.WaitGroup) {
		wg.Add(1)
		View(*e, c2)
		wg.Done()
	})(&wg)

	go (func(wg *sync.WaitGroup) {
		wg.Add(1)
		Drive(*e, c1)
		wg.Done()
	})(&wg)

	log.Println("Main: napping and then waiting.")
	time.Sleep(20)
	wg.Wait()
	log.Println("Main: done")
}
