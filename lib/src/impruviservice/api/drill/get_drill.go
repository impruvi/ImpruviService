package drills

import (
	drillDao "impruviService/dao/drill"
	drillFacade "impruviService/facade/drill"
)

type GetDrillRequest struct {
	DrillId string `json:"drillId"`
}

type GetDrillResponse struct {
	Drill *drillDao.DrillDB `json:"drill"`
}

func GetDrill(request *GetDrillRequest) (*GetDrillResponse, error) {
	drill, err := drillFacade.GetDrillById(request.DrillId)
	if err != nil {
		return nil, err
	}

	return &GetDrillResponse{Drill: drill}, nil
}
