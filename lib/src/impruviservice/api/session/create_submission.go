package session

import (
	"impruviService/exceptions"
	sessionFacade "impruviService/facade/session"
	"log"
)

type CreateSubmissionRequest struct {
	PlayerId              string `json:"playerId"`
	SessionNumber         int    `json:"sessionNumber"`
	DrillId               string `json:"drillId"`
	FileLocation          string `json:"fileLocation"`
	ThumbnailFileLocation string `json:"thumbnailFileLocation"`
}

func CreateSubmission(request *CreateSubmissionRequest) error {
	log.Printf("CreateSubmissionRequest: %+v\n", request)
	err := validateCreateSubmissionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CreateSubmissionRequest: %v\n", err)
		return err
	}

	return sessionFacade.CreateSubmission(request.PlayerId, request.SessionNumber, request.DrillId, request.FileLocation, request.ThumbnailFileLocation)
}

func validateCreateSubmissionRequest(request *CreateSubmissionRequest) error {
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
