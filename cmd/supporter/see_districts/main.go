package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app   = kingpin.New("see_distrincts", "A command-line app to to show districs for supporter for email address(es).")
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
	headers := []string{"FirstName",
		"LastName",
		"Email",
		"Result",
		"AddressLine1",
		"AddressLine2",
		"City",
		"State",
		"PostalCode",
		"FederalHouse",
		"StateHouse",
		"StateSenate",
	}
	fn := "email_districts.csv"
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.Write(headers)
	if err != nil {
		log.Fatal(err)
	}

	// from := goengage.Date("2000-01-01T00:00:00.000Z")
	// to := goengage.Date("2100-01-01T00:00:00.000Z")
	emails := strings.Split(*email, ",")

	payload := goengage.SupporterSearchPayload{
		Identifiers:    emails,
		IdentifierType: goengage.EmailAddressType,
		Offset:         0,
		Count:          e.Metrics.MaxBatchSize,
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
	for _, s := range resp.Payload.Supporters {
		email := goengage.FirstEmail(s)
		if email != nil {
			fmt.Printf("%-36s %-25v %v\n", s.SupporterID, *email, s.Result)
		} else { 
			fmt.Printf("%-36s %-25v %v\n", s.SupporterID, "(None)", s.Result)
		}
		if s.Result != goengage.Found {
			fmt.Printf("%v %v %v\n",
				s.FirstName,
				s.LastName,
				s.Result)
			continue
		}
		var r []string
		e := ""
		if email != nil {
			e = *email
		}
		for i := 0; i < len(headers); i++ {
			switch i {
			case 0:
				r = append(r, s.FirstName)
			case 1:
				r = append(r, s.LastName)
			case 2:
				r = append(r, e)
			case 3:
				r = append(r, s.Result)
			case 4:
				r = append(r, s.Address.AddressLine1)
			case 5:
				r = append(r, s.Address.AddressLine2)
			case 6:
				r = append(r, s.Address.City)
			case 7:
				r = append(r, s.Address.State)
			case 8:
				r = append(r, s.Address.PostalCode)
			case 9:
				r = append(r, s.Address.FederalHouseDistrict)
			case 10:
				r = append(r, s.Address.StateHouseDistrict)
			case 11:
				r = append(r, s.Address.StateSenateDistrict)
			}
		}
		w.Write(r)
	}
	rqt.Payload.Offset += rqt.Payload.Count
	w.Flush()
	fmt.Printf("Done.  Output is in %v\n", fn)
}
