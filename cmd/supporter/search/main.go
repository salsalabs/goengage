package main

import (
	"encoding/json"
	"fmt"

	"github.com/salsalabs/goengage"
)

const token = `wBTvk4rH5auTh4up8nOaVCcJBYWT3jr2Wk7QnlcOc4RlzTBx1sFmcTTI5go4M-lg_Jyh97x--zg4FwCCXx7Cmhnc_hRaAo_mk5pOloQtiOM`

func main() {
	rqt := goengage.SupSearchRequest{
		ModifiedFrom: "2016-09-01T00:00:00.000Z",
		ModifiedTo:   "2019-09-01T00:00:00.000Z",
		Offset:       0,
		Count:        20,
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     goengage.UatHost,
		Fragment: goengage.SupSearch,
		Token:    token,
		Request:  rqt,
		Response: &resp,
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
