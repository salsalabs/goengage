package main

import (
	"fmt"
	"os"
	"strings"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app   = kingpin.New("see-supporter", "A command-line app to to show supporters for an email.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		email = app.Flag("email", "Comma-separated list of email addresses to look up").Required().String()
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
	emails := strings.Split(*email, ",")
	for count == int32(e.Metrics.MaxBatchSize) {
		payload := goengage.SupporterSearchRequestPayload{
			Identifiers:    emails,
			IdentifierType: goengage.EmailAddressType,
			Offset:         offset,
			Count:          e.Metrics.MaxBatchSize,
		}
		rqt := goengage.SupporterSearchRequest{
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
		for i, s := range resp.Payload.Supporters {
			email := goengage.FirstEmail(s)
			e := ""
			if email != nil {
				e = *email
			}
			fmt.Printf("%2d %-20v %-20v %-25v %-25v %v\n",
				i+1,
				s.FirstName,
				s.LastName,
				s.CreatedDate,
				s.LastModified,
				e)
		}
	}
}
