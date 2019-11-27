package main

//Application scan the activities database from top to bottom and write them
//to the console.
import (
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seeBaseResponse(resp goengage.BaseResponse) {
	fmt.Println("\nActivities")
	for i, a := range resp.Payload.Activities {
		fmt.Printf("%2d %v %v %v %v %v %v %v\n",
			(i + 1),
			a.ActivityType,
			a.ActivityID,
			a.ActivityFormName,
			a.ActivityFormID,
			a.SupporterID,
			a.ActivityDate,
			a.LastModified)
	}
}

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
		//b, _ := json.MarshalIndent(n, "", "    ")
		//fmt.Printf("NetOp: %+v\n", string(b))

		err = n.Do()
		if err != nil {
			panic(err)
		}
		//b, _ = json.MarshalIndent(rqt, "", "    ")
		//fmt.Printf("Request: %+v\n", string(b))
		//b, _ = json.MarshalIndent(resp, "", "    ")
		//fmt.Printf("Response: %+v\n", string(b))
		seeBaseResponse(resp)
	}
}
