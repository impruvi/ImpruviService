package coach

import (
	coachDao "impruviService/dao/coach"
	coachFacade "impruviService/facade/coach"
)

type GetCoachRequest struct {
	CoachId string `json:"coachId"`
}

type GetCoachResponse struct {
	Coach *coachDao.CoachDB `json:"coach"`
}

func GetCoach(request *GetCoachRequest) (*GetCoachResponse, error) {
	coach, err := coachFacade.GetCoachById(request.CoachId)

	if err != nil {
		return nil, err
	}

	return &GetCoachResponse{Coach: coach}, nil
}
