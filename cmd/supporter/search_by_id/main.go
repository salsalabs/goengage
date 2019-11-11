package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"

	goengage "github.com/salsalabs/goengage/pkg"
	supporter "github.com/salsalabs/goengage/pkg/supporter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app     = kingpin.New("activity-search", "A command-line app to see emails for a list of supporter IDs.")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("csv", "CSV file with IDs.  Uses 'InternalID'.").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(*csvFile)
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)
	//records is an array of records.  Each record is
	//an array of strings with these offsets.
	//0 InternalID
	a, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	_ = f.Close()

	var lines []string
	for _, r := range a {
		if r[0] != "InternalID" {
			searchPayload := supporter.SupporterSearchPayload{
				IdentifierType: supporter.SupporterIDType,
				Identifiers:    []string{r[0]},
			}
			rqt := supporter.SupporterSearch{
				Header:  goengage.RequestHeader{},
				Payload: searchPayload,
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
			err = n.Do()
			if err != nil {
				panic(err)
			}

			for _, s := range resp.Payload.Supporters {
				x := goengage.FirstEmail(s)
				e := "(None)"
				if x != nil {
					e = *x
				}
				name := fmt.Sprintf("%s %s", s.FirstName, s.LastName)
				t := fmt.Sprintf("%-20s %-30s %-36v\n",
					name,
					e,
					s.CreatedDate)
				lines = append(lines, t)
			}
		}
	}
	sort.Strings(lines)
	for _, t := range lines {
		fmt.Printf(t)
	}
}
