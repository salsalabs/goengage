package main

//Program to display the groups to which a supporter belongs.

import (
	"fmt"
	"os"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//show the groups for a supporter.  You provide the primary key for the supporter.
func show(e *goengage.Environment, c chan goengage.Segment, k string) error {
	count := 0

	// Filter groups to keep the ones that have the supporter.
	for true {
		s, ok := <-c
		if !ok {
			break
		}
		sids := []string{k}

		rqt := goengage.SegSupporterSearchRequest{
			SegmentID:    s.SegmentID,
			SupporterIDs: sids,
			Offset:       0,
			Count:        e.Metrics.MaxBatchSize,
		}
		var resp goengage.SegSupporterSearchResult
		n := goengage.NetOp{
			Host:     e.Host,
			Endpoint: goengage.SegSupporterSearch,
			Method:   goengage.SearchMethod,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err := n.Do()
		if err != nil {
			return err
		}
		if len(resp.Supporters) > 0 {
			firstSupporter := resp.Supporters[0]
			if firstSupporter.Result == "FOUND" {
				if count == 0 {
					fmt.Printf("\nSupporter with key %v is in these groups:\n\n", k)
				}
				fmt.Printf("%s %-7s %s\n", s.SegmentID, s.Type, s.Name)
				count++
			}
		}
	}
	if count == 0 {
		fmt.Printf("\nSupporter with key %v is not any any groups.\n", k)
	}
	return nil
}

func main() {
	var (
		app          = kingpin.New("see-segments", "A command-line app to search for segments.")
		login        = app.Flag("login", "YAML file with API token").Required().String()
		supporterKEY = app.Flag("supporter-key", "Engage supporter key for the supporter").Required().String()
		count        = app.Flag("count", "Show number of members").Bool()
	)
	app.Parse(os.Args[1:])
	if len(*supporterKEY) == 0 {
		fmt.Println("Error: --supporter-key is REQUIRED.")
		os.Exit(1)
	}

	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	c := make(chan goengage.Segment)
	var wg sync.WaitGroup

	//Display segments that contain the supporter.  Panicking on error until we find an
	//elegant way to handle errors in a goroutine.
	go (func(e *goengage.Environment, c chan goengage.Segment, k string, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := show(e, c, k)
		if err != nil {
			panic(err)
		}
	})(e, c, *supporterKEY, &wg)
	fmt.Printf("Started show...")

	//Drive segments to the filter for a supporter.  Panicking on error until we find an
	//elegant way to handle errors in a goroutine.
	go (func(e *goengage.Environment, c chan goengage.Segment, count bool, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := goengage.AllSegments(e, count, c)
		if err != nil {
			panic(err)
		}
	})(e, c, *count, &wg)
	fmt.Println("Started AllSegments...")

	fmt.Println("Waiting patiently...")
	wg.Wait()
	fmt.Println("Done!")
}
