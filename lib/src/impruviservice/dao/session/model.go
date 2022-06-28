package session

import (
	"impruviService/awsclients/dynamoclient"
)

var dynamo = dynamoclient.GetClient()

const playerIdAttr = "playerId"
const sessionNumberAttr = "sessionNumber"

type Session struct {
	PlayerId      string   `json:"playerId"`
	SessionNumber int      `json:"sessionNumber"`
	Drills        []*Drill `json:"drills"`
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
