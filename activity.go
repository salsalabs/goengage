package goengage

//ActSearch is used to search for activities.
const ActSearch = "/api/integration/ext/v1/activities/search"

//Activity is a generic action that someone takes in Engage.
type Activity struct {
	ActivityType     string
	ActivityID       string
	ActivityFormName string
	ActivityFormID   string
	ActivityDate     string
	LastModified     string
}

//ActSearchRequest is used to ask for supporters.
type ActSearchRequest struct {
	ActivityIDs  []string `json:"activityIds"`
	ModifiedFrom string   `json:"modifiedFrom"`
	ModifiedTo   string   `json:"modifiedTo"`
	Offset       int32    `json:"offset"`
	Count        int32    `json:"count"`
	Type         string   `json:"type"`
}

//ActSearchResult is returned when supporters are found by a search.
type ActSearchResult struct {
	Payload struct {
		Count         int32         `json:"count"`
		Offset        int32         `json:"offset"`
		Total         int32         `json:"total"`
		SupActivities []SupActivity `json:"Activities"`
	} `json:"payload"`
}

//SupActivity shows when an activity occurred for a supporter.
type SupActivity struct {
	ActivityID       string
	ActivityFormName string
	ActivityFormID   string
	SupporterID      string
	ActivityDate     string
	ActivityType     string
	LastModified     string
	//CustomFieldValues []something
}
