package auth

import (
	"impruviService/exceptions"
	passwordResetCodeFacade "impruviService/facade/passwordresetcode"
	playerFacade "impruviService/facade/player"
)

type CompletePasswordResetRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func CompletePasswordReset(request *CompletePasswordResetRequest) error {
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
