package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/salsalabs/goengage"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type Input struct {
	InternalID     string
	ExternalID     string
	Email          string
	CorrectedEmail string
}

func NewInput(a []string) Input {
	i := Input{
		InternalID:     a[0],
		ExternalID:     a[1],
		Email:          a[2],
		CorrectedEmail: a[3],
	}
	return i
}

type InputMap map[string]Input

func main() {
	var (
		app     = kingpin.New("activity-search", "A command-line app to see emails for a list of supporter IDs.")
		login   = app.Flag("login", "YAML file with API token").Required().String()
		csvFile = app.Flag("csv", "CSV file with IDs.  Uses 'InternalID'.").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(*csvFile)
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)
	//records is an array of records.  Each record is
	//an array of strings with these offsets.
	//0 InternalID
	//1 ExternalID
	//2 Email
	//3 Corrected Email
	a, err := r.ReadAll()
	_ = f.Close()

    var records []Input
	inputMap := make(InputMap)
	for _, r := range a {
        if r[0] != "InternalID" {
            i := NewInput(r)
		    inputMap[i.Email] = i
            records = append(records, i)
        }
	}

	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}
	rqt := goengage.SupSearchIDRequest{
		IdentifierType: "EXTERNAL_ID",
		Offset:         0,
		Count:          m.MaxBatchSize,
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SupSearch,
		Method:   http.MethodPost,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}

	for rqt.Count > 0 {
		rOffset := int32(rqt.Offset)
		var identifiers []string
		for i := int32(0); i < m.MaxBatchSize && rOffset + i < int32(len(records)); i++ {
			x := rOffset + i
			identifiers = append(identifiers, records[x].ExternalID)
		}
        rqt.Count = int32(math.Min(float64(len(identifiers)), float64(m.MaxBatchSize)))
        if rqt.Count != 0 {
            rqt.Identifiers = identifiers
            fmt.Printf("Searching from offset %d for %d IDs\n", rqt.Offset, len(rqt.Identifiers))
            err := n.Do()
            if err != nil {
                panic(err)
            }

            count := int32(len(resp.Payload.Supporters))
            fmt.Printf("Read %d supporters from offset %d\n", count, rqt.Offset)
            for _, s := range resp.Payload.Supporters {
                e := goengage.FirstEmail(s)
                email := ""
                if e != nil && len(*e) > 0 {
                    email = *e
                }
                m, ok := inputMap[email]
                if !ok {
                    fmt.Printf("Warning: unable to find email %v in the map\n", e)
                } else {
                    fmt.Printf("%-20s %-30s -> %s\n", s.ExternalSystemID, email, m.CorrectedEmail)
                }
            }
        }
		rqt.Offset = rqt.Offset + rqt.Count
	}
}
