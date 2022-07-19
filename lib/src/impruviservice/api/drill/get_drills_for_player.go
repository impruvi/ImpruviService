package drills

import (
	drillsDao "impruviService/dao/drill"
	drillFacade "impruviService/facade/drill"
)

type GetDrillsForPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetDrillsForPlayerResponse struct {
	Drills []*drillsDao.DrillDB `json:"drills"`
}

func GetDrillsForPlayer(request *GetDrillsForPlayerRequest) (*GetDrillsForPlayerResponse, error) {
	drills, err := drillFacade.GetDrillsForPlayer(request.PlayerId)
	if err != nil {
		return nil, err
	}

	return &GetDrillsForPlayerResponse{Drills: drills}, nil
}
