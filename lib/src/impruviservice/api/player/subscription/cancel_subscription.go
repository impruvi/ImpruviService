package subscription

import (
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
)

type CancelSubscriptionRequest struct {
	Token string `json:"token"`
}

func CancelSubscription(request *CancelSubscriptionRequest) error {
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

	return stripeFacade.UpdateSubscriptionToCancelAtPeriodEnd(player.StripeCustomerId)
}
