package payment

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

// ProcessPayment processes all 3 types of "payments":
// 1. Free trial
// 2. One time purchase
// 3. Subscription
// Each of these payment types acts like a subscription meaning:
//   - it unlocks X number of sessions which will expire after 1 month
//   - if there is currently an active "subscription", that "subscription" will be set to not recur at the end of the month
//     and the new "subscription" will be set as the queued "subscription" (subscription here is in quotes to differentiate this payment
//     that acts like a one time purchase but is actually a subscription under the hood from the actual subscriptions.
//
// Key differences between these different payment types are:
// - Free trial:
//   - All sessions are free.
//   - If sessions are completed before the end of 1 month, the subscription will be automatically terminated.
//   - This subscription never auto-renews.
//
// - One time purchase:
//   - If sessions are completed before the end of 1 month, the subscription will be automatically terminated.
//   - This subscription never auto-renews.
//
// - Subscription:
//   - If sessions are completed before the end of 1 month, the subscription does not terminate.
//   - This subscription will auto-renew unless user explicitly cancels it
func ProcessPayment(token, paymentMethodId string, priceRef *model.PricingPlan) error {
	player, err := playerFacade.GetPlayerFromToken(token)
	if err != nil {
		return err
	}

	coach, err := coachFacade.GetCoachById(priceRef.CoachId)
	if err != nil {
		return err
	}
	if len(coach.IntroSessionDrills) < 4 {
		return exceptions.InvalidRequestError{Message: fmt.Sprintf("Coach: %+v does not have an intro session created.\n", coach.CoachId)}
	}

	hasSubscription, err := stripeFacade.HasSubscription(player.StripeCustomerId)
	if err != nil {
		return err
	}

	if hasSubscription {
		return setQueuedSubscription(player, priceRef, paymentMethodId)
	} else {
		return createSubscription(player, coach, priceRef, paymentMethodId)
	}
}

func setQueuedSubscription(player *playerFacade.Player, priceRef *model.PricingPlan, paymentMethodId string) error {
	err := stripeFacade.UpdateSubscriptionToCancelAtPeriodEnd(player.StripeCustomerId)
	if err != nil {
		log.Printf("Error updating subscription to cancel at period end: %v\n", err)
		return err
	}

	if paymentMethodId != "" {
		err = stripeFacade.AttachPaymentMethodIfNotExists(player.StripeCustomerId, paymentMethodId)
		if err != nil {
			log.Printf("Error attaching payment method: %v to customer: %v: %v\n", paymentMethodId, player.StripeCustomerId, err)
			return err
		}
	}

	player.QueuedSubscription = priceRef
	return playerFacade.UpdatePlayer(player)
}

func createSubscription(player *playerFacade.Player, coach *coachDao.CoachDB, priceRef *model.PricingPlan, paymentMethodId string) error {
	log.Printf("Creating subscription for player: %+v\n", player)
	player.CoachId = priceRef.CoachId
	err := playerFacade.UpdatePlayer(player)
	if err != nil {
		return err
	}
	player, err = stripeFacade.CreateSubscription(player, util.GetCurrentTimeEpochMillis(), paymentMethodId, priceRef)
	if err != nil {
		log.Printf("Error while creating subscription: %v\n", err)
		return err
	}

	shouldCreateIntroSession, err := isFirstSubscriptionWithCoach(player.StripeCustomerId, coach.CoachId)
	if err != nil {
		return err
	}
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
		PlayerId:      player.PlayerId,
		Drills:        drills,
		CoachId:       coach.CoachId,
		SessionNumber: latestSessionNumber + 1,
		//IsIntroSession: true, // TODO: do we even need this if we have trial plans?
		// this value must be >= the subscription start date value. Grab it from the subscription object itself
		// to ensure this is the case
		CreationDateEpochMillis:    subscription.CurrentPeriodStartDateEpochMillis + 1, // TODO: +1 here is only required due to a bug in mobile that should be fixed in next build
		LastUpdatedDateEpochMillis: subscription.CurrentPeriodStartDateEpochMillis + 1,
	})
}
