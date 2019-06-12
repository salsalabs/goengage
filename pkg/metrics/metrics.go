package metrics

import (
	"encoding/json"
	"net/http"

	"github.com/salsalabs/goengage/pkg"
)

//Command is used to retrieve runtime metrics.
const Command = "/api/integration/ext/v1/metrics"

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

//response is returned by Engage when asking for metrics.
type response struct {
	ID        string
	Timestamp string
	Header    struct {
		ProcessingTime int32  `json:"processingTime"`
		ServerID       string `jsin:"serverId"`
	}
	Payload MetricData
}

//Metrics reads metrics and returns them.
func Metrics(e goengage.EngEnv) (*MetricData, error) {
	body, err := e.Get(http.MethodGet, Command)
	if err != nil {
		return nil, err
	}
	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r.Payload, err
}
