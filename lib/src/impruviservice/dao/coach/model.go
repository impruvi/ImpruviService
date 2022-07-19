package coaches

import (
	"impruviService/model"
)

const coachIdAttr = "coachId"

type CoachDB struct {
	CoachId                    string       `json:"coachId"`
	FirstName                  string       `json:"firstName"`
	LastName                   string       `json:"lastName"`
	Email                      string       `json:"email"`
	Headshot                   *model.Media `json:"headshot"`
	Position                   string       `json:"position"`
	School                     string       `json:"school"`
	YouthClub                  string       `json:"youthClub"`
	About                      string       `json:"about"`
	CreationDateEpochMillis    int64        `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64        `json:"lastUpdatedDateEpochMillis"`
	NotificationId             string       `json:"notificationId"`
}
