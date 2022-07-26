package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	stripePaymentMethod "github.com/stripe/stripe-go/paymentmethod"
	stripePlan "github.com/stripe/stripe-go/plan"
	stripeSubscription "github.com/stripe/stripe-go/sub"
	"impruviService/exceptions"
	"log"
	"strconv"
)

func GetSubscription(stripeCustomerId string) (*Subscription, error) {
	if stripeCustomerId == "" {
		return nil, nil
	}
	subscriptions := make([]*Subscription, 0)

	iter := stripeSubscription.List(&stripe.SubscriptionListParams{
		Customer: stripeCustomerId,
	})
	for iter.Next() {
		subscription := iter.Subscription()
		log.Printf("Subscription: %+v\n", subscription)
		if len(subscription.Items.Data) < 1 {
			log.Printf("Less than 1 item in the subscription: %v\n", len(subscription.Items.Data))
			return nil, exceptions.ResourceNotFoundError{Message: fmt.Sprintf("Less than 1 item in the subscription: %v\n", len(subscription.Items.Data))}
		}
		if len(subscription.Items.Data) > 1 {
			log.Printf("More than 1 item in the subscription: %v\n", len(subscription.Items.Data))
		}
		item := subscription.Items.Data[0]

		stripeProductId := item.Plan.Product.ID
		stripePriceId := item.Plan.ID
		subscriptionPlan, err := GetSubscriptionPlan(stripeProductId, stripePriceId)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, &Subscription{
			StripeSubscriptionId:              subscription.ID,
			Plan:                              subscriptionPlan,
			CancelAtEndOfPeriod:               subscription.CancelAtPeriodEnd,
			CurrentPeriodStartDateEpochMillis: subscription.CurrentPeriodStart * 1000,
			CurrentPeriodEndDateEpochMillis:   subscription.CurrentPeriodEnd * 1000,
		})
	}

	if len(subscriptions) == 0 {
		return nil, exceptions.ResourceNotFoundError{Message: fmt.Sprintf("No subscription exists for customer: %v\n", stripeCustomerId)}
	}
	if len(subscriptions) > 1 {
		log.Printf("Customer: %v has more than one subscription\n", stripeCustomerId)
	}
	return subscriptions[0], nil
}

func GetSubscriptionPlan(stripeProductId, stripePriceId string) (*SubscriptionPlan, error) {
	plan, err := stripePlan.Get(stripePriceId, nil)
	if err != nil {
		log.Printf("Error while retreiving plan with planId: %v. err: %+v\n", stripePriceId, err)
		return nil, err
	}

	numberOfTrainings, err := strconv.Atoi(plan.Metadata["numberOfTrainings"])
	if err != nil {
		log.Printf("Error while getting number of trainings from plan. Metadata: %+v. Error: %v\n", plan.Metadata, err)
		return nil, err
	}

	return &SubscriptionPlan{
		StripeProductId:   stripeProductId,
		StripePriceId:     stripePriceId,
		CoachId:           plan.Metadata["coachId"],
		Type:              plan.Metadata["type"],
		NumberOfTrainings: numberOfTrainings,
		UnitAmount:        plan.Amount,
	}, nil
}

func GetPaymentMethods(stripeCustomerId string) ([]*PaymentMethod, error) {
	if stripeCustomerId == "" {
		return make([]*PaymentMethod, 0), nil
	}
	log.Printf("getting payment methods for customer: %v\n", stripeCustomerId)
	paymentMethodIds := make([]*PaymentMethod, 0)

	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(stripeCustomerId),
		Type:     stripe.String("card"),
	}
	log.Printf("params: %+v\n", params)
	i := stripePaymentMethod.List(params)
	if i.Err() != nil {
		return nil, i.Err()
	}
	for i.Next() {
		pm := i.PaymentMethod()
		log.Printf("payment method id: %+v\n", pm)
		paymentMethodIds = append(paymentMethodIds, &PaymentMethod{
			PaymentMethodId: pm.ID,
			Last4:           pm.Card.Last4,
			Brand:           string(pm.Card.Brand),
			ExpMonth:        pm.Card.ExpMonth,
			ExpYear:         pm.Card.ExpYear,
		})
	}

	return paymentMethodIds, nil
}
