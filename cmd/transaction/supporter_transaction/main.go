package main

//Application to create a CSV of transactions for a supporter. You
//provide credentials and a supporter_KEY.  This app writes a CSV of
//the supporter's transactions.

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Runtime holds the stuff that this app needs.
type Runtime struct {
	Env           *goengage.Environment
	IncludeCounts bool
	CSVFilename   string
	Logger        *goengage.UtilLogger
	SupporterID   string
}

func NewRuntime(e *goengage.Environment, f string, s string, v bool) (*Runtime, error) {
	rt := Runtime{
		Env:         e,
		CSVFilename: f,
		SupporterID: s,
	}
	if v {
		logger, err := goengage.NewUtilLogger()
		if err != nil {
			return nil, err
		}
		rt.Logger = logger
	}
	return &rt, nil
}

//Run finds and displays all transactions.
func Run(rt *Runtime) error {
	log.Println("Run: begin")
	f, err := os.Create(rt.CSVFilename)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	headers := []string{
		"CreatedDate",
		"TransactionID",
		"ActivityName",
		"SupporterID",
		"Amount",
		"TemplateID",
		"TransactionDate",
		"TransactionType",
		"DeductibleAmount",
		"FeesPaid",
		"Result",
	}
	err = writer.Write(headers)
	if err != nil {
		return err
	}

	count := rt.Env.Metrics.MaxBatchSize
	offset := int32(0)

	for count == rt.Env.Metrics.MaxBatchSize {
		payload := goengage.TransactionSearchRequestPayload{
			Identifiers:    []string{rt.SupporterID},
			IdentifierType: goengage.SupporterIDType,
			Count:          count,
			Offset:         offset,
		}
		fmt.Printf("Request payload: %+v\n", payload)
		rqt := goengage.TransactionSearchRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}

		var resp goengage.TransactionSearchResponse

		n := goengage.NetOp{
			Host:     rt.Env.Host,
			Method:   goengage.SearchMethod,
			Endpoint: goengage.SearchTransactionDetails,
			Token:    rt.Env.Token,
			Request:  &rqt,
			Response: &resp,
			Logger:   rt.Logger,
		}
		err := n.Do()
		if err != nil {
			return err
		}
		if offset%100 == 0 {
			log.Printf("Run: %6d: %2d of %6d\n",
				offset,
				len(resp.Payload.Transactions),
				resp.Payload.Total)
		}

		var cache [][]string
		for _, w := range resp.Payload.Transactions {
			s := w.DonationTransaction
			record := []string{
				s.CreatedDate,
				s.TransactionID,
				s.ActivityName,
				s.SupporterID,
				fmt.Sprintf("%v", s.Amount),
				s.TemplateID,
				s.TransactionDate,
				s.TransactionType,
				fmt.Sprintf("%v", s.DeductibleAmount),
				fmt.Sprintf("%v", s.FeesPaid),
				s.Result,
			}
			cache = append(cache, record)
		}
		err = writer.WriteAll(cache)
		if err != nil {
			return err
		}
		count = resp.Payload.Count
		offset += int32(count)
	}
	log.Printf("Run: end")
	return nil
}

//Program entry point.
func main() {
	var (
		app         = kingpin.New("supporter_transactions", "Creates a CSV of transactions for a supporter")
		login       = app.Flag("login", "YAML file with API token").Required().String()
		csvFile     = app.Flag("csv", "CSV filename to create").Default("supporter_transactions.csv").String()
		supporterID = app.Flag("supporterID", "Find transactions for this supporter").Required().String()
		verbose     = app.Flag("verbose", "Log all requests and responses to a file.  Verrrry noisy...").Bool()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		log.Fatalf("Error --login is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		log.Fatalf("main: %+v\n", e)
	}
	rt, err := NewRuntime(e, *csvFile, *supporterID, *verbose)
	if err != nil {
		log.Fatalf("main: %v\n", err)
	}
	err = Run(rt)
	if err != nil {
		log.Fatalf("main: %v\n", err)
	}
}
