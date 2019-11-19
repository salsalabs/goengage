package main

import (
	//"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/salsalabs/goengage/pkg"
	supporter "github.com/salsalabs/goengage/pkg/supporter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("supporter-search", "A command-line app to see all supporters.")
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
	db.AutoMigrate(&goengage.Supporter{})
	db.AutoMigrate(&goengage.Contact{})
	db.AutoMigrate(&goengage.CustomFieldValue{})

	rqtPayload := supporter.SupporterSearchPayload{
		ModifiedFrom: "2016-09-01T00:00:00.000Z",
		ModifiedTo:   "2019-09-01T00:00:00.000Z",
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
	}
	rqt := supporter.SupporterSearch{
		Header:  goengage.RequestHeader{},
		Payload: rqtPayload,
	}
	var resp supporter.SupporterSearchResults
	n := goengage.NetOp{
		Host:     e.Host,
		Endpoint: supporter.Search,
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
		count = int32(len(resp.Payload.Supporters))
		fmt.Printf("Read %d supporters from offset %d\n", count, rqt.Payload.Offset)
		rqt.Payload.Offset = rqt.Payload.Offset + count
		for _, s := range resp.Payload.Supporters {
			db.Create(s)
			if s.Contacts != nil {
				for _, c := range s.Contacts {
					db.Create(&c)
				}
			}
			if s.CustomFieldValues != nil {
				for _, c := range s.CustomFieldValues {
					db.Create(&c)
				}
			}
			e := goengage.FirstEmail(s)
			email := ""
			if e != nil && len(*e) > 0 {
				email = *e
			}
			fmt.Printf("%-20s %-20s %s\n", s.FirstName, s.LastName, email)
		}
	}
}
