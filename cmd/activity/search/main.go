package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/salsalabs/goengage"
)

//No Labels
const token = ``

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
const OutHeads = "SupporterID,CreatedDate,,Email,ActivityFormID,ActivityFormName"

//Line accepts a merge record and returns an Output Record.
//Note that the fields are in the same order as OutHeads.
func Line(m Merged) []string {
	email := "None"
	e := FirstEmail(m.Supporter)
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
func Lookup(in chan []goengage.SupActivity, out chan []Merged) {
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
			Host:     goengage.ProdHost,
			Fragment: goengage.SupSearch,
			Token:    token,
			Request:  &rqt,
			Response: &resp,
		}
		err := n.Search()
		if err != nil {
			panic(err)
		}
		var x []Merged
		for _, s := range resp.Payload.Supporters {
			if s.Result == "FOUND" {
				mg := Merged{
					Activity:  m[s.SupporterID],
					Supporter: s,
				}
				x = append(x, mg)
			} else {
				log.Printf("Lookup: %v status %v\n", FirstEmail(s), s.Result)
			}
		}
		out <- x
	}
}

//FirstEmail returns the first email address for the provided supporter.
//Returns nil if the supporter does not have an email.  (As if...)
func FirstEmail(s goengage.Supporter) *string {
	c := s.Contacts
	if c == nil || len(c) == 0 {
		return nil
	}
	for _, x := range c {
		if x.Type == "EMAIL" {
			email := x.Value
			return &email
		}
	}
	return nil
}

//View accepts a slice of merge records and writes them to disk
//in CSV format.
func View(in chan []Merged) {
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
func Drive(out chan []goengage.SupActivity) {
	// Search for all subscribe activities.  Retiurns a supporter ID
	// and activity information.
	rqt := goengage.ActSearchRequest{
		Offset:       0,
		Count:        50,
		Type:         "SUBSCRIBE",
		ModifiedFrom: "2010-01-01T00:00:00.000Z",
	}
	var resp goengage.ActSearchResult
	n := goengage.NetOp{
		Host:     goengage.ProdHost,
		Fragment: goengage.ActSearch,
		Token:    token,
		Request:  &rqt,
		Response: &resp,
	}
	err := n.Search()
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
	var wg sync.WaitGroup
	c1 := make(chan []goengage.SupActivity)
	c2 := make(chan []Merged)

	go (func(wg *sync.WaitGroup) {
		wg.Add(1)
		Lookup(c1, c2)
		wg.Done()
	})(&wg)

	go (func(wg *sync.WaitGroup) {
		wg.Add(1)
		View(c2)
		wg.Done()
	})(&wg)

	go (func(wg *sync.WaitGroup) {
		wg.Add(1)
		Drive(c1)
		wg.Done()
	})(&wg)

	log.Println("Main: napping and then waiting.")
	time.Sleep(20)
	wg.Wait()
	log.Println("Main: done")
}
