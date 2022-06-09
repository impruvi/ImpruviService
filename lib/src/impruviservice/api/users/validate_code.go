package users

import (
	"../../dao/coaches"
	"../../dao/invitationcodes"
	"../../dao/players"
	"../converter"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type ValidateCodeRequest struct {
	InvitationCode string `json:"invitationCode"`
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

	if invitationCodeEntry.UserType == invitationcodes.Coach {
		coach, err := coaches.GetCoachById(invitationCodeEntry.UserId)
		if err != nil {
			return converter.InternalServiceError("Error getting coach by coachId: %v\n", err)
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
		return converter.Success(ValidateCodeResponse{
			UserType: invitationCodeEntry.UserType,
			Player:   player,
		})
	}
}
