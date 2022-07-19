package session

import (
	sessionFacade "impruviService/facade/session"
)

type GetCoachSessionsRequest struct {
	CoachId string `json:"coachId"`
}

type GetCoachSessionsResponse struct {
	PlayerSessions []*sessionFacade.PlayerSessions `json:"playerSessions"`
}

func GetSessionForCoach(request *GetCoachSessionsRequest) (*GetCoachSessionsResponse, error) {
	playerSessions, err := sessionFacade.GetSessionsForCoach(request.CoachId)
	if err != nil {
		return nil, err
	}

	return &GetCoachSessionsResponse{
		PlayerSessions: playerSessions,
	}, nil
}
