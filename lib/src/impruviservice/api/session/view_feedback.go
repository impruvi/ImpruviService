package session

import (
	"impruviService/exceptions"
	sessionFacade "impruviService/facade/session"
	"log"
)

type ViewFeedbackRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
}

func ViewFeedback(request *ViewFeedbackRequest) error {
	log.Printf("ViewFeedbackRequest: %+v\n", request)
	err := validateViewFeedbackRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid ViewFeedbackRequest: %v\n", err)
		return err
	}

	return sessionFacade.ViewFeedback(request.PlayerId, request.SessionNumber)
}

func validateViewFeedbackRequest(request *ViewFeedbackRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	if request.SessionNumber <= 0 {
		return exceptions.InvalidRequestError{Message: "Valid session number must be provided"}
	}

	return nil
}
