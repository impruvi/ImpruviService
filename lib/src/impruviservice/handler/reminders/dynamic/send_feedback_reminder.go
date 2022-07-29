package dynamic

import (
	"encoding/json"
	sessionDao "impruviService/dao/session"
	notificationFacade "impruviService/facade/notification"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	"log"
)

func handleSendFeedbackReminders(request *dynamicReminderFacade.Input) (*dynamicReminderFacade.Input, error) {
	data, err := deserializeFeedbackEventData(request.Data)
	if err != nil {
		return nil, err
	}

	sessionDB, err := sessionDao.GetSession(data.PlayerId, data.SessionNumber)
	if err != nil {
		return nil, err
	}

	log.Printf("Session: %+v\n", sessionDB)
	if sessionDB.IsFeedbackComplete() {
		log.Printf("Feedback is complete. Not sending notifications")
		return &dynamicReminderFacade.Input{
			Completed: true,
		}, nil
	}

	var lastSubmissionUploadDateEpochMillis int64 = 0
	for _, drill := range sessionDB.Drills {
		if drill.Submission.UploadDateEpochMillis > lastSubmissionUploadDateEpochMillis {
			lastSubmissionUploadDateEpochMillis = drill.Submission.UploadDateEpochMillis
		}
	}

	hoursRemaining := getHoursRemaining(lastSubmissionUploadDateEpochMillis)
	log.Printf("%v hours remaining to provide feedback", hoursRemaining)

	if hoursRemaining <= 0 {
		err = notificationFacade.SendFeedbackOverdueNotifications(data.PlayerId, data.SessionNumber)
		if err != nil {
			return nil, err
		}
		return &dynamicReminderFacade.Input{
			Completed: true,
		}, nil
	}

	err = notificationFacade.SendFeedbackReminderNotifications(data.PlayerId, hoursRemaining)
	if err != nil {
		return nil, err
	}

	return &dynamicReminderFacade.Input{
		Completed:   false,
		WaitSeconds: getNewWaitSeconds(hoursRemaining, lastSubmissionUploadDateEpochMillis),
		Type:        request.Type,
		Data:        request.Data,
	}, nil
}

func deserializeFeedbackEventData(dataJSON string) (*dynamicReminderFacade.SendFeedbackReminderEventData, error) {
	var data dynamicReminderFacade.SendFeedbackReminderEventData
	err := json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		log.Printf("Failed to deserialize feedback event data: %v. json: %v\n", err, dataJSON)
		return nil, err
	}
	log.Printf("Deserialized feedback event data: %+v\n", data)
	return &data, nil
}
