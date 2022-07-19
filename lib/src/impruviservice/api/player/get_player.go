package player

import (
	"impruviService/dao/player"
	playerFacade "impruviService/facade/player"
)

type GetPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetPlayerResponse struct {
	Player *players.PlayerDB `json:"player"`
}

func GetPlayer(request *GetPlayerRequest) (*GetPlayerResponse, error) {
	player, err := playerFacade.GetPlayerById(request.PlayerId)

	if err != nil {
		return nil, err
	}

	return &GetPlayerResponse{Player: player}, nil
}
