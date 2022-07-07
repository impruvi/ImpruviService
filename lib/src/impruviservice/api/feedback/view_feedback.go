package feedback

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/session"
)

type ViewFeedbackRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
}

func ViewFeedback(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request ViewFeedbackRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = session.ViewFeedback(request.SessionNumber, request.PlayerId)
	if err != nil {
		return converter.InternalServiceError("Error while creating feedback: %v\n", err)
	}

	return converter.Success(nil)
}
