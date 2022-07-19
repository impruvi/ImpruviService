package session

import (
	sessionFacade "impruviService/facade/session"
)

type CreateSubmissionRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
	FileLocation  string `json:"fileLocation"`
}

func CreateSubmission(request *CreateSubmissionRequest) error {
	return sessionFacade.CreateSubmission(request.PlayerId, request.SessionNumber, request.DrillId, request.FileLocation)
}
