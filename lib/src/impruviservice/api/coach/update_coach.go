package coach

import (
	coachDao "impruviService/dao/coach"
	coachFacade "impruviService/facade/coach"
)

type UpdateCoachRequest struct {
	Coach *coachDao.CoachDB `json:"coach"`
}

func UpdateCoach(request *UpdateCoachRequest) error {
	return coachFacade.UpdateCoach(request.Coach)
}
