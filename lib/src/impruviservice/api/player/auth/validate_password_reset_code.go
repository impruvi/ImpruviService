package auth

import (
	"impruviService/exceptions"
	passwordResetCodeFacade "impruviService/facade/passwordresetcode"
	"log"
)

type ValidatePasswordResetCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func ValidatePasswordReset(request *ValidatePasswordResetCodeRequest) error {
	log.Printf("ValidatePasswordResetCodeRequest: %+v\n", request)
	err := validateValidatePasswordResetRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid ValidatePasswordResetCodeRequest: %v\n", err)
		return err
	}

	exists, err := passwordResetCodeFacade.Exists(request.Email, request.Code)
	if err != nil {
		return nil
	}
	if !exists {
		return exceptions.NotAuthorizedError{Message: "Invalid email/code combination"}
	}
	return nil
}

func validateValidatePasswordResetRequest(request *ValidatePasswordResetCodeRequest) error {
	if request.Email == "" {
		return exceptions.InvalidRequestError{Message: "Email cannot be null/empty"}
	}
	if request.Code == "" {
		return exceptions.InvalidRequestError{Message: "Code cannot be null/empty"}
	}

	return nil
}
