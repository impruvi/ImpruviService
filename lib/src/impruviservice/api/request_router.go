package api

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/coach"
	"impruviService/api/drills"
	"impruviService/api/feedback"
	"impruviService/api/invitationcode"
	"impruviService/api/player"
	"impruviService/api/session"
	"impruviService/api/submission"
	"impruviService/api/uploadurl"
	"impruviService/api/warmup"
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
	} else if request.Resource == "/player/update" {
		return *player.UpdatePlayer(&request), nil
	} else if request.Resource == "/coach/update" {
		return *coach.UpdateCoach(&request), nil
	} else if request.Resource == "/coach/get" {
		return *coach.GetCoach(&request), nil
	} else if request.Resource == "/sessions/player/get" {
		return *session.GetPlayerSessions(&request), nil
	} else if request.Resource == "/sessions/coach/get" {
		return *session.GetCoachSessions(&request), nil
	} else if request.Resource == "/sessions/submission/create" {
		return *submission.CreateSubmission(&request), nil
	} else if request.Resource == "/sessions/feedback/create" {
		return *feedback.CreateFeedback(&request), nil
	} else if request.Resource == "/sessions/create" {
		return *session.CreateSession(&request), nil
	} else if request.Resource == "/sessions/update" {
		return *session.UpdateSession(&request), nil
	} else if request.Resource == "/sessions/delete" {
		return *session.DeleteSession(&request), nil
	} else if request.Resource == "/drills/create" {
		return *drills.CreateDrill(&request), nil
	} else if request.Resource == "/drills/update" {
		return *drills.UpdateDrill(&request), nil
	} else if request.Resource == "/drills/delete" {
		return *drills.DeleteDrill(&request), nil
	} else if request.Resource == "/drills/coach/get" {
		return *drills.GetDrillsForCoach(&request), nil
	} else if request.Resource == "/get-video-upload-url" {
		return *uploadurl.GetVideoUploadUrl(&request), nil
	} else if request.Resource == "/get-video-thumbnail-upload-url" {
		return *uploadurl.GetVideoThumbnailUploadUrl(&request), nil
	} else if request.Resource == "/get-headshot-upload-url" {
		return *uploadurl.GetHeadshotUploadUrl(&request), nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("resource does not exist: %s", request.Path),
		StatusCode: http.StatusNotFound,
	}, nil
}
