package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	stripeCustomer "github.com/stripe/stripe-go/customer"
	stripePaymentMethod "github.com/stripe/stripe-go/paymentmethod"
	stripePlan "github.com/stripe/stripe-go/plan"
	stripeSubscription "github.com/stripe/stripe-go/sub"
	"impruviService/exceptions"
	"log"
	"strconv"
)

func GetSubscription(stripeCustomerId string) (*Subscription, error) {
	subscriptions, err := getSubscriptionsWithStatus(stripeCustomerId, string(stripe.SubscriptionStatusActive))
	if err != nil {
		return nil, err
	}

	if len(subscriptions) == 0 {
		return nil, exceptions.ResourceNotFoundError{Message: fmt.Sprintf("No active subscription exists for customer: %v\n", stripeCustomerId)}
	}
	if len(subscriptions) > 1 {
		log.Printf("Customer: %v has more than one active subscription\n", stripeCustomerId)
	}

	return subscriptions[0], nil
}

func ListSubscriptions(stripeCustomerId string) ([]*Subscription, error) {
	return getSubscriptionsWithStatus(stripeCustomerId, string(stripe.SubscriptionStatusAll))
}

func getSubscriptionsWithStatus(stripeCustomerId, status string) ([]*Subscription, error) {
	if stripeCustomerId == "" {
		return make([]*Subscription, 0), nil
	}
	subscriptions := make([]*Subscription, 0)

	iter := stripeSubscription.List(&stripe.SubscriptionListParams{
		Customer: stripeCustomerId,
		Status: status,
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
			log.Printf("Failed to get subscription plan: %v\n", err)
			return nil, err
		}

		recurrenceStartDateEpochMillis, err := strconv.ParseInt(subscription.Metadata["recurrenceStartDateEpochMillis"], 10, 64)
		if err != nil {
			log.Printf("Error getting recurrenceStartDateEpochMillis. Error: %v\n", err)
			return nil, err
		}

		subscriptions = append(subscriptions, &Subscription{
			StripeSubscriptionId:              subscription.ID,
			Plan:                              subscriptionPlan,
			CancelAtEndOfPeriod:               subscription.CancelAtPeriodEnd,
			CurrentPeriodStartDateEpochMillis: subscription.CurrentPeriodStart * 1000,
			CurrentPeriodEndDateEpochMillis:   subscription.CurrentPeriodEnd * 1000,
			PlayerId:                          subscription.Metadata["playerId"],
			RecurrenceStartDateEpochMillis:    recurrenceStartDateEpochMillis,
		})
	}

	return subscriptions, nil
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

	isTrial := false
	if isTrialString, ok := plan.Metadata["isTrial"]; ok {
		isTrial, err = strconv.ParseBool(isTrialString)
		if err != nil {
			log.Printf("Error while getting isTrial. Metadata: %+v. Error: %v\n", plan.Metadata, err)
			return nil, err
		}
	}

	return &SubscriptionPlan{
		StripeProductId:   stripeProductId,
		StripePriceId:     stripePriceId,
		CoachId:           plan.Metadata["coachId"],
		Type:              plan.Metadata["type"],
		NumberOfTrainings: numberOfTrainings,
		UnitAmount:        plan.Amount,
		IsTrial:           isTrial,
	}, nil
}

func GetPaymentMethods(stripeCustomerId string) ([]*PaymentMethod, error) {
	if stripeCustomerId == "" {
		return make([]*PaymentMethod, 0), nil
	}

	customer, err := stripeCustomer.Get(stripeCustomerId, nil)
	if err != nil {
		log.Printf("Error while getting customer by id: %v. error: %v\n", stripeCustomerId, err)
		return nil, err
	}
	log.Printf("Customer: %+v\n", customer)
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
			IsDefault:       getCustomerDefaultPaymentMethod(customer) == pm.ID,
		})
	}

	return paymentMethodIds, nil
}

func getCustomerDefaultPaymentMethod(customer *stripe.Customer) string {
	if customer == nil || customer.InvoiceSettings == nil || customer.InvoiceSettings.DefaultPaymentMethod == nil {
		return ""
	}
	return customer.InvoiceSettings.DefaultPaymentMethod.ID
}
