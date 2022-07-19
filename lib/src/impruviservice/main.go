package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"impruviService/api/coach"
	"impruviService/api/drill"
	"impruviService/api/inbox"
	"impruviService/api/invitationcode"
	"impruviService/api/player"
	"impruviService/api/session"
	"impruviService/api/uploadurl"
	"impruviService/api/warmup"
	"impruviService/handlers/notification"
	"impruviService/router"
	"log"
	"strings"
)

func main() {
	lambda.Start(Handler{})
}

type Handler struct{}

var requestRouter = router.RequestRouter{
	WarmupHandler: warmup.HandleWarmupEvent,
	Handlers: map[string]interface{}{
		"/invitation-code/validate":   invitationcode.ValidateCode,
		"/player/update":              player.UpdatePlayer,
		"/player/get":                 player.GetPlayer,
		"/player/inbox/get":           inbox.GetInboxForPlayer,
		"/coach/update":               coach.UpdateCoach,
		"/coach/get":                  coach.GetCoach,
		"/sessions/feedback/view":     session.ViewFeedback,
		"/sessions/player/get":        session.GetSessionsForPlayer,
		"/sessions/coach/get":         session.GetSessionForCoach,
		"/sessions/submission/create": session.CreateSubmission,
		"/sessions/feedback/create":   session.CreateFeedback,
		"/sessions/create":            session.CreateSession,
		"/sessions/update":            session.UpdateSession,
		"/sessions/delete":            session.DeleteSession,
		"/drills/create":              drills.CreateDrill,
		"/drills/update":              drills.UpdateDrill,
		"/drills/delete":              drills.DeleteDrill,
		"/drills/coach/get":           drills.GetDrillsForCoach,
		"/drills/player/get":          drills.GetDrillsForPlayer,
		"/media-upload-url/generate":  uploadurl.GetMediaUploadUrl,
	},
}

func (h Handler) Invoke(ctx context.Context, event []byte) ([]byte, error) {
	fmt.Printf("event: %s", string(event))
	lc, _ := lambdacontext.FromContext(ctx)
	if strings.Contains(lc.InvokedFunctionArn, "impruvi-service-api-handler") {
		return lambda.NewHandler(HandleAPIRequest).Invoke(ctx, event)
	} else if strings.Contains(lc.InvokedFunctionArn, "impruvi-service-notification-sender") {
		return lambda.NewHandler(HandleSendNotificationEvent).Invoke(ctx, event)
	} else {
		log.Printf(fmt.Sprintf("Unexpected lambda context: %s", lc.InvokedFunctionArn))
		return nil, errors.New(fmt.Sprintf("Unexpected lambda context: %s", lc.InvokedFunctionArn))
	}
}

func HandleAPIRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Resource: %s\n", request.Resource)

	return requestRouter.Route(request), nil
}

func HandleSendNotificationEvent() error {
	return notification.SendScheduledNotifications()
}
