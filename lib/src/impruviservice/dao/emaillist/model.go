package emaillist

const emailAttr = "email"

type EmailListSubscriptionDB struct {
	Email                      string `json:"email"`
	CreationDateEpochMillis    int64  `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64  `json:"lastUpdatedDateEpochMillis"`
}
