package session

import (
	"../../dao/drills"
	"../../dao/session"
)

type Session struct {
	UserId        string   `json:"userId"`
	SessionNumber int      `json:"sessionNumber"`
	Drills        []*Drill `json:"drills"`
}

type Drill struct {
	Drill           drills.Drill        `json:"drill"`
	Submission      *session.Submission `json:"submission"`
	Feedback        *session.Feedback   `json:"feedback"`
	Tips            []string            `json:"tips"`
	Repetitions     int                 `json:"repetitions"`
	DurationMinutes int                 `json:"durationMinutes"`
}
