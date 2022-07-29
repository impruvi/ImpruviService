package session

import (
	"impruviService/exceptions"
	sessionFacade "impruviService/facade/session"
	"log"
)

type CreateFeedbackRequest struct {
	PlayerId              string `json:"playerId"`
	SessionNumber         int    `json:"sessionNumber"`
	DrillId               string `json:"drillId"`
	FileLocation          string `json:"fileLocation"`
	ThumbnailFileLocation string `json:"thumbnailFileLocation"`
}

func CreateFeedback(request *CreateFeedbackRequest) error {
	log.Printf("CreateFeedbackRequest: %+v\n", request)
	err := validateCreateFeedbackRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CreateFeedbackRequest: %v\n", err)
		return err
	}

	return sessionFacade.CreateFeedback(request.PlayerId, request.SessionNumber, request.DrillId, request.FileLocation, request.ThumbnailFileLocation)
}

func validateCreateFeedbackRequest(request *CreateFeedbackRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	if request.DrillId == "" {
		return exceptions.InvalidRequestError{Message: "DrillId cannot be null/empty"}
	}
	if request.SessionNumber <= 0 {
		return exceptions.InvalidRequestError{Message: "Valid session number must be provided"}
	}
	if request.FileLocation == "" {
		return exceptions.InvalidRequestError{Message: "FileLocation cannot be null/empty"}
	}
	if request.ThumbnailFileLocation == "" {
		return exceptions.InvalidRequestError{Message: "ThumbnailFileLocation cannot be null/empty"}
	}

	return nil
}
