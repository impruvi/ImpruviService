package session

import (
	"impruviService/exceptions"
	sessionFacade "impruviService/facade/session"
	"log"
)

type GetSessionsForCoachRequest struct {
	CoachId string `json:"coachId"`
}

type GetSessionsForCoachResponse struct {
	PlayerSessions []*sessionFacade.PlayerSessions `json:"playerSessions"`
}

func GetSessionsForCoach(request *GetSessionsForCoachRequest) (*GetSessionsForCoachResponse, error) {
	log.Printf("GetSessionsForCoachRequest: %+v\n", request)
	err := validateGetSessionForCoachRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetSessionsForCoachRequest: %v\n", err)
		return nil, err
	}

	playerSessions, err := sessionFacade.GetSessionsForCoach(request.CoachId)
	if err != nil {
		return nil, err
	}

	return &GetSessionsForCoachResponse{
		PlayerSessions: playerSessions,
	}, nil
}

func validateGetSessionForCoachRequest(request *GetSessionsForCoachRequest) error {
	if request.CoachId == "" {
		return exceptions.InvalidRequestError{Message: "CoachId cannot be null/empty"}
	}

	return nil
}
