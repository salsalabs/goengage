package main

//Application scan the activities database from top to bottom and write them
//to the console.
import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("activity-see", "List all activities")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	types := []string{
		goengage.SubscriptionManagementType,
		goengage.SubscribeType,
		goengage.FundraiseType,
		goengage.PetitionType,
		goengage.TargetedLetterType,
		goengage.TicketedEventType,
		goengage.P2PEventType,
	}
	for _, r := range types {
		rqt := goengage.ActivityRequest{
			Type:         r,
			Offset:       0,
			Count:        e.Metrics.MaxBatchSize,
			ModifiedFrom: "2010-01-01T00:00:00.000Z",
		}
		var resp goengage.ActivityResponse
		n := goengage.NetOp{
			Host:     e.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.ActSearch,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err = n.Do()
		if err != nil {
			panic(err)
		}
		b, _ := json.MarshalIndent(rqt, "", "    ")
		fmt.Printf("Request: %+v\n", string(b))
		b, _ = json.MarshalIndent(resp, "", "    ")
		fmt.Printf("Response: %+v\n", string(b))
	}
}
