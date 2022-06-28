package session

import (
	"impruviService/dao/drills"
)

type FullSession struct {
	PlayerId      string       `json:"playerId"`
	SessionNumber int          `json:"sessionNumber"`
	Drills        []*FullDrill `json:"drills"`
}

// FullDrill is named as such as this object contains the combination of drill data from the session and drills
// table
type FullDrill struct {
	DrillId     string             `json:"drillId"`
	CoachId     string             `json:"coachId"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Category    string             `json:"category"` // DRIBBLING/WARMUP/SHOOTING/PASSING
	Equipment   []drills.Equipment `json:"equipment"`

	Submission               Media  `json:"submission"`
	Feedback                 Media  `json:"feedback"`
	Notes                    string `json:"notes"`
	EstimatedDurationMinutes int    `json:"estimatedDurationMinutes"`

	Demos Demos `json:"demos"`
}

type Demos struct {
	Front          Media `json:"front"`
	Side           Media `json:"side"`
	Close          Media `json:"close"`
	FrontThumbnail Media `json:"frontThumbnail"`
	SideThumbnail  Media `json:"sideThumbnail"`
	CloseThumbnail Media `json:"closeThumbnail"`
}

type Media struct {
	VideoUploadDateEpochMillis int64  `json:"videoUploadDateEpochMillis"` // only used in submission and feedback videos
	FileLocation               string `json:"fileLocation"`
}
