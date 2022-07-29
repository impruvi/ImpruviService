package dynamic

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	stepFunctionAccessor "impruviService/accessor/stepfunction"
	"impruviService/util"
	"log"
)

func StartFeedbackReminderStepFunctionExecution(data *SendFeedbackReminderEventData) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error while serializing send feedback reminder data: %+v. error: %v\n", data, err)
	}

	// invocationId must be less than or equal to 80 characters
	invocationId := fmt.Sprintf("%v-%v-%v-%v", FeedbackReminder, data.PlayerId, data.SessionNumber, uuid.New().String()[0:15])
	return startReminderStepFunctionExecution(invocationId, FeedbackReminder, string(bytes))
}

func StartCreateTrainingReminderStepFunctionExecution(data *CreateTrainingReminderEventData) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error while serializing send create training reminder data: %+v. error: %v\n", data, err)
		return err
	}

	invocationId := fmt.Sprintf("%v-%v-%v", CreateTrainingReminder, data.PlayerId, uuid.New().String()[0:15])
	return startReminderStepFunctionExecution(invocationId, CreateTrainingReminder, string(bytes))
}

func startReminderStepFunctionExecution(invocationId string, eventType EventType, data string) error {
	request := Input{
		WaitSeconds: util.TwelveHoursInSeconds, // send reminders with 12 hours remaining, then 1 hour remaining
		Type:        eventType,
		Data:        data,
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error while serializing request: %+v. error: %v\n", request, err)
		return err
	}

	return stepFunctionAccessor.StartExecution(invocationId, string(bytes))
}
