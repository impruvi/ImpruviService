package player

import (
	"impruviService/model"
)

type Player struct {
	PlayerId                   string                     `json:"playerId"`
	CoachId                    string                     `json:"coachId"`
	FirstName                  string                     `json:"firstName"`
	LastName                   string                     `json:"lastName"`
	Email                      string                     `json:"email"`
	Position                   string                     `json:"position"`
	AgeRange                   string                     `json:"ageRange"`
	AvailableEquipment         []string                   `json:"availableEquipment"`
	AvailableTrainingLocations []string                   `json:"availableTrainingLocations"`
	ShortTermGoal              string                     `json:"shortTermGoal"`
	LongTermGoal               string                     `json:"longTermGoal"`
	Headshot                   *model.Media               `json:"headshot"`
	CreationDateEpochMillis    int64                      `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64                      `json:"lastUpdatedDateEpochMillis"`
	NotificationId             string                     `json:"notificationId"`
	StripeCustomerId           string                     `json:"stripeCustomerId"`
	QueuedSubscription         *model.SubscriptionPlanRef `json:"queuedSubscription"`
}
