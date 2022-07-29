package drills

import (
	drillDao "impruviService/dao/drill"
	"impruviService/exceptions"
	drillFacade "impruviService/facade/drill"
	"log"
)

type UpdateDrillRequest struct {
	Drill *drillDao.DrillDB `json:"drill"`
}

func UpdateDrill(request *UpdateDrillRequest) error {
	log.Printf("GetDrillsForPlayerRequest: %+v\n", request)
	err := validateUpdateDrillRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid UpdateDrillRequest: %v\n", err)
		return err
	}

	return drillFacade.UpdateDrill(request.Drill)
}

func validateUpdateDrillRequest(request *UpdateDrillRequest) error {
	if request.Drill == nil {
		return exceptions.InvalidRequestError{Message: "Drill cannot be null/empty"}
	}
	if request.Drill.DrillId == "" {
		return exceptions.InvalidRequestError{Message: "DrillId cannot be null/empty"}
	}
	if request.Drill.Name == "" {
		return exceptions.InvalidRequestError{Message: "Name cannot be null/empty"}
	}
	if request.Drill.CoachId == "" {
		return exceptions.InvalidRequestError{Message: "CoachId cannot be null/empty"}
	}
	if request.Drill.Category == "" {
		return exceptions.InvalidRequestError{Message: "Category cannot be null/empty"}
	}
	if request.Drill.Description == "" {
		return exceptions.InvalidRequestError{Message: "Description cannot be null/empty"}
	}
	if request.Drill.Equipment == nil {
		return exceptions.InvalidRequestError{Message: "Equipment cannot be null"}
	}

	return nil
}
