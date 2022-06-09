package drills

import (
	drillsDao "../../dao/drills"
	"../../files"
	"../converter"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type GetDrillsForCoachRequest struct {
	CoachId string `json:"coachId"`
}

type GetDrillsForCoachResponse struct {
	Drills []*FullDrill `json:"drills"`
}

func GetDrillsForCoach(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetDrillsForCoachRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	drills, err := drillsDao.GetDrillsForCoach(request.CoachId)
	if err != nil {
		return converter.InternalServiceError("Error while getting drills for coach: %v. %v\n", request.CoachId, err)
	}

	fullDrills := getFullDrills(drills)

	return converter.Success(GetDrillsForCoachResponse{Drills: fullDrills})
}

func getFullDrills(drills []*drillsDao.Drill) []*FullDrill {
	fullDrills := make([]*FullDrill, 0)
	for _, drill := range drills {
		fullDrills = append(fullDrills, &FullDrill{
			DrillId:     drill.DrillId,
			CoachId:     drill.CoachId,
			Name:        drill.Name,
			Description: drill.Description,
			Category:    drill.Category,
			Equipment:   drill.Equipment,
			Demos: Demos{
				Front: Demo{FileLocation: files.GetDrillVideoFileLocation(drill.DrillId, files.Front).URL},
				Side:  Demo{FileLocation: files.GetDrillVideoFileLocation(drill.DrillId, files.Side).URL},
				Close: Demo{FileLocation: files.GetDrillVideoFileLocation(drill.DrillId, files.Close).URL},
			},
		})
	}

	return fullDrills
}
