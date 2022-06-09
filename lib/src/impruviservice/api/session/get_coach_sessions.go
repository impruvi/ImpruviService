package session

import (
	"../../dao/players"
	"../converter"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type GetCoachSessionsRequest struct {
	CoachId string `json:"coachId"`
}

type GetCoachSessionsResponse struct {
	PlayerSessions []*PlayerSessions `json:"playerSessions"`
}

type PlayerSessions struct {
	Player   *players.Player `json:"player"`
	Sessions []*FullSession  `json:"sessions"`
}

func GetCoachSessions(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetPlayerSessionsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	playersForCoach, err := players.GetPlayersForCoach(request.PlayerId)
	if err != nil {
		return converter.InternalServiceError("Error while players for coach: %v. %v\n", request.PlayerId, err)
	}
	playerSessions, err := getPlayerSessions(playersForCoach)
	if err != nil {
		return converter.InternalServiceError("Error while getting player sessions for players: %v. %v\n", playersForCoach, err)
	}

	return converter.Success(GetCoachSessionsResponse{
		PlayerSessions: playerSessions,
	})
}

func getPlayerSessions(playersForCoach []*players.Player) ([]*PlayerSessions, error) {
	playerSessions := make([]*PlayerSessions, 0)
	for _, player := range playersForCoach {
		sessionsWithDrills, err := getFullSessionsForPlayer(player.PlayerId)
		if err != nil {
			return nil, err
		}

		playerSessions = append(playerSessions, &PlayerSessions{
			Player:   player,
			Sessions: sessionsWithDrills,
		})
	}
	return playerSessions, nil
}
