package invitationcodes

import (
	"impruviService/model"
)

const invitationCodeAttr = "invitationCode"

type InvitationCodeEntryDB struct {
	InvitationCode string         `json:"invitationCode"`
	UserId         string         `json:"userId"`
	UserType       model.UserType `json:"userType"`
}
