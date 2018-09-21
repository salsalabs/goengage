package main

import (
	"fmt"

	"github.com/salsalabs/goengage"
)

const token = "wBTvk4rH5auTh4up8nOaVCcJBYWT3jr2Wk7QnlcOc4Qa7dvkgaDBGK6pP3hUaneP_aw0vGveE3XqDEfXSBIsQy7slH24kQ_SZVlojNYkNrg"

func main() {
	e := goengage.EngEnv{
		Host:  "hq.uat.igniteaction.net",
		Token: token}
	fmt.Printf("EngEnv is %+v\n", e)
	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Metrics: %+v\n", m)
}
