package auth

import (
	"impruviService/exceptions"
	passwordResetCodeFacade "impruviService/facade/passwordresetcode"
	playerFacade "impruviService/facade/player"
	"log"
)

type CompletePasswordResetRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func CompletePasswordReset(request *CompletePasswordResetRequest) error {
	log.Printf("CompletePasswordResetRequest: %+v\n", request)
	err := validateCompletePasswordResetRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CompletePasswordResetRequest: %v\n", err)
		return err
	}

	exists, err := passwordResetCodeFacade.Exists(request.Email, request.Code)
	if err != nil {
		return err
	}

	if !exists {
		return exceptions.NotAuthorizedError{Message: "Invalid email/code combination for password reset"}
	}

	player, err := playerFacade.GetPlayerByEmail(request.Email)
	if err != nil {
		return err
	}
	_, err = playerFacade.UpdatePassword(player.PlayerId, request.Password)
	return err
}

func validateCompletePasswordResetRequest(request *CompletePasswordResetRequest) error {
	if request.Email == "" {
		return exceptions.InvalidRequestError{Message: "Email cannot be null/empty"}
	}
	if request.Code == "" {
		return exceptions.InvalidRequestError{Message: "Code cannot be null/empty"}
	}
	if request.Password == "" {
		return exceptions.InvalidRequestError{Message: "Password cannot be null/empty"}
	}
	return nil
}
