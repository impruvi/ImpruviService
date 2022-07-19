package notification

import (
	"fmt"
	expoAccessor "impruviService/accessor/expo"
	snsAccessor "impruviService/accessor/sns"
	coachFacade "impruviService/facade/coach"
	playerFacade "impruviService/facade/player"
	"log"
)

func SendFeedbackNotifications(playerId string) error {
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return err
	}
	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	snsAccessor.SendTextToSystem(fmt.Sprintf("Coach %v submitted feedback on your session!", coach.FirstName))
	if player.NotificationId != "" {
		log.Printf("Sending push notification for feedback to: %v. %v\n", player.FirstName, player.NotificationId)
		expoAccessor.SendPushNotification(
			fmt.Sprintf("Coach %v submitted feedback!", coach.FirstName),
			fmt.Sprintf("Review your feedback before the next session"),
			player.NotificationId)
	} else {
		log.Printf("Not sending push notification for feedback")
	}
	return nil
}

func SendSubmissionNotifications(playerId string) error {
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return err
	}
	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	snsAccessor.SendTextToSystem(fmt.Sprintf("%v completed a session!", player.FirstName))
	if coach.NotificationId != "" {
		log.Printf("Sending push notification for submission to: %v. %v\n", coach.CoachId, coach.NotificationId)
		expoAccessor.SendPushNotification(
			fmt.Sprintf("%v completed a session!", player.FirstName),
			fmt.Sprintf("You have 24 hours to submit feedback"),
			coach.NotificationId)
	} else {
		log.Printf("Not sending push notification for submission")
	}

	return nil
}
