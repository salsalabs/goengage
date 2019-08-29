package main

import (
	"encoding/csv"
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("see-segment-census", "A command-line app to display segment names and supporters.")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("group_census.csv")
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)

	a, err := goengage.AllSegmentCensus(e, true)
	if err != nil {
		panic(err)
	}
	for _, s := range a {
		sName := s.Segment.Name
		fmt.Println(sName)
	}
	for _, s := range a {
		sName := s.Segment.Name
		fmt.Println(sName)
		for _, u := range s.Supporters {
			email := "?"
			if len(u.Contacts) > 0 {
				for _, c := range u.Contacts {
					if c.Type == goengage.ContactTypeEmail {
						email = c.Value
					}
				}
				r := []string{sName, email}
				err := w.Write(r)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%v,%v\n", sName, email)
			}
		}
	}
	w.Flush()
	f.Close()
	fmt.Println("Output can be found in group_census.csv")
}
