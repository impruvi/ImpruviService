package invitationcodes

import (
	"impruviService/awsclients/dynamoclient"
	"impruviService/model"
)

var dynamo = dynamoclient.GetClient()

const invitationCodeAttr = "invitationCode"

type InvitationCodeEntry struct {
	InvitationCode string         `json:"invitationCode"`
	UserId         string         `json:"userId"`
	UserType       model.UserType `json:"userType"`
}
