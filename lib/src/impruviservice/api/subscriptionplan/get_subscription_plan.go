package subscriptionplan

import (
	stripeFacade "impruviService/facade/stripe"
)

type GetSubscriptionPlanRequest struct {
	StripeProductId string `json:"stripeProductId"`
	StripePriceId   string `json:"stripePriceId"`
}

type GetSubscriptionPlanResponse struct {
	SubscriptionPlan *stripeFacade.SubscriptionPlan `json:"subscriptionPlan"`
}

func GetSubscriptionPlan(request *GetSubscriptionPlanRequest) (*GetSubscriptionPlanResponse, error) {
	plan, err := stripeFacade.GetSubscriptionPlan(request.StripeProductId, request.StripePriceId)
	if err != nil {
		return nil, err
	}
	return &GetSubscriptionPlanResponse{SubscriptionPlan: plan}, nil
}
