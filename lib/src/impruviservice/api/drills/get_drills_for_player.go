package drills

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	playerFacade "impruviService/facade/player"
)

type GetDrillsForPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetDrillsForPlayerResponse struct {
	Drills []*FullDrill `json:"drills"`
}

func GetDrillsForPlayer(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetDrillsForPlayerRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	drills, err := playerFacade.GetDrillsForPlayer(request.PlayerId)
	if err != nil {
		return converter.InternalServiceError("Error while getting drills for player: %v. %v\n", request.PlayerId, err)
	}

	fullDrills := getFullDrills(drills)

	return converter.Success(GetDrillsForPlayerResponse{Drills: fullDrills})
}