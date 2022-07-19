package player

import (
	players "impruviService/dao/player"
)

func GetPlayerById(playerId string) (*players.PlayerDB, error) {
	player, err := players.GetPlayerById(playerId)
	if err != nil {
		return nil, err
	}
	return player, nil
}
