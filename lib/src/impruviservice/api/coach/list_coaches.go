package coach

import (
	coachDao "impruviService/dao/coach"
	coachFacade "impruviService/facade/coach"
)

type ListCoachesRequest struct {
	Limit int `json:"limit"`
}

type ListCoachesResponse struct {
	Coaches []*coachDao.CoachDB `json:"coaches"`
}

func ListCoaches(request *ListCoachesRequest) (*ListCoachesResponse, error) {
	coaches, err := coachFacade.ListCoaches(request.Limit)
	if err != nil {
		return nil, err
	}

	return &ListCoachesResponse{Coaches: coaches}, nil
}
