package invitationcodes

import (
	"impruviService/awsclients/dynamoclient"
<<<<<<< HEAD
	"impruviService/model"
=======
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
)

var dynamo = dynamoclient.GetClient()

const invitationCodeAttr = "invitationCode"

type InvitationCodeEntry struct {
	InvitationCode string         `json:"invitationCode"`
	UserId         string         `json:"userId"`
	UserType       model.UserType `json:"userType"`
}
