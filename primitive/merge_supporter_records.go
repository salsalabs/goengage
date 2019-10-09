package goengage

//Result value constants.
const (
	//Found will be reported for the destination supporter if no updates were
	//specified to be performed.
	Found = "FOUND"
	//Update will be reported for the destination supporter if updates were
	//specified. It will also be reported on the main payload if the merge
	//operation was successful.
	Update = "UPDATE"
	//NotFound will be reported for the destination or source supporter if the
	//provided id(s) do not exist.
	NotFound = "NOT_FOUND"
	//Deleted will be reported for the source supporter on a successful merge.
	Deleted = "DELETED"
	//ValidationError will be reported on the main payload if either the source
	//or the destination supporter is not found, or a request to update the
	//destination was specified and validation errors occurred during that
	//update.
	ValidationError = "VALIDATION_ERROR"
	//SystemError if the merge could not be completed.
	SystemError = "SYSTEM_ERROR"
)

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
