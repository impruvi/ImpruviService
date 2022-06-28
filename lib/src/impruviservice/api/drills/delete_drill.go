package drills

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/drills"
)

type DeleteDrillRequest struct {
	DrillId string `json:"drillId"`
}

func DeleteDrill(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request DeleteDrillRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = drills.DeleteDrill(request.DrillId)
	if err != nil {
		return converter.InternalServiceError("Error while deleting drill with id: %v. %v\n", request.DrillId, err)
	}

	return converter.Success(nil)
}
