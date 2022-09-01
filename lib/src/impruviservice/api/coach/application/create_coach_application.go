package application

import (
	coachApplicationDao "impruviService/dao/coachapplication"
	"impruviService/exceptions"
	coachApplicationFacade "impruviService/facade/coachapplication"
	"log"
)

type CreateCoachApplicationRequest struct {
	Application *coachApplicationDao.CoachApplicationDB `json:"application"`
}

func CreateCoachApplication(request *CreateCoachApplicationRequest) error {
	log.Printf("CreateCoachApplicationRequest: %+v\n", request)
	err := validateCreateCoachApplicationRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CreateCoachApplicationRequest: %v\n", err)
		return err
	}
	return coachApplicationFacade.CreateApplication(request.Application)
}

func validateCreateCoachApplicationRequest(request *CreateCoachApplicationRequest) error {
	if request.Application == nil {
		return exceptions.InvalidRequestError{Message: "Application cannot be null/empty"}
	}
	if request.Application.Email == "" {
		return exceptions.InvalidRequestError{Message: "Email cannot be null/empty"}
	}
	if request.Application.FirstName == "" {
		return exceptions.InvalidRequestError{Message: "FirstName cannot be null/empty"}
	}
	if request.Application.LastName == "" {
		return exceptions.InvalidRequestError{Message: "LastName cannot be null/empty"}
	}
	if request.Application.Experience == "" {
		return exceptions.InvalidRequestError{Message: "Experience cannot be null/empty"}
	}
	if request.Application.PhoneNumber == "" {
		return exceptions.InvalidRequestError{Message: "PhoneNumber cannot be null/empty"}
	}

	return nil
}
