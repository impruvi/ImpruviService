package player

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	"log"
)

type GetPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetPlayerResponse struct {
	Player *playerFacade.Player `json:"player"`
}

func GetPlayer(request *GetPlayerRequest) (*GetPlayerResponse, error) {
	log.Printf("GetPlayerRequest: %+v\n", request)
	err := validateGetPlayerRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetPlayerRequest: %v\n", err)
		return nil, err
	}

	player, err := playerFacade.GetPlayerById(request.PlayerId)

	if err != nil {
		return nil, err
	}

	return &GetPlayerResponse{Player: player}, nil
}

func validateGetPlayerRequest(request *GetPlayerRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}

	return nil
}
