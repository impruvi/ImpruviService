package inbox

type InboxEntryType string

const (
	NewSessionAdded  InboxEntryType = "NEW_SESSION_ADDED"
	FeedbackProvided InboxEntryType = "FEEDBACK_PROVIDED"
	Message          InboxEntryType = "MESSAGE"
)

type InboxEntry struct {
	CreationDateEpochMillis int64             `json:"creationDateEpochMillis"`
	Actor                   *InboxEntryActor  `json:"actor"`
	Type                    InboxEntryType    `json:"type"`
	Metadata                map[string]string `json:"metadata"`
}

type InboxEntryActor struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	ImageFileLocation string `json:"image"`
}
