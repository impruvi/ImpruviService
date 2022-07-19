package player

import (
	playerDao "impruviService/dao/player"
	"log"
)

func UpdatePlayer(player *playerDao.PlayerDB) error {
	currentPlayer, err := playerDao.GetPlayerById(player.PlayerId)
	if err != nil {
		return err
	}
	player.NotificationId = currentPlayer.NotificationId

	err = playerDao.PutPlayer(player)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePlayerNotificationId(player *playerDao.PlayerDB, notificationId string) *playerDao.PlayerDB {
	var newPlayer = player
	newPlayer.NotificationId = notificationId
	err := playerDao.PutPlayer(player)
	if err != nil {
		// don't error if we can't update push notification id
		log.Printf("Error while updating player notification id: %v\n", err)
		return player
	}
	log.Printf("Updated player id's %v push notification token", player.PlayerId)
	return newPlayer
}
