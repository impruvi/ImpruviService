package coach

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/coaches"
	coachFacade "impruviService/facade/coach"
)

type UpdateCoachRequest struct {
	Coach *coaches.Coach `json:"coach"`
}

func UpdateCoach(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request UpdateCoachRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = coachFacade.UpdateCoach(request.Coach)
	if err != nil {
		return converter.InternalServiceError("Error getting invitation code entry: %v\n", err)
	}

	return converter.Success(nil)
}
