package coach

import (
<<<<<<< HEAD
<<<<<<< HEAD
	"impruviService/dao/coaches"
	"impruviService/api/converter"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
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

	err = coaches.PutCoach(request.Coach)
	if err != nil {
		return converter.InternalServiceError("Error getting invitation code entry: %v\n", err)
=======
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
>>>>>>> b2c6df1ca043c348ab7faab66c2a8cad9aaf1762
	}

	return converter.Success(nil)
}
