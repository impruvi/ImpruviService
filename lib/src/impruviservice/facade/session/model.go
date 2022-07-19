package session

import (
	drillDao "impruviService/dao/drill"
	playerDao "impruviService/dao/player"
	"impruviService/model"
)

type PlayerSessions struct {
	Player   *playerDao.PlayerDB `json:"player"`
	Sessions []*Session          `json:"sessions"`
}

type Session struct {
	PlayerId                   string          `json:"playerId"`
	Name                       string          `json:"name"`
	SessionNumber              int             `json:"sessionNumber"`
	Drills                     []*SessionDrill `json:"drills"`
	CreationDateEpochMillis    int64           `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64           `json:"lastUpdatedDateEpochMillis"`
	HasViewedFeedback          bool            `json:"hasViewedFeedback"`
}

type SessionDrill struct {
	DrillId     string                 `json:"drillId"`
	CoachId     string                 `json:"coachId"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"` // DRIBBLING/WARMUP/SHOOTING/PASSING
	Equipment   []drillDao.EquipmentDB `json:"equipment"`
	Demos       *drillDao.DemosDB      `json:"demos"`

	Submission *model.Media `json:"submission"`
	Feedback   *model.Media `json:"feedback"`
	Notes      string       `json:"notes"`
}
