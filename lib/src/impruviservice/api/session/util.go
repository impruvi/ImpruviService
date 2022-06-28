package session

import (
	"impruviService/dao/drills"
	"impruviService/dao/session"
	"impruviService/files"
	"log"
)

func getFullSessionsForPlayer(playerId string) ([]*FullSession, error) {
	sessions, err := session.GetSessions(playerId)
	if err != nil {
		return nil, err
	}
	log.Printf("Sessions: %v\n", sessions)

	return getFullSessions(sessions)
}

func getFullSessions(sessions []*session.Session) ([]*FullSession, error) {
	log.Printf("Getting full sessions: %v\n", sessions)
	fullSessions := make([]*FullSession, 0)
	for _, sess := range sessions {
		fullSession, err := getFullSession(sess)
		log.Printf("full session: %v\n", fullSession)
		if err != nil {
			return nil, err
		}
		fullSessions = append(fullSessions, fullSession)
	}
	return fullSessions, nil
}

func getFullSession(sess *session.Session) (*FullSession, error) {
	drillIds := getDrillIds(sess.Drills)
	drillDetails, err := drills.BatchGetDrills(drillIds)
	if err != nil {
		return nil, err
	}

	fullDrills := make([]*FullDrill, 0)
	for _, sessionDrill := range sess.Drills {
		drill := drillDetails[sessionDrill.DrillId]
		var feedback Media
		if sessionDrill.Feedback != nil && sessionDrill.Feedback.VideoUploadDateEpochMillis > 0 {
			feedback = Media{
				VideoUploadDateEpochMillis: sessionDrill.Feedback.VideoUploadDateEpochMillis,
				FileLocation:               files.GetFeedbackVideoFileLocation(sess.PlayerId, sess.SessionNumber, drill.DrillId).URL,
			}
		}
		var submission Media
		if sessionDrill.Submission != nil && sessionDrill.Submission.VideoUploadDateEpochMillis > 0 {
			submission = Media{
				VideoUploadDateEpochMillis: sessionDrill.Submission.VideoUploadDateEpochMillis,
				FileLocation:               files.GetSubmissionVideoFileLocation(sess.PlayerId, sess.SessionNumber, drill.DrillId).URL,
			}
		}
		fullDrills = append(fullDrills, &FullDrill{
			DrillId:                  drill.DrillId,
			CoachId:                  drill.CoachId,
			Name:                     drill.Name,
			Description:              drill.Description,
			Category:                 drill.Category,
			Equipment:                drill.Equipment,
			Submission:               submission,
			Feedback:                 feedback,
			Notes:                    sessionDrill.Notes,
			EstimatedDurationMinutes: sessionDrill.EstimatedDurationMinutes,
			Demos: Demos{
				Front:          Media{FileLocation: files.GetDemoVideoFileLocation(drill.DrillId, files.Front).URL},
				Side:           Media{FileLocation: files.GetDemoVideoFileLocation(drill.DrillId, files.Side).URL},
				Close:          Media{FileLocation: files.GetDemoVideoFileLocation(drill.DrillId, files.Close).URL},
				FrontThumbnail: Media{FileLocation: files.GetDemoVideoThumbnailFileLocation(drill.DrillId, files.Front).URL},
				SideThumbnail:  Media{FileLocation: files.GetDemoVideoThumbnailFileLocation(drill.DrillId, files.Side).URL},
				CloseThumbnail: Media{FileLocation: files.GetDemoVideoThumbnailFileLocation(drill.DrillId, files.Close).URL},
			},
		})
	}

	return &FullSession{
		PlayerId:      sess.PlayerId,
		SessionNumber: sess.SessionNumber,
		Drills:        fullDrills,
	}, nil
}

func getDrillIds(drills []*session.Drill) []string {
	drillIds := make([]string, 0)
	for _, drill := range drills {
		drillIds = append(drillIds, drill.DrillId)
	}
	return drillIds
}

// check if session drills are completed
func SessionSubmissionComplete(session *session.Session) bool {
	for _, drill := range session.Drills {
		if !drillSubmissionComplete(drill) {
			return false
		}
	}
	return true
}

// check if drill is completed
func drillSubmissionComplete(drill *session.Drill) bool {
	return submissionComplete(drill.Submission)
}

// check if drill has uploaded video
func submissionComplete(submission *session.Submission) bool {
	return &submission.VideoUploadDateEpochMillis != nil
}

// check if session feedback is completed for all drills
func SessionFeedbackComplete(session *session.Session) bool {
	for _, drill := range session.Drills {
		if !drillFeedbackComplete(drill) {
			return false
		}
	}
	return true
}

// check if feedback is completed for a drill
func drillFeedbackComplete(drill *session.Drill) bool {
	return feedbackComplete(drill.Feedback)
}

// check if feedback has an uploaded video
func feedbackComplete(feedback *session.Feedback) bool {
	return &feedback.VideoUploadDateEpochMillis != nil
}
