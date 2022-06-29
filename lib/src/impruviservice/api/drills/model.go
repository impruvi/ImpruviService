package drills

import (
	"impruviService/dao/drills"
<<<<<<< HEAD
	"impruviService/model"
=======
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
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
	Front          model.Media `json:"front"`
	Side           model.Media `json:"side"`
	Close          model.Media `json:"close"`
	FrontThumbnail model.Media `json:"frontThumbnail"`
	SideThumbnail  model.Media `json:"sideThumbnail"`
	CloseThumbnail model.Media `json:"closeThumbnail"`
}
