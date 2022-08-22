package session

import (
	sessionDao "impruviService/dao/session"
	stripeFacade "impruviService/facade/stripe"
)

func GetNumberOfSessionsCreatedForPlan(subscription *stripeFacade.Subscription, sessions []*sessionDao.SessionDB) int {
	numberOfSessionsCreatedAfter := 0
	for _, session := range sessions {
		if !session.IsIntroSession && isSessionInPlan(subscription, session) {
			numberOfSessionsCreatedAfter += 1
		}
	}

	return numberOfSessionsCreatedAfter
}

func HasCompletedPlan(subscription *stripeFacade.Subscription, sessions []*sessionDao.SessionDB) bool {
	numberOfSessionsCompletedAfter := 0
	for _, session := range sessions {
		if !session.IsIntroSession && isSessionInPlan(subscription, session) && session.IsFeedbackComplete() {
			numberOfSessionsCompletedAfter += 1
		}
	}

	return numberOfSessionsCompletedAfter == subscription.Plan.NumberOfTrainings
}

func isSessionInPlan(subscription *stripeFacade.Subscription, session *sessionDao.SessionDB) bool {
	return session.CreationDateEpochMillis > subscription.CurrentPeriodStartDateEpochMillis
}