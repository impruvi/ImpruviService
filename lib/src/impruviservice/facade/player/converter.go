package player

import (
	playerDao "impruviService/dao/player"
)

func convertAll(playerDBs []*playerDao.PlayerDB) []*Player {
	players := make([]*Player, 0)
	for _, playerDB := range playerDBs {
		players = append(players, convert(playerDB))
	}
	return players
}

func convert(playerDB *playerDao.PlayerDB) *Player {
	return &Player{
		PlayerId:                   playerDB.PlayerId,
		CoachId:                    playerDB.CoachId,
		StripeCustomerId:           playerDB.StripeCustomerId,
		FirstName:                  playerDB.FirstName,
		LastName:                   playerDB.LastName,
		Email:                      playerDB.Email,
		Headshot:                   playerDB.Headshot,
		Position:                   playerDB.Position,
		AgeRange:                   playerDB.AgeRange,
		AvailableEquipment:         playerDB.AvailableEquipment,
		AvailableTrainingLocations: playerDB.AvailableTrainingLocations,
		ShortTermGoal:              playerDB.ShortTermGoal,
		LongTermGoal:               playerDB.LongTermGoal,
		CreationDateEpochMillis:    playerDB.CreationDateEpochMillis,
		LastUpdatedDateEpochMillis: playerDB.LastUpdatedDateEpochMillis,
		NotificationId:             playerDB.NotificationId,
		QueuedSubscription:         playerDB.QueuedSubscription,
	}
}

func reverseConvert(player *Player) *playerDao.PlayerDB {
	return &playerDao.PlayerDB{
		PlayerId:                   player.PlayerId,
		CoachId:                    player.CoachId,
		StripeCustomerId:           player.StripeCustomerId,
		FirstName:                  player.FirstName,
		LastName:                   player.LastName,
		Email:                      player.Email,
		Headshot:                   player.Headshot,
		Position:                   player.Position,
		AgeRange:                   player.AgeRange,
		AvailableEquipment:         player.AvailableEquipment,
		AvailableTrainingLocations: player.AvailableTrainingLocations,
		ShortTermGoal:              player.ShortTermGoal,
		LongTermGoal:               player.LongTermGoal,
		CreationDateEpochMillis:    player.CreationDateEpochMillis,
		LastUpdatedDateEpochMillis: player.LastUpdatedDateEpochMillis,
		NotificationId:             player.NotificationId,
		QueuedSubscription:         player.QueuedSubscription,
	}
}
