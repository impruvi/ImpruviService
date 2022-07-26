package coach

import (
	coachDao "impruviService/dao/coach"
)

func GetCoachById(coachId string) (*coachDao.CoachDB, error) {
	coach, err := coachDao.GetCoachById(coachId)
	if err != nil {
		return nil, err
	}
	return coach, nil
}

func ListCoaches(limit int) ([]*coachDao.CoachDB, error) {
	coaches, err := coachDao.ListCoaches()
	if err != nil {
		return nil, err
	}
	if limit < 0 {
		return coaches, nil
	}
	return coaches[0:limit], nil
}
