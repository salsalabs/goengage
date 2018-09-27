package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/salsalabs/goengage"
)

//No Labels
const token = `sLw8A6soxe-TiccQVB22QiFHSvw1HlYiEQ8aUdfAYaUYLWUf0okqAaXonEexfP_VSxzmEfg6ifh9jIPUKoIjVBTw2BoPpmfvX5yArYmRaXY6mV6gjpQUWu-y5glDm_esgnLODGZZPEmMxKRSS8tqTA`

//Merged is a supporter and the acivity used to subscribe the supporter.
type Merged struct {
	Activity  goengage.SupActivity
	Supporter goengage.Supporter
}

//Lookup accepts a SupActivity and finds the supporter record.  Output goes to the
//merged queue for downstream process.
func Lookup(in chan goengage.SupActivity, out chan Merged) {
	fmt.Printf("Lookup: start")
	for {
		sa, ok := <-in
		if !ok {
			fmt.Printf("Lookup done!")
			close(in)
			return
		}
		//fmt.Printf("Lookkup: received %+v\n", sa)
		rqt := goengage.SupSearchRequest{
			Identifiers:    []string{sa.SupporterID},
			IdentifierType: "SUPPORTER_ID",
			Offset:         0,
			Count:          1,
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
		if resp.Payload.Count != 0 {
			m := Merged{
				Activity:  sa,
				Supporter: resp.Payload.Supporters[0],
			}
			out <- m
		}
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

//View accepts a merge record and displays it.  Or writes it a disk.  Or something.
func View(in chan Merged) {
	fmt.Printf("Merge: start")
	for {
		m, ok := <-in
		if !ok {
			fmt.Printf("View done!")
			close(in)
			return
		}
		//fmt.Printf("View: Received %+v", m)
		e := FirstEmail(m.Supporter)
		email := "(None)"
		if e != nil {
			email = *e
		}
		fmt.Printf("%v,%v,%v,%v,%v,%v\n", m.Supporter.SupporterID,
			email,
			m.Supporter.CreatedDate,
			m.Activity.ActivityFormID,
			m.Activity.ActivityFormName,
			m.Activity.ActivityDate)
	}
}

//Drive finds all subscribe activities and pushes them onto a queue.
func Drive(out chan goengage.SupActivity) {
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
		fmt.Printf("Read %d activities from offset %d\n", count, rqt.Offset)
		rqt.Offset = rqt.Offset + count
		for _, a := range resp.Payload.SupActivities {
			//fmt.Printf("Drive: pushing %s %s %s\n", a.SupporterID, a.ActivityID, a.ActivityFormName)
			out <- a
		}
	}
}

//main starts all of the go routines and then waits for them to finish.
func main() {
	var wg sync.WaitGroup
	c1 := make(chan goengage.SupActivity)
	c2 := make(chan Merged)

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

	fmt.Println("Main: napping and then waiting.")
	time.Sleep(20)
	wg.Wait()
	fmt.Println("Main: done")
}
