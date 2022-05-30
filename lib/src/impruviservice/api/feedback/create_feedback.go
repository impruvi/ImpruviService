package feedback

import (
	"../../dao/session"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type CreateFeedbackRequest struct {
	UserId string `json:"userId"`
	SessionNumber int `json:"sessionNumber"`
	DrillId string `json:"drillId"`
	FileLocation string `json:"fileLocation"`
}

func CreateFeedback(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request CreateFeedbackRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	err = session.CreateFeedback(request.SessionNumber, request.UserId, request.DrillId, request.FileLocation)
	if err != nil {
		log.Printf("Error while creating feedback: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
	}
}
