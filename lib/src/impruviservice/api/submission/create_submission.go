package submission

import (
	"../../dao/session"
	"../../dao/users"
	"../../notification"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type CreateSubmissionRequest struct {
	UserId        string `json:"userId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
	FileLocation  string `json:"fileLocation"`
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

	user, err := users.GetUserById(request.UserId)
	if err == nil {
		notification.Notify(fmt.Sprintf("%v submitted a video!", user.Name))
	} else {
		log.Printf("Error while getting user by id: %v\n", err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
	}
}
