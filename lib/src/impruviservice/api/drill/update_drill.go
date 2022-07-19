package drills

import (
	drillDao "impruviService/dao/drill"
	drillFacade "impruviService/facade/drill"
)

type UpdateDrillRequest struct {
	Drill *drillDao.DrillDB `json:"drill"`
}

func UpdateDrill(request *UpdateDrillRequest) error {
	return drillFacade.UpdateDrill(request.Drill)
}
