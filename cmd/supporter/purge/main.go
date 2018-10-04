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

	"github.com/salsalabs/goengage"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//deleteAfter shows the last valid date for supporters.
const deleteAfter = "2016-12-31T23:59:59.999Z"

//stack gets the maximum number of records, then pushes them onto
//the channel in MaxBatchSize offsets.
func stack(e *goengage.EngEnv, d chan int32) {
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
	fmt.Printf("stack: pushing increments of %d up to %d\n", m.MaxBatchSize, resp.Payload.Total)
	for i := int32(0); i <= resp.Payload.Total; i += m.MaxBatchSize {
		d <- i
	}
	fmt.Printf("stack: done after %d\n", resp.Payload.Total)
	close(d)
}

//pack reads offset from a channel and pushes batches of supporters onto
//the other channel channel.  Errors are noisy and fatal.
func pack(e *goengage.EngEnv, d chan int32, c chan []goengage.Supporter) {
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
			fmt.Println("pack done!")
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
	close(c)
}

//cack accepts batches of supporters from the channel and deletes them.
func cack(e *goengage.EngEnv, c chan []goengage.Supporter, b chan int32) {
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
			fmt.Println("cack done!")
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

//main is the standard entry point for Go applications.
func main() {
	var (
		app   = kingpin.New("delete-supporters", "A command-line app to DELETE ENGAGE SUPPORTERS.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		yes   = app.Flag("yes", "Yes, I understand that this program will utterly and completely remove Engage supporters.").Required().Bool()
	)
	app.Parse(os.Args[1:])
	if !*yes {
		fmt.Printf("You made a good choice.  Supporters won't be deleted.\n")
		return
	} else {
		fmt.Println("***")
		fmt.Printf("*** Alrighty, then.  You supplied --yes, supporters will be deleted.\n")
		fmt.Println("***")
	}

	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	c := make(chan []goengage.Supporter, 100)
	d := make(chan int32, 200)
	b := make(chan int32, 1000)
	var wg sync.WaitGroup

	go (func(b chan int32, wg *sync.WaitGroup) {
		wg.Add(1)
		yack(b)
		wg.Done()
	})(b, &wg)
	for i := 0; i < 5; i++ {
		go (func(c chan []goengage.Supporter, d chan int32, e *goengage.EngEnv, wg *sync.WaitGroup) {
			wg.Add(1)
			pack(e, d, c)
			wg.Done()
		})(c, d, e, &wg)
	}
	for i := 0; i < 5; i++ {
		go (func(c chan []goengage.Supporter, b chan int32, e *goengage.EngEnv, wg *sync.WaitGroup) {
			wg.Add(1)
			cack(e, c, b)
			wg.Done()
		})(c, b, e, &wg)
	}
	go (func(d chan int32, e *goengage.EngEnv, wg *sync.WaitGroup) {
		wg.Add(1)
		stack(e, d)
		wg.Done()
	})(d, e, &wg)

	time.Sleep(3000)
	fmt.Println("Main: waiting...")
	wg.Wait()
}
