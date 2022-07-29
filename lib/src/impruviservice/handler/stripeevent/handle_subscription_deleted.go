package stripeevent

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"impruviService/exceptions"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	stripeFacade "impruviService/facade/stripe"
	"log"
)

func handleSubscriptionDeleted(subscription *stripe.Subscription) error {
	log.Printf("Subscription deleted: %v\n", subscription)

	playerId := subscription.Metadata["playerId"]
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		log.Printf("Error while getting player by id: %v. error: %v\n", playerId, err)
	}

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

		err = stripeFacade.CreateSubscription(player, defaultPaymentMethod.PaymentMethodId, player.QueuedSubscription)
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
		err = notificationFacade.SendSubscriptionDidNotRenewNotifications(player)
		if err != nil {
			log.Printf("Error while updating subscription ended notifications: %v\n", err)
			return err
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
