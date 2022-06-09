package players

import (
	"../../awsclients/dynamoclient"
)

var dynamo = dynamoclient.GetClient()

const coachIdIndexName = "coachId-index"
const playerIdAttr = "playerId"
const coachIdAttr = "coachId"

type Player struct {
	PlayerId  string `json:"playerId"`
	CoachId   string `json:"coachId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
