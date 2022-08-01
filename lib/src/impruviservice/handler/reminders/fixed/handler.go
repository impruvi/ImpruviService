package fixed

import (
	playerDao "impruviService/dao/player"
	sessionDao "impruviService/dao/session"
	notificationFacade "impruviService/facade/notification"
	"log"
)

func HandleSendFixedRemindersEvent() (interface{}, error) {
	log.Printf("Sending fixed reminder notifications")

	// Send reminder to every player that hasn't completed their trainings for the billing period
	players, err := playerDao.ListPlayers()
	if err != nil {
		log.Printf("Failed to list players: %v\n", err)
		return nil, err
	}

	for _, player := range players {
		shouldSendReminders, err := hasOutstandingSessions(player.PlayerId)
		if err != nil {
			log.Printf("Error while checking if player has outstanding sessions: %v\n", err)
			return nil, err
		}

		if shouldSendReminders {
			err = notificationFacade.SendSubmissionReminderNotifications(player)
			if err != nil {
				log.Printf("Error while sending submission reminder notifications for player: %v\n", player)
				return nil, err
			}
		}
	}

	return nil, nil
}

func hasOutstandingSessions(playerId string) (bool, error) {
	sessions, err := sessionDao.GetSessions(playerId)
	if err != nil {
		log.Printf("Failed to get sessions for playerId: %+v. error: %v\n", playerId, err)
		return false, err
	}

	for _, session := range sessions {
		if !session.IsSubmissionComplete() {
			return true, nil
		}
	}

	return false, nil
}
