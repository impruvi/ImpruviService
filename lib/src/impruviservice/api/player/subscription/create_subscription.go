package subscription

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"impruviService/model"
	"log"
)

type CreateSubscriptionRequest struct {
	Token               string                     `json:"token"` // TODO: move token into header and pass in second arg to func
	PaymentMethodId     string                     `json:"paymentMethodId"`
	SubscriptionPlanRef *model.SubscriptionPlanRef `json:"subscriptionPlanRef"`
}

func CreateSubscription(request *CreateSubscriptionRequest) error {
	player, err := playerFacade.GetPlayerFromToken(request.Token)
	if err != nil {
		return err
	}

	hasSubscription, err := alreadyHasSubscription(player)
	if err != nil {
		return err
	}

	if hasSubscription {
		err := stripeFacade.UpdateSubscriptionToCancelAtPeriodEnd(player.StripeCustomerId)
		if err != nil {
			return err
		}
		player.QueuedSubscription = request.SubscriptionPlanRef
		return playerFacade.UpdatePlayer(player)
	} else {
		log.Printf("Creating subscription for player: %+v\n", player)
		player.CoachId = request.SubscriptionPlanRef.CoachId
		log.Printf("Setting coachId: %v\n", player)
		err = playerFacade.UpdatePlayer(player)
		if err != nil {
			return err
		}
		return stripeFacade.CreateSubscription(player, request.PaymentMethodId, request.SubscriptionPlanRef)
	}
}

func alreadyHasSubscription(player *playerFacade.Player) (bool, error) {
	if player.StripeCustomerId == "" {
		return false, nil
	}

	_, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		if _, ok := err.(exceptions.ResourceNotFoundError); ok {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
