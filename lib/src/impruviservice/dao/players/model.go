package players

import (
	"impruviService/awsclients/dynamoclient"
)

var dynamo = dynamoclient.GetClient()

const coachIdIndexName = "coachId-index"
const playerIdAttr = "playerId"
const coachIdAttr = "coachId"

type Player struct {
	PlayerId                   string   `json:"playerId"`
	CoachId                    string   `json:"coachId"`
	FirstName                  string   `json:"firstName"`
	LastName                   string   `json:"lastName"`
	Email                      string   `json:"email"`
	Availability               []string `json:"availability"`
	NotificationId             string   `json:"notificationId"`
	CreationDateEpochMillis    int64    `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64    `json:"lastUpdatedDateEpochMillis"`
}
