package player

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
<<<<<<< HEAD
	"impruviService/dao/players"
	playerFacade "impruviService/facade/players
=======
	"impruviService/dao/player"
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
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
<<<<<<< HEAD
		return converter.InternalServiceError("Error while creating session: %v. %v\n", request.Session, err))
=======
		return converter.InternalServiceError("Error while creating session: %v. %v\n", request.Session, err)
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
	}

	return converter.Success(nil)
}
