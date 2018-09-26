package main

import (
	//"encoding/json"
	"fmt"

	"github.com/salsalabs/goengage"
)

const token = ``

func main() {
	rqt := goengage.SupSearchRequest{
		ModifiedFrom: "2016-09-01T00:00:00.000Z",
		ModifiedTo:   "2019-09-01T00:00:00.000Z",
		Offset:       0,
		Count:        20,
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     goengage.ProdHost,
		Fragment: goengage.SupSearch,
		Token:    token,
		Request:  &rqt,
		Response: &resp,
	}
	count := int32(rqt.Count)
	for count > 0 {
		err := n.Search()
		if err != nil {
			panic(err)
		}
		count = int32(len(resp.Payload.Supporters))
		fmt.Printf("Read %d supporters from offset %d\n", count, rqt.Offset)
		rqt.Offset = rqt.Offset + count
		for _, s := range resp.Payload.Supporters {
			fmt.Printf("%-20s %-20s\n", s.FirstName, s.LastName)
		}
	}
}
