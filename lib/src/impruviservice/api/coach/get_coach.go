package coach

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/coaches"
	coachFacade "impruviService/facade/coach"
)

type GetCoachRequest struct {
	CoachId string `json:"coachId"`
}

type GetCoachResponse struct {
	Coach *coaches.Coach `json:"coach"`
}

func GetCoach(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetCoachRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	coach, err := coachFacade.GetCoachById(request.CoachId)

	if err != nil {
		return converter.InternalServiceError("Error getting invitation code entry: %v\n", err)
	}

	return converter.Success(GetCoachResponse{Coach: coach})
}
