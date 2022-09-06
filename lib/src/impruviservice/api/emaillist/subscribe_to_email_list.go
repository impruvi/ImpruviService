package emaillist

import (
	"impruviService/exceptions"
	emailListFacade "impruviService/facade/emaillist"
	"log"
)

type SubscribeToEmailListRequest struct {
	Email string `json:"email"`
}

func SubscribeToEmailList(request *SubscribeToEmailListRequest) error {
	log.Printf("SubscribeToEmailListRequest: %+v\n", request)
	err := validateSubscribeToEmailListRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid SubscribeToEmailListRequest: %v\n", err)
		return err
	}

	return emailListFacade.Subscribe(request.Email)
}

func validateSubscribeToEmailListRequest(request *SubscribeToEmailListRequest) error {
	if request.Email == "" {
		return exceptions.InvalidRequestError{Message: "Email cannot be null/empty"}
	}

	return nil
}
