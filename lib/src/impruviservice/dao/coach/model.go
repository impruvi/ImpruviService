package coaches

import (
	"impruviService/model"
)

const coachIdAttr = "coachId"

type CoachDB struct {
	CoachId                    string                       `json:"coachId"`
	FirstName                  string                       `json:"firstName"`
	LastName                   string                       `json:"lastName"`
	Email                      string                       `json:"email"`
	Headshot                   *model.Media                 `json:"headshot"`
	CardImage                  *model.Media                 `json:"cardImage"`
	BackgroundImage            *model.Media                 `json:"backgroundImage"`
	Location                   string                       `json:"location"`
	Position                   string                       `json:"position"`
	School                     string                       `json:"school"`
	Team                       string                       `json:"team"`
	YouthClub                  string                       `json:"youthClub"`
	About                      string                       `json:"about"`
	CreationDateEpochMillis    int64                        `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64                        `json:"lastUpdatedDateEpochMillis"`
	NotificationId             string                       `json:"notificationId"`
	SubscriptionPlanRefs       []*model.SubscriptionPlanRef `json:"subscriptionPlanRefs"`
}
