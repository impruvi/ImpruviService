package session

import (
	"../converter"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

type GetPlayerSessionsRequest struct {
	PlayerId string `json:"playerId"`
}

type GetPlayerSessionsResponse struct {
	Sessions []*FullSession `json:"sessions"`
}

func GetPlayerSessions(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetPlayerSessionsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	fullSessions, err := getFullSessionsForPlayer(request.PlayerId)
	if err != nil {
		return converter.InternalServiceError("Error while getting full sessions for player with id: %v, %v\n", request.PlayerId, err)
	}
	log.Printf("Full sessions: %v\n", fullSessions)

	return converter.Success(GetPlayerSessionsResponse{
		Sessions: fullSessions,
	})
}
