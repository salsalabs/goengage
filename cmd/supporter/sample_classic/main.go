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
	fmt.Printf("\n%+v\n", resp)
}
