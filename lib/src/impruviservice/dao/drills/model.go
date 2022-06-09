package drills

import (
	"../../awsclients/dynamoclient"
)

var dynamo = dynamoclient.GetClient()

const coachIdIndexName = "coachId-index"
const coachIdAttr = "coachId"
const drillIdAttr = "drillId"

type Drill struct {
	DrillId     string      `json:"drillId"`
	CoachId     string      `json:"coachId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Category    string      `json:"category"` // DRIBBLING/WARMUP/SHOOTING/PASSING
	Equipment   []Equipment `json:"equipment"`
}

type Equipment struct {
	EquipmentType string      `json:"equipmentType"` // BALL/CONE/SPACE/GOAL
	Requirement   Requirement `json:"requirement"`
}

type Requirement struct {
	RequirementType string     `json:"requirementType"` // COUNT/DIMENSION
	Count           int        `json:"count"`
	Dimensions      Dimensions `json:"dimension"`
}

// Dimensions in yards
type Dimensions struct {
	Width  int `json:"width"`
	Length int `json:"length"`
}
