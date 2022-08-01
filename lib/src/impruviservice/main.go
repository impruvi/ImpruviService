package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stripe/stripe-go"
	"impruviService/api/coach"
	appCompatibility "impruviService/api/compatibility"
	"impruviService/api/drill"
	"impruviService/api/inbox"
	"impruviService/api/invitationcode"
	"impruviService/api/player"
	playerAuth "impruviService/api/player/auth"
	playerSubscription "impruviService/api/player/subscription"
	"impruviService/api/session"
	"impruviService/api/subscriptionplan"
	"impruviService/api/uploadurl"
	"impruviService/api/warmup"
	"impruviService/handler/mediaconvertevent"
	"impruviService/handler/reminders/dynamic"
	"impruviService/handler/reminders/fixed"
	"impruviService/router"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// Set your secret key. Remember to switch to your live secret key in production.
	// See your keys here: https://dashboard.stripe.com/apikeys
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func main() {
	lambda.Start(Handler{})
}

type Handler struct{}

var requestRouter = router.RequestRouter{
	WarmupHandler: warmup.HandleWarmupEvent,
	Handlers: map[string]interface{}{
		"/invitation-code/validate":            invitationcode.ValidateInvitationCode,
		"/subscription-plan/get":               subscriptionplan.GetSubscriptionPlan,
		"/player/sign-in":                      playerAuth.SignIn,
		"/player/sign-up/initiate":             playerAuth.InitiateSignUp,
		"/player/sign-up/complete":             playerAuth.CompleteSignUp,
		"/player/password-reset/initiate":      playerAuth.InitiatePasswordReset,
		"/player/password-reset/validate-code": playerAuth.ValidatePasswordReset,
		"/player/password-reset/complete":      playerAuth.CompletePasswordReset,
		"/player/payment-methods/get":          playerSubscription.GetPaymentMethods,
		"/player/subscription/re-activate":     playerSubscription.ReactivateSubscription,
		"/player/subscription/create":          playerSubscription.CreateSubscription,
		"/player/subscription/get":             playerSubscription.GetSubscription,
		"/player/subscription/cancel":          playerSubscription.CancelSubscription,
		"/player/update":                       player.UpdatePlayer,
		"/player/get":                          player.GetPlayer,
		"/player/inbox/get":                    inbox.GetInboxForPlayer,
		"/coaches/list":                        coach.ListCoaches,
		"/coach/update":                        coach.UpdateCoach,
		"/coach/get":                           coach.GetCoach,
		"/coach/players-and-subscriptions/get": coach.GetPlayersAndSubscriptions,
		"/sessions/feedback/view":              session.ViewFeedback,
		"/sessions/player/get":                 session.GetSessionsForPlayer,
		"/sessions/coach/get":                  session.GetSessionsForCoach,
		"/sessions/submission/create":          session.CreateSubmission,
		"/sessions/feedback/create":            session.CreateFeedback,
		"/sessions/get":                        session.GetSession,
		"/sessions/create":                     session.CreateSession,
		"/sessions/update":                     session.UpdateSession,
		"/sessions/delete":                     session.DeleteSession,
		"/drills/get":                          drills.GetDrill,
		"/drills/create":                       drills.CreateDrill,
		"/drills/update":                       drills.UpdateDrill,
		"/drills/delete":                       drills.DeleteDrill,
		"/drills/coach/get":                    drills.GetDrillsForCoach,
		"/drills/player/get":                   drills.GetDrillsForPlayer,
		"/media-upload-url/generate":           uploadurl.GetMediaUploadUrl,
		"/app-version/is-compatible":           appCompatibility.IsAppVersionCompatible,
	},
}

func (h Handler) Invoke(ctx context.Context, event []byte) ([]byte, error) {
	fmt.Printf("event: %s", string(event))
	lc, _ := lambdacontext.FromContext(ctx)
	if strings.Contains(lc.InvokedFunctionArn, "impruvi-service-api-handler") {
		return lambda.NewHandler(HandleAPIRequest).Invoke(ctx, event)
	} else if strings.Contains(lc.InvokedFunctionArn, "impruvi-service-fixed-reminder-notification-sender") {
		return lambda.NewHandler(fixed.HandleSendFixedRemindersEvent).Invoke(ctx, event)
	} else if strings.Contains(lc.InvokedFunctionArn, "impruvi-service-dynamic-reminder-notification-sender") {
		return lambda.NewHandler(dynamic.HandleSendDynamicRemindersEvent).Invoke(ctx, event)
	} else if strings.Contains(lc.InvokedFunctionArn, "impruvi-service-mediaconvert-event") {
		log.Printf("In impruvi-service-mediaconvert-event")
		return lambda.NewHandler(mediaconvertevent.HandleMediaConvertEvent).Invoke(ctx, event)
	} else {
		log.Printf(fmt.Sprintf("Unexpected lambda context: %s", lc.InvokedFunctionArn))
		return nil, errors.New(fmt.Sprintf("Unexpected lambda context: %s", lc.InvokedFunctionArn))
	}
}

func HandleAPIRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Resource: %s\n", request.Resource)

	response := requestRouter.Route(request)
	response.Headers = map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
	}
	return response, nil
}
