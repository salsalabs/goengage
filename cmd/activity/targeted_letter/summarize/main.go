package main

//Application scan the activities database from top to bottom and write them
//to the console.
import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	goengage "github.com/salsalabs/goengage/pkg"
	activity "github.com/salsalabs/goengage/pkg/activity"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seeBaseResponse(resp activity.BaseResponse) {
	var unsorted []string
	for _, a := range resp.Payload.Activities {
		date := a.ActivityDate.Format("2006-01-02")
		s := fmt.Sprintf("%-10v %-52v %-15v\n",
			//(i + 1),
			date,
			a.ActivityFormName,
			a.SupporterID)
		unsorted = append(unsorted, s)
	}
	sort.Strings(unsorted)
	for i, s := range unsorted {
		fmt.Printf("%02d %s", i+1, s)
	}
}

func process(e *goengage.Environment, offset int32) (int32, error) {
	payload := activity.ActivityRequestPayload{
		Type:   activity.TargetedLetterType,
		Offset: int32(offset),
		Count:  e.Metrics.MaxBatchSize,
		//Pagination does *not* work with activityFormIds at this writing.
		ModifiedFrom: "2019-09-01T00:00:00.0Z",
		ModifiedTo:   "2019-10-31T00:00:00.0Z",
	}
	rqt := activity.ActivityRequest{
		Header:  goengage.Header{},
		Payload: payload,
	}
	var resp activity.BaseResponse
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   goengage.SearchMethod,
		Endpoint: activity.Search,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err := n.Do()
	if err != nil {
		return 0, err
	}
	// fmt.Printf("process:  offset: %2d, requested: %2d, total: %2d, returned %2d\n",
	// 	offset,
	// 	rqt.Payload.Count,
	// 	resp.Payload.Total,
	// 	resp.Payload.Count)
	b, _ := json.MarshalIndent(rqt, "", "    ")
	fmt.Printf("Request: %+v\n", string(b))
	b, _ = json.MarshalIndent(resp, "", "    ")
	fmt.Printf("Response: %+v\n", string(b))
	fmt.Println(strings.Repeat("=", 70))
	// seeBaseResponse(resp)
	return resp.Payload.Count, nil
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
	offset := int32(0)
	count := int32(e.Metrics.MaxBatchSize)
	for count > 0 {
		count, err = process(e, offset)
		if err != nil {
			panic(err)
		}
		offset += count
	}
}
