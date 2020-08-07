package main

//Application scan the activities database from top to bottom and write them
//to the console.
import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seeBaseResponse(resp goengage.BaseResponse, writer *csv.Writer) {
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
	writer.Flush()
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
	f, err := os.Create(*csvFile)
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
			fmt.Printf("Payload: %+v\n", payload)
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
			seeBaseResponse(resp, writer)
			count = resp.Payload.Count
			offset += count
		}
	}
}
