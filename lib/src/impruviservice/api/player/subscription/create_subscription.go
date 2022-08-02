package subscription

import (
	"fmt"
	coachDao "impruviService/dao/coach"
	sessionDao "impruviService/dao/session"
	"impruviService/exceptions"
	coachFacade "impruviService/facade/coach"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	sessionFacade "impruviService/facade/session"
	stripeFacade "impruviService/facade/stripe"
	"impruviService/model"
	"impruviService/util"
	"log"
)

type CreateSubscriptionRequest struct {
	Token               string                     `json:"token"` // TODO: move token into header and pass in second arg to func
	PaymentMethodId     string                     `json:"paymentMethodId"`
	SubscriptionPlanRef *model.SubscriptionPlanRef `json:"subscriptionPlanRef"`
}

func CreateSubscription(request *CreateSubscriptionRequest) error {
	log.Printf("CreateSubscriptionRequest: %+v\n", request)
	err := validateCreateSubscriptionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CreateSubscriptionRequest: %v\n", err)
		return err
	}

	player, err := playerFacade.GetPlayerFromToken(request.Token)
	if err != nil {
		return err
	}

	coach, err := coachFacade.GetCoachById(request.SubscriptionPlanRef.CoachId)
	if err != nil {
		return err
	}
	if len(coach.IntroSessionDrills) < 4 {
		return exceptions.InvalidRequestError{Message: fmt.Sprintf("Coach: %+v does not have an intro session created.\n", coach.CoachId)}
	}

	hasSubscription, err := alreadyHasSubscription(player)
	if err != nil {
		return err
	}

	if hasSubscription {
		err = stripeFacade.UpdateSubscriptionToCancelAtPeriodEnd(player.StripeCustomerId)
		if err != nil {
			log.Printf("Error updating subscription to cancel at period end: %v\n", err)
			return err
		}
		err = stripeFacade.AttachPaymentMethodIfNotExists(player.StripeCustomerId, request.PaymentMethodId)
		if err != nil {
			log.Printf("Error attaching payment method: %v to customer: %v: %v\n", request.PaymentMethodId, player.StripeCustomerId, err)
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
		err = stripeFacade.CreateSubscription(player, util.GetCurrentTimeEpochMillis(), request.PaymentMethodId, request.SubscriptionPlanRef)
		if err != nil {
			log.Printf("Error while creating subscription: %v\n", err)
			return err
		}

		err = createIntroSession(player, coach)
		if err != nil {
			log.Printf("Error while creating initial session: %v\n", err)
			return err
		}

		err = notificationFacade.SendSubscriptionCreatedNotifications(player)
		if err != nil {
			log.Printf("Error while sending subscription created notifications: %v\n", err)
			return err
		}
		err = dynamicReminderFacade.StartCreateTrainingReminderStepFunctionExecution(&dynamicReminderFacade.CreateTrainingReminderEventData{PlayerId: player.PlayerId})
		if err != nil {
			log.Printf("Error while starting subscription created step function execution: %v\n", err)
			return err
		}

		return nil
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

func createIntroSession(player *playerFacade.Player, coach *coachDao.CoachDB) error {
	coach, err := coachFacade.GetCoachById(player.CoachId)
	if err != nil {
		return err
	}

	drills := make([]*sessionDao.SessionDrillDB, 0)
	for _, drill := range coach.IntroSessionDrills {
		drills = append(drills, &sessionDao.SessionDrillDB{
			DrillId: drill.DrillId,
			Notes:   drill.Notes,
		})
	}
	return sessionFacade.CreateSession(&sessionDao.SessionDB{
		PlayerId:       player.PlayerId,
		Drills:         drills,
		IsIntroSession: true,
	})
}

func validateCreateSubscriptionRequest(request *CreateSubscriptionRequest) error {
	if request.Token == "" {
		return exceptions.InvalidRequestError{Message: "Token cannot be null/empty"}
	}
	if request.PaymentMethodId == "" {
		return exceptions.InvalidRequestError{Message: "PaymentMethodId cannot be null/empty"}
	}
	if request.SubscriptionPlanRef == nil || request.SubscriptionPlanRef.CoachId == "" || request.SubscriptionPlanRef.StripeProductId == "" || request.SubscriptionPlanRef.StripePriceId == "" {
		return exceptions.InvalidRequestError{Message: "SubscriptionPlanRef must have a coachId, productId, and priceId"}
	}

	return nil
}
