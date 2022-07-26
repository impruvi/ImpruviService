package subscription

import (
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
)

type ReactivateSubscriptionRequest struct {
	Token string `json:"token"`
}

func ReactivateSubscription(request *ReactivateSubscriptionRequest) error {
	player, err := playerFacade.GetPlayerFromToken(request.Token)
	if err != nil {
		return err
	}

	return stripeFacade.ReactivateSubscription(player.StripeCustomerId)
}