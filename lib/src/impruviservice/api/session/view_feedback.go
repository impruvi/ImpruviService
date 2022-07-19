package session

import (
	sessionFacade "impruviService/facade/session"
)

type ViewFeedbackRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
}

func ViewFeedback(request *ViewFeedbackRequest) error {
	return sessionFacade.ViewFeedback(request.PlayerId, request.SessionNumber)
}
