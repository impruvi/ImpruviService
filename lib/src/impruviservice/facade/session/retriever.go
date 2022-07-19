package session

import (
	playersDao "impruviService/dao/player"
	sessionDao "impruviService/dao/session"
)

func GetSessionsForPlayers(players []*playersDao.PlayerDB) ([]*PlayerSessions, error) {
	playerSessions := make([]*PlayerSessions, 0)
	for _, player := range players {
		sessions, err := GetSessionsForPlayer(player.PlayerId)
		if err != nil {
			return nil, err
		}

		playerSessions = append(playerSessions, &PlayerSessions{
			Player:   player,
			Sessions: sessions,
		})
	}
	return playerSessions, nil
}

func GetSessionsForPlayer(playerId string) ([]*Session, error) {
	sessionDBs, err := sessionDao.GetSessions(playerId)
	if err != nil {
		return nil, err
	}

	return convertAll(sessionDBs)
}

func GetSessionsForCoach(coachId string) ([]*PlayerSessions, error) {
	players, err := playersDao.GetPlayersForCoach(coachId)
	if err != nil {
		return nil, err
	}
	playerSessions, err := GetSessionsForPlayers(players)
	if err != nil {
		return playerSessions, err
	}

	return playerSessions, nil
}
