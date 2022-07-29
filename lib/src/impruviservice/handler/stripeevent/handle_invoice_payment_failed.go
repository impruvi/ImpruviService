package stripeevent

import (
	"github.com/stripe/stripe-go"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"log"
)

func handleInvoicePaymentFailed(invoice *stripe.Invoice) error {
	log.Printf("Invoice payment failed: %v\n", invoice)

	playerId := invoice.Subscription.Metadata["playerId"]
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		log.Printf("Error while getting player by id: %v. error: %v\n", playerId, err)
		return err
	}

	err = stripeFacade.CancelSubscription(player.StripeCustomerId)
	if err != nil {
		log.Printf("Failed to cancel subscription due to invoice payment failing: %v\n", err)
		return err
	}

	err = notificationFacade.SendSubscriptionRenewalFailureNotifications(player)
	if err != nil {
		log.Printf("Failed to send subscription renewal failed notifications to player: %+v. error: %v\n", player, err)
		return err
	}

	return nil
}
