package session

import (
	"impruviService/exceptions"
	sessionFacade "impruviService/facade/session"
	"log"
)

type GetSessionsForPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetSessionsForPlayerResponse struct {
	Sessions []*sessionFacade.Session `json:"sessions"`
}

func GetSessionsForPlayer(request *GetSessionsForPlayerRequest) (*GetSessionsForPlayerResponse, error) {
	log.Printf("GetSessionsForPlayerRequest: %+v\n", request)
	err := validateGetSessionsForPlayerRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetSessionsForPlayerRequest: %v\n", err)
		return nil, err
	}

	sessions, err := sessionFacade.GetSessionsForPlayer(request.PlayerId)
	if err != nil {
		return nil, err
	}

	return &GetSessionsForPlayerResponse{
		Sessions: sessions,
	}, nil
}

func validateGetSessionsForPlayerRequest(request *GetSessionsForPlayerRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}

	return nil
}
