package session

import (
	"impruviService/awsclients/dynamoclient"
	"impruviService/model"
)

var dynamo = dynamoclient.GetClient()

const playerIdAttr = "playerId"
const sessionNumberAttr = "sessionNumber"

type Session struct {
	PlayerId                   string      `json:"playerId"`
	SessionNumber              int         `json:"sessionNumber"`
	Name                       string      `json:"name"`
	Drills                     []*Drill    `json:"drills"`
	Date                       *model.Date `json:"date"`
	CreationDateEpochMillis    int64       `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64       `json:"lastUpdatedDateEpochMillis"`
	HasViewedFeedback          bool        `json:"hasViewedFeedback"`
}

type Drill struct {
	DrillId                  string      `json:"drillId"`
	Submission               *Submission `json:"submission"`
	Feedback                 *Feedback   `json:"feedback"`
	Notes                    string      `json:"notes"`
	EstimatedDurationMinutes int         `json:"estimatedDurationMinutes"`
}

type Submission struct {
	VideoUploadDateEpochMillis int64 `json:"videoUploadDateEpochMillis"`
}

type Feedback struct {
	VideoUploadDateEpochMillis int64 `json:"videoUploadDateEpochMillis"`
}
