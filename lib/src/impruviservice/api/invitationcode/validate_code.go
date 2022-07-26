package invitationcode

import (
	"fmt"
	"impruviService/dao/coach"
	"impruviService/dao/invitationcode"
	"impruviService/exceptions"
	coachFacade "impruviService/facade/coach"
	playerFacade "impruviService/facade/player"
	"impruviService/model"
	"log"
	"strings"
)

type ValidateCodeRequest struct {
	InvitationCode string `json:"invitationCode"`
	ExpoPushToken  string `json:"expoPushToken"`
}

type ValidateCodeResponse struct {
	UserType model.UserType       `json:"userType"`
	Player   *playerFacade.Player `json:"player"`
	Coach    *coaches.CoachDB     `json:"coach"`
}

func ValidateCode(request *ValidateCodeRequest) (*ValidateCodeResponse, error) {
	request.InvitationCode = strings.TrimSpace(request.InvitationCode)

	invitationCodeEntry, err := invitationcodes.GetInvitationCodeEntry(request.InvitationCode)

	if err != nil {
		log.Printf("Error: %v\n", err)
		if _, ok := err.(exceptions.ResourceNotFoundError); ok {
			return nil, exceptions.NotAuthorizedError{Message: fmt.Sprintf("Invalid invitation code: %v\n", request.InvitationCode)}
		} else {
			return nil, err
		}
	}
	log.Printf("Invitation code entry: %v\n", invitationCodeEntry)

	if invitationCodeEntry.UserType == model.Coach {
		coach, err := coachFacade.GetCoachById(invitationCodeEntry.UserId)
		if err != nil {
			return nil, err
		}
		log.Printf("Coach: %v\n", coach)
		if coach.NotificationId != request.ExpoPushToken {
			coach, err = coachFacade.UpdateNotificationId(coach, request.ExpoPushToken)
			if err != nil {
				return nil, err
			}
		}
		return &ValidateCodeResponse{
			UserType: invitationCodeEntry.UserType,
			Coach:    coach,
		}, nil
	} else {
		player, err := playerFacade.GetPlayerById(invitationCodeEntry.UserId)
		if err != nil {
			return nil, err
		}
		log.Printf("Player: %v\n", player)
		if player.NotificationId != request.ExpoPushToken {
			player, err = playerFacade.UpdateNotificationId(player.PlayerId, request.ExpoPushToken)
			if err != nil {
				return nil, err
			}
		}
		return &ValidateCodeResponse{
			UserType: invitationCodeEntry.UserType,
			Player:   player,
		}, nil
	}
}
