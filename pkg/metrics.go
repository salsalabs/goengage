package goengage

//MetricsCommand is used to retrieve runtime metrics.
const MetricsCommand = "/api/integration/ext/v1/metrics"

//MetricData contains the measurable stsuff in Engage.
type MetricData struct {
	RateLimit                      int32  `json:"rateLimit"`
	MaxBatchSize                   int32  `json:"maxBatchSize"`
	CurrentRateLimit               int32  `json:"currentRateLimit"`
	TotalAPICalls                  int32  `json:"totalAPICalls"`
	LastAPICall                    string `json:"lastAPICall"`
	TotalAPICallFailures           int32  `json:"totalAPICallFailures"`
	LastAPICallFailure             string `json:"lastAPICallFailure"`
	ActivityFundraise              int32  `json:"activityFundraise"`
	ActivityP2PEvent               int32  `json:"activityP2PEvent"`
	ActivityPetition               int32  `json:"activityPetition"`
	ActivitySubscribe              int32  `json:"activitySubscribe"`
	ActivitySubscriptionManagement int32  `json:"activitySubscriptionManagement"`
	ActivityTargetedLetter         int32  `json:"activityTargetedLetter"`
	ActivityTicketedEvent          int32  `json:"activityTicketedEvent"`
	OfflineDonationAdd             int32  `json:"offlineDonationAdd"`
	OfflineDonationUpdate          int32  `json:"offlineDonationUpdate"`
	SegmentAdd                     int32  `json:"segmentAdd"`
	SegmentDelete                  int32  `json:"segmentDelete"`
	SegmentRead                    int32  `json:"segmentRead"`
	SegmentUpdate                  int32  `json:"segmentUpdate"`
	SegmentAssignmentAdd           int32  `json:"segmentAssignmentAdd"`
	SegmentAssignmentDelete        int32  `json:"segmentAssignmentDelete"`
	SegmentAssignmentRead          int32  `json:"segmentAssignmentRead"`
	SegmentAssignmentUpdate        int32  `json:"segmentAssignmentUpdate"`
	SupporterAdd                   int32  `json:"supporterAdd"`
	SupporterDelete                int32  `json:"supporterDelete"`
	SupporterRead                  int32  `json:"supporterRead"`
	SupporterUpdate                int32  `json:"supporterUpdate"`
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
