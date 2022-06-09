package invitationcodes

import (
	"../../awsclients/dynamoclient"
)

var dynamo = dynamoclient.GetClient()

const invitationCodeAttr = "invitationCode"

type UserType string

const (
	Player UserType = "PLAYER"
	Coach           = "COACH"
)

type InvitationCodeEntry struct {
	InvitationCode string   `json:"invitationCode"`
	UserId         string   `json:"userId"`
	UserType       UserType `json:"userType"`
}
