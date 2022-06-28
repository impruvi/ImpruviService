package session

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/session"
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

	return converter.Success(nil)
}
