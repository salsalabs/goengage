package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("activity-search", "A command-line app to see all supporters.")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}
	rqt := goengage.SupSearchRequest{
		ModifiedFrom: "2016-09-01T00:00:00.000Z",
		ModifiedTo:   "2019-09-01T00:00:00.000Z",
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SupSearch,
		Method:   goengage.SearchMethod,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	count := int32(rqt.Count)
	for count > 0 {
		fmt.Printf("Searching from offset %d\n", rqt.Offset)
		err := n.Do()
		if err != nil {
			panic(err)
		}
		count = int32(len(resp.Payload.Supporters))
		fmt.Printf("Read %d supporters from offset %d\n", count, rqt.Offset)
		rqt.Offset = rqt.Offset + count
		for _, s := range resp.Payload.Supporters {
			e := goengage.FirstEmail(s)
			email := ""
			if e != nil && len(*e) > 0 {
				email = *e
			}
			fmt.Printf("%-20s %-20s %s\n", s.FirstName, s.LastName, email)
		}
	}
}
