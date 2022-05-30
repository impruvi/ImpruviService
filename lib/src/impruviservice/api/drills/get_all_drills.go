package drills

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type GetAllDrillsRequest struct {
	Code string `json:"code"`
}

type GetAllDrillsResponse struct {
	UserId string `json:"userId"`
}

func GetAllDrills(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetAllDrillsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	// query users table by invitation code

	rspBody, err := json.Marshal(GetAllDrillsResponse{
		UserId: "", // TODO: get the actual UserId
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

