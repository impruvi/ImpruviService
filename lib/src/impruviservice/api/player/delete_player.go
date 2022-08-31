package player

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	"log"
)

type DeletePlayerRequest struct {
	PlayerId string `json:"playerId"`
}

func DeletePlayer(request *DeletePlayerRequest) error {
	log.Printf("DeletePlayerRequest: %+v\n", request)
	err := validateDeletePlayerRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid DeletePlayerRequest: %v\n", err)
		return err
	}

	return playerFacade.DeletePlayer(request.PlayerId)
}

func validateDeletePlayerRequest(request *DeletePlayerRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}

	return nil
}
