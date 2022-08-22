package session

import (
	playersDao "impruviService/dao/player"
	sessionDao "impruviService/dao/session"
)


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
	playerSessions, err := getPlayerSessionsForCoach(players, coachId)
	if err != nil {
		return nil, err
	}



	return playerSessions, nil
}

func getPlayerSessionsForCoach(players []*playersDao.PlayerDB, coachId string) ([]*PlayerSessions, error) {
	playerSessions := make([]*PlayerSessions, 0)
	for _, player := range players {
		sessions, err := GetSessionsForPlayer(player.PlayerId)
		if err != nil {
			return nil, err
		}

		sessionsWithCoach := make([]*Session, 0)
		for _, session := range sessions {
			if session.CoachId == coachId {
				sessionsWithCoach = append(sessionsWithCoach, session)
			}
		}
		playerSessions = append(playerSessions, &PlayerSessions{
			Player:   player,
			Sessions: sessionsWithCoach,
		})
	}
	return playerSessions, nil
}