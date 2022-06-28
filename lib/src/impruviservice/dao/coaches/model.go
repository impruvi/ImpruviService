package coaches

import (
	"impruviService/awsclients/dynamoclient"
)

var dynamo = dynamoclient.GetClient()

const coachIdAttr = "coachId"

type Coach struct {
	CoachId        string `json:"coachId"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	NotificationId string `json:"notificationId""`
}
