package drills

import (
	drillsDao "impruviService/dao/drill"
	"impruviService/exceptions"
	drillFacade "impruviService/facade/drill"
	"log"
)

type GetDrillsForPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetDrillsForPlayerResponse struct {
	Drills []*drillsDao.DrillDB `json:"drills"`
}

func GetDrillsForPlayer(request *GetDrillsForPlayerRequest) (*GetDrillsForPlayerResponse, error) {
	log.Printf("GetDrillsForPlayerRequest: %+v\n", request)
	err := validateGetDrillsForPlayerRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetDrillsForPlayerRequest: %v\n", err)
		return nil, err
	}

	drills, err := drillFacade.GetDrillsForPlayer(request.PlayerId)
	if err != nil {
		return nil, err
	}

	return &GetDrillsForPlayerResponse{Drills: drills}, nil
}

func validateGetDrillsForPlayerRequest(request *GetDrillsForPlayerRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	return nil
}
