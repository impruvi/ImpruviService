package users

import (
	"../../dao/users"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type GetAllUsersRequest struct {
}

type GetAllUsersResponse struct {
	Users []*users.User `json:"users"`
}

func GetAllUsers(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetAllUsersRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	allUsers, err := users.GetAllUsers()
	if err != nil {
		log.Printf("Error while fetching all users: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	rspBody, err := json.Marshal(GetAllUsersResponse{
		Users: allUsers,
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
