package drills

import (
	"impruviService/exceptions"
	coachFacade "impruviService/facade/coach"
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

	drill, err := drillFacade.GetDrillById(request.DrillId)
	if err != nil {
		log.Printf("Unexpected error while getting drill: %v\n", drill)
		return err
	}
	coach, err := coachFacade.GetCoachById(drill.CoachId)
	if err != nil {
		log.Printf("Unexpected error while getting coach: %v\n", drill)
		return err
	}
	for _, introSessionDrill := range coach.IntroSessionDrills {
		if introSessionDrill.DrillId == drill.DrillId {
			return exceptions.InvalidRequestError{Message: "Intro session drills cannot be deleted"}
		}
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
