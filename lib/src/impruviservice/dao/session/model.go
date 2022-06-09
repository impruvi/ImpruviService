package session

import (
	"../../awsclients/dynamoclient"
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
	DrillId                  string       `json:"drillId"`
	Submission               *Submission  `json:"submission"`
	Feedback                 *Feedback    `json:"feedback"`
	Notes                    string       `json:"notes"`
	EstimatedDurationMinutes int          `json:"estimatedDurationMinutes"`
	Prescription             Prescription `json:"prescription"`
}

type Prescription struct {
	Type  string `json:"type"` // REPETITIONS/DURATION_MINUTES
	Value int    `json:"value"`
}

type Submission struct {
	VideoUploadDateEpochMillis int64 `json:"videoUploadDateEpochMillis"`
}

type Feedback struct {
	VideoUploadDateEpochMillis int64 `json:"videoUploadDateEpochMillis"`
}
