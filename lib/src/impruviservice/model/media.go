package model

type Media struct {
	UploadDateEpochMillis int64  `json:"uploadDateEpochMillis"`
	FileLocation          string `json:"fileLocation"`
}

func (m *Media) IsPresent() bool {
	return m.FileLocation != ""
}
