package stripeevent

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"impruviService/exceptions"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	stripeFacade "impruviService/facade/stripe"
	"impruviService/util"
	"log"
	"strconv"
)

func handleSubscriptionDeleted(subscription *stripe.Subscription) error {
	log.Printf("Subscription deleted: %v\n", subscription)

	playerId := subscription.Metadata["playerId"]
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		log.Printf("Error while getting player by id: %v. error: %v\n", playerId, err)
		return err
	}
	currentRecurrenceStartDateEpochMillis, err := strconv.ParseInt(subscription.Metadata["recurrenceStartDateEpochMillis"], 10, 64)
	if err != nil {
		log.Printf("Error while getting recurrence start date. Error: %v\n", err)
		return err
	}
	coachId := subscription.Metadata["coachId"]

	if player.QueuedSubscription != nil {
		defaultPaymentMethod, err := getDefaultPaymentMethod(player.StripeCustomerId)
		if err != nil {
			log.Printf("Error while getting default payment method for player: %+v. error: %v\n", player, err)
			return err
		}
		if defaultPaymentMethod == nil {
			log.Printf("No payment method for player: %+v. error: %v\n", player, err)
			return exceptions.ResourceNotFoundError{Message: fmt.Sprintf("No default payment method for player: %v. Cannot update subscription.\n", player.PlayerId)}
		}

		var recurrenceStartDateEpochMillis int64
		if player.QueuedSubscription.CoachId == coachId {
			recurrenceStartDateEpochMillis = currentRecurrenceStartDateEpochMillis
		} else {
			recurrenceStartDateEpochMillis = util.GetCurrentTimeEpochMillis()
		}
		err = stripeFacade.CreateSubscription(player, recurrenceStartDateEpochMillis, defaultPaymentMethod.PaymentMethodId, player.QueuedSubscription)
		if err != nil {
			log.Printf("Error while updating subscription to queued subscription: %v\n", err)
			return err
		}

		player.QueuedSubscription = nil
		err = playerFacade.UpdatePlayer(player)
		if err != nil {
			log.Printf("Error while setting player queued subscription to nil: %v\n", err)
			return err
		}

		err = notificationFacade.SendSubscriptionRenewedNotifications(player, true)
		if err != nil {
			log.Printf("Error while updating subscription renewal notifications: %v\n", err)
			return err
		}

		err = dynamicReminderFacade.StartCreateTrainingReminderStepFunctionExecution(&dynamicReminderFacade.CreateTrainingReminderEventData{PlayerId: player.PlayerId})
		if err != nil {
			log.Printf("Error while starting subscription created step function execution: %v\n", err)
			return err
		}

		return nil
	} else {
		player.CoachId = ""
		err = playerFacade.UpdatePlayer(player)
		if err != nil {
			log.Printf("Error while removing coachId from player: %v\n", err)
			return err
		}

		isTrial := false
		if isTrialString, ok := subscription.Plan.Metadata["isTrial"]; ok {
			isTrial, err = strconv.ParseBool(isTrialString)
			log.Printf("Error while getting isTrial. Metadata: %+v. Error: %v\n", subscription.Plan.Metadata, err)
			return err
		}

		if isTrial {
			err = notificationFacade.SendTrialEndedNotifications(player)
			if err != nil {
				log.Printf("Error while sending trial ended notifications: %v\n", err)
				return err
			}
		} else {
			err = notificationFacade.SendSubscriptionDidNotRenewNotifications(player)
			if err != nil {
				log.Printf("Error while sending subscription ended notifications: %v\n", err)
				return err
			}
		}
	}

	return nil
}

func getDefaultPaymentMethod(stripeCustomerId string) (*stripeFacade.PaymentMethod, error) {
	paymentMethods, err := stripeFacade.GetPaymentMethods(stripeCustomerId)
	if err != nil {
		return nil, err
	}

	for _, paymentMethod := range paymentMethods {
		if paymentMethod.IsDefault {
			return paymentMethod, nil
		}
	}

	return nil, nil
}
