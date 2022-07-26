package player

import (
	"fmt"
	"github.com/google/uuid"
	playerDao "impruviService/dao/player"
	"impruviService/exceptions"
	"impruviService/util"
	"log"
)

func CreatePlayer(player *Player, password string) (string, string, error) {
	doesPlayerAlreadyExistWithEmail, isActive, err := doesPlayerWithEmailExist(player.Email)
	if err != nil {
		return "", "", err
	}
	if doesPlayerAlreadyExistWithEmail && isActive {
		return "", "", exceptions.ResourceAlreadyExistsError{Message: fmt.Sprintf("Player with email: %v already exists\n", player.Email)}
	}

	log.Printf("Active player with email: %v does not already exist. creating...", player.Email)

	playerDB := reverseConvert(player)
	if !doesPlayerAlreadyExistWithEmail {
		playerDB.PlayerId = uuid.New().String()
	} else {
		currentPlayerDB, err := GetPlayerByEmail(player.Email)
		if err != nil {
			return "", "", err
		}
		playerDB.PlayerId = currentPlayerDB.PlayerId
	}

	currentTimeMillis := util.GetCurrentTimeEpochMillis()
	playerDB.Password = password
	playerDB.IsActive = false
	playerDB.ActivationCode = util.GenerateVerificationCode()
	playerDB.CreationDateEpochMillis = currentTimeMillis
	playerDB.LastUpdatedDateEpochMillis = currentTimeMillis

	log.Printf("Player: %+v\n...", player)

	err = playerDao.PutPlayer(playerDB)
	if err != nil {
		return "", "", err
	}
	return playerDB.PlayerId, playerDB.ActivationCode, nil
}

func ActivatePlayer(playerId, code string) (*Player, error) {
	playerDB, err := playerDao.GetPlayerById(playerId)
	if err != nil {
		return nil, err
	}
	if playerDB.ActivationCode != code {
		return nil, exceptions.NotAuthorizedError{Message: "Invalid activation code"}
	}

	playerDB.IsActive = true
	playerDB.LastUpdatedDateEpochMillis = util.GetCurrentTimeEpochMillis()
	err = playerDao.PutPlayer(playerDB)
	if err != nil {
		return nil, err
	}
	return convert(playerDB)
}

func UpdatePlayer(player *Player) error {
	currentPlayer, err := playerDao.GetPlayerById(player.PlayerId)
	if err != nil {
		return err
	}
	newPlayer := reverseConvert(player)
	newPlayer.NotificationId = currentPlayer.NotificationId
	newPlayer.Password = currentPlayer.Password
	newPlayer.LastUpdatedDateEpochMillis = util.GetCurrentTimeEpochMillis()
	newPlayer.IsActive = currentPlayer.IsActive

	err = playerDao.PutPlayer(newPlayer)
	if err != nil {
		return err
	}

	return nil
}

func UpdateNotificationId(playerId, notificationId string) (*Player, error) {
	playerDB, err := playerDao.GetPlayerById(playerId)
	if err != nil {
		return nil, err
	}
	playerDB.NotificationId = notificationId
	playerDB.LastUpdatedDateEpochMillis = util.GetCurrentTimeEpochMillis()
	err = playerDao.PutPlayer(playerDB)
	if err != nil {
		log.Printf("Error while updating player notification id: %v\n", err)
		return nil, err
	}
	return convert(playerDB)
}

func UpdatePassword(playerId string, password string) (*Player, error) {
	playerDB, err := playerDao.GetPlayerById(playerId)
	if err != nil {
		return nil, err
	}
	playerDB.Password = password
	playerDB.LastUpdatedDateEpochMillis = util.GetCurrentTimeEpochMillis()
	playerDB.IsActive = true
	err = playerDao.PutPlayer(playerDB)
	if err != nil {
		log.Printf("Error while updating player password: %v\n", err)
		return nil, err
	}
	return convert(playerDB)
}

func doesPlayerWithEmailExist(email string) (bool, bool, error) {
	player, err := playerDao.GetPlayerByEmail(email)
	if err != nil {
		if _, ok := err.(exceptions.ResourceNotFoundError); ok {
			return false, false, nil
		} else {
			return false, false, err
		}
	}

	return true, player.IsActive, nil
}
