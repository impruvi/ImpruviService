package drills

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	drillsDao "impruviService/dao/drills"
	"impruviService/files"
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
				Front:          Media{FileLocation: files.GetDemoVideoFileLocation(drill.DrillId, files.Front).URL},
				Side:           Media{FileLocation: files.GetDemoVideoFileLocation(drill.DrillId, files.Side).URL},
				Close:          Media{FileLocation: files.GetDemoVideoFileLocation(drill.DrillId, files.Close).URL},
				FrontThumbnail: Media{FileLocation: files.GetDemoVideoThumbnailFileLocation(drill.DrillId, files.Front).URL},
				SideThumbnail:  Media{FileLocation: files.GetDemoVideoThumbnailFileLocation(drill.DrillId, files.Side).URL},
				CloseThumbnail: Media{FileLocation: files.GetDemoVideoThumbnailFileLocation(drill.DrillId, files.Close).URL},
			},
		})
	}

	return fullDrills
}
