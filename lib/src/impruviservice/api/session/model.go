package session

import (
	"../../dao/drills"
	"../../model"
)

type FullSession struct {
	PlayerId      string       `json:"playerId"`
	Name          string       `json:"name"`
	Date          *model.Date  `json:"date"`
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

	Submission               model.Media `json:"submission"`
	Feedback                 model.Media `json:"feedback"`
	Notes                    string      `json:"notes"`
	EstimatedDurationMinutes int         `json:"estimatedDurationMinutes"`

	Demos Demos `json:"demos"`
}

type Demos struct {
	Front          model.Media `json:"front"`
	Side           model.Media `json:"side"`
	Close          model.Media `json:"close"`
	FrontThumbnail model.Media `json:"frontThumbnail"`
	SideThumbnail  model.Media `json:"sideThumbnail"`
	CloseThumbnail model.Media `json:"closeThumbnail"`
}
