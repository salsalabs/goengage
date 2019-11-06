package goengage

import (
	"time"

	goengage "github.com/salsalabs/goengage/pkg"
)

//PetitionResponse is returned when the request type is "PETITION".
type PetitionResponse struct {
	Header  goengage.Header         `json:"header,omitempty"`
	Payload PetitionResponsePayload `json:"payload,omitempty"`
}

//Petition contains information about a petition being signed.
//Note that PetitionActivity starts with the contents of BaseActivity...
type Petition struct {
	ActivityID               string    `json:"activityId,omitempty"`
	ActivityFormName         string    `json:"activityFormName,omitempty"`
	ActivityFormID           string    `json:"activityFormId,omitempty"`
	SupporterID              string    `json:"supporterId,omitempty"`
	ActivityDate             time.Time `json:"activityDate,omitempty"`
	ActivityType             string    `json:"activityType,omitempty"`
	LastModified             time.Time `json:"lastModified,omitempty"`
	Comment                  string    `json:"comment,omitempty"`
	ModerationState          string    `json:"moderationState,omitempty"`
	DisplaySignaturePublicly bool      `json:"displaySignaturePublicly,omitempty"`
	DisplayCommentPublicly   bool      `json:"displayCommentPublicly,omitempty"`
}

//PetitionResponsePayload contains the data returned for a PETITION search.
type PetitionResponsePayload struct {
	Total      int        `json:"total,omitempty"`
	Offset     int        `json:"offset,omitempty"`
	Count      int32        `json:"count,omitempty"`
	Activities []Petition `json:"activities,omitempty"`
}

//TargetedLetterResponse is returned when the request is "TARGETED_LETTERS".
type TargetedLetterResponse struct {
	Header  goengage.Header               `json:"header,omitempty"`
	Payload TargetedLetterResponsePayload `json:"payload,omitempty"`
}

//Target describes the recipient of a targeted letter or call.
type Target struct {
	TargetID            string `json:"targetId,omitempty"`
	TargetName          string `json:"targetName,omitempty"`
	TargetTitle         string `json:"targetTitle,omitempty"`
	PoliticalParty      string `json:"politicalParty,omitempty"`
	TargetType          string `json:"targetType,omitempty"`
	State               string `json:"state,omitempty"`
	DistrictID          string `json:"districtId,omitempty"`
	DistrictName        string `json:"districtName,omitempty"`
	Role                string `json:"role,omitempty"`
	SentEmail           bool   `json:"sentEmail,omitempty"`
	SentFacebook        bool   `json:"sentFacebook,omitempty"`
	SentTwitter         bool   `json:"sentTwitter,omitempty"`
	MadeCall            bool   `json:"madeCall,omitempty"`
	CallDurationSeconds int64  `json:"callDurationSeconds,omitempty"`
	CallResult          string `json:"callResult,omitempty"`
}

//Letter contains the contents sent by a supporters to one or more Targets.
type Letter struct {
	Name               string   `json:"name,omitempty"`
	Subject            string   `json:"subject,omitempty"`
	Message            string   `json:"message,omitempty"`
	AdditionalComment  string   `json:"additionalComment,omitempty"`
	SubjectWasModified bool     `json:"subjectWasModified,omitempty"`
	MessageWasModified bool     `json:"messageWasModified,omitempty"`
	Targets            []Target `json:"targets,omitempty"`
}

//TargetedLetter describes the action taken for a targeted letter.
type TargetedLetter struct {
	ActivityID       string    `json:"activityId,omitempty"`
	ActivityFormName string    `json:"activityFormName,omitempty"`
	ActivityFormID   string    `json:"activityFormId,omitempty"`
	SupporterID      string    `json:"supporterId,omitempty"`
	ActivityDate     time.Time `json:"activityDate,omitempty"`
	ActivityType     string    `json:"activityType,omitempty"`
	LastModified     time.Time `json:"lastModified,omitempty"`
	Letters          []Letter  `json:"letters,omitempty"`
}

//TargetedLetterResponsePayload  contains the data returned for a TARGETED_LETTER search.
type TargetedLetterResponsePayload struct {
	Total      int              `json:"total,omitempty"`
	Offset     int              `json:"offset,omitempty"`
	Count      int32              `json:"count,omitempty"`
	Activities []TargetedLetter `json:"activities,omitempty"`
}
