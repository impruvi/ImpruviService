package session

import (
	sessionFacade "impruviService/facade/session"
)

type DeleteSessionRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
}

func DeleteSession(request *DeleteSessionRequest) error {
	return sessionFacade.DeleteSession(request.SessionNumber, request.PlayerId)
}
