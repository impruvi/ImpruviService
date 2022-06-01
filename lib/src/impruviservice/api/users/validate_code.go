package users

import (
	"../../dao/users"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type ValidateCodeRequest struct {
	InvitationCode string `json:"invitationCode"`
}

type ValidateCodeResponse struct {
	User *users.User `json:"user"`
}

func ValidateCode(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request ValidateCodeRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	user, err := users.GetUserByInvitationCode(request.InvitationCode)
	if err != nil {
		log.Printf("Error getting user by invitation code: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	rspBody, err := json.Marshal(ValidateCodeResponse{
		User: user,
	})
	if err != nil {
		log.Printf("Error while marshalling response: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &events.APIGatewayProxyResponse{
		Body:       string(rspBody),
		StatusCode: http.StatusAccepted,
	}
}
