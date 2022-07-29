package player

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	"log"
)

type UpdatePlayerRequest struct {
	Player *playerFacade.Player `json:"player"`
}

func UpdatePlayer(request *UpdatePlayerRequest) error {
	log.Printf("UpdatePlayerRequest: %+v\n", request)
	err := validateUpdatePlayerRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid UpdatePlayerRequest: %v\n", err)
		return err
	}

	return playerFacade.UpdatePlayer(request.Player)
}

func validateUpdatePlayerRequest(request *UpdatePlayerRequest) error {
	if request.Player == nil {
		return exceptions.InvalidRequestError{Message: "Player cannot be null/empty"}
	}
	if request.Player.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	if request.Player.FirstName == "" {
		return exceptions.InvalidRequestError{Message: "FirstName cannot be null/empty"}
	}
	if request.Player.LastName == "" {
		return exceptions.InvalidRequestError{Message: "LastName cannot be null/empty"}
	}

	return nil
}
