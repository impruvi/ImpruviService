package session

import (
	"../../dao/drills"
	"../../dao/session"
)

func getSessionsWithDrillsForUser(userId string) ([]*Session, error) {
	sessions, err := session.GetSessions(userId)
	if err != nil {
		return nil, err
	}

	return getSessionsWithDrills(sessions)
}

func getSessionsWithDrills(sessions []*session.Session) ([]*Session, error) {
	sessionsWithDrills := make([]*Session, 0)
	for _, sess := range sessions {
		sessionWithDrill, err := getSessionWithDrills(sess)
		if err != nil {
			return nil, err
		}
		sessionsWithDrills = append(sessionsWithDrills, sessionWithDrill)
	}
	return sessionsWithDrills, nil
}

func getSessionWithDrills(sess *session.Session) (*Session, error) {
	drillIds := getDrillIds(sess.Drills)
	drillDetails, err := drills.BatchGetDrills(drillIds)
	if err != nil {
		return nil, err
	}

	fullDrills := make([]*Drill, 0)
	for _, drill := range sess.Drills {
		fullDrills = append(fullDrills, &Drill{
			Drill:           *drillDetails[drill.DrillId],
			Submission:      drill.Submission,
			Feedback:        drill.Feedback,
			Tips:            drill.Tips,
			Repetitions:     drill.Repetitions,
			DurationMinutes: drill.DurationMinutes,
		})
	}

	return &Session{
		UserId:        sess.UserId,
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
