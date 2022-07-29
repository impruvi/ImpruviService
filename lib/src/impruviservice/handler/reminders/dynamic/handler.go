package dynamic

import (
	"context"
	"errors"
	"fmt"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	"log"
)

func HandleSendDynamicRemindersEvent(_ context.Context, request dynamicReminderFacade.Input) (dynamicReminderFacade.Input, error) {
	log.Printf("Sending dynamic reminder notifications: %+v\n", request)

	var response *dynamicReminderFacade.Input
	var err error
	if request.Type == dynamicReminderFacade.FeedbackReminder {
		response, err = handleSendFeedbackReminders(&request)
	} else if request.Type == dynamicReminderFacade.CreateTrainingReminder {
		response, err = handleSendCreateTrainingReminders(&request)
	} else {
		log.Printf("Unrecognized request type: %v\n", request.Type)
		return dynamicReminderFacade.Input{}, errors.New(fmt.Sprintf("Unrecognized request type: %v\n", request.Type))
	}

	if err != nil {
		log.Printf("Error while sending dynamic reminders for event: %+v. error: %v\n", request, err)
		return dynamicReminderFacade.Input{}, err
	}

	return *response, nil
}
