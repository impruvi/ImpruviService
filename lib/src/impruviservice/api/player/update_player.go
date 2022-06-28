package player

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/players"
	playerFacade "impruviService/facade/players
)

type UpdatePlayerRequest struct {
	Player *player.Player `json:"player"`
}

func UpdatePlayer(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request CreatePlayerRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = player.PutPlayer(request.Coach)
	if err != nil {
		return converter.InternalServiceError("Error while creating session: %v. %v\n", request.Session, err))
	}

	return converter.Success(nil)
}
