package main

import (
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	enggoengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("see-segments", "A command-line app to search for segments.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		name  = app.Flag("name", "Group name").Required().String()
		desc  = app.Flag("description", "Group description").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	rqt := enggoengage.SegmentUpsertRequest{
        Payload: {
            Segments: [ {
                Name: name,
                Description: desc,
            }]
	}
}
