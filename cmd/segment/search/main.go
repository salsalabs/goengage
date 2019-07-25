package main

import (
	"fmt"
	"os"

	"github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("see-segments", "A command-line app to search for segments.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		fast  = app.Flag("fast", "Don't show number of members").Default("false").Bool()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}
	a, err := goengage.AllSegments(e, m, !*fast)
	if err != nil {
		panic(err)
	}
	dashes := "----------------------------------------------------------"
	fh := "%-36v %-40v %-10v %7v %-8v %v\n"
	fl := "%-36v %-40v %-10v %7d %-8v %v\n"
	fmt.Println()
	fmt.Printf(fh, "SegmentID", "Name", "Type", "Members", "ExtID", "Description")
	fmt.Printf(fh, dashes[0:36], dashes[0:40], dashes[0:10], dashes[0:7], dashes[0:8], dashes[0:25])
	for _, s := range a {
		name := s.Name
		if len(name) > 40 {
			name = name[0:37] + "..."
		}
		fmt.Printf(fl,
			s.SegmentID,
			name,
			s.Type,
			s.TotalMembers,
			s.ExternalSystemID,
			s.Description)
	}
	fmt.Println()
}
