package session

import (
	"errors"
	"fmt"
	sessionDao "impruviService/dao/session"
	"impruviService/exceptions"
	notificationFacade "impruviService/facade/notification"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	"impruviService/model"
	"impruviService/util"
	"log"
)

func CreateSession(session *sessionDao.SessionDB) error {
	latestSessionNumber, err := sessionDao.GetLatestSessionNumber(session.PlayerId)
	if err != nil {
		return err
	}

	currentTimeMillis := util.GetCurrentTimeEpochMillis()
	session.SessionNumber = latestSessionNumber + 1
	session.CreationDateEpochMillis = currentTimeMillis
	session.LastUpdatedDateEpochMillis = currentTimeMillis
	return sessionDao.PutSession(session)
}

func UpdateSession(session *sessionDao.SessionDB) error {
	currentSession, err := sessionDao.GetSession(session.PlayerId, session.SessionNumber)
	if err != nil {
		return err
	}

	session.LastUpdatedDateEpochMillis = util.GetCurrentTimeEpochMillis()
	session.CreationDateEpochMillis = currentSession.CreationDateEpochMillis
	session.HasViewedFeedback = currentSession.HasViewedFeedback

	return sessionDao.PutSession(session)
}

func DeleteSession(sessionNumber int, playerId string) error {
	err := sessionDao.DeleteSession(sessionNumber, playerId)
	if err != nil {
		return err
	}

	return decrementAllSessionsAbove(sessionNumber, playerId)
}

func ViewFeedback(playerId string, sessionNumber int) error {
	session, err := sessionDao.GetSession(playerId, sessionNumber)
	if err != nil {
		return err
	}

	session.LastUpdatedDateEpochMillis = util.GetCurrentTimeEpochMillis()
	session.HasViewedFeedback = true

	return sessionDao.PutSession(session)
}

func CreateFeedback(playerId string, sessionNumber int, drillId, fileLocation, thumbnailFileLocation string) error {
	session, err := sessionDao.GetSession(playerId, sessionNumber)
	if err != nil {
		return err
	}

	drill, err := findDrill(session.Drills, drillId)
	if err != nil {
		return err
	}
	currentTimeEpochMillis := util.GetCurrentTimeEpochMillis()
	drill.Feedback = &model.Media{
		UploadDateEpochMillis: currentTimeEpochMillis,
		FileLocation:          fileLocation,
	}
	drill.FeedbackThumbnail = &model.Media{
		UploadDateEpochMillis: currentTimeEpochMillis,
		FileLocation:          thumbnailFileLocation,
	}
	err = sessionDao.PutSession(session)
	if err != nil {
		return err
	}

	if session.IsFeedbackComplete() {
		err = notificationFacade.SendFeedbackNotifications(playerId)
		if err != nil {
			log.Printf("Error sending feedback notifications: %v\n", err)
			return err
		}
	}

	return nil
}

func CreateSubmission(playerId string, sessionNumber int, drillId, fileLocation, thumbnailFileLocation string) error {
	session, err := sessionDao.GetSession(playerId, sessionNumber)
	if err != nil {
		return err
	}

	drill, err := findDrill(session.Drills, drillId)
	if err != nil {
		return err
	}
	currentTimeEpochMillis := util.GetCurrentTimeEpochMillis()
	drill.Submission = &model.Media{
		UploadDateEpochMillis: currentTimeEpochMillis,
		FileLocation:          fileLocation,
	}
	drill.SubmissionThumbnail = &model.Media{
		UploadDateEpochMillis: currentTimeEpochMillis,
		FileLocation:          thumbnailFileLocation,
	}

	err = sessionDao.PutSession(session)
	if err != nil {
		return err
	}

	if session.IsSubmissionComplete() {
		log.Printf("Completed session: %v\n", session)
		err = notificationFacade.SendSubmissionNotifications(playerId)
		if err != nil {
			log.Printf("Error while sending notifications on submission: %v\n", err)
			return err
		}
		err = dynamicReminderFacade.StartFeedbackReminderStepFunctionExecution(&dynamicReminderFacade.SendFeedbackReminderEventData{
			PlayerId:      session.PlayerId,
			SessionNumber: session.SessionNumber,
		})
		if err != nil {
			log.Printf("Error while starting feedback reminder step function execution: %v\n", err)
			return err
		}
	}

	return nil
}

func findDrill(drills []*sessionDao.SessionDrillDB, drillId string) (*sessionDao.SessionDrillDB, error) {
	for _, drill := range drills {
		if drill.DrillId == drillId {
			return drill, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("drill %v does not exist in drills: %v\n.", drillId, drills))
}

func decrementAllSessionsAbove(sessionNumber int, playerId string) error {

	for true {
		sessionNumber += 1

		sess, err := sessionDao.GetSession(playerId, sessionNumber)
		if err != nil {
			_, ok := err.(exceptions.ResourceNotFoundError)
			if !ok {
				return err
			} else {
				return nil
			}
		}

		err = sessionDao.DeleteSession(sess.SessionNumber, playerId)
		if err != nil {
			return err
		}

		sess.SessionNumber = sess.SessionNumber - 1
		err = sessionDao.PutSession(sess)
		if err != nil {
			return err
		}
	}

	return nil
}
