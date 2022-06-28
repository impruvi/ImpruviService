package coaches

import (
	"impruviService/awsclients/dynamoclient"
	"impruviService/model"
)

var dynamo = dynamoclient.GetClient()

const coachIdAttr = "coachId"

type Coach struct {
	CoachId        string       `json:"coachId"`
	FirstName      string       `json:"firstName"`
	LastName       string       `json:"lastName"`
	Email          string       `json:"email"`
	Headshot       *model.Media `json:"headshot"`
	Position       string       `json:"position"`
	School         string       `json:"school"`
	YouthClub      string       `json:"youthClub"`
	About          string       `json:"about"`
	NotificationId string       `json:"notificationId"`
}
