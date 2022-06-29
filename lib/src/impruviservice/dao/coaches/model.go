package coaches

import (
	"impruviService/awsclients/dynamoclient"
<<<<<<< HEAD
	"impruviService/model"
=======
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
)

var dynamo = dynamoclient.GetClient()

const coachIdAttr = "coachId"

type Coach struct {
<<<<<<< HEAD
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
=======
	CoachId        string `json:"coachId"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	NotificationId string `json:"notificationId""`
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
}
