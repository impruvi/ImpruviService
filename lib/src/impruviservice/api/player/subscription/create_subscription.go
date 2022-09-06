package subscription

import (
	"impruviService/exceptions"
	paymentProcessor "impruviService/facade/payment"
	"impruviService/model"
	"log"
)

type CreateSubscriptionRequest struct {
	Token               string             `json:"token"` // TODO: move token into header and pass in second arg to func
	PaymentMethodId     string             `json:"paymentMethodId"`
	SubscriptionPlanRef *model.PricingPlan `json:"subscriptionPlanRef"`
}

func CreateSubscription(request *CreateSubscriptionRequest) error {
	log.Printf("CreateSubscriptionRequest: %+v\n", request)
	err := validateCreateSubscriptionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CreateSubscriptionRequest: %v\n", err)
		return err
	}

	return paymentProcessor.ProcessPayment(request.Token, request.PaymentMethodId, request.SubscriptionPlanRef)
}

func validateCreateSubscriptionRequest(request *CreateSubscriptionRequest) error {
	if request.Token == "" {
		return exceptions.InvalidRequestError{Message: "Token cannot be null/empty"}
	}
	if request.SubscriptionPlanRef == nil || request.SubscriptionPlanRef.CoachId == "" || request.SubscriptionPlanRef.StripeProductId == "" || request.SubscriptionPlanRef.StripePriceId == "" {
		return exceptions.InvalidRequestError{Message: "PricingPlan must have a coachId, productId, and priceId"}
	}

	return nil
}
