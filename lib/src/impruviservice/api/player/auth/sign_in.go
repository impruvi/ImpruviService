package auth

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	"log"
)

type SignInRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	ExpoPushToken string `json:"expoPushToken"`
}

type SignInResponse struct {
	Token  string               `json:"token"`
	Player *playerFacade.Player `json:"player"`
}

func SignIn(request *SignInRequest) (*SignInResponse, error) {
	log.Printf("SignInRequest: %+v\n", request)
	err := validateSignInRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid SignInRequest: %v\n", err)
		return nil, err
	}

	doesPasswordMatch, err := playerFacade.DoesPasswordMatch(request.Email, request.Password)
	if err != nil {
		if _, ok := err.(exceptions.ResourceNotFoundError); ok {
			return nil, exceptions.NotAuthorizedError{Message: "Invalid email/password combination"}
		} else {
			return nil, err
		}
	}

	if !doesPasswordMatch {
		return nil, exceptions.NotAuthorizedError{Message: "Invalid email/password combination"}
	}

	player, err := playerFacade.GetPlayerByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	token, err := playerFacade.GenerateToken(player.PlayerId)
	if err != nil {
		return nil, err
	}
	if player.NotificationId != request.ExpoPushToken {
		player, err = playerFacade.UpdateNotificationId(player.PlayerId, request.ExpoPushToken)
		if err != nil {
			return nil, err
		}
	}

	return &SignInResponse{
		Player: player,
		Token:  token,
	}, nil
}

func validateSignInRequest(request *SignInRequest) error {
	if request.Email == "" {
		return exceptions.InvalidRequestError{Message: "Email cannot be null/empty"}
	}
	if request.Password == "" {
		return exceptions.InvalidRequestError{Message: "Password cannot be null/empty"}
	}

	return nil
}
