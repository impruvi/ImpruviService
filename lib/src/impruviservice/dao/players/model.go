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
<<<<<<< HEAD
	PlayerId     string   `json:"playerId"`
	CoachId      string   `json:"coachId"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Email        string   `json:"email"`
	Availability []string `json:"availability"`
=======
=======
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
	PlayerId       string `json:"playerId"`
	CoachId        string `json:"coachId"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	NotificationId string `json:"notificationId""`
<<<<<<< HEAD
>>>>>>> b2c6df1 (push notification changes)
=======
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
}
