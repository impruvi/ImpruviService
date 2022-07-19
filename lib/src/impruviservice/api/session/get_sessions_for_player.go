package session

import (
	sessionFacade "impruviService/facade/session"
)

type GetPlayerSessionsRequest struct {
	PlayerId string `json:"playerId"`
}

type GetPlayerSessionsResponse struct {
	Sessions []*sessionFacade.Session `json:"sessions"`
}

func GetSessionsForPlayer(request *GetPlayerSessionsRequest) (*GetPlayerSessionsResponse, error) {
	sessions, err := sessionFacade.GetSessionsForPlayer(request.PlayerId)
	if err != nil {
		return nil, err
	}

	return &GetPlayerSessionsResponse{
		Sessions: sessions,
	}, nil
}
