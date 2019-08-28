package main

//Application scan the activities database from top to bottom and write them
//to the console.
import (
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
		//goengage.SubscriptionManagementType,
		//goengage.SubscribeType,
		//goengage.FundraiseType,
		goengage.PetitionType,
		//goengage.TargetedLetterType,
		//goengage.TicketedEventType,
		//goengage.P2PEventType,
	}

	for _, r := range types {
		rqt := goengage.ActivityRequest{
			Type:         r,
			Offset:       0,
			Count:        e.Metrics.MaxBatchSize,
			ModifiedFrom: "2010-01-01T00:00:00.000Z",
			//ModifiedTo:   "2020-12-31T23:59:59.000Z",
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
		fmt.Printf("\nActivity Type: %v\n", r)
		fmt.Printf("\nRequest: %+v\n", n)
		err = n.Do()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Response: %+v\n", resp)
	}
}
