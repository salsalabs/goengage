package goengage

import "time"

//MetricsCommand is used to retrieve runtime metrics.
const MetricsCommand = "/api/integration/ext/v1/metrics"

//There metrics command does not require request JSON.

//MetricsResponse wraps the results of the metrics call.
type MetricsResponse struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Header    Header    `json:"header"`
	Payload   Metrics   `json:"payload"`
}

//Metrics data returned via the Metrics call.
//See https://help.salsalabs.com/hc/en-us/articles/224531208-General-Use#understanding-available-calls-remaining
type Metrics struct {
	RateLimit                      int       `json:"rateLimit"`
	MaxBatchSize                   int       `json:"maxBatchSize"`
	SupporterRead                  int       `json:"supporterRead"`
	SupporterAdd                   int       `json:"supporterAdd"`
	SupporterDelete                int       `json:"supporterDelete"`
	SupporterUpdate                int       `json:"supporterUpdate"`
	SegmentRead                    int       `json:"segmentRead"`
	SegmentAdd                     int       `json:"segmentAdd"`
	SegmentDelete                  int       `json:"segmentDelete"`
	SegmentUpdate                  int       `json:"segmentUpdate"`
	SegmentAssignmentRead          int       `json:"segmentAssignmentRead"`
	SegmentAssignmentAdd           int       `json:"segmentAssignmentAdd"`
	SegmentAssignmentUpdate        int       `json:"segmentAssignmentUpdate"`
	SegmentAssignmentDelete        int       `json:"segmentAssignmentDelete"`
	OfflineDonationAdd             int       `json:"offlineDonationAdd"`
	OfflineDonationUpdate          int       `json:"offlineDonationUpdate"`
	ActivityTicketedEvent          int       `json:"activityTicketedEvent"`
	ActivityP2PEvent               int       `json:"activityP2PEvent"`
	ActivitySubscribe              int       `json:"activitySubscribe"`
	ActivityFundraise              int       `json:"activityFundraise"`
	ActivityTargetedLetter         int       `json:"activityTargetedLetter"`
	ActivityPetition               int       `json:"activityPetition"`
	ActivitySubscriptionManagement int       `json:"activitySubscriptionManagement"`
	LastAPICall                    time.Time `json:"lastAPICall"`
	TotalAPICalls                  int       `json:"totalAPICalls"`
	TotalAPICallFailures           int       `json:"totalAPICallFailures"`
	CurrentRateLimit               int       `json:"currentRateLimit"`
}
