package coachapplication

import (
	coachApplicationDao "impruviService/dao/coachapplication"
)

func CreateApplication(application *coachApplicationDao.CoachApplicationDB) error {
	return coachApplicationDao.CreateApplication(application)
}
