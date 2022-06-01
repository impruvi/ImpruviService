package drills

import (
	"../../dao/drills"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type GetAllDrillsRequest struct {}

type GetAllDrillsResponse struct {
	Drills []*drills.Drill `json:"drills"`
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

	allDrills, err := drills.GetAllDrills()
	if err != nil {
		log.Printf("Error while getting all drills: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	rspBody, err := json.Marshal(GetAllDrillsResponse{
		Drills: allDrills,
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

