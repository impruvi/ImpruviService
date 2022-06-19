package users

import (
	"../../dao/coaches"
	"../../dao/invitationcodes"
	"../../dao/players"
	"../../model"
	"../converter"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

type ValidateCodeRequest struct {
	InvitationCode string `json:"invitationCode"`
}

type ValidateCodeResponse struct {
	UserType model.UserType  `json:"userType"` // PLAYER/COACH
	Player   *players.Player `json:"player"`
	Coach    *coaches.Coach  `json:"coach"`
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

	if invitationCodeEntry.UserType == model.Coach {
		coach, err := coaches.GetCoachById(invitationCodeEntry.UserId)
		if err != nil {
			return converter.InternalServiceError("Error getting coach by coachId: %v\n", err)
		}
		log.Printf("Coach: %v\n", coach)
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
		return converter.Success(ValidateCodeResponse{
			UserType: invitationCodeEntry.UserType,
			Player:   player,
		})
	}
}
