package coachapplication

const applicationIdAttr = "applicationId"

type CoachApplicationDB struct {
	ApplicationId string `json:"applicationId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phoneNumber"`
	Experience    string `json:"experience"`
}
