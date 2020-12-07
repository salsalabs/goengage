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
		app   = kingpin.New("see_districts", "A command-line app to write supporters and state districts to a CSV.")
		login = app.Flag("login", "YAML file with API token").Required().String()
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
	fn := "supporters_and_districts.csv"
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	headers := []string{"FirstName",
		"LastName",
		"Email",
		"Phone",
		"AddressLine1",
		"AddressLine2",
		"City",
		"State",
		"PostalCode",
		"StateHouse",
		"StateSenate",
	}
	err = w.Write(headers)
	if err != nil {
		log.Fatal(err)
	}

	withDistricts := int32(0)
	count := e.Metrics.MaxBatchSize
	offset := int32(0)
	for count == e.Metrics.MaxBatchSize {
		payload := goengage.SupporterSearchPayload{
			ModifiedFrom: "2006-01-01T00:00:00.000Z",
			Offset:       offset,
			Count:        e.Metrics.MaxBatchSize,
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
				// emailTemp := ""
				// if email != nil {
				// 	emailTemp = *email
				// }
				// log.Printf("%v %v %v %v\n",
				// 	s.FirstName,
				// 	s.LastName,
				// 	emailTemp,
				// 	s.Result)
				continue
			}
			if len(s.Address.StateHouseDistrict) == 0 && len(s.Address.StateSenateDistrict) == 0 {
				// log.Printf("%v %v %v %v\n",
				// 	s.FirstName,
				// 	s.LastName,
				// 	emailTemp,
				// 	"no districts")
				continue
			}
			phone := goengage.FirstPhone(s)
			var a []string
			e := ""
			if email != nil {
				e = *email
			}
			p := ""
			if phone != nil {
				p = *phone
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
					a = append(a, p)
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
					a = append(a, s.Address.StateHouseDistrict)
				case 10:
					a = append(a, s.Address.StateSenateDistrict)
				}
			}
			w.Write(a)
			withDistricts++
		}
		count = resp.Payload.Count
		offset += count
		log.Printf("Status: %6d of %6d with districts out of  %6d total\n", withDistricts, offset, resp.Payload.Total)

	}
	w.Flush()
	log.Printf("Done.  Output is in %v\n", fn)
}
