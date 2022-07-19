package players

import (
	"impruviService/model"
)

const coachIdIndexName = "coachId-index"
const playerIdAttr = "playerId"
const coachIdAttr = "coachId"

type PlayerDB struct {
	PlayerId                   string       `json:"playerId"`
	CoachId                    string       `json:"coachId"`
	FirstName                  string       `json:"firstName"`
	LastName                   string       `json:"lastName"`
	Email                      string       `json:"email"`
	Headshot                   *model.Media `json:"headshot"`
	CreationDateEpochMillis    int64        `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64        `json:"lastUpdatedDateEpochMillis"`
	NotificationId             string       `json:"notificationId"`

	// TODO: move this out
	Subscription *Subscription `json:"subscription"`
}

// TODO: move this to model once we have subscriptions up and running
type Subscription struct {
	CurrentPeriodStartDateEpochMillis int64 `json:"currentPeriodStartDateEpochMillis"`
	CurrentPeriodEndDateEpochMillis   int64 `json:"currentPeriodEndDateEpochMillis"`
	NumberOfSessions                  int   `json:"numberOfSessions"`
}
