package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	stripeCustomer "github.com/stripe/stripe-go/customer"
	stripePaymentMethod "github.com/stripe/stripe-go/paymentmethod"
	stripeProduct "github.com/stripe/stripe-go/product"
	stripeSubscription "github.com/stripe/stripe-go/sub"
	playerFacade "impruviService/facade/player"
	"impruviService/model"
	"log"
)

func CancelSubscription(stripeCustomerId string) error {
	if stripeCustomerId == "" {
		return nil
	}

	subscription, err := GetSubscription(stripeCustomerId)
	if err != nil {
		log.Printf("Error while getting subscription to cancel: %v\n", err)
		return err
	}

	updatedSubscription, err := stripeSubscription.Cancel(subscription.StripeSubscriptionId, nil)
	if err != nil {
		log.Printf("Error while cancelling subscription: %v\n", err)
	}

	log.Printf("Updated subscription after cancelling: %+v\n", updatedSubscription)
	return nil
}

func UpdateSubscriptionToCancelAtPeriodEnd(stripeCustomerId string) error {
	if stripeCustomerId == "" {
		return nil
	}
	subscription, err := GetSubscription(stripeCustomerId)
	if err != nil {
		log.Printf("Error while getting subscription to cancel: %v\n", err)
		return err
	}

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	log.Printf("Subscription params: %+v\n", params)
	updatedSubscription, err := stripeSubscription.Update(subscription.StripeSubscriptionId, params)
	if err != nil {
		log.Printf("Error updating subscription to cancel at the end of the period: %v\n", err)
		return err
	}
	log.Printf("updatedSubscription: %+v\n", updatedSubscription)
	return err
}

func ReactivateSubscription(stripeCustomerId string) error {
	if stripeCustomerId == "" {
		return nil
	}
	subscription, err := GetSubscription(stripeCustomerId)
	if err != nil {
		log.Printf("Error while getting subscription to reactivate: %v\n", err)
		return err
	}

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}
	log.Printf("Subscription params: %+v\n", params)
	updatedSubscription, err := stripeSubscription.Update(subscription.StripeSubscriptionId, params)
	if err != nil {
		log.Printf("Error updating subscription to not cancel at the end of the period: %v\n", err)
		return err
	}
	log.Printf("updatedSubscription: %+v\n", updatedSubscription)
	return err
}

func CreateSubscription(player *playerFacade.Player, paymentMethodId string, subscriptionPlanRef *model.SubscriptionPlanRef) error {
	customer, err := getOrCreateCustomer(player)
	if err != nil {
		return err
	}

	log.Printf("Customer: %+v\n", customer)
	err = AttachPaymentMethodIfNotExists(customer.ID, paymentMethodId)
	if err != nil {
		return err
	}

	return subscribeToPlan(player.PlayerId, customer.ID, subscriptionPlanRef)
}

func getOrCreateCustomer(player *playerFacade.Player) (*stripe.Customer, error) {
	if player.StripeCustomerId != "" {
		// customer already exists
		log.Printf("Stripe customerId already exists: %v\n", player.StripeCustomerId)
		return stripeCustomer.Get(player.StripeCustomerId, nil)
	} else {
		// create customer
		params := &stripe.CustomerParams{
			Name:  stripe.String(fmt.Sprintf("%s %s", player.FirstName, player.LastName)),
			Email: stripe.String(player.Email),
		}
		log.Printf("Stripe customerId does not already exist. Creating with params: %+v\n", params)
		customer, _ := stripeCustomer.New(params)
		log.Printf("Created customer: %+v\n", customer)

		// store customerId along with player object
		player.StripeCustomerId = customer.ID
		err := playerFacade.UpdatePlayer(player)
		if err != nil {
			return nil, err
		}

		return customer, nil
	}
}

func AttachPaymentMethodIfNotExists(stripeCustomerId, paymentMethodId string) error {
	// Check if payment method is already attached to customer
	existingPaymentMethods, err := GetPaymentMethods(stripeCustomerId)
	if err != nil {
		return err
	}
	for _, existingPaymentMethod := range existingPaymentMethods {
		if existingPaymentMethod.PaymentMethodId == paymentMethodId {
			log.Printf("Customer: %v already had payment method attached: %v\n", stripeCustomerId, paymentMethodId)
			return nil
		}
	}

	// Attach PaymentMethod
	paymentMethod, err := stripePaymentMethod.Attach(
		paymentMethodId,
		&stripe.PaymentMethodAttachParams{
			Customer: stripe.String(stripeCustomerId),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Payment method: %v\n", paymentMethod)

	// Update invoice settings default
	customer, err := stripeCustomer.Update(
		stripeCustomerId,
		&stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(paymentMethod.ID),
			},
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Updated customer: %+v\n", customer)

	return nil
}

func subscribeToPlan(playerId, stripeCustomerId string, subscriptionPlanRef *model.SubscriptionPlanRef) error {
	// TODO: we can probably remove the below
	product, err := stripeProduct.Get(subscriptionPlanRef.StripeProductId, nil)
	if err != nil {
		return err
	}

	log.Printf("Product: %v\n", product)

	// Create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomerId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(subscriptionPlanRef.StripePriceId),
			},
		},
	}
	subscriptionParams.AddMetadata("playerId", playerId)
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	log.Printf("subscriptionParams: %v\n", subscriptionParams)
	subscription, err := stripeSubscription.New(subscriptionParams)
	if err != nil {
		return err
	}
	log.Printf("Subscription: %v\n", subscription)
	if subscription.Status != stripe.SubscriptionStatusActive {
		log.Printf("Subscription status is not active!") // TODO: notify us of unexpected event
	}
	return nil
}
