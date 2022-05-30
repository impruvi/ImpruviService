package session

import (
	"../../dao/session"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type GetSessionsRequest struct {
	UserId string `json:"userId"`
}

type GetSessionsResponse struct {
	Sessions []*session.Session `json:"sessions"`
}

func GetSessions(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetSessionsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	sessions, err := session.GetSessions(request.UserId)

	rspBody, err := json.Marshal(GetSessionsResponse{
		Sessions: sessions,
	})
	if err != nil {
		log.Printf("Error while marshalling response: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &events.APIGatewayProxyResponse{
		Body:       string(rspBody),
		StatusCode: http.StatusAccepted,
	}
}
