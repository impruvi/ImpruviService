package feedback

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

type CreateFeedbackRequest struct {
	UserId        string `json:"userId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
	FileLocation  string `json:"fileLocation"`
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

	user, err := users.GetUserById(request.UserId)
	coach, err := users.GetUserById(user.CoachUserId)
	if err == nil {
		notification.Notify(fmt.Sprintf("%v submitted feedback!", coach.Name))
	} else {
		log.Printf("Error while sending text notification: %v\n", err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
	}
}
