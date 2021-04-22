package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Runtime is used to pass arguments between methods.
type Runtime struct {
	E         *goengage.Environment
	StartDate string
	EndDate   string
	W         *csv.Writer
}

//Zippopatamus is the record that's returned for a postalcode lookup.
type Zippopatamus struct {
	PostCode            string `json:"post code"`
	Country             string `json:"country"`
	CountryAbbreviation string `json:"country abbreviation"`
	Places              []struct {
		PlaceName         string `json:"place name"`
		Longitude         string `json:"longitude"`
		State             string `json:"state"`
		StateAbbreviation string `json:"state abbreviation"`
		Latitude          string `json:"latitude"`
	} `json:"places"`
}

//Fix an address record using a call to Zippopatm.us.  Returns the modified
//address record.  Errors trigger panics.
func fixWithZip(a *goengage.Address) (*goengage.Address, error) {
	url := fmt.Sprintf("https://api.zippopotam.us/us/%s", a.PostalCode)
	// log.Printf("fixWithZip url: '%s'", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// log.Printf("fixWithZip body: '%s'", string(body))
	if err != nil {
		return a, err
	}
	// Can get back an empty struct ("{}") here...
	if len(body) > 2 {
		var z Zippopatamus
		err = json.Unmarshal(body, &z)
		if err != nil {
			return a, err
		}
		// log.Printf("fixWithZip z: %v\n", z)
		a.City = z.Places[0].PlaceName
		a.State = z.Places[0].StateAbbreviation
	}
	return a, nil
}

//Process one supporter record by fixing city and/or state using a lookup from
//zippopatm.us.  Outputs a CSV row if either the city or state changes. Errors
//trigger panics.
func process(rt Runtime, s goengage.Supporter) {
	e := ""
	email := goengage.FirstEmail(s)
	if email != nil {
		e = *email
	}
	a := s.Address
	if a == nil {
		a = &goengage.Address{
			City:       "",
			State:      "",
			PostalCode: "",
		}
	}
	zipCheck := len(a.PostalCode) != 0
	cityStateCheck := (len(a.City) == 0 || len(a.State) == 0)
	countryCheck := len(a.Country) == 0 || (len(a.Country) != 0 && strings.ToUpper(a.Country) != "US")
	if zipCheck && cityStateCheck && countryCheck {
		c := a.City
		st := a.State
		fixWithZip(a)
		cityModified := a.City != c
		stateModified := a.State != st
		if cityModified || stateModified {
			record := []string{
				s.SupporterID,
				e,
				a.City,
				c,
				a.State,
				st,
				a.PostalCode,
				a.Country,
			}
			err := rt.W.Write(record)
			log.Printf("%v\n", record)
			if err != nil {
				panic(err)
			}
		}
	}
}

//Drives the process by processing all supporters. Errors panic.
func drive(rt Runtime) {
	count := int32(rt.E.Metrics.MaxBatchSize)
	offset := int32(0)
	for count == int32(rt.E.Metrics.MaxBatchSize) {
		payload := goengage.SupporterSearchPayload{
			IdentifierType: goengage.EmailAddressType,
			Offset:         offset,
			Count:          rt.E.Metrics.MaxBatchSize,
			ModifiedFrom:   rt.StartDate,
			ModifiedTo:     rt.EndDate,
		}
		rqt := goengage.SupporterSearch{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}
		var resp goengage.SupporterSearchResults
		n := goengage.NetOp{
			Host:     rt.E.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.SearchSupporter,
			Token:    rt.E.Token,
			Request:  &rqt,
			Response: &resp,
		}
		if offset%1000 == int32(0) {
			log.Printf("main: %5d\n", offset)
		}
		err := n.Do()
		if err != nil {
			panic(err)
		}
		count = resp.Payload.Count
		for _, s := range resp.Payload.Supporters {
			process(rt, s)
		}
		offset += count
	}
}

//Program entry point.  Look for supporters in a last_modified range.
//No values means forever.
func main() {
	var (
		app       = kingpin.New("ZIP City State Lookup", "Use Zippotam.us to find missing states and cities by postalCode")
		login     = app.Flag("login", "YAML file with API token").Required().String()
		startDate = app.Flag("start", "Last modified start").Default("2001-01-01T00:00:00.000Z").String()
		endDate   = app.Flag("end", "Last modified end").Default("2101-01-01T00:00:00.000Z").String()
		csvFile   = app.Flag("csv", "CSV to receive modified records").Default("zip_city_state_fixes.csv").String()
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

	f, err := os.Create(*csvFile)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	h := strings.Split("SupporterID,Email,City,OriginalCity,State,OriginalState,PostalCode,Country", ",")
	w.Write(h)

	rt := Runtime{
		E:         e,
		StartDate: *startDate,
		EndDate:   *endDate,
		W:         w,
	}
	drive(rt)
	close(w)
}
