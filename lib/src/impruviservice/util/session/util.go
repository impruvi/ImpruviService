package session

import (
	sessionDao "impruviService/dao/session"
	stripeFacade "impruviService/facade/stripe"
)

func GetNumberOfSessionsCreatedForPlan(subscription *stripeFacade.Subscription, sessions []*sessionDao.SessionDB) int {
	subscriptionStartDateEpochMillis := subscription.CurrentPeriodStartDateEpochMillis
	numberOfSessionsCreatedAfter := 0
	for _, session := range sessions {
		if !session.IsIntroSession && session.CreationDateEpochMillis > subscriptionStartDateEpochMillis {
			numberOfSessionsCreatedAfter += 1
		}
	}

	return numberOfSessionsCreatedAfter
}
