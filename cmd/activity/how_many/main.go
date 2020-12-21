package main

//Application scan the activities database from top to bottom and write them
//to the console.
import (
	"log"
	"math"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("how-many", "See number of activities")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	types := []string{
		goengage.SubscriptionManagementType,
		goengage.SubscriptionType,
		goengage.FundraiseType,
		goengage.PetitionType,
		goengage.TargetedLetterType,
		goengage.TicketedEventType,
		goengage.P2PEventType,
	}
	for _, r := range types {
		offset := int32(0)
		payload := goengage.ActivityRequestPayload{
			Type:         r,
			Offset:       offset,
			Count:        0,
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
		passes := int32(math.Ceil(float64(resp.Payload.Total) / float64(e.Metrics.MaxBatchSize)))
		log.Printf("%-27s %6d %6d\n", r, resp.Payload.Total, passes)
	}
}
