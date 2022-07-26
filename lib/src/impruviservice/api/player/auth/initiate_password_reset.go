package auth

import (
	"fmt"
	sesAccessor "impruviService/accessor/ses"
	passwordResetCodeFacade "impruviService/facade/passwordresetcode"
	"log"
)

type InitiatePasswordResetRequest struct {
	Email string `json:"email"`
}

func InitiatePasswordReset(request *InitiatePasswordResetRequest) error {
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
