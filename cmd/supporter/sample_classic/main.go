package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/salsalabs/godig"
	"github.com/salsalabs/goengage"
	"gopkg.in/alecthomas/kingpin.v2"
)

//App to read a number of supporter records from Salsa and
//write them to Engage.

func xform(c map[string]string) goengage.Supporter {
	// I can't find a place in engage to store job-related info.
	// leaving it out of this test.

	s := goengage.Supporter{
		FirstName:        c["First_Name"],
		LanguageCode:     c["Language_Code"],
		LastName:         c["Last_Name"],
		MiddleName:       c["MI"],
		Timezone:         c["Timezone"],
		Title:            c["Title"],
		Status:           c["Receive_Email"],
		ExternalSystemID: c["supporter_KEY"],
	}

	f := false
	af := []string{
		"AddressLine1",
		"AddressLine2",
		"City",
		"State",
		"Country",
		"PostalCode",
	}
	for _, k := range af {
		f = f || len(c[k]) > 0
	}
	if f {
		s.Address = goengage.Address{
			AddressLine1: c["Street"],
			AddressLine2: c["Street_2"],
			City:         c["City"],
			State:        c["State"],
			Country:      c["Country"],
			PostalCode:   c["Zip"],
		}
	}

	am := map[string]string{
		"Email":      "EMAIL",
		"Phone":      "HOME_PHONE",
		"Cell_Phone": "CELL_PHONE",
		"WorkPhone":  "WORK_PHONE",
	}
	as := map[string]string{
		"Email":      "OPT_IN",
		"Phone":      "",
		"Cell_Phone": "",
		"WorkPhone":  "",
	}

	var contacts []goengage.Contact
	for _, k := range af {
		if len(c[k]) > 0 {
			contact := goengage.Contact{
				Type:   am[k],
				Value:  c[k],
				Status: as[k],
			}
			contacts = append(contacts, contact)
		}
	}
	if len(contacts) > 0 {
		s.Contacts = contacts
	}
	return s
}

func get(api *godig.API, m *goengage.MetricData) ([]map[string]string, error) {
	t := api.Supporter()
	c := []string{
		"Email IS NOT EMPTY",
		"Email LIKE %@%.%",
		"Receive_Email>0",
	}
	crit := strings.Join(c, "&conmdition")
	x, err := t.ManyMap(int32(0), int(m.MaxBatchSize), crit)
	return x, err
}

func main() {

	var (
		app    = kingpin.New("engexport", "Classic-to-Engage exporter.")
		cLogin = app.Flag("login", "YAML file with Classic credentials").Required().String()
		eLogin = app.Flag("token", "YAML file with Engage API token").Required().String()
	)
	app.Parse(os.Args[1:])
	api, err := (godig.YAMLAuth(*cLogin))
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}
	e, err := goengage.Credentials(*eLogin)
	if err != nil {
		panic(err)
	}
	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}

	x, err := get(api, m)
	if err != nil {
		panic(err)
	}

	var supporters []goengage.Supporter
	for _, c := range x {
		s := xform(c)
		supporters = append(supporters, s)
	}

	rqt := goengage.SupUpsertRequest{}
	rqt.Payload.Supporters = supporters

	var resp goengage.SupUpsertResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SegSearch,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}

	// This is WRONG but it makes the compile work.
	err = n.Search()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%+v\n", resp)
}
