package drills

import (
	"impruviService/exceptions"
	drillFacade "impruviService/facade/drill"
	"log"
)

type DeleteDrillRequest struct {
	DrillId string `json:"drillId"`
}

func DeleteDrill(request *DeleteDrillRequest) error {
	log.Printf("DeleteDrillRequest: %+v\n", request.DrillId)
	err := validateDeleteDrillRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid DeleteDrillRequest: %v\n", err)
		return err
	}

	err = drillFacade.DeleteDrill(request.DrillId)
	if err != nil {
		log.Printf("Unexpected error while deleting drill: %v\n", err)
	}
	return err
}

func validateDeleteDrillRequest(request *DeleteDrillRequest) error {
	if request.DrillId == "" {
		return exceptions.InvalidRequestError{Message: "DrillId cannot be null/empty"}
	}
	return nil
}
