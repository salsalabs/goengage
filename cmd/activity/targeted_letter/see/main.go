package main

//Application to find and detail petition signatures.
import (
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	activity "github.com/salsalabs/goengage/pkg/activity"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seeTargetedLetterResponse(resp activity.TargetedLetterResponse) {
	fmt.Println("\nHeader")
	fmt.Printf("\tProcessingTime: %v\n", resp.Header.ProcessingTime)
	fmt.Printf("\tServerID: %v\n", resp.Header.ServerID)

	fmt.Println("\nPayload")
	fmt.Printf("\tTotal: %v\n", resp.Payload.Total)
	fmt.Printf("\tOffset: %v\n", resp.Payload.Offset)
	fmt.Printf("\tCount: %v\n", resp.Payload.Count)
	fmt.Printf("\tLength: %v\n", len(resp.Payload.Activities))

	fmt.Println("\nTargetedLetters")
	for i, a := range resp.Payload.Activities {
		fmt.Printf("\nTargetedLetter %d\n", i)
		fmt.Printf("ActivityID: %v\n", a.ActivityID)
		fmt.Printf("ActivityFormName: %v\n", a.ActivityFormName)
		fmt.Printf("ActivityFormID: %v\n", a.ActivityFormID)
		fmt.Printf("SupporterID: %v\n", a.SupporterID)
		fmt.Printf("ActivityDate: %v\n", a.ActivityDate)
		fmt.Printf("ActivityType: %v\n", a.ActivityType)
		fmt.Printf("LastModified: %v\n", a.LastModified)

		fmt.Println("Letters")
		for j, letter := range a.Letters {
			fmt.Printf("\n\tLetter %d\n", j)
			fmt.Printf("\tName: %v\n", letter.Name)
			fmt.Printf("\tSubject: %v\n", letter.Subject)
			fmt.Printf("\tMessage: %v\n", letter.Message)
			fmt.Printf("\tAdditionalComment: %v\n", letter.AdditionalComment)
			fmt.Printf("\tSubjectWasModified: %v\n", letter.SubjectWasModified)
			fmt.Printf("\tMessageWasModified: %v\n", letter.MessageWasModified)

			fmt.Println("\tTargets")
			for k, t := range letter.Targets {
				fmt.Printf("\n\t\tTarget %d\n", k)
				fmt.Printf("\t\tTargetID: %v\n", t.TargetID)
				fmt.Printf("\t\tTargetName: %v\n", t.TargetName)
				fmt.Printf("\t\tTargetTitle: %v\n", t.TargetTitle)
				fmt.Printf("\t\tPoliticalParty: %v\n", t.PoliticalParty)
				fmt.Printf("\t\tTargetType: %v\n", t.TargetType)
				fmt.Printf("\t\tState: %v\n", t.State)
				fmt.Printf("\t\tDistrictID: %v\n", t.DistrictID)
				fmt.Printf("\t\tDistrictName: %v\n", t.DistrictName)
				fmt.Printf("\t\tRole: %v\n", t.Role)
				fmt.Printf("\t\tSentEmail: %v\n", t.SentEmail)
				fmt.Printf("\t\tSentFacebook: %v\n", t.SentFacebook)
				fmt.Printf("\t\tSentTwitter: %v\n", t.SentTwitter)
				fmt.Printf("\t\tMadeCall: %v\n", t.MadeCall)
				fmt.Printf("\t\tCallDurationSeconds: %v\n", t.CallDurationSeconds)
				fmt.Printf("\t\tCallResult: %v\n", t.CallResult)
			}
		}
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
	rqt := activity.ActivityRequest{
		Type:         activity.TargetedLetterType,
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
		ModifiedFrom: "2000-01-01T00:00:00.000Z",
	}
	var resp activity.TargetedLetterResponse
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   goengage.SearchMethod,
		Endpoint: activity.Search,
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
	seeTargetedLetterResponse(resp)
}
