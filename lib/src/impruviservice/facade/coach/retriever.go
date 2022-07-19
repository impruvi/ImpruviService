package coach

import (
	"impruviService/dao/coach"
)

func GetCoachById(coachId string) (*coaches.CoachDB, error) {
	coach, err := coaches.GetCoachById(coachId)
	if err != nil {
		return nil, err
	}
	return coach, nil
}
