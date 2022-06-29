package session

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/session"
<<<<<<< HEAD
	playerFacade "impruviService/facade/player"
	"log"
=======
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
)

type DeleteSessionRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
}

func DeleteSession(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request DeleteSessionRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = session.DeleteSession(request.SessionNumber, request.PlayerId)
	if err != nil {
		return converter.InternalServiceError("Error while deleting session number: %v. for player: %v, %v\n", request.SessionNumber, request.PlayerId, err)
	}
	err = playerFacade.UpdateSessionDatesForPlayerByPlayerId(request.PlayerId)
	if err != nil {
		return converter.InternalServiceError("Error while updating session dates for player: %v. %v\n", request.PlayerId, err)
	}

	log.Printf("Deleted session: %v, for player %v\n", request.SessionNumber, request.PlayerId)

	return converter.Success(nil)
}
