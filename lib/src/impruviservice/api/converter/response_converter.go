package converter

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

func Success(rsp any) *events.APIGatewayProxyResponse {
	rspBody, err := json.Marshal(rsp)
	if err != nil {
		return InternalServiceError("Error while marshalling response: %v\n", err)
	}

	return &events.APIGatewayProxyResponse{
		Body:       string(rspBody),
		StatusCode: http.StatusAccepted,
	}
}

func BadRequest(format string, v ...any) *events.APIGatewayProxyResponse {
	msg := fmt.Sprintf(format, v)
	log.Println(msg)
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       msg,
	}
}

func InternalServiceError(format string, v ...any) *events.APIGatewayProxyResponse {
	msg := fmt.Sprintf(format, v)
	log.Println(msg)
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "An unexpected error occurred",
	}
}

func NotAuthorizedError(format string, v ...any) *events.APIGatewayProxyResponse {
	msg := fmt.Sprintf(format, v)
	log.Println(msg)
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusForbidden,
		Body:       msg,
	}
}
