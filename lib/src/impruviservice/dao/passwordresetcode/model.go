package passwordresetcode

const emailAttr = "email"
const creationDateEpochMillisAttr = "creationDateEpochMillis"

type PasswordResetCodeEntryDB struct {
	Email                   string `json:"email"`
	CreationDateEpochMillis int64  `json:"creationDateEpochMillis"`
	Code                    string `json:"code"`
}
