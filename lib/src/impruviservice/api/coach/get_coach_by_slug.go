package coach

import (
	coachDao "impruviService/dao/coach"
	"impruviService/exceptions"
	coachFacade "impruviService/facade/coach"
	"log"
)

type GetCoachBySlugRequest struct {
	Slug string `json:"slug"`
}

type GetCoachBySlugResponse struct {
	Coach *coachDao.CoachDB `json:"coach"`
}

func GetCoachBySlug(request *GetCoachBySlugRequest) (*GetCoachBySlugResponse, error) {
	log.Printf("GetCoachBySlugRequest: %+v\n", request)
	err := validateGetCoachBySlugRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetCoachBySlugRequest: %v\n", err)
		return nil, err
	}

	coach, err := coachFacade.GetCoachBySlug(request.Slug)

	if err != nil {
		return nil, err
	}

	return &GetCoachBySlugResponse{Coach: coach}, nil
}

func validateGetCoachBySlugRequest(request *GetCoachBySlugRequest) error {
	if request.Slug == "" {
		return exceptions.InvalidRequestError{Message: "Slug cannot be null/empty"}
	}
	return nil
}
