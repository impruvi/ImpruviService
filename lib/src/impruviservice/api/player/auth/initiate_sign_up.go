package auth

import (
	"fmt"
	sesAccessor "impruviService/accessor/ses"
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
)

type InitiateSignUpRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type InitiateSignUpResponse struct {
	PlayerId string `json:"playerId"`
}

func InitiateSignUp(request *InitiateSignUpRequest) (*InitiateSignUpResponse, error) {
	playerId, activationCode, err := playerFacade.CreatePlayer(&playerFacade.Player{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
	}, request.Password)
	if err != nil {
		if _, ok := err.(exceptions.ResourceAlreadyExistsError); ok {
			return nil, exceptions.InvalidRequestError{Message: "Invalid email/password combination"}
		} else {
			return nil, err
		}
	}

	err = sesAccessor.SendEmail(
		request.Email,
		"Complete your Impruvi sign up",
		fmt.Sprintf("<div>Your verification code is %s</div>", activationCode),
		fmt.Sprintf("Your verification code is %s", activationCode))

	if err != nil {
		return nil, err
	}

	return &InitiateSignUpResponse{
		PlayerId: playerId,
	}, nil
}
