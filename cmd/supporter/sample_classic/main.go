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

//get retrieves active supporters from Salsa Classic.  The number of records
//is limited to the Engage instance's maximum batch size.
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
	// Reminder: true will show all Classic URLs and response bodies.
	api.Verbose = false

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
		s := goengage.SupXform(c)
		supporters = append(supporters, s)
	}

	rqt := goengage.SupUpsertRequest{}
	rqt.Supporters = supporters

	var resp goengage.SupUpsertResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SupUpsert,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}

	err = n.Upsert()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nUpsert supporter results")
	for _, s := range resp.Payload.Supporters {
		//if s.Result != "INSERTED" && s.Result != "UPDATED" {
		fmt.Printf("%-10v %v\n", s.ExternalSystemID, s.Result)
		for _, c := range s.Contacts {
			fmt.Printf("%10v %10v %20v %10v\n", "", c.Type, c.Value, c.Status)
			if len(c.Errors) > 0 {
				for _, e := range c.Errors {
					fmt.Printf("%10v %10v %v Code: %v\n", "", "", "", e.Code)
					fmt.Printf("%10v %10v %v Message: %v\n", "", "", "", e.Message)
					fmt.Printf("%10v %10v %v Details: %v\n", "", "", "", e.Details)
				}
			}
		}
		//}
		if s.Address != nil {
			a := *s.Address
			if len(a.Errors) > 0 {
				fmt.Printf("%v %v %v, %v %v\n", a.AddressLine1, a.City, a.State, a.PostalCode, a.Country)
				for _, e := range a.Errors {
					fmt.Printf("%10v %10v %v Code: %v\n", "", "", "", e.Code)
					fmt.Printf("%10v %10v %v Message: %v\n", "", "", "", e.Message)
					fmt.Printf("%10v %10v %v Details: %v\n", "", "", "", e.Details)
				}

			}
		}
	}
}
