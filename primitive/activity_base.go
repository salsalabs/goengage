package goengage

import "time"

//BaseResponse is the set of common fields returned for all activities.
//Some activities (like SUBSCRIBE or SUBSCRIPTION_MANAGEMENT) only return
//ActivityBase.  Other activities, like donations, events and P2P, return
//data appended to the base.
type BaseResponse struct {
	Header  Header              `json:"header,omitempty"`
	Payload BaseResponsePayload `json:"payload,omitempty"`
}

//BaseActivity returns activity information from SUBSCRIBE or
//SUBSCRIPTION_MANAGEMENT requests.  Note that BaseActivity is actually
//contained in the other activity result objects.
type BaseActivity struct {
	ActivityType     string    `json:"activityType,omitempty"`
	ActivityID       string    `json:"activityId,omitempty"`
	ActivityFormName string    `json:"activityFormName,omitempty"`
	ActivityFormID   string    `json:"activityFormId,omitempty"`
	SupporterID      string    `json:"supporterId,omitempty"`
	ActivityDate     time.Time `json:"activityDate,omitempty"`
	LastModified     time.Time `json:"lastModified,omitempty"`
}

//BaseResponsePayload contains the data returned by a SUBSCRIBE or
//SUBSCRIPTION_MANAGEMENT requests.
type BaseResponsePayload struct {
	Total      int            `json:"total,omitempty"`
	Offset     int            `json:"offset,omitempty"`
	Count      int            `json:"count,omitempty"`
	Activities []BaseActivity `json:"activities,omitempty"`
}
