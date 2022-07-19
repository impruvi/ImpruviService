package drills

import (
	drillFacade "impruviService/facade/drill"
	"log"
)

type DeleteDrillRequest struct {
	DrillId string `json:"drillId"`
}

func DeleteDrill(request *DeleteDrillRequest) error {
	log.Printf("Deleting drill: %v\n", request.DrillId)
	err := drillFacade.DeleteDrill(request.DrillId)
	if err != nil {
		log.Printf("Unexpected error while deleting drill: %v\n", err)
	}
	return err
}
