package dynamic

type EventType string

const (
	FeedbackReminder       EventType = "FEEDBACK_REMINDER"
	CreateTrainingReminder EventType = "CREATE_TRAINING_REMINDER"
)

type Input struct {
	Type        EventType `json:"type"`
	Data        string    `json:"data"`
	WaitSeconds int64     `json:"waitSeconds"`
	Completed   bool      `json:"completed"`
}

type CreateTrainingReminderEventData struct {
	PlayerId string `json:"playerId"`
}

type SendFeedbackReminderEventData struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
}
