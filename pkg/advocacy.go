package goengage

// Petition contains information about a petition being signed.
// Note that PetitionActivity starts with the contents of BaseActivity...
type Petition struct {
	BaseActivity
	Comment                  string `json:"comment,omitempty"`
	ModerationState          string `json:"moderationState,omitempty"`
	DisplaySignaturePublicly bool   `json:"displaySignaturePublicly,omitempty"`
	DisplayCommentPublicly   bool   `json:"displayCommentPublicly,omitempty"`
}

// Target describes the recipient of a targeted letter or call.
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

// Letter contains the contents sent by a supporters to one or more Targets.
type Letter struct {
	Name               string   `json:"name,omitempty"`
	Subject            string   `json:"subject,omitempty"`
	Message            string   `json:"message,omitempty"`
	AdditionalComment  string   `json:"additionalComment,omitempty"`
	SubjectWasModified bool     `json:"subjectWasModified,omitempty"`
	MessageWasModified bool     `json:"messageWasModified,omitempty"`
	Targets            []Target `json:"targets,omitempty"`
}

// TargetedLetter describes the action taken for a targeted letter.
type TargetedLetter struct {
	BaseActivity
	Letters []Letter `json:"letters,omitempty"`
}

// PetitionResponse is returned when the request type is "PETITION".
type PetitionResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Total      int32      `json:"total,omitempty"`
		Offset     int32      `json:"offset,omitempty"`
		Count      int32      `json:"count,omitempty"`
		Activities []Petition `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}

// TargetedLetterResponse is returned when the request is "TARGETED_LETTERS".
type TargetedLetterResponse struct {
	Header  Header `json:"header,omitempty"`
	Payload struct {
		Total      int32            `json:"total,omitempty"`
		Offset     int32            `json:"offset,omitempty"`
		Count      int32            `json:"count,omitempty"`
		Activities []TargetedLetter `json:"activities,omitempty"`
	} `json:"payload,omitempty"`
}
