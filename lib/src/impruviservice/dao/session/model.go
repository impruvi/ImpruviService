package session

import (
	"impruviService/model"
	"log"
)

const playerIdAttr = "playerId"
const sessionNumberAttr = "sessionNumber"

type SessionDB struct {
	PlayerId                   string            `json:"playerId"`
	SessionNumber              int               `json:"sessionNumber"`
	Drills                     []*SessionDrillDB `json:"drills"`
	CreationDateEpochMillis    int64             `json:"creationDateEpochMillis"`
	LastUpdatedDateEpochMillis int64             `json:"lastUpdatedDateEpochMillis"`
	HasViewedFeedback          bool              `json:"hasViewedFeedback"`
	IsIntroSession             bool              `json:"isIntroSession"`
}

type SessionDrillDB struct {
	DrillId             string       `json:"drillId"`
	Submission          *model.Media `json:"submission"`
	SubmissionThumbnail *model.Media `json:"submissionThumbnail"`
	Feedback            *model.Media `json:"feedback"`
	FeedbackThumbnail   *model.Media `json:"feedbackThumbnail"`
	Notes               string       `json:"notes"`
}

func (s *SessionDB) IsSubmissionComplete() bool {
	log.Printf("checking if submission complete for session: %v\n", s)
	for _, drill := range s.Drills {
		if drill.Submission == nil || !drill.Submission.IsPresent() {
			return false
		}
	}
	return true
}

func (s *SessionDB) IsFeedbackComplete() bool {
	log.Printf("checking if feedback complete for session: %v\n", s)
	for _, drill := range s.Drills {
		if drill.Feedback == nil || !drill.Feedback.IsPresent() {
			log.Printf("Feedback is not complete")
			return false
		}
	}
	log.Printf("Feedback is complete")
	return true
}
