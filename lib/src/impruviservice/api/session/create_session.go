package session

import (
	sessionDao "impruviService/dao/session"
	sessionFacade "impruviService/facade/session"
)

type CreateSessionRequest struct {
	Session *sessionDao.SessionDB `json:"session"`
}

func CreateSession(request *CreateSessionRequest) error {
	return sessionFacade.CreateSession(request.Session)
}
