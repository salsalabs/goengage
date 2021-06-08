package main

//Application to accept a segmentId and output the supporters that belong
//to the segment.  Output includes a list of the other segments that a
//supporter belongs to.  Produces a CSV of supporter_KEY, Email, Groups.

import (
	"log"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Run finds and displays all segments.
func Run(env *goengage.Environment) error {
	log.Println("Run: begin")

	count := env.Metrics.MaxBatchSize
	offset := int32(0)

	payload := goengage.SegmentSearchRequestPayload{
		// IncludeMemberCounts: false,
	}

	rqt := goengage.SegmentSearchRequest{
		Header:  goengage.RequestHeader{},
		Payload: payload,
	}

	for count == env.Metrics.MaxBatchSize {
		payload.Offset = offset
		payload.Count = count
		var resp goengage.SegmentSearchResponse

		n := goengage.NetOp{
			Host:     env.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.SearchSegment,
			Token:    env.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err := n.Do()
		if err != nil {
			return err
		}
		if offset%500 == 0 {
			log.Printf("Run: %6d: %2d of %6d\n",
				offset,
				len(resp.Payload.Segments),
				resp.Payload.Total)
		}

		log.Printf("payload: %+v", resp.Payload)
		for _, s := range resp.Payload.Segments {
			log.Printf("%-36s %-50s %6d ",
				s.SegmentID,
				s.Name,
				s.TotalMembers)
		}
		count = resp.Payload.Count
		offset += int32(count)
	}
	return nil
}

//Program entry point.
func main() {
	var (
		app   = kingpin.New("one_segment_xref", "Lists segments for an organization.")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		log.Fatalf("Error --login is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("Error: %+v\n", e)
	}
	err = Run(e)
	if err != nil {
		panic(err)
	}
}
