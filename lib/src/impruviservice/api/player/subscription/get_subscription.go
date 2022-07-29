package subscription

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"log"
)

type GetSubscriptionRequest struct {
	PlayerId string `json:"playerId"` // TODO: eventually pass token rather than playerId
}

type GetSubscriptionResponse struct {
	Subscription *stripeFacade.Subscription `json:"subscription"`
}

func GetSubscription(request *GetSubscriptionRequest) (*GetSubscriptionResponse, error) {
	log.Printf("GetSubscriptionRequest: %+v\n", request)
	err := validateGetSubscriptionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetSubscriptionRequest: %v\n", err)
		return nil, err
	}

	player, err := playerFacade.GetPlayerById(request.PlayerId)
	if err != nil {
		return nil, err
	}

	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		return nil, err
	}
	return &GetSubscriptionResponse{Subscription: subscription}, nil
}

func validateGetSubscriptionRequest(request *GetSubscriptionRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}

	return nil
}
