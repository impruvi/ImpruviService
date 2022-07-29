package session

import (
	"impruviService/exceptions"
	sessionFacade "impruviService/facade/session"
	"log"
)

type DeleteSessionRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
}

func DeleteSession(request *DeleteSessionRequest) error {
	log.Printf("DeleteSessionRequest: %+v\n", request)
	err := validateDeleteSessionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid DeleteSessionRequest: %v\n", err)
		return err
	}

	return sessionFacade.DeleteSession(request.SessionNumber, request.PlayerId)
}

func validateDeleteSessionRequest(request *DeleteSessionRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	if request.SessionNumber <= 0 {
		return exceptions.InvalidRequestError{Message: "Valid session number must be provided"}
	}

	return nil
}
