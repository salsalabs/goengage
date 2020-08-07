package main

//Application scan the activities database from top to bottom and write them
//to the console.
import (
	"encoding/csv"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seeBaseResponse(resp goengage.BaseResponse, writer *csv.Writer) {
	for _, a := range resp.Payload.Activities {
		record := []string{
			a.SupporterID,
			a.PersonName,
			a.PersonEmail,
			a.ActivityType,
			a.ActivityFormName,
			a.ActivityFormID,
		}
		err := writer.Write(record)
		if err != nil {
			panic(err)
		}
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
		goengage.SubscriptionType,
		// goengage.FundraiseType,
		// goengage.PetitionType,
		// goengage.TargetedLetterType,
		// goengage.TicketedEventType,
		// goengage.P2PEventType,
	}
	f, err := os.Create(*csvFile)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(f)
	for _, r := range types {
		offset := int32(0)
		count := int32(e.Metrics.MaxBatchSize)
		for count > 0 {
			payload := goengage.ActivityRequestPayload{
				Type:         r,
				Offset:       0,
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
			// b, _ := json.MarshalIndent(n, "", "    ")
			// fmt.Printf("NetOp: %+v\n", string(b))

			err = n.Do()
			if err != nil {
				panic(err)
			}
			//b, _ = json.MarshalIndent(rqt, "", "    ")
			//fmt.Printf("Request: %+v\n", string(b))
			//b, _ = json.MarshalIndent(resp, "", "    ")
			//fmt.Printf("Response: %+v\n", string(b))
			seeBaseResponse(resp, writer)
			count = resp.Payload.Count
			offset += count
		}
	}
}
