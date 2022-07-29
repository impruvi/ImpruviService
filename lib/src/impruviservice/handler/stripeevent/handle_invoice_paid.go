package stripeevent

import (
	"github.com/stripe/stripe-go"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	stripeFacade "impruviService/facade/stripe"
	"log"
)

func handleInvoicePaid(invoice *stripe.Invoice) error {
	log.Printf("Invoice paid: %+v\n", invoice)

	if invoice.BillingReason == stripe.InvoiceBillingReasonSubscriptionCycle {
		subscription, err := stripeFacade.GetSubscription(invoice.Customer.ID)
		if err != nil {
			log.Printf("Error while getting subscription for customer: %v. error: %v\n", invoice.Customer.ID, err)
			return err
		}
		log.Printf("Subscription: %+v\n", subscription)
		player, err := playerFacade.GetPlayerById(subscription.PlayerId)
		if err != nil {
			log.Printf("Error while getting player by id: %v. error: %v\n", subscription.PlayerId, err)
			return err
		}
		log.Printf("Player: %+v\n", player)

		// send us and coach notification
		err = notificationFacade.SendSubscriptionRenewedNotifications(player, false)
		if err != nil {
			log.Printf("Error while sending subscription renewed notifications for player: %+v. error %v\n", player, err)
			return err
		}
		err = dynamicReminderFacade.StartCreateTrainingReminderStepFunctionExecution(&dynamicReminderFacade.CreateTrainingReminderEventData{PlayerId: player.PlayerId})
		if err != nil {
			log.Printf("Error while starting subscription created step function execution: %v\n", err)
			return err
		}

		return nil
	} else if invoice.BillingReason == stripe.InvoiceBillingReasonSubscriptionCreate {
		// Do nothing, we should have already sent notifications synchronously while creating the subscription.
	} else {
		log.Printf("Unexpected invoice paid billing reason: %v\n", invoice.BillingReason)
		return notificationFacade.SendUnexpectedInvoiceBillingReasonNotifications(invoice)
	}

	return nil
}
