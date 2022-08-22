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

		if request.PaymentMethodId != "" {
			err = stripeFacade.AttachPaymentMethodIfNotExists(player.StripeCustomerId, request.PaymentMethodId)
			if err != nil {
				log.Printf("Error attaching payment method: %v to customer: %v: %v\n", request.PaymentMethodId, player.StripeCustomerId, err)
				return err
			}
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

		player, err = playerFacade.GetPlayerById(player.PlayerId)
		if err != nil {
			log.Printf("Error while getting player after creating subscription: %v\n", err)
			return err
		}

		subscriptionPlan, err := stripeFacade.GetSubscriptionPlan(request.SubscriptionPlanRef.StripeProductId, request.SubscriptionPlanRef.StripePriceId)
		if err != nil {
			log.Printf("Error while getting subscription plan: %v\n", err)
			return err
		}
		if subscriptionPlan.IsTrial {
			err = stripeFacade.UpdateSubscriptionToCancelAtPeriodEnd(player.StripeCustomerId)
			if err != nil {
				log.Printf("Error while updating trial subscription to cancel at period end: %v\n", err)
				return err
			}
		}

		shouldCreateIntroSession, err := isFirstSubscriptionWithCoach(player.StripeCustomerId, coach.CoachId)
		if err != nil {
			return err
		}
		log.Printf("shouldCreateIntroSession: %v\n", shouldCreateIntroSession)

		if shouldCreateIntroSession {
			err = createIntroSession(player, coach)
			if err != nil {
				log.Printf("Error while creating initial session: %v\n", err)
				return err
			}
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

func isFirstSubscriptionWithCoach(stripeCustomerId, coachId string) (bool, error) {
	log.Printf("isFirstSubscriptionWithCoach. stripeCustomerId: %v. coachId: %v\n", stripeCustomerId, coachId)
	subscriptionHistory, err := stripeFacade.ListSubscriptions(stripeCustomerId)
	if err != nil {
		return false, err
	}
	log.Printf("subscriptionHistory: %v\n", subscriptionHistory)

	subscriptionsWithCoach := 0
	for _, subscription := range subscriptionHistory {
		if subscription.Plan.CoachId == coachId {
			subscriptionsWithCoach += 1
		}
	}

	return subscriptionsWithCoach == 1, nil
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

	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		log.Printf("Error getting subscription: %v\n", err)
	}

	latestSessionNumber, err := sessionDao.GetLatestSessionNumber(player.PlayerId)
	if err != nil {
		return err
	}

	return sessionDao.PutSession(&sessionDao.SessionDB{
		PlayerId:       player.PlayerId,
		Drills:         drills,
		CoachId: coach.CoachId,
		SessionNumber:  latestSessionNumber + 1,
		//IsIntroSession: true, // TODO: do we even need this if we have trial plans?
		// this value must be >= the subscription start date value. Grab it from the subscription object itself
		// to ensure this is the case
		CreationDateEpochMillis:    subscription.CurrentPeriodStartDateEpochMillis + 1, // TODO: +1 here is only required due to a bug in mobile that should be fixed in next build
		LastUpdatedDateEpochMillis: subscription.CurrentPeriodStartDateEpochMillis + 1,
	})
}

func validateCreateSubscriptionRequest(request *CreateSubscriptionRequest) error {
	if request.Token == "" {
		return exceptions.InvalidRequestError{Message: "Token cannot be null/empty"}
	}
	if request.SubscriptionPlanRef == nil || request.SubscriptionPlanRef.CoachId == "" || request.SubscriptionPlanRef.StripeProductId == "" || request.SubscriptionPlanRef.StripePriceId == "" {
		return exceptions.InvalidRequestError{Message: "SubscriptionPlanRef must have a coachId, productId, and priceId"}
	}

	return nil
}
