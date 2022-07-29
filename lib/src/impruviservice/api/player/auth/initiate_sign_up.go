package auth

import (
	"fmt"
	sesAccessor "impruviService/accessor/ses"
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	"log"
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
	log.Printf("InitiateSignUpRequest: %+v\n", request)
	err := validateInitiateSignUpRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid InitiateSignUpRequest: %v\n", err)
		return nil, err
	}

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

func validateInitiateSignUpRequest(request *InitiateSignUpRequest) error {
	if request.Email == "" {
		return exceptions.InvalidRequestError{Message: "Email cannot be null/empty"}
	}
	if request.Password == "" {
		return exceptions.InvalidRequestError{Message: "Password cannot be null/empty"}
	}
	if request.FirstName == "" {
		return exceptions.InvalidRequestError{Message: "FirstName cannot be null/empty"}
	}
	if request.LastName == "" {
		return exceptions.InvalidRequestError{Message: "LastName cannot be null/empty"}
	}

	return nil
}
