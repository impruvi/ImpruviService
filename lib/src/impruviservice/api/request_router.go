package api

import (
	"./feedback"
	"./session"
	"./submission"
	"./uploadurl"
	"./users"
	"./warmup"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

func RouteRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("request body: %v", request.Body)
	log.Printf("request path: %v", request.Path)
	log.Printf("request resource: %v", request.Resource)

	if request.Body == "WARM_UP_EVENT" {
		warmup.HandleWarmupEvent()
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusAccepted,
		}, nil
	}

	if request.Resource == "/validate-invitation-code" {
		return *users.ValidateCode(&request), nil
	} else if request.Resource == "/player/get-sessions" {
		return *session.GetPlayerSessions(&request), nil
	} else if request.Resource == "/coach/get-sessions" {
		return *session.GetCoachSessions(&request), nil
	} else if request.Resource == "/get-video-upload-url" {
		return *uploadurl.GetVideoUploadUrl(&request), nil
	} else if request.Resource == "/create-submission" {
		return *submission.CreateSubmission(&request), nil
	} else if request.Resource == "/create-feedback" {
		return *feedback.CreateFeedback(&request), nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("resource does not exist: %s", request.Path),
		StatusCode: http.StatusNotFound,
	}, nil
}
