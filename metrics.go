package goengage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//MetricData contains the measurable stsuff in Engage.
type MetricData struct {
	RateLimit                      int32  `json:"rateLimit"`
	MaxBatchSize                   int32  `json:"maxBatchSize"`
	CurrentRateLimit               int32  `json:"currentRateLimit"`
	SupporterRead                  int32  `json:"supporterRead"`
	TotalAPICalls                  int32  `json:"totalAPICalls"`
	LastAPICall                    string `json:"lastAPICall"`
	TotalAPICallFailures           int32  `json:"totalAPICallFailures"`
	LastAPICallFailure             string `json:"lastAPICallFailure"`
	SupporterReads1                int32  `json:"supporterRead"`
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

//MetricReturn is returned by Engage when asking for metrics.
type MetricReturn struct {
	ID        string
	Timestamp string
	Header    struct {
		ProcessingTime int32  `json:"processingTime"`
		ServerID       string `jsin:"serverId"`
	}
	Payload MetricData
}

//Measure reads metrics and returns them.
func (e EngEnv) Metrics() (*MetricData, error) {
	u, _ := url.Parse("/api/integration/ext/v1/metrics")
	x := fmt.Sprintf("https://%v", e.Host)
	b, _ := url.Parse(x)
	t := b.ResolveReference(u)
	fmt.Printf("Meterics URL is %v", t)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, t.String(), nil)
	req.Header.Set("authToken", e.Token)
	var body []byte
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Printf("body: %v\n", string(body))
	var m MetricReturn
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MetricReturn is: %+v\n", m)
	return &m.Payload, err
}
