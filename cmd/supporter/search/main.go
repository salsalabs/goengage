package main

import (
	"encoding/json"
	"fmt"

	"github.com/salsalabs/goengage"
)

const token = "wBTvk4rH5auTh4up8nOaVCcJBYWT3jr2Wk7QnlcOc4Qa7dvkgaDBGK6pP3hUaneP_aw0vGveE3XqDEfXSBIsQy7slH24kQ_SZVlojNYkNrg"

func main() {
	rqt := goengage.SupSearchRequest{
		ModifiedFrom: "2010-09-01T00:00:00.00Z",
		ModifiedTo:   "2010-09-01T00:00:00.00Z",
		Offset:       0,
		Count:        20,
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     goengage.UatHost,
		Fragment: goengage.SupSearch,
		Token:    token,
		Request:  rqt,
		Response: resp,
	}
	err := n.Search()
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response\n%v\n", string(b))

}
