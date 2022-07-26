package auth

import (
	"impruviService/exceptions"
	passwordResetCodeFacade "impruviService/facade/passwordresetcode"
)

type ValidatePasswordResetCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func ValidatePasswordReset(request *ValidatePasswordResetCodeRequest) error {
	exists, err := passwordResetCodeFacade.Exists(request.Email, request.Code)
	if err != nil {
		return nil
	}
	if !exists {
		return exceptions.NotAuthorizedError{Message: "Invalid email/code combination"}
	}
	return nil
}
