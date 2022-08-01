package notification

import (
	"fmt"
	"github.com/stripe/stripe-go"
	expoAccessor "impruviService/accessor/expo"
	sesAccessor "impruviService/accessor/ses"
	snsAccessor "impruviService/accessor/sns"
	playerDao "impruviService/dao/player"
	coachFacade "impruviService/facade/coach"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	emailTemplateProvider "impruviService/provider/emailtemplate"
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

func SendFeedbackReminderNotifications(playerId string, hoursRemaining int) error {
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return err
	}

	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	snsAccessor.SendTextToSystem(fmt.Sprintf("Coach %v %v has %v hour(s) remain to provide feedback", coach.FirstName, coach.LastName, hoursRemaining))
	if coach.NotificationId != "" {
		log.Printf("Sending push notification for feedback reminder to: %v. %v\n", coach.CoachId, coach.NotificationId)
		var hourText string
		if hoursRemaining > 1 {
			hourText = "hours"
		} else {
			hourText = "hour"
		}
		expoAccessor.SendPushNotification(
			fmt.Sprintf("Provide feedback on %v %v's training", player.FirstName, player.LastName),
			fmt.Sprintf("You have %v %v remaining to submit feedback", hoursRemaining, hourText),
			coach.NotificationId)
	} else {
		log.Printf("Not sending push notification for feedback reminder")
	}

	return nil
}

func SendFeedbackOverdueNotifications(playerId string, sessionNumber int) error {
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return err
	}

	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	snsAccessor.SendTextToSystem(fmt.Sprintf("Coach %v %v failed to provide feedback in 24 hours for %v %v on session %v", coach.FirstName, coach.LastName, player.FirstName, player.LastName, sessionNumber))
	return nil
}

func SendCreateTrainingPlanReminderNotifications(playerId string, hoursRemaining int) error {
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return err
	}

	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	snsAccessor.SendTextToSystem(fmt.Sprintf("Coach %v %v has %v hour(s) remain to create a training plan", coach.FirstName, coach.LastName, hoursRemaining))
	if coach.NotificationId != "" {
		log.Printf("Sending push notification for create training reminder to: %v. %v\n", coach.CoachId, coach.NotificationId)
		var hourText string
		if hoursRemaining > 1 {
			hourText = "hours"
		} else {
			hourText = "hour"
		}
		expoAccessor.SendPushNotification(
			fmt.Sprintf("Create a training plan for %v %v", player.FirstName, player.LastName),
			fmt.Sprintf("You have %v %v remaining to create the training plan", hoursRemaining, hourText),
			coach.NotificationId)
	} else {
		log.Printf("Not sending push notification for create training reminder")
	}

	return nil
}

func SendCreateTrainingPlanOverdueNotifications(playerId string) error {
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return err
	}

	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	snsAccessor.SendTextToSystem(fmt.Sprintf("Coach %v %v failed to create training plan for %v %v in 24 hours", coach.FirstName, coach.LastName, player.FirstName, player.LastName))
	return nil
}

func SendSubmissionReminderNotifications(player *playerDao.PlayerDB) error {
	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	if player.NotificationId != "" {
		log.Printf("Sending push notification reminder to complete training to: %v. %v\n", player, player.NotificationId)
		expoAccessor.SendPushNotification(
			fmt.Sprintf("Complete your next training!"),
			fmt.Sprintf("Hey %v, complete your training session to receive feedback from coach %v %v!", player.FirstName, coach.FirstName, coach.LastName),
			player.NotificationId)
	} else {
		log.Printf("Not sending push notification reminder to complete training")
	}

	return nil
}

func SendSubscriptionCreatedNotifications(player *playerFacade.Player) error {
	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		return err
	}

	snsAccessor.SendTextToSystem(fmt.Sprintf("%v %v subscribed to %v %v's %v plan!", player.FirstName, player.LastName, coach.FirstName, coach.LastName, subscription.Plan.Type))
	if coach.NotificationId != "" {
		log.Printf("Sending push notification for subscription creation to: %v. %v\n", coach.CoachId, coach.NotificationId)
		expoAccessor.SendPushNotification(
			fmt.Sprintf("%v %v subscribed to your %v plan!", player.FirstName, player.LastName, subscription.Plan.Type),
			fmt.Sprintf("You have 24 hours to create %v training sessions for them", subscription.Plan.NumberOfTrainings),
			coach.NotificationId)
	} else {
		log.Printf("Not sending push notification for subscription")
	}
	return sesAccessor.SendEmail(
		coach.Email,
		"NEW SUBSCRIBER, ACTION REQUIRED",
		emailTemplateProvider.GetSubscriptionCreatedCoachEmailHtml(player, coach),
		emailTemplateProvider.GetSubscriptionCreatedCoachEmailText(player, coach))
}

func SendTrainingPlanCreatedNotifications(player *playerFacade.Player, numberOfSessions int) error {
	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	snsAccessor.SendTextToSystem(fmt.Sprintf("%v %v created trainings for %v %v", coach.FirstName, coach.LastName, player.FirstName, player.LastName))
	if player.NotificationId != "" {
		log.Printf("Sending push notification for training plan created to: %v. %v\n", coach.CoachId, coach.NotificationId)
		expoAccessor.SendPushNotification(
			fmt.Sprintf("Coach %v has created your training plan", coach.LastName),
			fmt.Sprintf("You have %v new training sessions. Get started!", numberOfSessions),
			player.NotificationId)
	} else {
		log.Printf("Not sending push notification for training plan created")
	}
	return sesAccessor.SendEmail(
		coach.Email,
		fmt.Sprintf("Coach %v has created your training plan", coach.LastName),
		fmt.Sprintf("<div>You have %v new training sessions. Get started!</div>", numberOfSessions),
		fmt.Sprintf("You have %v new training sessions. Get started!", numberOfSessions))
}

func SendSubscriptionRenewedNotifications(player *playerFacade.Player, isSubscriptionUpdated bool) error {
	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		return err
	}

	if isSubscriptionUpdated {
		err = snsAccessor.SendTextToSystem(fmt.Sprintf("%v %v's subscription was renewed to updated plan.", player.FirstName, player.LastName))
	} else {
		err = snsAccessor.SendTextToSystem(fmt.Sprintf("%v %v's subscription was autorenewed.", player.FirstName, player.LastName))
	}
	if err != nil {
		log.Printf("Failed to send system text message on subscription renewal")
		return err
	}

	if coach.NotificationId != "" {
		log.Printf("Sending push notification for subscription renewal to: %v. %v\n", coach.CoachId, coach.NotificationId)
		expoAccessor.SendPushNotification(
			fmt.Sprintf("%v %v renewed their subscription!", player.FirstName, player.LastName),
			fmt.Sprintf("You have 24 hours to create %v training sessions for them", subscription.Plan.NumberOfTrainings),
			coach.NotificationId)
	} else {
		log.Printf("Not sending push notification for subscription renewal")
	}
	return sesAccessor.SendEmail(
		coach.Email,
		fmt.Sprintf("%v %v renewed their subscription!", player.FirstName, player.LastName),
		fmt.Sprintf("<div>You have 24 hours to create %v training sessions for them</div>", subscription.Plan.NumberOfTrainings),
		fmt.Sprintf("You have 24 hours to create %v training sessions for them", subscription.Plan.NumberOfTrainings))
}

func SendSubscriptionCancelledNotification(player *playerFacade.Player) error {
	return snsAccessor.SendTextToSystem(fmt.Sprintf("%v %v's subscription has been cancelled.", player.FirstName, player.LastName))
}

func SendSubscriptionDidNotRenewNotifications(player *playerFacade.Player) error {
	return snsAccessor.SendTextToSystem(fmt.Sprintf("%v %v's subscription has ended and is not renewed.", player.FirstName, player.LastName))
}

func SendUnexpectedInvoiceBillingReasonNotifications(invoice *stripe.Invoice) error {
	return snsAccessor.SendTextToSystem(fmt.Sprintf("Unexpected billing reason: %v for invoice: %v\n.", invoice.BillingReason, invoice))
}

func SendSubscriptionRenewalFailureNotifications(player *playerFacade.Player) error {
	err := snsAccessor.SendTextToSystem(fmt.Sprintf("%v %v's subscription failed to renew.", player.FirstName, player.LastName))
	if err != nil {
		log.Printf("Failed to send system text message on subscription renewal failure")
		return err
	}

	return sesAccessor.SendEmail(
		player.Email,
		"Your Impruvi subscription failed to renew",
		"<div>Your Impruvi subscription failed to renew. You will need to create a new subscription on <a href=\"https://impruviapp.com\">impruviapp.com</a></div>",
		"Your Impruvi subscription failed to renew. You will need to create a new subscription on impruviapp.com")
}

func SendUnhandledStripeEventTypeNotifications(eventType string) error {
	return snsAccessor.SendTextToSystem(fmt.Sprintf("Unhandled stripe event type: %v\n.", eventType))
}
