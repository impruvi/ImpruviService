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

func UpdateNotificationId(coach *coachDao.CoachDB, notificationId string) (*coachDao.CoachDB, error) {
	coach.NotificationId = notificationId
	err := coachDao.PutCoach(coach)
	if err != nil {
		log.Printf("Error while updating coach notification id: %v\n", err)
		return coach, err
	}
	return coach, nil
}
