package drills

import (
	drillDao "impruviService/dao/drill"
	"impruviService/exceptions"
	drillFacade "impruviService/facade/drill"
	"log"
)

type GetDrillRequest struct {
	DrillId string `json:"drillId"`
}

type GetDrillResponse struct {
	Drill *drillDao.DrillDB `json:"drill"`
}

func GetDrill(request *GetDrillRequest) (*GetDrillResponse, error) {
	log.Printf("GetDrillRequest: %+v\n", request)
	err := validateGetDrillRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetDrillRequest: %v\n", err)
		return nil, err
	}

	drill, err := drillFacade.GetDrillById(request.DrillId)
	if err != nil {
		return nil, err
	}

	return &GetDrillResponse{Drill: drill}, nil
}

func validateGetDrillRequest(request *GetDrillRequest) error {
	if request.DrillId == "" {
		return exceptions.InvalidRequestError{Message: "DrillId cannot be null/empty"}
	}
	return nil
}
