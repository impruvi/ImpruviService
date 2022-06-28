package session

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/session"
)

type CreateSessionRequest struct {
	Session *session.Session `json:"session"`
}

func CreateSession(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request CreateSessionRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = session.CreateSession(request.Session)
	if err != nil {
		return converter.InternalServiceError("Error while creating session: %v. %v\n", request.Session, err)
	}

	return converter.Success(nil)
}
