package drills

import (
	drillsDao "impruviService/dao/drill"
	drillFacade "impruviService/facade/drill"
)

type GetDrillsForCoachRequest struct {
	CoachId string `json:"coachId"`
}

type GetDrillsForCoachResponse struct {
	Drills []*drillsDao.DrillDB `json:"drills"`
}

func GetDrillsForCoach(request *GetDrillsForCoachRequest) (*GetDrillsForCoachResponse, error) {
	drills, err := drillFacade.GetDrillsForCoach(request.CoachId)
	if err != nil {
		return nil, err
	}

	return &GetDrillsForCoachResponse{Drills: drills}, nil
}
