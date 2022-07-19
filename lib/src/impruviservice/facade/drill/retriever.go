package drill

import (
	drillsDao "impruviService/dao/drill"
	sessionDao "impruviService/dao/session"
)

func GetDrillsForPlayer(playerId string) ([]*drillsDao.DrillDB, error) {
	sessions, err := sessionDao.GetSessions(playerId)
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

	drillsForPlayerMap, err := drillsDao.BatchGetDrills(drillIds)
	if err != nil {
		return nil, err
	}

	drillsForPlayer := make([]*drillsDao.DrillDB, 0)
	for _, drill := range drillsForPlayerMap {
		drillsForPlayer = append(drillsForPlayer, drill)
	}

	return drillsForPlayer, nil
}

func GetDrillsForCoach(coachId string) ([]*drillsDao.DrillDB, error) {
	drills, err := drillsDao.GetDrillsForCoach(coachId)
	if err != nil {
		return nil, err
	}
	activeDrills := make([]*drillsDao.DrillDB, 0)
	for _, drill := range drills {
		if !drill.IsDeleted {
			activeDrills = append(activeDrills, drill)
		}
	}
	return activeDrills, nil
}
