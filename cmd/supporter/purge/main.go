package main

//Search for all supporters and cack 'em. Supporters will be deleted
//if they were added after 12-Dec-2016.
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//deleteAfter shows the last valid date for supporters.
const deleteAfter = "2016-12-31T23:59:59.999Z"

//stack gets the maximum number of records, then pushes them onto
//the channel in MaxBatchSize offsets.  Note that the user can override
//the maximum number...
func stack(e *goengage.EngEnv, d chan int32, max *int32) {
	fmt.Printf("stack started\n")
	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}

	rqt := goengage.SupSearchRequest{
		ModifiedFrom: deleteAfter,
		Offset:       0,
		Count:        1,
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SupSearch,
		Method:   http.MethodPost,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err = n.Do()
	if err != nil {
		panic(err)
	}
	most := resp.Payload.Total
	if max != nil {
		most = *max
	}
	fmt.Printf("stack: pushing increments of %d up to %d\n", m.MaxBatchSize, most)
	for i := int32(0); i <= most; i += m.MaxBatchSize {
		d <- i
	}
	fmt.Printf("stack: done after %d\n", most)
	close(d)
}

//pack reads offset from a channel and pushes batches of supporters onto
//the other channel channel.  Errors are noisy and fatal.
func pack(e *goengage.EngEnv, d chan int32, c chan []goengage.Supporter, done chan bool) {
	fmt.Printf("pack started\n")
	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}

	rqt := goengage.SupSearchRequest{
		ModifiedFrom: deleteAfter,
		Offset:       0,
		Count:        m.MaxBatchSize,
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SupSearch,
		Method:   http.MethodPost,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}

	for {
		offset, ok := <-d
		if !ok {
			break
		}
		rqt.Offset = offset
		rqt.Count = m.MaxBatchSize
		// fmt.Printf("Reading %d supporters from offset %d\n", rqt.Count, rqt.Offset)
		err := n.Do()
		if err != nil {
			panic(err)
		}
		c <- resp.Payload.Supporters
	}
	done <- true
}

//cack accepts batches of supporters from the channel and deletes them.
func cack(e *goengage.EngEnv, c chan []goengage.Supporter, b chan int32, done chan bool) {
	fmt.Printf("cack started\n")
	dRqt := goengage.SupDeleteRequest{}
	dResp := goengage.SupDeleteResult{}
	nDel := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SupDelete,
		Method:   http.MethodDelete,
		Token:    e.Token,
		Request:  &dRqt,
		Response: &dResp,
	}

	for {
		p, ok := <-c
		if !ok {
			break
		}
		var a []goengage.DeletingSupporters
		for _, x := range p {
			d := goengage.DeletingSupporters{SupporterID: x.SupporterID}
			a = append(a, d)
		}
		dRqt.Supporters = a
		err := nDel.Do()
		if err != nil {
			panic(err)
		}

		//for _, s := range dResp.Payload.Supporters {
		//	fmt.Printf("%s %s\n", s.SupporterID, s.Result)
		//}
		b <- int32(len(dResp.Payload.Supporters))
	}
	done <- true
}

//yack keeps a running total of the supporters processed.
func yack(e chan int32) {
	fmt.Printf("yack started\n")
	t := int32(0)
	for {
		x, ok := <-e
		if !ok {
			return
		}
		t += x
		log.Printf("yack: %d\n", t)
	}
}

//pWait listens to a channel for 'n' messages.  When the n-th message arrives,
//pWait closed the other channel.
func pWait(b chan bool, c chan []goengage.Supporter, n int32) {
	fmt.Println("pWait started")
	for {
		_, ok := <-b
		if !ok {
			fmt.Printf("pWait: channel closed with %d remaining\n", n)
			break
		}
		n--
		fmt.Printf("pWait: waiting for %d\n", n)
		if n <= 0 {
			fmt.Println("pWait: done on count")
			break
		}
	}
	close(c)
}

//cWait listens to a channel for 'n' messages.  When the n-th message arrives,
//cWait closed the other channel.
func cWait(b chan bool, c chan int32, n int32) {
	fmt.Println("cWait started")
	for {
		_, ok := <-b
		if !ok {
			fmt.Printf("cWait: channel closed with %d remaining\n", n)
			break
		}
		n--
		fmt.Printf("pWait: waiting for %d\n", n)
		if n <= 0 {
			fmt.Println("cWait: done on count")
			break
		}
	}
	close(c)
}

//main is the standard entry point for Go applications.
func main() {
	var (
		app   = kingpin.New("delete-supporters", "A command-line app to DELETE ENGAGE SUPPORTERS.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		max   = app.Flag("max", "Delete no more than this many supporter").Int32()
		yes   = app.Flag("yes", "Yes, I understand that this program will utterly and completely remove Engage supporters.").Required().Bool()
	)
	app.Parse(os.Args[1:])
	if !*yes {
		fmt.Printf("You made a good choice.  Supporters won't be deleted.\n")
		return
	}
	fmt.Println("***")
	fmt.Printf("*** Alrighty, then.  You supplied --yes, supporters will be deleted.\n")
	fmt.Println("***")

	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	c := make(chan []goengage.Supporter, 100)
	d := make(chan int32, 200)
	b := make(chan int32, 1000)
	cDone := make(chan bool)
	pDone := make(chan bool)
	var wg sync.WaitGroup

	go (func(b chan int32, wg *sync.WaitGroup) {
		wg.Add(1)
		yack(b)
		wg.Done()
	})(b, &wg)

	go (func(b chan bool, c chan []goengage.Supporter, n int32, wg *sync.WaitGroup) {
		wg.Add(1)
		pWait(b, c, n)
		wg.Done()
	})(pDone, c, 5, &wg)

	for i := 0; i < 5; i++ {
		go (func(c chan []goengage.Supporter, d chan int32, e *goengage.EngEnv, wg *sync.WaitGroup) {
			wg.Add(1)
			pack(e, d, c, pDone)
			wg.Done()
		})(c, d, e, &wg)
	}

	go (func(b chan bool, c chan int32, n int32, wg *sync.WaitGroup) {
		wg.Add(1)
		cWait(b, c, n)
		wg.Done()
	})(cDone, b, 5, &wg)

	for i := 0; i < 5; i++ {
		go (func(c chan []goengage.Supporter, b chan int32, e *goengage.EngEnv, wg *sync.WaitGroup) {
			wg.Add(1)
			cack(e, c, b, cDone)
			wg.Done()
		})(c, b, e, &wg)
	}
	go (func(d chan int32, e *goengage.EngEnv, m *int32, wg *sync.WaitGroup) {
		wg.Add(1)
		stack(e, d, max)
		wg.Done()
	})(d, e, max, &wg)

	time.Sleep(3000)
	fmt.Println("Main: waiting...")
	wg.Wait()
}
