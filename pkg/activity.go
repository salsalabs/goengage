package goengage

import "time"

//ActivityRequest is used to retrieve activities from Engage.
//Note that ActivityRequest can be used to retrieve activities based
//on three types of criteria: activity IDs, activity form IDs, modified
//date range.  Choose one and provide the necessary data.  The remainder
//will be ignored when the request is sent to Engage.
type ActivityRequest struct {
	Header  RequestHeader          `json:"header"`
	Payload ActivityRequestPayload `json:"payload"`
}

//ActivityRequestPayload specifies the activities to return.
type ActivityRequestPayload struct {
	Type            string   `json:"type,omitempty"`
	Offset          int32    `json:"offset"`
	Count           int32    `json:"count"`
	ActivityIDs     []string `json:"activityIds,omitempty"`
	ActivityFormIDs []string `json:"activityFormIds,omitempty"`
	ModifiedFrom    string   `json:"modifiedFrom,omitempty"`
	ModifiedTo      string   `json:"modifiedTo,omitempty"`
}

//BaseActivity returns activity information from SUBSCRIBE or
//SUBSCRIPTION_MANAGEMENT requests.  Note that Base is actually
//contained in the other activity result objects.
type BaseActivity struct {
	ActivityType     string     `json:"activityType,omitempty"`
	ActivityID       string     `json:"activityId,omitempty"`
	ActivityFormName string     `json:"activityFormName,omitempty"`
	ActivityFormID   string     `json:"activityFormId,omitempty"`
	SupporterID      string     `json:"supporterId,omitempty"`
	ActivityDate     *time.Time `json:"activityDate,omitempty"`
	LastModified     *time.Time `json:"lastModified,omitempty"`
}

//BaseResponse is the set of common fields returned for all activities.
//Some activities (like SUBSCRIBE or SUBSCRIPTION_MANAGEMENT) only return
//ActivityBase.  Other activities, like donations, events and P2P, return
//data appended to the base.
type BaseResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Total      int32          `json:"total,omitempty"`
		Offset     int32          `json:"offset,omitempty"`
		Count      int32          `json:"count,omitempty"`
		Activities []BaseActivity `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}
