package auth

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	"log"
)

type CompleteSignUpRequest struct {
	PlayerId string `json:"playerId"`
	Code     string `json:"code"`
	ExpoPushToken string `json:"expoPushToken"`
}

type CompleteSignUpResponse struct {
	Token  string               `json:"token"`
	Player *playerFacade.Player `json:"player"`
}

func CompleteSignUp(request *CompleteSignUpRequest) (*CompleteSignUpResponse, error) {
	log.Printf("CompleteSignUpRequest: %+v\n", request)
	err := validateCompleteSignUpRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CompleteSignUpRequest: %v\n", err)
		return nil, err
	}

	player, err := playerFacade.ActivatePlayer(request.PlayerId, request.Code)
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

	return &CompleteSignUpResponse{
		Player: player,
		Token:  token,
	}, nil
}

func validateCompleteSignUpRequest(request *CompleteSignUpRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	if request.Code == "" {
		return exceptions.InvalidRequestError{Message: "Code cannot be null/empty"}
	}
	return nil
}
