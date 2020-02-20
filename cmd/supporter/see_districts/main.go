package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app     = kingpin.New("see_distrincts", "A command-line app to to show districs for supporter for email address(es).")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("csv", "Comma-separated file of email addresses to look up.  Email must be first.").Required().String()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		fmt.Println("Error --login is required.")
		os.Exit(1)
	}
	if csvFile == nil || len(*csvFile) == 0 {
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

	g, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	r := csv.NewReader(g)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read %d rows from %s\n", len(rows), *csvFile)
	//Clip headers.
	rows = rows[1:len(rows)]

	//Do in legal sized batches.  Each batch is a list of
	//string slices.  Dereference those to get a list of emails.
	for i := 0; i < len(rows); i += int(e.Metrics.MaxBatchSize) {
		var emails []string
		batch := rows[i : i+int(e.Metrics.MaxBatchSize)]
		log.Printf("Batch at offset %d contains %d emails.", i, len(batch))
		for _, v := range batch {
			if len(v) > 0 {
				emails = append(emails, v[0])
			}
		}

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

			if s.Result != goengage.Found {
				fmt.Printf("%v %v %v %v\n",
					s.FirstName,
					s.LastName,
					email,
					s.Result)
				continue
			}
			var a []string
			e := ""
			if email != nil {
				e = *email
			}
			for i := 0; i < len(headers); i++ {
				switch i {
				case 0:
					a = append(a, s.FirstName)
				case 1:
					a = append(a, s.LastName)
				case 2:
					a = append(a, e)
				case 3:
					a = append(a, s.Result)
				case 4:
					a = append(a, s.Address.AddressLine1)
				case 5:
					a = append(a, s.Address.AddressLine2)
				case 6:
					a = append(a, s.Address.City)
				case 7:
					a = append(a, s.Address.State)
				case 8:
					a = append(a, s.Address.PostalCode)
				case 9:
					a = append(a, s.Address.FederalHouseDistrict)
				case 10:
					a = append(a, s.Address.StateHouseDistrict)
				case 11:
					a = append(a, s.Address.StateSenateDistrict)
				}
			}
			w.Write(a)
		}
	}
	w.Flush()
	fmt.Printf("Done.  Output is in %v\n", fn)
}
