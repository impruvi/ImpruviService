package session

import (
	sessionDao "impruviService/dao/session"
	"impruviService/exceptions"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	sessionFacade "impruviService/facade/session"
	stripeFacade "impruviService/facade/stripe"
	sessionUtil "impruviService/util/session"
	"log"
)

type CreateSessionRequest struct {
	Session *sessionDao.SessionDB `json:"session"`
}

func CreateSession(request *CreateSessionRequest) error {
	log.Printf("CreateSessionRequest: %+v\n", request)
	err := validateCreateSessionRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid CreateSessionRequest: %v\n", err)
		return err
	}

	err = sessionFacade.CreateSession(request.Session)
	if err != nil {
		log.Printf("Failed to create session: %v\n", err)
		return err
	}

	err = sendNotificationsIfAllTrainingSessionsForPlanAreCreated(request.Session.PlayerId)
	if err != nil {
		log.Printf("Error while sending training plan created notification: %v\n", err)
		return err
	}

	return nil
}

func sendNotificationsIfAllTrainingSessionsForPlanAreCreated(playerId string) error {
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return err
	}

	sessions, err := sessionDao.GetSessions(player.PlayerId)
	if err != nil {
		return err
	}

	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		return err
	}

	numberOfSessionsCreatedForPlan := sessionUtil.GetNumberOfSessionsCreatedForPlan(subscription, sessions)
	if numberOfSessionsCreatedForPlan >= subscription.Plan.NumberOfTrainings {
		err = notificationFacade.SendTrainingPlanCreatedNotifications(player, numberOfSessionsCreatedForPlan)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateCreateSessionRequest(request *CreateSessionRequest) error {
	if request.Session == nil {
		return exceptions.InvalidRequestError{Message: "Session cannot be null/empty"}
	}
	if request.Session.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	if request.Session.Drills == nil || len(request.Session.Drills) < 4 {
		return exceptions.InvalidRequestError{Message: "You must provide at least 4 drills"}
	} else {
		for _, drill := range request.Session.Drills {
			if drill.DrillId == "" {
				return exceptions.InvalidRequestError{Message: "DrillId cannot be null/empty"}
			}
		}
	}

	return nil
}
