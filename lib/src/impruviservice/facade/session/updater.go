package session

import (
	"errors"
	"fmt"
	mediaConvertAccessor "impruviService/accessor/mediaconvert"
	sessionDao "impruviService/dao/session"
	"impruviService/exceptions"
	notificationFacade "impruviService/facade/notification"
	playerFacade "impruviService/facade/player"
	dynamicReminderFacade "impruviService/facade/reminder/dynamic"
	stripeFacade "impruviService/facade/stripe"
	"impruviService/model"
	"impruviService/util"
	sessionUtil "impruviService/util/session"
	"log"
)

func CreateSession(session *sessionDao.SessionDB) error {
	latestSessionNumber, err := sessionDao.GetLatestSessionNumber(session.PlayerId)
	if err != nil {
		return err
	}

	coachId, err := getCoachId(session.PlayerId)
	if err != nil {
		return err
	}
	session.CoachId = coachId

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
	coachId, err := getCoachId(session.PlayerId)
	if err != nil {
		return err
	}
	session.CoachId = coachId

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
	if drill.HasFeedback() {
		log.Printf("Drill already has feedback")
		return nil
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

	err = startFeedbackMediaConversion(session.PlayerId, drillId, session.SessionNumber, fileLocation)
	if err != nil {
		return err
	}

	if session.IsFeedbackComplete() {
		err = notificationFacade.SendFeedbackNotifications(playerId)
		if err != nil {
			log.Printf("Error sending feedback notifications: %v\n", err)
			return err
		}

		player, err := playerFacade.GetPlayerById(playerId)
		if err != nil {
			log.Printf("Error while getting player: %v. Error: %v\n", playerId, err)
			return err
		}

		justCompletedTrial, err := hasCompletedTrial(player)
		if err != nil {
			return err
		}

		if justCompletedTrial {
			player.CoachId = ""
			err = playerFacade.UpdatePlayer(player)
			if err != nil {
				log.Printf("Error while removing coachId from player: %v\n", err)
				return err
			}

			err = stripeFacade.CancelSubscription(player.StripeCustomerId)
			if err != nil {
				log.Printf("Error while cancelling subscription for player: %+v. Error: %v\n", player, err)
				return err
			}

			err = notificationFacade.SendTrialEndedNotifications(player)
			if err != nil {
				log.Printf("Error while sending trial ended notifications: %v\n", err)
				return err
			}
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
	if drill.HasSubmission() {
		log.Printf("Drill already has submission")
		return nil
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

	err = startSubmissionMediaConversion(session.PlayerId, drillId, session.SessionNumber, fileLocation)
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

func getCoachId(playerId string) (string, error){
	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return "", err
	}
	return player.CoachId, nil
}

func findDrill(drills []*sessionDao.SessionDrillDB, drillId string) (*sessionDao.SessionDrillDB, error) {
	for _, drill := range drills {
		if drill.DrillId == drillId {
			return drill, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("drill %v does not exist in drills: %v\n.", drillId, drills))
}

func startSubmissionMediaConversion(playerId, drillId string, sessionNumber int, fileLocation string) error {
	return mediaConvertAccessor.StartJob(fileLocation, &mediaConvertAccessor.Metadata{
		Type: mediaConvertAccessor.SubmissionVideo,
		SubmissionVideoMetadata: &mediaConvertAccessor.SubmissionVideoMetadata{
			DrillId:       drillId,
			PlayerId:      playerId,
			SessionNumber: sessionNumber,
		},
	})
}

func startFeedbackMediaConversion(playerId, drillId string, sessionNumber int, fileLocation string) error {
	return mediaConvertAccessor.StartJob(fileLocation, &mediaConvertAccessor.Metadata{
		Type: mediaConvertAccessor.FeedbackVideo,
		FeedbackVideoMetadata: &mediaConvertAccessor.FeedbackVideoMetadata{
			DrillId:       drillId,
			PlayerId:      playerId,
			SessionNumber: sessionNumber,
		},
	})
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

func hasCompletedTrial(player *playerFacade.Player) (bool, error) {
	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		return false, err
	}

	if !subscription.Plan.IsTrial {
		return false, nil
	}

	sessions, err := sessionDao.GetSessions(player.PlayerId)
	if err != nil {
		return false, err
	}

	return sessionUtil.HasCompletedPlan(subscription, sessions), nil
}