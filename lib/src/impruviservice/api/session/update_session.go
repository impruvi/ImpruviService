package session

import (
	sessionDao "impruviService/dao/session"
	sessionFacade "impruviService/facade/session"
)

type UpdateSessionRequest struct {
	Session *sessionDao.SessionDB `json:"session"`
}

func UpdateSession(request *UpdateSessionRequest) error {
	return sessionFacade.UpdateSession(request.Session)
}
