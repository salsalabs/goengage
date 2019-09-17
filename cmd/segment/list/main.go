package main

import (
	"fmt"
	"os"
	"sync"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//show reads a channel of segments (groups) and displays information about each.
func show(c chan goengage.Segment) {
	dashes := "----------------------------------------------------------"
	fh := "%-36v %-40v %-10v %7v %-8v %v\n"
	fl := "%-36v %-40v %-10v %7d %-8v %v\n"
	fmt.Println()
	fmt.Printf(fh, "SegmentID", "Name", "Type", "Members", "ExtID", "Description")
	fmt.Printf(fh, dashes[0:36], dashes[0:40], dashes[0:10], dashes[0:7], dashes[0:8], dashes[0:25])
	for true {
		s, ok := <-c
		if !ok {
			fmt.Println()
			return
		}
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

func main() {
	var (
		app   = kingpin.New("see-segments", "A command-line app to search for segments.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		count = app.Flag("count", "Show number of supporters (expensive)").Bool()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	c := make(chan goengage.Segment)
	var wg sync.WaitGroup

	go (func(c chan goengage.Segment, wg *sync.WaitGroup) {
		wg.Add(1)
		show(c)
		wg.Done()
	})(c, &wg)

	go (func(c chan goengage.Segment, wg *sync.WaitGroup, count bool) {
		wg.Add(1)
		err := goengage.AllSegments(e, count, c)
		if err != nil {
			panic(err)
		}
		wg.Done()
	})(c, &wg, *count)
	wg.Wait()
}
