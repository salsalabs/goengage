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
func get(api *godig.API, offset int32, m *goengage.MetricData) ([]map[string]string, error) {
	t := api.Supporter()
	c := []string{
		"Email IS NOT EMPTY",
		"Email LIKE %@%.%",
		"Receive_Email>0",
	}
	crit := strings.Join(c, "&conmdition")
	x, err := t.ManyMap(offset, int(m.MaxBatchSize), crit)
	log.Printf("get: offset %7d, count %d\n", offset, len(x))
	return x, err
}

//note writes a line to a file and dies if there's an error.
func note(s string, f *os.File) {
	_, err := f.WriteString(s)
	if err != nil {
		log.Fatalf("Error writing to log, %v\n", err)
	}
}

//see shows the interesting bits of a supporter record.
func see(s goengage.Supporter, f *os.File) {
	note("\n", f)
	note(fmt.Sprintf("%v %v\n", s.ExternalSystemID, s.Result), f)
	for _, c := range s.Contacts {
		note(fmt.Sprintf("%v Contact %v, %v, %v\n", s.ExternalSystemID, c.Type, c.Value, c.Status), f)
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				note(fmt.Sprintf("%v Contact Code: %v\n", s.ExternalSystemID, e.Code), f)
				note(fmt.Sprintf("%v Contact Message: %v\n", s.ExternalSystemID, e.Message), f)
				note(fmt.Sprintf("%v Contact Details: %v\n", s.ExternalSystemID, e.Details), f)
			}
		}
	}
	//}
	if s.Address != nil {
		a := *s.Address
		if len(a.Errors) > 0 {
			note(fmt.Sprintf(`%v Address "%v" "%v" "%v", "%v" "%v"\n`, s.ExternalSystemID, a.AddressLine1, a.City, a.State, a.PostalCode, a.Country), f)
			for _, e := range a.Errors {
				note(fmt.Sprintf("%v Address Code: %v\n", s.ExternalSystemID, e.Code), f)
				note(fmt.Sprintf("%v Address Message: %v\n", s.ExternalSystemID, e.Message), f)
				note(fmt.Sprintf("%v Address Details: %v\n", s.ExternalSystemID, e.Details), f)
			}

		}
	}
}

//main is the entry point for Go applications.
func main() {
	var (
		app    = kingpin.New("engexport", "Classic-to-Engage exporter.")
		eLogin = app.Flag("login", "YAML file with Engage API token").Required().String()
		cLogin = app.Flag("classic", "YAML file with Classic credentials").Required().String()
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
		log.Fatalf("Main: %v\n", err)
	}
	m, err := e.Metrics()
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}

	f, err := os.Create("transfer_results.txt")
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}

	count := m.MaxBatchSize
	offset := int32(0)
	//Make this "count > 0" to copy everything from classic
	//to engage in MaxBatchSize chunks.  Will be slow...
	for offset < 100 {
		x, err := get(api, offset, m)
		if err != nil {
			log.Fatalf("Main: %v\n", err)
		}

		count = int32(len(x))
		offset = offset + count
		if count == 0 {
			break
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
			log.Fatalf("Main: %v\n", err)
		}
		for _, s := range resp.Payload.Supporters {
			see(s, f)
		}
	}
	_ = f.Close()
	fmt.Println("Transfer result details can be found in transfer_results.txt")
}
