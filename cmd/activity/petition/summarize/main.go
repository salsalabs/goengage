package main

//Application reads all petitions and shows the petition
//and the list of action takers.
import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seePetitionResponse(resp goengage.PetitionResponse) {
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

func process(e *goengage.Environment, writer *csv.Writer, offset int32) (int32, error) {
	payload := goengage.ActivityRequestPayload{
		Type:   goengage.PetitionType,
		Offset: int32(offset),
		Count:  e.Metrics.MaxBatchSize,
		//Pagination does *not* work with activityFormIds at this writing.
		ModifiedFrom: "2006-09-01T00:00:00.0Z",
		ModifiedTo:   "2022-10-31T00:00:00.0Z",
	}
	rqt := goengage.ActivityRequest{
		Header:  goengage.RequestHeader{RefID: "cmd/activity/petition/summarize"},
		Payload: payload,
	}
	var resp goengage.PetitionResponse
	logger, err := goengage.NewUtilLogger()
	if err != nil {
		return 0, err
	}
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   goengage.SearchMethod,
		Endpoint: goengage.SearchActivity,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
		Logger:   logger,
	}
	err = n.Do()
	if err != nil {
		return 0, err
	}
	fmt.Printf("process:  offset: %2d, requested: %2d, total: %2d, returned %2d\n",
		offset,
		rqt.Payload.Count,
		resp.Payload.Total,
		resp.Payload.Count)

	for _, r := range resp.Payload.Activities {
		record := []string{
			r.SupporterID,
			r.PersonName,
			r.PersonEmail,
			r.ActivityType,
			r.ActivityFormName,
		}
		err := writer.Write(record)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
	return resp.Payload.Count, nil
}

func main() {
	var (
		app     = kingpin.New("activity-see", "List all activities")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("output", "CSVf file for results").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(*csvFile)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(f)
	offset := int32(0)
	count := int32(e.Metrics.MaxBatchSize)
	for count > 0 {
		count, err = process(e, writer, offset)
		if err != nil {
			panic(err)
		}
		offset += count
	}
}
