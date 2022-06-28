package players

import (
	"impruviService/awsclients/dynamoclient"
)

var dynamo = dynamoclient.GetClient()

const coachIdIndexName = "coachId-index"
const playerIdAttr = "playerId"
const coachIdAttr = "coachId"

type Player struct {
<<<<<<< HEAD
	PlayerId     string   `json:"playerId"`
	CoachId      string   `json:"coachId"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Email        string   `json:"email"`
	Availability []string `json:"availability"`
=======
	PlayerId       string `json:"playerId"`
	CoachId        string `json:"coachId"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	NotificationId string `json:"notificationId""`
>>>>>>> b2c6df1 (push notification changes)
}
