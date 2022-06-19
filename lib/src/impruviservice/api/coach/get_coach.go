package coach

import (
	"../../dao/coaches"
	"../../files"
	"../../model"
	"../converter"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
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

	coach, err := coaches.GetCoachById(request.CoachId)
	if coach.Headshot.UploadDateEpochMillis > 0 {
		coach.Headshot.FileLocation = files.GetHeadshotFileLocation(model.Coach, request.CoachId).URL
	}

	if err != nil {
		return converter.InternalServiceError("Error getting invitation code entry: %v\n", err)
	}

	return converter.Success(GetCoachResponse{Coach: coach})
}
