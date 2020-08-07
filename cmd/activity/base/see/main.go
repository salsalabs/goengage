package main

//Application scan the activities database from top to bottom and write them
//to the console.
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//handle retrieves responses from the channel, formats them, and
//writes them to the handle's own CSV file.
func handle(c chan goengage.BaseResponse, writer *csv.Writer, id int) {
	log.Printf("handle-%d: begin\n", id)
	for true {
		resp, ok := <-c
		if !ok {
			break
		}
		var cache [][]string
		for _, a := range resp.Payload.Activities {
			date := strings.Split(fmt.Sprintf("%v", a.ActivityDate), " ")[0]
			record := []string{
				a.SupporterID,
				a.PersonName,
				a.PersonEmail,
				a.ActivityType,
				date,
			}
			cache = append(cache, record)
		}
		err := writer.WriteAll(cache)
		if err != nil {
			panic(err)
		}
		log.Printf("handle-%d: write %d\n", id, len(cache))
		writer.Flush()
	}
	log.Printf("handle-%d: end\n", id)
}

//startHandler creates a handler that reads from a channel of responses
//and writes to the 'n'th output file. Output files have "-n" just before
//the dot that separates the name from the extension (whatever-1.csv,
//whatever-2.csv, etc.)  Errors panic.
func startHandler(c chan goengage.BaseResponse, filename string, n int) {
	parts := strings.Split(filename, ".")
	csvFile := fmt.Sprintf("%s-%d.%s", parts[0], n, parts[1])
	f, err := os.Create(csvFile)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(f)
	headers := []string{
		"SupporterID",
		"PersonName",
		"PersonEmail",
		"ActivityType",
		"ActivityDate",
	}
	err = writer.Write(headers)
	if err != nil {
		panic(err)
	}
	handle(c, writer, n)
}

func main() {
	var (
		app     = kingpin.New("activity-see", "List all activities")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("output", "CSVf file for results").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	types := []string{
		// goengage.SubscriptionManagementType,
		//goengage.SubscriptionType,
		// goengage.FundraiseType,
		goengage.PetitionType,
		goengage.TargetedLetterType,
		// goengage.TicketedEventType,
		// goengage.P2PEventType,
	}
	c := make(chan goengage.BaseResponse, 1000)
	var wg sync.WaitGroup
	for i := 1; i < 6; i++ {
		go func(c chan goengage.BaseResponse, filename string, id int, wg *sync.WaitGroup) {
			wg.Add(1)
			startHandler(c, *csvFile, id)
			wg.Done()
		}(c, *csvFile, i, &wg)
		log.Printf("main: started handler %d\n", i)
	}

	// Listeners are all ready.  Start the talker.
	go func(e *goengage.Environment, c chan<- goengage.BaseResponse, wg *sync.WaitGroup) {
		wg.Add(1)
		for _, r := range types {
			offset := int32(0)
			count := int32(e.Metrics.MaxBatchSize)
			for count == int32(e.Metrics.MaxBatchSize) {
				payload := goengage.ActivityRequestPayload{
					Type:         r,
					Offset:       offset,
					Count:        e.Metrics.MaxBatchSize,
					ModifiedFrom: "2000-01-01T00:00:00.000Z",
				}
				rqt := goengage.ActivityRequest{
					Header:  goengage.RequestHeader{},
					Payload: payload,
				}
				var resp goengage.BaseResponse
				n := goengage.NetOp{
					Host:     e.Host,
					Method:   goengage.SearchMethod,
					Endpoint: goengage.SearchActivity,
					Token:    e.Token,
					Request:  &rqt,
					Response: &resp,
				}
				err = n.Do()
				if err != nil {
					panic(err)
				}
				c <- resp
				count = resp.Payload.Count
				offset += count
				log.Printf("main: offset %d\n", offset)
			}
		}
		close(c)
		wg.Done()
	}(e, c, &wg)
	log.Print("main: started talker")
	log.Print("main: waiting...")
	wg.Wait()
	log.Printf("main: done  Look for output files like '%s'\n", *csvFile)
}
