package players

import (
	"impruviService/model"
)

const coachIdIndexName = "coachId-index"
const emailIndexName = "email-index"
const playerIdAttr = "playerId"
const coachIdAttr = "coachId"
const emailAttr = "email"

type PlayerDB struct {
	PlayerId                   string             `json:"playerId"`
	CoachId                    string             `json:"coachId"`
	StripeCustomerId           string             `json:"stripeCustomerId"`
	FirstName                  string             `json:"firstName"`
	LastName                   string             `json:"lastName"`
	Email                      string             `json:"email"`
	Password                   string             `json:"password"`
	Headshot                   *model.Media       `json:"headshot"`
	Position                   string             `json:"position"`
	AgeRange                   string             `json:"ageRange"`
	AvailableEquipment         []string           `json:"availableEquipment"`
	AvailableTrainingLocations []string           `json:"availableTrainingLocations"`
	ShortTermGoal              string             `json:"shortTermGoal"`
	LongTermGoal               string             `json:"longTermGoal"`
	CreationDateEpochMillis    int64              `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64              `json:"lastUpdatedDateEpochMillis"`
	NotificationId             string             `json:"notificationId"`
	QueuedSubscription         *model.PricingPlan `json:"queuedSubscription"`

	// when user has initiated signup but not yet confirmed their email, isActive is set to false
	IsActive       bool   `json:"isActive"`
	ActivationCode string `json:"activationCode"`
}
