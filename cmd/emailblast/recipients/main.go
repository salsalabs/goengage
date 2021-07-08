package main

// Quick and dirty application to read all email blasts and write blast
// and recipient activity to CSVs.  Need this done today.  No performance
// tricks -- just getting the job done.
import (
	"encoding/csv"
	"log"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	"gopkg.in/alecthomas/kingpin.v2"
)

//Runtime is the internal data store for this app.
type Runtime struct {
	Env             *goengage.Environment
	PublishedFrom   string
	RecipientsFile  *csv.Writer
	ConversionsFile *csv.Writer
}

//Blasts reads email blasts and passes blasts off to the
//blast detail reader.
func (rt *Runtime) Blasts() error {
	count := rt.Env.Metrics.MaxBatchSize
	offset := int32(0)
	for count == rt.Env.Metrics.MaxBatchSize {
		payload := goengage.EmailBlastSearchRequestPayload{
			PublishedFrom: rt.PublishedFrom,
			Offset:        offset,
			Count:         count,
			Type:          goengage.Email,
		}
		rqt := goengage.EmailBlastSearchRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}
		var resp goengage.EmailBlastSearchResponse

		n := goengage.NetOp{
			Host:     rt.Env.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.EmailBlastSearch,
			Token:    rt.Env.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err := n.Do()
		if err != nil {
			return err
		}
		for _, s := range resp.Payload.EmailActivities {
			err = rt.OneBlast(s)
			if err != nil {
				return err
			}
		}
		offset = offset + count
		count = resp.Payload.Count
	}
	return nil
}

//Details reads the activity version of the blast and writes recipients
//and conversions.
func (rt *Runtime) OneBlast(r goengage.EmailActivity) error {
	log.Printf("OneBlast: %s %s\n", r.ID, r.Name)
	count := rt.Env.Metrics.MaxBatchSize
	offset := int32(0)
	cursor := ""
	// for count == rt.Env.Metrics.MaxBatchSize {
	log.Printf("OneBlast: ID: %s, cursor: %s, offset: %d, count: %d\n", id, cursor, offset, count)
	payload := goengage.IndivualBlastRequestPayload{
		ID:     id,
		Offset: offset,
		Count:  count,
		Type:   goengage.Email,
	}
	if len(cursor) > 0 {
		payload.Cursor = cursor
	}

	rqt := goengage.IndivualBlastRequest{
		Header:  goengage.RequestHeader{},
		Payload: payload,
	}
	var resp goengage.InvidualBlastResponse

	n := goengage.NetOp{
		Host:     rt.Env.Host,
		Method:   goengage.SearchMethod,
		Endpoint: goengage.IndividualBlastSearch,
		Token:    rt.Env.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err := n.Do()
	if err != nil {
		return err
	}

	log.Printf("OneBlast: response %+v\n", resp)
	for _, s := range resp.Payload.IndividualEmailActivityData {
		log.Printf("OneBlast: individual activity %+v\n", s)
		err = rt.Recipients(id, s.RecipientsData)
		if err != nil {
			return err
		}
		// cursor = s.Cursor
		// if len(cursor) == 0 {
		// 	log.Printf("OneBlast: cursor is empty\n")
		// 	count = 0
		// } else {
		// count = int32(len(s.RecipientsData.Recipients))
		// }
		offset = offset + count
	}
	offset = offset + count
	// }
	return nil
}

//Recipients formats and writes email and conversion activity.
func (rt *Runtime) Recipients(blastId string, r goengage.SingleBlastRecipientsData) error {
	log.Printf("Recipients: blast ID: %s, %d recipients", blastId, len(r.Recipients))
	return nil
}

//Program entry point.
func main() {
	var (
		app             = kingpin.New("recipients", "Write recipient and conversion data for blasts sent after a date")
		login           = app.Flag("login", "YAML file with API token").Required().String()
		recipientsFile  = app.Flag("recipients", "CSV filename to recipient info").Default("recipients.csv").String()
		conversionsFile = app.Flag("conversions", "CSV filename to store conversion info").Default("conversion.csv").String()
		publishedFrom   = app.Flag("published-from", "Engage-formatted start date").Default("2021-03-12T00:00:00.000Z").String()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		log.Fatalf("Error --login is required.")
		os.Exit(1)
	}
	if recipientsFile == nil || len(*recipientsFile) == 0 {
		log.Fatalf("Error --blast-csv is required.")
		os.Exit(1)
	}
	if conversionsFile == nil || len(*conversionsFile) == 0 {
		log.Fatalf("Error --csv is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("Error %v\n", err)
		os.Exit(1)
	}
	f1, err := os.Create(*recipientsFile)
	if err != nil {
		log.Fatalf("Error %v\n", err)
		os.Exit(1)
	}
	defer f1.Close()

	f2, err := os.Create(*conversionsFile)
	if err != nil {
		log.Fatalf("Error %v\n", err)
		os.Exit(1)
	}
	defer f2.Close()

	rtx := Runtime{
		Env:             e,
		PublishedFrom:   *publishedFrom,
		RecipientsFile:  csv.NewWriter(f1),
		ConversionsFile: csv.NewWriter(f2),
	}
	rt := &rtx

	rt.Blasts()

}
