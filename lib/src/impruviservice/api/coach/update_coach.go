package coach

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
)

type UpdateCoachRequest struct {
	Coach *coach.Coach `json:"coach"`
}

func UpdateCoach(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request CreateCoachRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = coach.PutCoach(request.Coach)
	if err != nil {
		return converter.InternalServiceError("Error while creating session: %v. %v\n", request.Session, err)
	}

	return converter.Success(nil)
}
