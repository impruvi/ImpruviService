package coach

import (
	coachDao "impruviService/dao/coach"
	coachFacade "impruviService/facade/coach"
	"log"
)

type ListCoachesRequest struct {
	Limit int `json:"limit"`
}

type ListCoachesResponse struct {
	Coaches []*coachDao.CoachDB `json:"coaches"`
}

func ListCoaches(request *ListCoachesRequest) (*ListCoachesResponse, error) {
	log.Printf("ListCoachesRequest: %+v\n", request)
	coaches, err := coachFacade.ListCoaches(request.Limit)
	if err != nil {
		return nil, err
	}
	availableCoaches := make([]*coachDao.CoachDB, 0)
	for _, coach := range coaches {
		if len(coach.IntroSessionDrills) >= 4 {
			availableCoaches = append(availableCoaches, coach)
		} else {
			log.Printf("Coach: %v, does not have an intro session available.", coach.CoachId)
		}
	}

	return &ListCoachesResponse{Coaches: availableCoaches}, nil
}
