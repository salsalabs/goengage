package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/salsalabs/goengage/pkg"
	supporter "github.com/salsalabs/goengage/pkg/supporter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app   = kingpin.New("see-supporter", "A command-line app to to show supporters for an email.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		email = app.Flag("email", "Email address to look up").Required().String()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		fmt.Println("Error --login is required.")
		os.Exit(1)
	}
	if email == nil || len(*email) == 0 {
		fmt.Println("Error --email is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	count := int32(e.Metrics.MaxBatchSize)
	offset := int32(0)
	// from := goengage.Date("2000-01-01T00:00:00.000Z")
	// to := goengage.Date("2100-01-01T00:00:00.000Z")
	for count == int32(e.Metrics.MaxBatchSize) {
		payload := supporter.SupporterSearchPayload{
			Identifiers:    []string{*email},
			IdentifierType: supporter.EmailAddressType,
			Offset:         offset,
			Count:          e.Metrics.MaxBatchSize,
		}
		rqt := supporter.SupporterSearch{
			Header:  goengage.Header{},
			Payload: payload,
		}
		var resp supporter.SupporterSearchResults
		n := goengage.NetOp{
			Host:     e.Host,
			Method:   goengage.SearchMethod,
			Endpoint: supporter.Search,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err = n.Do()
		if err != nil {
			panic(err)
		}
		b, _ := json.MarshalIndent(rqt, "", "    ")
		fmt.Printf("Request:\n%v\n", string(b))
		b, _ = json.MarshalIndent(resp, "", "    ")
		fmt.Printf("Response:\n%v\n", string(b))
		count = resp.Payload.Count
		for i, s := range resp.Payload.Supporters {
			email := goengage.FirstEmail(s)
			fmt.Printf("%2d %-20v %-20v %v\n",
				i+1,
				s.FirstName,
				s.LastName,
				email)
			b, _ := json.MarshalIndent(s, "", "    ")
			fmt.Println(string(b))
		}
	}
}
