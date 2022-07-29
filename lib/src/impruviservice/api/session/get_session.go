package session

import (
	sessionDao "impruviService/dao/session"
	"impruviService/exceptions"
	"log"
)

type GetSessionRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
}

type GetSessionResponse struct {
	Session *sessionDao.SessionDB `json:"session"`
}

func GetSession(request *GetSessionRequest) (*GetSessionResponse, error) {
	log.Printf("GetSessionRequest: %+v\n", request)
	err := validateGetSessionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetSessionRequest: %v\n", err)
		return nil, err
	}

	session, err := sessionDao.GetSession(request.PlayerId, request.SessionNumber)
	if err != nil {
		return nil, err
	}

	return &GetSessionResponse{Session: session}, nil
}

func validateGetSessionRequest(request *GetSessionRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	if request.SessionNumber <= 0 {
		return exceptions.InvalidRequestError{Message: "Valid session number must be provided"}
	}

	return nil
}
