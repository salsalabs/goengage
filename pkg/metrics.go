package goengage

import "time"

//MetricsCommand is used to retrieve runtime metrics.
const MetricsCommand = "/api/integration/ext/v1/metrics"

//There metrics command does not require request JSON.

//MetricsResponse wraps the results of the metrics call.
type MetricsResponse struct {
	ID        string    `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	Header    Header    `json:"header"`
	Payload   Metrics   `json:"payload"`
}

//Metrics data returned via the Metrics call.
//See https://help.salsalabs.com/hc/en-us/articles/224531208-General-Use#understanding-available-calls-remaining
type Metrics struct {
	RateLimit                      int32     `json:"rateLimit"`
	MaxBatchSize                   int32     `json:"maxBatchSize"`
	SupporterRead                  int32     `json:"supporterRead"`
	SupporterAdd                   int32     `json:"supporterAdd"`
	SupporterDelete                int32     `json:"supporterDelete"`
	SupporterUpdate                int32     `json:"supporterUpdate"`
	SegmentRead                    int32     `json:"segmentRead"`
	SegmentAdd                     int32     `json:"segmentAdd"`
	SegmentDelete                  int32     `json:"segmentDelete"`
	SegmentUpdate                  int32     `json:"segmentUpdate"`
	SegmentAssignmentRead          int32     `json:"segmentAssignmentRead"`
	SegmentAssignmentAdd           int32     `json:"segmentAssignmentAdd"`
	SegmentAssignmentUpdate        int32     `json:"segmentAssignmentUpdate"`
	SegmentAssignmentDelete        int32     `json:"segmentAssignmentDelete"`
	OfflineDonationAdd             int32     `json:"offlineDonationAdd"`
	OfflineDonationUpdate          int32     `json:"offlineDonationUpdate"`
	ActivityTicketedEvent          int32     `json:"activityTicketedEvent"`
	ActivityP2PEvent               int32     `json:"activityP2PEvent"`
	ActivitySubscribe              int32     `json:"activitySubscribe"`
	ActivityFundraise              int32     `json:"activityFundraise"`
	ActivityTargetedLetter         int32     `json:"activityTargetedLetter"`
	ActivityPetition               int32     `json:"activityPetition"`
	ActivitySubscriptionManagement int32     `json:"activitySubscriptionManagement"`
	LastAPICall                    *time.Time `json:"lastAPICall"`
	TotalAPICalls                  int32     `json:"totalAPICalls"`
	TotalAPICallFailures           int32     `json:"totalAPICallFailures"`
	CurrentRateLimit               int32     `json:"currentRateLimit"`
}
