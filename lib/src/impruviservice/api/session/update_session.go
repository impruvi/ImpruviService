package session

import (
	"../../dao/session"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type UpdateSessionsRequest struct {
	UserId string `json:"userId"`
	SessionNumber int `json:"sessionNumber"`
	Session *session.Session `json:"session"`
}

type UpdateSessionResponse struct {
	UserId string `json:"userId"`
}

func UpdateSession(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request UpdateSessionsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	err = session.UpdateSession(request.SessionNumber, request.UserId, request.Session)
	if err != nil {
		log.Printf("Error while updating session: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
	}
}

