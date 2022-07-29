package dynamic

import (
	"encoding/json"
	sessionDao "impruviService/dao/session"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	stripeFacade "impruviService/facade/stripe"
	sessionUtil "impruviService/util/session"
	"log"
)

func handleSendCreateTrainingReminders(request *dynamicReminderFacade.Input) (*dynamicReminderFacade.Input, error) {
	data, err := deserializeCreateTrainingEventData(request.Data)
	if err != nil {
		return nil, err
	}

	player, err := playerFacade.GetPlayerById(data.PlayerId)
	if err != nil {
		return nil, err
	}

	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		return nil, err
	}

	isTrainingPlanCreated, err := hasCreatedTrainingPlan(player, subscription)
	if err != nil {
		log.Printf("Error while checking if create training reminder should be sent: %v\n", err)
		return nil, err
	}

	if isTrainingPlanCreated {
		log.Printf("Training plan is created. Not sending notifications")
		return &dynamicReminderFacade.Input{
			Completed: true,
		}, nil
	}

	hoursRemaining := getHoursRemaining(subscription.CurrentPeriodStartDateEpochMillis)
	log.Printf("%v hours remaining to create training plan", hoursRemaining)

	if hoursRemaining <= 0 {
		err = notificationFacade.SendCreateTrainingPlanOverdueNotifications(data.PlayerId)
		if err != nil {
			return nil, err
		}
		return &dynamicReminderFacade.Input{
			Completed: true,
		}, nil
	}

	err = notificationFacade.SendCreateTrainingPlanReminderNotifications(data.PlayerId, hoursRemaining)
	if err != nil {
		return nil, err
	}

	return &dynamicReminderFacade.Input{
		Completed:   false,
		WaitSeconds: getNewWaitSeconds(hoursRemaining, subscription.CurrentPeriodStartDateEpochMillis),
		Type:        request.Type,
		Data:        request.Data,
	}, nil
}

func hasCreatedTrainingPlan(player *playerFacade.Player, subscription *stripeFacade.Subscription) (bool, error) {
	sessions, err := sessionDao.GetSessions(player.PlayerId)
	if err != nil {
		return false, err
	}

	numberOfSessionCreatedForPlan := sessionUtil.GetNumberOfSessionsCreatedForPlan(subscription, sessions)
	return numberOfSessionCreatedForPlan >= subscription.Plan.NumberOfTrainings, nil
}

func deserializeCreateTrainingEventData(dataJSON string) (*dynamicReminderFacade.CreateTrainingReminderEventData, error) {
	var data dynamicReminderFacade.CreateTrainingReminderEventData
	err := json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		log.Printf("Failed to deserialize create training event data: %v. json: %v\n", err, dataJSON)
		return nil, err
	}
	log.Printf("Deserialized create training event data: %+v\n", data)
	return &data, nil
}
