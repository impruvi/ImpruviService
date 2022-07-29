package stripe

type Subscription struct {
	StripeSubscriptionId              string            `json:"stripeSubscriptionId"`
	CancelAtEndOfPeriod               bool              `json:"cancelAtEndOfPeriod"`
	Plan                              *SubscriptionPlan `json:"plan"`
	CurrentPeriodStartDateEpochMillis int64             `json:"currentPeriodStartDateEpochMillis"`
	CurrentPeriodEndDateEpochMillis   int64             `json:"currentPeriodEndDateEpochMillis"`
	PlayerId                          string            `json:"playerId"`
}

type SubscriptionPlan struct {
	StripeProductId   string `json:"stripeProductId"`
	StripePriceId     string `json:"stripePriceId"`
	Type              string `json:"type"`
	CoachId           string `json:"coachId"`
	NumberOfTrainings int    `json:"numberOfTrainings"`
	UnitAmount        int64  `json:"unitAmount"`
}

type PaymentMethod struct {
	PaymentMethodId string `json:"paymentMethodId"`
	Last4           string `json:"last4"`
	Brand           string `json:"brand"`
	ExpMonth        uint64 `json:"expMonth"`
	ExpYear         uint64 `json:"expYear"`
	IsDefault       bool   `json:"isDefault"`
}
