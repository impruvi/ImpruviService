package drills

import (
	drillDao "impruviService/dao/drill"
	"impruviService/exceptions"
	drillFacade "impruviService/facade/drill"
	"log"
)

type CreateDrillRequest struct {
	Drill *drillDao.DrillDB `json:"drill"`
}

type CreateDrillResponse struct {
	Drill *drillDao.DrillDB `json:"drill"`
}

func CreateDrill(request *CreateDrillRequest) (*CreateDrillResponse, error) {
	log.Printf("CreateDrillRequest: %+v\n", request)
	err := validateCreateDrillRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CreateDrillRequest: %v\n", err)
		return nil, err
	}
	drill, err := drillFacade.CreateDrill(request.Drill)
	if err != nil {
		return nil, err
	}
	return &CreateDrillResponse{Drill: drill}, nil
}

func validateCreateDrillRequest(request *CreateDrillRequest) error {
	if request.Drill == nil {
		return exceptions.InvalidRequestError{Message: "Drill cannot be null/empty"}
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
