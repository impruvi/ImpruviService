package drills

import (
	"impruviService/model"
)

const coachIdIndexName = "coachId-index"
const coachIdAttr = "coachId"
const drillIdAttr = "drillId"

type DemoAngle string

const (
	FrontAngle DemoAngle = "FRONT"
	SideAngle  DemoAngle = "SIDE"
	CloseAngle DemoAngle = "CLOSE"
)

type DrillDB struct {
	DrillId     string         `json:"drillId"`
	CoachId     string         `json:"coachId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"` // DRIBBLING/WARMUP/SHOOTING/PASSING
	Equipment   []*EquipmentDB `json:"equipment"`
	Demos       *DemosDB       `json:"demos"`
	IsDeleted   bool           `json:"isDeleted"`
}

type EquipmentDB struct {
	EquipmentType string         `json:"equipmentType"` // BALL/CONE/SPACE/GOAL
	Requirement   *RequirementDB `json:"requirement"`
}

type RequirementDB struct {
	RequirementType string        `json:"requirementType"` // COUNT/DIMENSION
	Count           int           `json:"count"`
	Dimensions      *DimensionsDB `json:"dimension"`
}

// DimensionsDB in yards
type DimensionsDB struct {
	Width  int `json:"width"`
	Length int `json:"length"`
}

type DemosDB struct {
	Front          *model.Media `json:"front"`
	Side           *model.Media `json:"side"`
	Close          *model.Media `json:"close"`
	FrontThumbnail *model.Media `json:"frontThumbnail"`
	SideThumbnail  *model.Media `json:"sideThumbnail"`
	CloseThumbnail *model.Media `json:"closeThumbnail"`
}
