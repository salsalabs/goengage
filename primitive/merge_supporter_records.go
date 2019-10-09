package goengage

//MergeSupporterRecordsRequest specifies both source and destination supporters to be merged.
type MergeSupporterRecordsRequest struct {
	Payload struct {
		Destination struct {
			ReadOnly    bool   `json:"readOnly"`
			SupporterID string `json:"supporterId"`
		} `json:"destination"`
		Source struct {
			SupporterID string `json:"supporterId"`
		} `json:"source"`
	} `json:"payload"`
}

//MergeSupporterRecordsResponse shows the results of the merge request.
type MergeSupporterRecordsResponse struct {
	Payload struct {
		Destination struct {
			ReadOnly    bool   `json:"readOnly"`
			SupporterID string `json:"supporterId"`
			Result      string `json:"result"`
		} `json:"destination"`
		Source struct {
			SupporterID string `json:"supporterId"`
			Result      string `json:"result"`
		} `json:"source"`
		Result string `json:"result"`
	} `json:"payload"`
}
