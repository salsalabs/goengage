package main

// Search for all supporters and whack 'em.  NO SUPPORTERS WILL SURVIVE.
import (
	//"encoding/json"
	"fmt"
	"os"

	"github.com/salsalabs/goengage"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("delete-supporters", "A command-line app to DELETE ENGAGE SUPPORTERS.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		yes   = app.Flag("yes", "Yes, I understand that this program will utterly and completely remove Engage supporters.").Required().Bool()
	)
	app.Parse(os.Args[1:])
	if !*yes {
		fmt.Printf("You made a good choice.  Supporters won't be deleted.\n")
		return
	} else {
		fmt.Println("***")
		fmt.Printf("*** Alrighty, then.  You supplied --yes, supporters will be deleted.\n")
		fmt.Println("***")
	}

	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	rqt := goengage.SupSearchRequest{
		ModifiedFrom: "2016-09-01T00:00:00.000Z",
		ModifiedTo:   "2019-09-01T00:00:00.000Z",
		Offset:       0,
		Count:        20,
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SupSearch,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}

	dRqt := goengage.SupDeleteRequest{}
	dResp := goengage.SupDeleteResult{}
	nDel := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SupDelete,
		Token:    e.Token,
		Request:  &dRqt,
		Response: &dResp,
	}

	maxCount := int32(100)
	count := int32(rqt.Count)
	for count < maxCount && count > 0 {
		fmt.Printf("Deleting from offset %d, %d remain.\n", rqt.Offset, maxCount)
		err = n.Search()
		if err != nil {
			panic(err)
		}
		count = int32(len(resp.Payload.Supporters))
		maxCount = maxCount - count
		fmt.Printf("Deleting %d supporters from offset %d\n", count, rqt.Offset)
		rqt.Offset = rqt.Offset + count

		var a goengage.DeletingSupporters
		for _, x := range resp.Payload.Supporters {
			a = append(a, x.SupporterID)
		}
		dRqt.Supporters = a
		err = nDel.Delete()
		if err != nil {
			panic(err)
		}

		for _, s := range resp.Payload.Supporters {
			fmt.Printf("%s %s\n", s.SupporterID, s.Result)
		}
	}
}
