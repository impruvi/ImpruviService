package session

import (
	sessionFacade "impruviService/facade/session"
)

type CreateFeedbackRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
	FileLocation  string `json:"fileLocation"`
}

func CreateFeedback(request *CreateFeedbackRequest) error {
	return sessionFacade.CreateFeedback(request.PlayerId, request.SessionNumber, request.DrillId, request.FileLocation)
}
