package player

import (
	playerFacade "impruviService/facade/player"
)

type UpdatePlayerRequest struct {
	Player *playerFacade.Player `json:"player"`
}

func UpdatePlayer(request *UpdatePlayerRequest) error {
	return playerFacade.UpdatePlayer(request.Player)
}
