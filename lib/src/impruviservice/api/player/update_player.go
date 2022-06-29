package player

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/players"
	playerFacade "impruviService/facade/player"
)

type UpdatePlayerRequest struct {
	Player *players.Player `json:"player"`
}

func UpdatePlayer(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request UpdatePlayerRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = playerFacade.UpdatePlayer(request.Player)
	if err != nil {
		return converter.InternalServiceError("Error getting invitation code entry: %v\n", err)
	}

	return converter.Success(nil)
}
