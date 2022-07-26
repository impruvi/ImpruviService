package subscription

import (
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
)

type GetSubscriptionRequest struct {
	Token string `json:"token"`
}

type GetSubscriptionResponse struct {
	Subscription *stripeFacade.Subscription `json:"subscription"`
}

func GetSubscription(request *GetSubscriptionRequest) (*GetSubscriptionResponse, error) {
	player, err := playerFacade.GetPlayerFromToken(request.Token)
	if err != nil {
		return nil, err
	}

	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		return nil, err
	}
	return &GetSubscriptionResponse{Subscription: subscription}, nil
}
