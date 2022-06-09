package drills

import (
	"../../dao/drills"
)

// FullDrill is named as such as this object contains the combination of drill data from the session and drills
// table
type FullDrill struct {
	DrillId     string             `json:"drillId"`
	CoachId     string             `json:"coachId"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Category    string             `json:"category"` // DRIBBLING/WARMUP/SHOOTING/PASSING
	Equipment   []drills.Equipment `json:"equipment"`

	Demos Demos `json:"demos"`
}

type Demos struct {
	Front Demo `json:"front"`
	Side  Demo `json:"side"`
	Close Demo `json:"close"`
}

type Demo struct {
	FileLocation string `json:"fileLocation"`
}
