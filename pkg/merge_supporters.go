package goengage

// Destination describes the target of the mail merge.  Note that "Result" is
// only provided in the response.
type Destination struct {
	ReadOnly    bool   `json:"readOnly"`
	SupporterID string `json:"supporterId"`
	Result      string `json:"result,omitempty"`
}

// Source describes the source of a mail merge.  This record's contents will
// be merged into the destination.  If the merge is successful, then the
// source record is removed.  That will appear in the Result field.  Note that
// the result field is only provided in the response.
type Source struct {
	SupporterID string `json:"supporterId"`
	Result      string `json:"result,omitempty"`
}

// MergeSupporterRecordsRequest specifies both source and destination supporters to be merged.
type MergeSupporterRecordsRequest struct {
	Header  RequestHeader `json:"header,omitempty"`
	Payload struct {
		Destination Destination `json:"destination"`
		Source      Source      `json:"source"`
	} `json:"payload"`
}

// MergeSupporterRecordsResponse shows the results of the merge request.
type MergeSupporterRecordsResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Destination Destination `json:"destination"`
		Source      Source      `json:"source"`
		Result      string      `json:"result"`
	} `json:"payload"`
}
