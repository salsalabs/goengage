package metrics

import (
	"fmt"

	"github.com/salsalabs/goengage/pkg"
)

const (
	//Command is used to retrieve runtime metrics.
	Command = "/api/integration/ext/v1/metrics"
	//MetricsMethod is the HTTP method used to retrieve metrics.
	MetricsMethod = "GET"
)

//MetricData contains the measurable stsuff in Engage.
type MetricData struct {
	RateLimit                      int32  `json:"rateLimit"`
	MaxBatchSize                   int32  `json:"maxBatchSize"`
	CurrentRateLimit               int32  `json:"currentRateLimit"`
	TotalAPICalls                  int32  `json:"totalAPICalls"`
	LastAPICall                    string `json:"lastAPICall"`
	TotalAPICallFailures           int32  `json:"totalAPICallFailures"`
	LastAPICallFailure             string `json:"lastAPICallFailure"`
	SupporterRead                  int32  `json:"supporterRead"`
	SupporterAdd                   int32  `json:"supporterAdd"`
	SupporterUpdate                int32  `json:"supporterUpdate"`
	SupporterDelete                int32  `json:"supporterDelete"`
	ActivityEvent                  int32  `json:"activityEvent"`
	ActivitySubscribe              int32  `json:"activitySubscribe"`
	ActivityFundraise              int32  `json:"activityFundraise"`
	ActivityTargetedLetter         int32  `json:"activityTargetedLetter"`
	ActivityPetition               int32  `json:"activityPetition"`
	ActivitySubscriptionManagement int32  `json:"activitySubscriptionManagement"`
}

//Metrics reads metrics and returns them.
func Metrics(e goengage.EngEnv) (*MetricData, error) {
	m := MetricData{}
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   MetricsMethod,
		Fragment: Command,
		Token:    e.Token,
		Response: &m,
	}
	fmt.Printf("NetOp is %+v\n", n)
	err := n.Do()
	return &m, err
}
