package profile

import (
	"impruviService/dao/drills"
	"impruviService/dao/session"
)

func GetDrillsForPlayer(playerId string) ([]*drills.Drill, error) {
	sessions, err := session.GetSessions(playerId)
	if err != nil {
		return nil, err
	}

	drillIdsSet := make(map[string]bool, 0)
	for _, sess := range sessions {
		for _, drill := range sess.Drills {
			drillIdsSet[drill.DrillId] = true
		}
	}

	drillIds := make([]string, 0)
	for drillId, _ := range drillIdsSet {
		drillIds = append(drillIds, drillId)
	}

	drillsForPlayerMap, err := drills.BatchGetDrills(drillIds)
	if err != nil {
		return nil, err
	}

	drillsForPlayer := make([]*drills.Drill, 0)
	for _, drill := range drillsForPlayerMap {
		drillsForPlayer = append(drillsForPlayer, drill)
	}

	return drillsForPlayer, nil
}
