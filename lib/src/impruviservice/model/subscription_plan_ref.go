package model

type SubscriptionPlanRef struct {
	CoachId         string `json:"coachId"`
	StripeProductId string `json:"stripeProductId"`
	StripePriceId   string `json:"stripePriceId"`
}
