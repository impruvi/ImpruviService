package coaches

import (
	"impruviService/model"
)

const coachIdAttr = "coachId"
const slugAttr = "slug"
const slugIndexName = "slug-index"

type CoachDB struct {
	CoachId                    string                 `json:"coachId"`
	Slug                       string                 `json:"slug"`
	FirstName                  string                 `json:"firstName"`
	LastName                   string                 `json:"lastName"`
	Email                      string                 `json:"email"`
	Headshot                   *model.Media           `json:"headshot"`
	CardImagePortrait          *model.Media           `json:"cardImagePortrait"`
	CardImageLandscape         *model.Media           `json:"cardImageLandscape"`
	BackgroundImage            *model.Media           `json:"backgroundImage"`
	HeaderImage                *model.Media           `json:"headerImage"`
	FocusAreas                 []string               `json:"focusAreas"`
	Location                   string                 `json:"location"`
	Position                   string                 `json:"position"`
	School                     string                 `json:"school"`
	Team                       string                 `json:"team"`
	YouthClub                  string                 `json:"youthClub"`
	About                      string                 `json:"about"`
	CreationDateEpochMillis    int64                  `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64                  `json:"lastUpdatedDateEpochMillis"`
	NotificationId             string                 `json:"notificationId"`
	IntroSessionDrills         []*IntroSessionDrillDB `json:"introSessionDrills"`
	PricingPlans               []*model.PricingPlan   `json:"pricingPlans"`
	SubscriptionPlanRefs       []*model.PricingPlan   `json:"subscriptionPlanRefs"` //TODO: deprecate this
}

type IntroSessionDrillDB struct {
	DrillId string `json:"drillId"`
	Notes   string `json:"notes"`
}
