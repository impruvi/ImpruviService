package auth

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token  string               `json:"token"`
	Player *playerFacade.Player `json:"player"`
}

func SignIn(request *SignInRequest) (*SignInResponse, error) {
	doesPasswordMatch, err := playerFacade.DoesPasswordMatch(request.Email, request.Password)
	if err != nil {
		return nil, err
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

	return &SignInResponse{
		Player: player,
		Token:  token,
	}, nil
}
