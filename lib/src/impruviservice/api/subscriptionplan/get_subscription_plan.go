package subscriptionplan

import (
	"impruviService/exceptions"
	stripeFacade "impruviService/facade/stripe"
	"log"
)

type GetSubscriptionPlanRequest struct {
	StripeProductId string `json:"stripeProductId"`
	StripePriceId   string `json:"stripePriceId"`
}

type GetSubscriptionPlanResponse struct {
	SubscriptionPlan *stripeFacade.SubscriptionPlan `json:"subscriptionPlan"`
}

func GetSubscriptionPlan(request *GetSubscriptionPlanRequest) (*GetSubscriptionPlanResponse, error) {
	log.Printf("GetSubscriptionPlanRequest: %+v\n", request)
	err := validateGetSubscriptionPlanRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetSubscriptionPlanRequest: %v\n", err)
		return nil, err
	}

	plan, err := stripeFacade.GetSubscriptionPlan(request.StripeProductId, request.StripePriceId)
	if err != nil {
		return nil, err
	}
	return &GetSubscriptionPlanResponse{SubscriptionPlan: plan}, nil
}

func validateGetSubscriptionPlanRequest(request *GetSubscriptionPlanRequest) error {
	if request.StripeProductId == "" {
		return exceptions.InvalidRequestError{Message: "StripeProductId cannot be null/empty"}
	}
	if request.StripePriceId == "" {
		return exceptions.InvalidRequestError{Message: "StripePriceId cannot be null/empty"}
	}

	return nil
}
