package coach

import (
	coachDao "impruviService/dao/coach"
	"impruviService/exceptions"
	coachFacade "impruviService/facade/coach"
	"log"
)

type UpdateCoachRequest struct {
	Coach *coachDao.CoachDB `json:"coach"`
}

func UpdateCoach(request *UpdateCoachRequest) error {
	log.Printf("UpdateCoachRequest: %+v\n", request)
	err := validateUpdateCoachRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid UpdateCoachRequest: %v\n", err)
		return err
	}
	return coachFacade.UpdateCoach(request.Coach)
}

func validateUpdateCoachRequest(request *UpdateCoachRequest) error {
	if request.Coach == nil {
		return exceptions.InvalidRequestError{Message: "Coach cannot be null/empty"}
	}
	if request.Coach.CoachId == "" {
		return exceptions.InvalidRequestError{Message: "CoachId cannot be null/empty"}
	}
	if request.Coach.CreationDateEpochMillis == 0 {
		return exceptions.InvalidRequestError{Message: "CreationDateEpochMillis cannot be null/empty"}
	}
	if request.Coach.LastUpdatedDateEpochMillis == 0 {
		return exceptions.InvalidRequestError{Message: "LastUpdatedDateEpochMillis cannot be null/empty"}
	}
	if request.Coach.FirstName == "" {
		return exceptions.InvalidRequestError{Message: "FirstName cannot be null/empty"}
	}
	if request.Coach.LastName == "" {
		return exceptions.InvalidRequestError{Message: "LastName cannot be null/empty"}
	}
	if request.Coach.About == "" {
		return exceptions.InvalidRequestError{Message: "About cannot be null/empty"}
	}
	if request.Coach.Location == "" {
		return exceptions.InvalidRequestError{Message: "Location cannot be null/empty"}
	}
	if request.Coach.Team == "" {
		return exceptions.InvalidRequestError{Message: "Team cannot be null/empty"}
	}
	if request.Coach.Headshot == nil || !request.Coach.Headshot.IsPresent() {
		return exceptions.InvalidRequestError{Message: "Headshot must be provided"}
	}
	if request.Coach.BackgroundImage == nil || !request.Coach.BackgroundImage.IsPresent() {
		return exceptions.InvalidRequestError{Message: "BackgroundImage must be provide"}
	}
	if request.Coach.CardImagePortrait == nil || !request.Coach.CardImagePortrait.IsPresent() {
		return exceptions.InvalidRequestError{Message: "CardImagePortrait must be provided"}
	}
	if request.Coach.CardImageLandscape == nil || !request.Coach.CardImageLandscape.IsPresent() {
		return exceptions.InvalidRequestError{Message: "CardImageLandscape must be provided"}
	}
	if request.Coach.FocusAreas == nil || len(request.Coach.FocusAreas) == 0 {
		return exceptions.InvalidRequestError{Message: "Coach must have at least 1 focus area"}
	}
	if request.Coach.SubscriptionPlanRefs == nil || len(request.Coach.SubscriptionPlanRefs) == 0 {
		return exceptions.InvalidRequestError{Message: "Coach must have at least 1 subscription plan"}
	}

	return nil
}
