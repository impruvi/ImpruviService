package coach

import (
	"impruviService/dao/coaches"
)

func UpdateCoach(coach *coaches.Coach) error {
	currentCoach, err := coaches.GetCoachById(coach.CoachId)
	if err != nil {
		return err
	}
	coach.NotificationId = currentCoach.NotificationId

	return coaches.PutCoach(coach)
}
