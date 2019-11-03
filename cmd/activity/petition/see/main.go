package main

//Application to find and detail petition signatures.
import (
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seePetitionResponse(resp goengage.PetitionResponse) {
	fmt.Println("\nHeader")
	fmt.Printf("\tProcessingTime: %v\n", resp.Header.ProcessingTime)
	fmt.Printf("\tServerID: %v\n", resp.Header.ServerID)

	fmt.Println("\nPayload")
	fmt.Printf("\tTotal: %v\n", resp.Payload.Total)
	fmt.Printf("\tOffset: %v\n", resp.Payload.Offset)
	fmt.Printf("\tCount: %v\n", resp.Payload.Count)
	fmt.Printf("\tLength: %v\n", len(resp.Payload.Activities))

	fmt.Println("\nPetitions")
	for i, a := range resp.Payload.Activities {
		fmt.Printf("\n\tPetition %d\n", i)
		fmt.Printf("\tActivityID: %v\n", a.ActivityID)
		fmt.Printf("\tActivityFormName: %v\n", a.ActivityFormName)
		fmt.Printf("\tActivityFormID: %v\n", a.ActivityFormID)
		fmt.Printf("\tSupporterID: %v\n", a.SupporterID)
		fmt.Printf("\tActivityDate: %v\n", a.ActivityDate)
		fmt.Printf("\tActivityType: %v\n", a.ActivityType)
		fmt.Printf("\tLastModified: %v\n", a.LastModified)
		fmt.Printf("\tComment: %v\n", a.Comment)
		fmt.Printf("\tModerationState: %v\n", a.ModerationState)
		fmt.Printf("\tDisplaySignaturePublicly: %v\n", a.DisplaySignaturePublicly)
		fmt.Printf("\tDisplayCommentPublicly: %v\n", a.DisplayCommentPublicly)
	}
}

func main() {
	var (
		app   = kingpin.New("activity-see", "List all activities")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	rqt := goengage.ActivityRequest{
		Type:         goengage.PetitionType,
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
		ModifiedFrom: "2010-01-01T00:00:00.000Z",
	}
	var resp goengage.PetitionResponse
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   goengage.SearchMethod,
		Endpoint: goengage.ActSearch,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	//b, _ := json.MarshalIndent(n, "", "    ")
	//fmt.Printf("NetOp: %+v\n", string(b))

	err = n.Do()
	if err != nil {
		panic(err)
	}
	//b, _ = json.MarshalIndent(rqt, "", "    ")
	//fmt.Printf("Request: %+v\n", string(b))
	//b, _ = json.MarshalIndent(resp, "", "    ")
	//fmt.Printf("Response: %+v\n", string(b))
	seePetitionResponse(resp)
}
