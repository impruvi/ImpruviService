package users

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/coaches"
	"impruviService/dao/invitationcodes"
	"impruviService/dao/players"
	"log"
)

type ValidateCodeRequest struct {
	InvitationCode string `json:"invitationCode"`
	ExpoPushToken  string `json:"expoPushToken"`
}

type ValidateCodeResponse struct {
	UserType invitationcodes.UserType `json:"userType"` // PLAYER/COACH
	Player   *players.Player          `json:"player"`
	Coach    *coaches.Coach           `json:"coach"`
}

func ValidateCode(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request ValidateCodeRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	invitationCodeEntry, err := invitationcodes.GetInvitationCodeEntry(request.InvitationCode)
	if err != nil {
		return converter.InternalServiceError("Error getting invitation code entry: %v\n", err)
	}
	log.Printf("Invitation code entry: %v\n", invitationCodeEntry)

	if invitationCodeEntry.UserType == invitationcodes.Coach {
		coach, err := coaches.GetCoachById(invitationCodeEntry.UserId)
		if err != nil {
			return converter.InternalServiceError("Error getting coach by coachId: %v\n", err)
		}
		log.Printf("Coach: %v\n", coach)
		if coach.NotificationId != request.ExpoPushToken {
			coach = updateCoachNotificationId(coach, request.ExpoPushToken)
		}
		return converter.Success(ValidateCodeResponse{
			UserType: invitationCodeEntry.UserType,
			Coach:    coach,
		})
	} else {
		player, err := players.GetPlayerById(invitationCodeEntry.UserId)
		if err != nil {
			return converter.InternalServiceError("Error getting player by playerId: %v\n", err)
		}
		log.Printf("Player: %v\n", player)
		if player.NotificationId != request.ExpoPushToken {
			player = updatePlayerNotificationId(player, request.ExpoPushToken)
		}
		return converter.Success(ValidateCodeResponse{
			UserType: invitationCodeEntry.UserType,
			Player:   player,
		})
	}
}

func updateCoachNotificationId(coach *coaches.Coach, notificationId string) *coaches.Coach {
	var newCoach = coach
	newCoach.NotificationId = notificationId
	err := coaches.PutCoach(newCoach)
	if err != nil {
		// don't error if we can't update push notification id
		log.Printf("Error while updating coach notification id: %v\n", err)
		return coach
	}
	log.Printf("Updated coach id's %v push notification token", coach.CoachId)
	return newCoach
}

func updatePlayerNotificationId(player *players.Player, notificationId string) *players.Player {
	var newPlayer = player
	newPlayer.NotificationId = notificationId
	err := players.PutPlayer(player)
	if err != nil {
		// don't error if we can't update push notification id
		log.Printf("Error while updating player notification id: %v\n", err)
		return player
	}
	log.Printf("Updated player id's %v push notification token", player.PlayerId)
	return newPlayer
}
