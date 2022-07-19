package drills

import (
	drillDao "impruviService/dao/drill"
	drillFacade "impruviService/facade/drill"
)

type CreateDrillRequest struct {
	Drill *drillDao.DrillDB `json:"drill"`
}

type CreateDrillResponse struct {
	Drill *drillDao.DrillDB `json:"drill"`
}

func CreateDrill(request *CreateDrillRequest) (*CreateDrillResponse, error) {
	drill, err := drillFacade.CreateDrill(request.Drill)
	if err != nil {
		return nil, err
	}
	return &CreateDrillResponse{Drill: drill}, nil
}
