package session

import (
	"../../dao/drills"
	"../../dao/session"
	"../../files"
)

func getFullSessionsForPlayer(playerId string) ([]*FullSession, error) {
	sessions, err := session.GetSessions(playerId)
	if err != nil {
		return nil, err
	}

	return getFullSessions(sessions)
}

func getFullSessions(sessions []*session.Session) ([]*FullSession, error) {
	sessionsWithDrills := make([]*FullSession, 0)
	for _, sess := range sessions {
		sessionWithDrill, err := getFullSession(sess)
		if err != nil {
			return nil, err
		}
		sessionsWithDrills = append(sessionsWithDrills, sessionWithDrill)
	}
	return sessionsWithDrills, nil
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
		fullDrills = append(fullDrills, &FullDrill{
			DrillId:                  drill.DrillId,
			CoachId:                  drill.CoachId,
			Name:                     drill.Name,
			Description:              drill.Description,
			Category:                 drill.Category,
			Equipment:                drill.Equipment,
			Submission:               sessionDrill.Submission,
			Feedback:                 sessionDrill.Feedback,
			Notes:                    sessionDrill.Notes,
			EstimatedDurationMinutes: sessionDrill.EstimatedDurationMinutes,
			Prescription:             sessionDrill.Prescription,
			Demos: Demos{
				Front: Demo{FileLocation: files.GetDrillVideoFileLocation(drill.DrillId, files.Front).URL},
				Side:  Demo{FileLocation: files.GetDrillVideoFileLocation(drill.DrillId, files.Side).URL},
				Close: Demo{FileLocation: files.GetDrillVideoFileLocation(drill.DrillId, files.Close).URL},
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
