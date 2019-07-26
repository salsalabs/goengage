package main

//Program to display the groups to which a supporter belongs.

import (
	"fmt"
	"net/http"
	"os"

	"github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app          = kingpin.New("see-segments", "A command-line app to search for segments.")
		login        = app.Flag("login", "YAML file with API token").Required().String()
		supporterKEY = app.Flag("supporter-key", "Engage supporter key for the supporter").Required().String()
		count        = app.Flag("count", "Show number of members").Bool()
	)
	app.Parse(os.Args[1:])
	if len(*supporterKEY) == 0 {
		fmt.Println("Error: --supporterKEY is REQUIRED.")
		os.Exit(1)
	}

	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}
	//Get all groups.
	a, err := goengage.AllSegments(e, m, *count)
	if err != nil {
		panic(err)
	}
	// Filter groups to keep the ones that have the supporter.
	var b []goengage.Segment
	for _, s := range a {
		sids := []string{*supporterKEY}

		rqt := goengage.SegSupporterSearchRequest{
			SegmentID:    s.SegmentID,
			SupporterIDs: sids,
			Offset:       0,
			Count:        m.MaxBatchSize,
		}
		var resp goengage.SegSupporterSearchResult
		n := goengage.NetOp{
			Host:     e.Host,
			Fragment: goengage.SegSupporterSearch,
			Method:   http.MethodPost,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err = n.Do()
		if err != nil {
			panic(err)
		}
		firstSupporter := resp.Payload.Supporters[0]
		if firstSupporter.Result == "FOUND" {
			b = append(b, s)
		}
	}
	if len(b) == 0 {
		fmt.Printf("\nSupporter with key %v is not any any groups.", *supporterKEY)
	} else {
		fmt.Printf("\nSupporter with key %v is in these groups:\n\n", *supporterKEY)
		for _, s := range b {
			if *count {
				fmt.Printf("%s %-7s %4d %s\n", s.SegmentID, s.Type, s.TotalMembers, s.Name)

			} else {
				fmt.Printf("%s %-7s %s\n", s.SegmentID, s.Type, s.Name)
			}
		}
	}
}
