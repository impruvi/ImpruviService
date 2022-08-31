package router

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/exceptions"
	"impruviService/handler/stripeevent"
	"log"
	"net/http"
	"reflect"
)

type RequestRouter struct {
	WarmupHandler interface{}
	Handlers      map[string]interface{}
}

func (r *RequestRouter) Route(apiRequest events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	//log.Printf("request body: %v", apiRequest.Body)
	log.Printf("request resource: %v", apiRequest.Resource)
	log.Printf("headers: %v\n", apiRequest.Headers)

	if apiRequest.Body == "WARM_UP_EVENT" {
		reflect.ValueOf(r.WarmupHandler).Call([]reflect.Value{})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusAccepted,
		}
	}

	// TODO: don't make a special case for stripe events
	if apiRequest.Resource == "/stripe-event" {
		return *stripeevent.HandleStripeEvent(&apiRequest)
	}

	// find the appropriate handler based on the HTTP resource path
	handler, ok := r.Handlers[apiRequest.Resource]
	if !ok {
		log.Printf("No handler for resource: %v\n", apiRequest.Resource)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("resource does not exist: %s", apiRequest.Resource),
			StatusCode: http.StatusNotFound,
		}
	}

	// unmarshal the request body into the appropriate type
	// here we attempt to unmarshal the request body JSON into the type of the request handlers first argument
	request := reflect.New(reflect.ValueOf(handler).Type().In(0))
	var err = json.Unmarshal([]byte(apiRequest.Body), request.Interface())
	if err != nil {
		log.Printf("Unable to unmarshal body in to request: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("Bad request body: %v for resource: %v\n", apiRequest.Body, apiRequest.Resource),
		}
	}

	// invoke the request handler
	res := reflect.ValueOf(handler).Call([]reflect.Value{request.Elem()})

	// convert the response
	// each request handler can either return
	// - a response and an error
	// - just an error
	// So here we check if the handler is returning just one parameter or two. If the handler
	// returns only one parameter, that parameter is assumed to be of type `error`. If the handler
	// returns two arguments they are assumed to be of type `responseObject` and `error`.
	if len(res) == 1 {
		if res[0].Interface() != nil {
			// The handler returned an error
			err = res[0].Interface().(error)
			log.Printf("Error: %v\n", err)
			return convertError(err)
		} else {
			// The handler did not return an error
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusAccepted,
			}
		}
	} else {
		if res[1].Interface() != nil {
			// The handler returned an error
			err = res[1].Interface().(error)
			log.Printf("Error: %v\n", err)
			return convertError(err)
		} else {
			// The handler did not return an error
			// marshal the response object to JSON
			resBytes, err := json.Marshal(res[0].Interface())
			if err != nil {
				log.Printf("[ERROR] Error: %v\n", err)
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Body:       "An unexpected error occurred",
				}
			}

			log.Printf("Response: %v\n", string(resBytes))
			return events.APIGatewayProxyResponse{
				Body:       string(resBytes),
				StatusCode: http.StatusAccepted,
			}
		}
	}
}

func convertError(err error) events.APIGatewayProxyResponse {
	if _, ok := err.(exceptions.NotAuthorizedError); ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusForbidden,
			Body:       err.Error(),
		}
	} else if _, ok := err.(exceptions.InvalidRequestError); ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}
	} else if _, ok := err.(exceptions.ResourceNotFoundError); ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       err.Error(),
		}
	} else {
		log.Printf("[ERROR] Error: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "An unexpected error occurred",
		}
	}
}
