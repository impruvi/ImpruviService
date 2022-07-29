package auth

import (
	"fmt"
	sesAccessor "impruviService/accessor/ses"
	"impruviService/exceptions"
	passwordResetCodeFacade "impruviService/facade/passwordresetcode"
	"log"
)

type InitiatePasswordResetRequest struct {
	Email string `json:"email"`
}

func InitiatePasswordReset(request *InitiatePasswordResetRequest) error {
	log.Printf("InitiatePasswordResetRequest: %+v\n", request)
	err := validateInitiatePasswordResetRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid InitiatePasswordResetRequest: %v\n", err)
		return err
	}

	code, err := passwordResetCodeFacade.CreateCode(request.Email)
	if err != nil {
		return err
	}

	log.Printf("Created code: %v\n", code)
	return sesAccessor.SendEmail(
		request.Email,
		"Reset your Impruvi password",
		fmt.Sprintf("<div>Your verification code is %s</div>", code),
		fmt.Sprintf("Your verification code is %s", code))
}

func validateInitiatePasswordResetRequest(request *InitiatePasswordResetRequest) error {
	if request.Email == "" {
		return exceptions.InvalidRequestError{Message: "Email cannot be null/empty"}
	}
	return nil
}
