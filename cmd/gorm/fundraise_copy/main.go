package main

import (
	//"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/salsalabs/goengage/pkg"
	activity "github.com/salsalabs/goengage/pkg/activity"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("gorm-activity-copy", "A command-line app to copy fundraising activities to SQLite via GORM")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&activity.Fundraise{})
	db.AutoMigrate(&activity.Transaction{})

	rqtPayload := activity.ActivityRequestPayload{
		Type:         activity.FundraiseType,
		ModifiedFrom: "2010-09-01T00:00:00.000Z",
		ModifiedTo:   "2020-09-01T00:00:00.000Z",
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
	}
	rqt := activity.ActivityRequest{
		Header:  goengage.RequestHeader{},
		Payload: rqtPayload,
	}
	var resp activity.FundraiseResponse
	n := goengage.NetOp{
		Host:     e.Host,
		Endpoint: activity.Search,
		Method:   goengage.SearchMethod,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	count := int32(rqt.Payload.Count)
	for count > 0 {
		fmt.Printf("Searching from offset %d\n", rqt.Payload.Offset)
		err := n.Do()
		if err != nil {
			panic(err)
		}
		count = int32(len(resp.Payload.Activities))
		fmt.Printf("Read %d activities from offset %d\n", count, rqt.Payload.Offset)
		rqt.Payload.Offset = rqt.Payload.Offset + count
		fmt.Printf("%-36s %-36s %-10s %-10s %7s %7s %5s\n",
			"ActivityID",
			"ActivityDate",
			"ActivityType",
			"DonationType",
			"TotalReceivedAmount",
			"RecurringAmount",
			"OneTimeAmount")

		for _, s := range resp.Payload.Activities {

			db.Create(s)
			if len(s.Transactions) != 0 {
				for _, c := range s.Transactions {
					db.Create(&c)
				}
				fmt.Printf("%-36s %36s %-10s %-10s %7.2f %7.2f %7.2f\n",
					s.ActivityID,
					s.ActivityDate,
					s.ActivityType,
					s.DonationType,
					s.TotalReceivedAmount,
					s.RecurringAmount,
					s.OneTimeAmount)
			}
		}
	}
}
