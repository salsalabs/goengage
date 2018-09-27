package main

//Application scan the activities database from top to bottom.  Shows
//all activities and all supporters.
import (
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
	ActivityType     string
	ActivityFormID   string
	ActivityFormName string
	DeltaT           string
}

//OutHeads are the headers for Out.  Sets the order so that the CSV
//output is consistent.
const OutHeads = "SupporterID,CreatedDate,Email,ActivityType,ActivityFormID,ActivityFormName"

//Line accepts a merge record and returns an Output Record.
//Note that the fields are in the same order as OutHeads.
func Line(m Merged) []string {
	email := "None"
	e := goengage.FirstEmail(m.Supporter)
	if e != nil {
		email = *e
	}
	d := strings.Split(m.Activity.ActivityDate, "T")[0]
	var a []string
	// a = append(a, m.Supporter.SupporterID)
	a = append(a, email)
	a = append(a, d)
	a = append(a, m.Activity.ActivityType)
	//a = append(a, m.Activity.ActivityFormID)
	a = append(a, m.Activity.ActivityFormName)

	//Let's see how much time elapsed between the activity
	//and when the supporter was created.
	t1 := goengage.Date(m.Activity.ActivityDate)
	t2 := goengage.Date(m.Supporter.CreatedDate)
	diff := t1.Sub(t2)
	a = append(a, diff.String())

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

//View accepts a slice of merge records and writes them to the console.
func View(e goengage.EngEnv, in chan []Merged) {
	log.Println("Merge: start")
	for {
		m, ok := <-in
		if !ok {
			log.Println("View done!")
			return
		}
		for _, r := range m {
			x := Line(r)
			fmt.Printf("%-40v %v\n", x[0], strings.Join(x[1:], "\t"))
		}
	}
}

//Drive finds all subscribe activities and pushes them onto a queue.
func Drive(e goengage.EngEnv, out chan []goengage.SupActivity) {
	// Use the metrics to determine buffer size.
	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}

	log.Printf("Drive: max size is %d, we're using %d\n", m.MaxBatchSize, 20)

	// Search for all subscribe activities.  Retiurns a supporter ID
	// and activity information.
	rqt := goengage.ActSearchRequest{
		Offset:       0,
		Count:        20,
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
