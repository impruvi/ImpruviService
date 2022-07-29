package coach

import (
	coachDao "impruviService/dao/coach"
	"impruviService/exceptions"
	coachFacade "impruviService/facade/coach"
	"log"
)

type GetCoachRequest struct {
	CoachId string `json:"coachId"`
}

type GetCoachResponse struct {
	Coach *coachDao.CoachDB `json:"coach"`
}

func GetCoach(request *GetCoachRequest) (*GetCoachResponse, error) {
	log.Printf("GetCoachRequest: %+v\n", request)
	err := validateGetCoachRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetCoachRequest: %v\n", err)
		return nil, err
	}

	coach, err := coachFacade.GetCoachById(request.CoachId)

	if err != nil {
		return nil, err
	}

	return &GetCoachResponse{Coach: coach}, nil
}

func validateGetCoachRequest(request *GetCoachRequest) error {
	if request.CoachId == "" {
		return exceptions.InvalidRequestError{Message: "CoachId cannot be null/empty"}
	}
	return nil
}
