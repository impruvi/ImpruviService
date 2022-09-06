package onetimepurchase

import (
	"impruviService/exceptions"
	paymentProcessor "impruviService/facade/payment"
	"impruviService/model"
	"log"
)

type CreateOneTimePurchaseRequest struct {
	Token           string             `json:"token"`
	PaymentMethodId string             `json:"paymentMethodId"`
	PriceRef        *model.PricingPlan `json:"priceRef"`
}

func CreateOneTimePurchase(request *CreateOneTimePurchaseRequest) error {
	log.Printf("CreateOneTimePurchaseRequest: %+v\n", request)
	err := validateCreateOneTimePurchaseRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CreateOneTimePurchaseRequest: %v\n", err)
		return err
	}

	return paymentProcessor.ProcessPayment(request.Token, request.PaymentMethodId, request.PriceRef)
}

func validateCreateOneTimePurchaseRequest(request *CreateOneTimePurchaseRequest) error {
	if request.Token == "" {
		return exceptions.InvalidRequestError{Message: "Token cannot be null/empty"}
	}
	if request.PriceRef == nil || request.PriceRef.CoachId == "" || request.PriceRef.StripeProductId == "" || request.PriceRef.StripePriceId == "" {
		return exceptions.InvalidRequestError{Message: "PricingPlan must have a coachId, productId, and priceId"}
	}

	return nil
}
