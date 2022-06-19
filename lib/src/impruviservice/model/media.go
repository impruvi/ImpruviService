package model

type Media struct {
	UploadDateEpochMillis int64  `json:"uploadDateEpochMillis"`
	FileLocation          string `json:"fileLocation"`
}
