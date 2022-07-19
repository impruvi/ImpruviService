package coach

import (
	coachDao "impruviService/dao/coach"
	"log"
)

func UpdateCoach(coach *coachDao.CoachDB) error {
	currentCoach, err := coachDao.GetCoachById(coach.CoachId)
	if err != nil {
		return err
	}
	coach.NotificationId = currentCoach.NotificationId

	return coachDao.PutCoach(coach)
}

func UpdateCoachNotificationId(coach *coachDao.CoachDB, notificationId string) *coachDao.CoachDB {
	coach.NotificationId = notificationId
	err := coachDao.PutCoach(coach)
	if err != nil {
		// don't error if we can't update push notification id
		log.Printf("Error while updating coach notification id: %v\n", err)
		return coach
	}
	log.Printf("Updated coach id's %v push notification token: %v\n", coach.CoachId, notificationId)
	return coach
}
