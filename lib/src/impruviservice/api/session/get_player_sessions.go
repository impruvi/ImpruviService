package session

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type GetPlayerSessionsRequest struct {
	UserId string `json:"userId"`
}

type GetPlayerSessionsResponse struct {
	Sessions []*Session `json:"sessions"`
}

func GetPlayerSessions(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetPlayerSessionsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	sessionsWithDrills, err := getSessionsWithDrillsForUser(request.UserId)
	if err != nil {
		log.Printf("Error while getting drills for sessions for userId: %v, %v\n", request.UserId, err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	rspBody, err := json.Marshal(GetPlayerSessionsResponse{
		Sessions: sessionsWithDrills,
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
