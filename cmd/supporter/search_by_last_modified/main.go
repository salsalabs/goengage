package main

import (
	"encoding/csv"
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Program entry point.  Look for supporters in a last_modified range.
//No values means forever.
func main() {
	var (
		app       = kingpin.New("see-supporter", "A command-line app to to show supporters for an email.")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		startDate = app.Flag("start", "Start of the date range").Default("2001-01-01T00:00:00.000Z").String()
		endDate   = app.Flag("end", "End of the date range").Default("2101-01-01T00:00:00.000Z").String()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		fmt.Println("Error --login is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("last_modified.csv")
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)

	count := int32(e.Metrics.MaxBatchSize)
	offset := int32(0)
	for count == int32(e.Metrics.MaxBatchSize) {
		payload := goengage.SupporterSearchPayload{
			IdentifierType: goengage.EmailAddressType,
			Offset:         offset,
			Count:          e.Metrics.MaxBatchSize,
			ModifiedFrom:   *startDate,
			ModifiedTo:     *endDate,
		}
		rqt := goengage.SupporterSearch{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}
		var resp goengage.SupporterSearchResults
		n := goengage.NetOp{
			Host:     e.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.SearchSupporter,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err = n.Do()
		if err != nil {
			panic(err)
		}
		count = resp.Payload.Count
		for _, s := range resp.Payload.Supporters {
			email := goengage.FirstEmail(s)
			lastModified := fmt.Sprintf("%v", s.LastModified)
			e := ""
			if email != nil {
				e = *email
			}
			record := []string{
				s.SupporterID,
				s.FirstName,
				s.LastName,
				lastModified,
				e,
			}
			err = w.Write(record)
			if err != nil {
				panic(err)
			}
		}
		offset += count
	}
}
