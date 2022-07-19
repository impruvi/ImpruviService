package player

import (
	"impruviService/dao/player"
	playerFacade "impruviService/facade/player"
)

type UpdatePlayerRequest struct {
	Player *players.PlayerDB `json:"player"`
}

func UpdatePlayer(request *UpdatePlayerRequest) error {
	return playerFacade.UpdatePlayer(request.Player)
}
