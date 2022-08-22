package session

import (
	drillsDao "impruviService/dao/drill"
	sessionDao "impruviService/dao/session"
)

func convertAll(sessions []*sessionDao.SessionDB) ([]*Session, error) {
	fullSessions := make([]*Session, 0)
	for _, sess := range sessions {
		fullSession, err := convert(sess)
		if err != nil {
			return nil, err
		}
		fullSessions = append(fullSessions, fullSession)
	}
	return fullSessions, nil
}

func convert(sess *sessionDao.SessionDB) (*Session, error) {
	drillIds := getDrillIds(sess.Drills)
	drillDetails, err := drillsDao.BatchGetDrills(drillIds)
	if err != nil {
		return nil, err
	}

	drills := make([]*SessionDrill, 0)
	for _, sessionDrill := range sess.Drills {
		drill := drillDetails[sessionDrill.DrillId]
		drills = append(drills, &SessionDrill{
			DrillId:             drill.DrillId,
			CoachId:             drill.CoachId,
			Name:                drill.Name,
			Description:         drill.Description,
			Category:            drill.Category,
			Equipment:           drill.Equipment,
			Demos:               drill.Demos,
			Submission:          sessionDrill.Submission,
			SubmissionThumbnail: sessionDrill.SubmissionThumbnail,
			Feedback:            sessionDrill.Feedback,
			FeedbackThumbnail:   sessionDrill.FeedbackThumbnail,
			Notes:               sessionDrill.Notes,
		})
	}

	return &Session{
		PlayerId:                   sess.PlayerId,
		SessionNumber:              sess.SessionNumber,
		CoachId: sess.CoachId,
		Drills:                     drills,
		CreationDateEpochMillis:    sess.CreationDateEpochMillis,
		LastUpdatedDateEpochMillis: sess.LastUpdatedDateEpochMillis,
		HasViewedFeedback:          sess.HasViewedFeedback,
		IsIntroSession:             sess.IsIntroSession,
	}, nil
}

func getDrillIds(drills []*sessionDao.SessionDrillDB) []string {
	drillIds := make([]string, 0)
	for _, drill := range drills {
		drillIds = append(drillIds, drill.DrillId)
	}
	return drillIds
}
