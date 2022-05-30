package main

import (
	"./api"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"log"
	"strings"
)

func main() {
	lambda.StartHandler(Handler{})
}

type Handler struct{}

func (h Handler) Invoke(ctx context.Context, event []byte) ([]byte, error) {
	fmt.Printf("event: %s", string(event))
	lc, _ := lambdacontext.FromContext(ctx)
	if strings.Contains(lc.InvokedFunctionArn, "impruvi-service-api-handler") {
		return lambda.NewHandler(HandleAPIRequest).Invoke(ctx, event)
	} else {
		log.Printf(fmt.Sprintf("Unexpected lambda context: %s", lc.InvokedFunctionArn))
		return nil, errors.New(fmt.Sprintf("Unexpected lambda context: %s", lc.InvokedFunctionArn))
	}
}

func HandleAPIRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Resource: %s\n", request.Resource)

	return api.RouteRequest(request)
}
