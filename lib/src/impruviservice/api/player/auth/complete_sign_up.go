package auth

import (
	playerFacade "impruviService/facade/player"
)

type CompleteSignUpRequest struct {
	PlayerId string `json:"playerId"`
	Code     string `json:"code"`
}

type CompleteSignUpResponse struct {
	Token  string               `json:"token"`
	Player *playerFacade.Player `json:"player"`
}

func CompleteSignUp(request *CompleteSignUpRequest) (*CompleteSignUpResponse, error) {
	player, err := playerFacade.ActivatePlayer(request.PlayerId, request.Code)
	if err != nil {
		return nil, err
	}

	token, err := playerFacade.GenerateToken(player.PlayerId)
	if err != nil {
		return nil, err
	}

	return &CompleteSignUpResponse{
		Player: player,
		Token:  token,
	}, nil
}
