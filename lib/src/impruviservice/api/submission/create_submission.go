package submission

import (
	"../../dao/session"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type CreateSubmissionRequest struct {
	UserId string `json:"userId"`
	SessionNumber int `json:"sessionNumber"`
	DrillId string `json:"drillId"`
	FileLocation string `json:"fileLocation"`
}

func CreateSubmission(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request CreateSubmissionRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	err = session.CreateSubmission(request.SessionNumber, request.UserId, request.DrillId, request.FileLocation)
	if err != nil {
		log.Printf("Error while creating submission: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
	}
}

