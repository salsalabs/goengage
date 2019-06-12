package main

import (
	"fmt"
	"os"

	"github.com/salsalabs/goengage/pkg"
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

	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}
	dashes := "----------------------------------------------------------"
	fmt.Println()
	fmt.Printf("%-30v %v\n", "Setting", "Value")
	fmt.Printf("%-30v %v\n", dashes[0:30], dashes[0:25])
	fmt.Printf("%-30v %v\n", "RateLimit", m.RateLimit)
	fmt.Printf("%-30v %v\n", "MaxBatchSize", m.MaxBatchSize)
	fmt.Printf("%-30v %v\n", "CurrentRateLimit", m.CurrentRateLimit)
	fmt.Printf("%-30v %v\n", "TotalAPICalls", m.TotalAPICalls)
	fmt.Printf("%-30v %v\n", "LastAPICall", m.LastAPICall)
	fmt.Printf("%-30v %v\n", "TotalAPICallFailures", m.TotalAPICallFailures)
	fmt.Printf("%-30v %v\n", "LastAPICallFailure", m.LastAPICallFailure)
	fmt.Printf("%-30v %v\n", "SupporterRead", m.SupporterRead)
	fmt.Printf("%-30v %v\n", "SupporterAdd", m.SupporterAdd)
	fmt.Printf("%-30v %v\n", "SupporterUpdate", m.SupporterUpdate)
	fmt.Printf("%-30v %v\n", "SupporterDelete", m.SupporterDelete)
	fmt.Printf("%-30v %v\n", "ActivityEvent", m.ActivityEvent)
	fmt.Printf("%-30v %v\n", "ActivitySubscribe", m.ActivitySubscribe)
	fmt.Printf("%-30v %v\n", "ActivityFundraise", m.ActivityFundraise)
	fmt.Printf("%-30v %v\n", "ActivityTargetedLetter", m.ActivityTargetedLetter)
	fmt.Printf("%-30v %v\n", "ActivityPetition", m.ActivityPetition)
	fmt.Printf("%-30v %v\n", "ActivitySubscriptionManagement", m.ActivitySubscriptionManagement)
	fmt.Println()
}
