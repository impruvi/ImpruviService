package session

import (
	sessionDao "impruviService/dao/session"
	"impruviService/exceptions"
	sessionFacade "impruviService/facade/session"
	"log"
)

type UpdateSessionRequest struct {
	Session *sessionDao.SessionDB `json:"session"`
}

func UpdateSession(request *UpdateSessionRequest) error {
	log.Printf("UpdateSessionRequest: %+v\n", request)
	err := validateUpdateSessionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid UpdateSessionRequest: %v\n", err)
		return err
	}

	return sessionFacade.UpdateSession(request.Session)
}

func validateUpdateSessionRequest(request *UpdateSessionRequest) error {
	if request.Session == nil {
		return exceptions.InvalidRequestError{Message: "Session cannot be null/empty"}
	}
	if request.Session.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	if request.Session.SessionNumber <= 0 {
		return exceptions.InvalidRequestError{Message: "Valid session number must be provided"}
	}
	if request.Session.Drills == nil || len(request.Session.Drills) < 4 {
		return exceptions.InvalidRequestError{Message: "You must provide at least 4 drills"}
	} else {
		for _, drill := range request.Session.Drills {
			if drill.DrillId == "" {
				return exceptions.InvalidRequestError{Message: "DrillId cannot be null/empty"}
			}
		}
	}

	return nil
}
