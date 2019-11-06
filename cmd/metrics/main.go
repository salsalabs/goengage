package main

import (
	"fmt"
	"os"
	"strings"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("metrics", "A command-line app to display the current Engage metrics for a token.")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	m := e.Metrics
	fmt.Println()
	fmt.Printf("%-30v %v\n", "Setting", "Value")
	fmt.Printf("%-30v %v\n", strings.Repeat("-", 30), strings.Repeat("-", 25))
	fmt.Printf("%-30v %v\n", "RateLimit", m.RateLimit)
	fmt.Printf("%-30v %v\n", "MaxBatchSize", e.Metrics.MaxBatchSize)
	fmt.Printf("%-30v %v\n", "SupporterRead", m.SupporterRead)
	fmt.Printf("%-30v %v\n", "SupporterAdd", m.SupporterAdd)
	fmt.Printf("%-30v %v\n", "SupporterDelete", m.SupporterDelete)
	fmt.Printf("%-30v %v\n", "SupporterUpdate", m.SupporterUpdate)
	fmt.Printf("%-30v %v\n", "SegmentRead", m.SegmentRead)
	fmt.Printf("%-30v %v\n", "SegmentAdd", m.SegmentAdd)
	fmt.Printf("%-30v %v\n", "SegmentDelete", m.SegmentDelete)
	fmt.Printf("%-30v %v\n", "SegmentUpdate", m.SegmentUpdate)
	fmt.Printf("%-30v %v\n", "SegmentAssignmentRead", m.SegmentAssignmentRead)
	fmt.Printf("%-30v %v\n", "SegmentAssignmentAdd", m.SegmentAssignmentAdd)
	fmt.Printf("%-30v %v\n", "SegmentAssignmentDelete", m.SegmentAssignmentDelete)
	fmt.Printf("%-30v %v\n", "SegmentAssignmentUpdate", m.SegmentAssignmentUpdate)
	fmt.Printf("%-30v %v\n", "OfflineDonationAdd", m.OfflineDonationAdd)
	fmt.Printf("%-30v %v\n", "OfflineDonationUpdate", m.OfflineDonationUpdate)
	fmt.Printf("%-30v %v\n", "ActivityTicketedEvent", m.ActivityTicketedEvent)
	fmt.Printf("%-30v %v\n", "ActivityP2PEvent", m.ActivityP2PEvent)
	fmt.Printf("%-30v %v\n", "ActivityPetition", m.ActivityPetition)
	fmt.Printf("%-30v %v\n", "ActivitySubscribe", m.ActivitySubscribe)
	fmt.Printf("%-30v %v\n", "ActivityFundraise", m.ActivityFundraise)
	fmt.Printf("%-30v %v\n", "ActivityTargetedLetter", m.ActivityTargetedLetter)
	fmt.Printf("%-30v %v\n", "ActivitySubscriptionManagement", m.ActivitySubscriptionManagement)
	fmt.Printf("%-30v %v\n", "LastAPICall", m.LastAPICall)
	fmt.Printf("%-30v %v\n", "TotalAPICalls", m.TotalAPICalls)
	fmt.Printf("%-30v %v\n", "TotalAPICallFailures", m.TotalAPICallFailures)
	fmt.Printf("%-30v %v\n", "CurrentRateLimit", m.CurrentRateLimit)
	fmt.Println()
}
