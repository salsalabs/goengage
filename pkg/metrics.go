package goengage

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

//FragMetrics is used to retrieve runtime metrics.
const FragMetrics = "/api/integration/ext/v1/metrics"

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

//MetResponse is returned by Engage when asking for metrics.
type MetResponse struct {
	ID        string
	Timestamp string
	Header    struct {
		ProcessingTime int32  `json:"processingTime"`
		ServerID       string `jsin:"serverId"`
	}
	Payload MetricData
}

func (e EngEnv) get(method string, fragment string) ([]byte, error) {
	u, _ := url.Parse(fragment)
	u.Scheme = "https"
	u.Host = e.Host
	client := &http.Client{}
	req, _ := http.NewRequest(method, u.String(), nil)
	req.Header.Set("authToken", e.Token)
	var body []byte
	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return body, err
}

//Metrics reads metrics and returns them.
func (e EngEnv) Metrics() (*MetricData, error) {
	body, err := e.get(http.MethodGet, FragMetrics)
	if err != nil {
		return nil, err
	}
	var m MetResponse
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	return &m.Payload, err
}