package session

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/players"
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
	var request GetCoachSessionsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	playersForCoach, err := players.GetPlayersForCoach(request.CoachId)
	if err != nil {
		return converter.InternalServiceError("Error while getting players for coach: %v. %v\n", request.CoachId, err)
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
