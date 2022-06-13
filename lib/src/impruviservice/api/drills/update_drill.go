package drills

import (
	"../../dao/drills"
	"../converter"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type UpdateDrillRequest struct {
	Drill *drills.Drill `json:"drill"`
}

func UpdateDrill(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request CreateDrillRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = drills.PutDrill(request.Drill)
	if err != nil {
		return converter.InternalServiceError("Error while updating drill: %v. %v\n", request.Drill, err)
	}

	return converter.Success(nil)
}
