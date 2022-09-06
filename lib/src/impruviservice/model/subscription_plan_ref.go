package model

type PriceType string

const (
	OneTimePurchase PriceType = "OneTimePurchase"
	Trial           PriceType = "Trial"
	Subscription    PriceType = "Subscription"
)

type PricingPlan struct {
	CoachId               string    `json:"coachId"`
	StripeProductId       string    `json:"stripeProductId"`
	StripePriceId         string    `json:"stripePriceId"`
	Type                  PriceType `json:"type"`
	UnitAmountPerTraining int       `json:"unitAmountPerTraining"`
	NumberOfTrainings     int       `json:"numberOfTrainings"`
}
