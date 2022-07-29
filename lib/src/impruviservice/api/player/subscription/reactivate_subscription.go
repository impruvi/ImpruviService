package subscription

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"log"
)

type ReactivateSubscriptionRequest struct {
	Token string `json:"token"` // TODO: move token into header and pass in second arg to func
}

func ReactivateSubscription(request *ReactivateSubscriptionRequest) error {
	log.Printf("ReactivateSubscriptionRequest: %+v\n", request)
	err := validateReactivateSubscriptionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid ReactivateSubscriptionRequest: %v\n", err)
		return err
	}

	player, err := playerFacade.GetPlayerFromToken(request.Token)
	if err != nil {
		return err
	}

	return stripeFacade.ReactivateSubscription(player.StripeCustomerId)
}

func validateReactivateSubscriptionRequest(request *ReactivateSubscriptionRequest) error {
	if request.Token == "" {
		return exceptions.InvalidRequestError{Message: "Token cannot be null/empty"}
	}

	return nil
}
