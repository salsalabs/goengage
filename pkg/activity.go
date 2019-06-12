package goengage

//ActSearch is used to search for activities.
const ActSearch = "/api/integration/ext/v1/activities/search"

//ActMethod is the method needed to do activity searches
const ActMethod = "POST"

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
	ActivityIDs  []string `json:"activityIds,omitempty"`
	ModifiedFrom string   `json:"modifiedFrom,omitempty"`
	ModifiedTo   string   `json:"modifiedTo,omitempty"`
	Offset       int32    `json:"offset,omitempty"`
	Count        int32    `json:"count,omitempty"`
	Type         string   `json:"type,omitempty"`
}

//ActSearchResult is returned when supporters are found by a search.
type ActSearchResult struct {
	Payload struct {
		Count         int32         `json:"count,omitempty"`
		Offset        int32         `json:"offset,omitempty"`
		Total         int32         `json:"total,omitempty"`
		SupActivities []SupActivity `json:"Activities,omitempty"`
	} `json:"payload,omitempty"`
}

//SupActivity shows when an activity occurred for a supporter.
type SupActivity struct {
	ActivityID       string `json:"activityID,omitempty"`
	ActivityFormName string `json:"activityFormName,omitempty"`
	ActivityFormID   string `json:"activityFormID,omitempty"`
	SupporterID      string `json:"supporterID,omitempty"`
	ActivityDate     string `json:"activityDate,omitempty"`
	ActivityType     string `json:"activityType,omitempty"`
	LastModified     string `json:"lastModified,omitempty"`
	//CustomFieldValues []something
}
