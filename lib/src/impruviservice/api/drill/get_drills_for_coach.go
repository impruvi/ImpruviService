package drills

import (
	drillsDao "impruviService/dao/drill"
	"impruviService/exceptions"
	drillFacade "impruviService/facade/drill"
	"log"
)

type GetDrillsForCoachRequest struct {
	CoachId string `json:"coachId"`
}

type GetDrillsForCoachResponse struct {
	Drills []*drillsDao.DrillDB `json:"drills"`
}

func GetDrillsForCoach(request *GetDrillsForCoachRequest) (*GetDrillsForCoachResponse, error) {
	log.Printf("GetDrillsForCoachRequest: %+v\n", request)
	err := validateGetDrillsForCoachRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetDrillsForCoachRequest: %v\n", err)
		return nil, err
	}

	drills, err := drillFacade.GetDrillsForCoach(request.CoachId)
	if err != nil {
		return nil, err
	}

	return &GetDrillsForCoachResponse{Drills: drills}, nil
}

func validateGetDrillsForCoachRequest(request *GetDrillsForCoachRequest) error {
	if request.CoachId == "" {
		return exceptions.InvalidRequestError{Message: "CoachId cannot be null/empty"}
	}
	return nil
}
