package subscription

import (
	"impruviService/exceptions"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"log"
)

type CancelSubscriptionRequest struct {
	Token string `json:"token"` // TODO: move token into header and pass in second arg to func
}

func CancelSubscription(request *CancelSubscriptionRequest) error {
	log.Printf("CancelSubscriptionRequest: %+v\n", request)
	err := validateCancelSubscriptionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CancelSubscriptionRequest: %v\n", err)
		return err
	}

	player, err := playerFacade.GetPlayerFromToken(request.Token)
	if err != nil {
		return err
	}

	if player.QueuedSubscription != nil {
		player.QueuedSubscription = nil
		err = playerFacade.UpdatePlayer(player)
		if err != nil {
			return err
		}
	}

	err = stripeFacade.UpdateSubscriptionToCancelAtPeriodEnd(player.StripeCustomerId)
	if err != nil {
		return err
	}

	return notificationFacade.SendSubscriptionCancelledNotification(player)
}

func validateCancelSubscriptionRequest(request *CancelSubscriptionRequest) error {
	if request.Token == "" {
		return exceptions.InvalidRequestError{Message: "Token cannot be null/empty"}
	}

	return nil
}
