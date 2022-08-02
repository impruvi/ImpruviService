package inbox

import (
	"fmt"
	sessionDao "impruviService/dao/session"
	coachFacade "impruviService/facade/coach"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"log"
	"sort"
)

func GetInboxForPlayer(playerId string) ([]*InboxEntry, error) {
	sessions, err := sessionDao.GetSessions(playerId)
	if err != nil {
		return nil, err
	}

	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return nil, err
	}
	if player.CoachId == "" {
		return make([]*InboxEntry, 0), nil
	}

	subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
	if err != nil {
		log.Printf("Error getting subscription: %v\n", err)
		return nil, err
	}

	coach, err := coachFacade.GetCoachById(player.CoachId)
	var headshotFileLocation = ""
	if coach.Headshot != nil {
		headshotFileLocation = coach.Headshot.FileLocation
	}
	coachActor := &InboxEntryActor{
		FirstName:         fmt.Sprintf(coach.FirstName),
		LastName:          fmt.Sprintf(coach.LastName),
		ImageFileLocation: headshotFileLocation,
	}

	entries := make([]*InboxEntry, 0)
	entries = append(entries, &InboxEntry{
		CreationDateEpochMillis: subscription.RecurrenceStartDateEpochMillis,
		Actor:                   coachActor,
		Type:                    Message,
		Metadata: map[string]string{
			"messageContent": fmt.Sprintf("Hey %v, welcome to Impruvi! I'm excited to get started training.", player.FirstName),
		},
	})

	for _, sess := range sessions {
		entries = append(entries, &InboxEntry{
			CreationDateEpochMillis: sess.CreationDateEpochMillis,
			Actor:                   coachActor,
			Type:                    NewSessionAdded,
			Metadata: map[string]string{
				"sessionNumber": fmt.Sprintf("%v", sess.SessionNumber),
			},
		})
		var feedbackDate int64 = 0
		for _, drill := range sess.Drills {
			if drill.Feedback != nil && drill.Feedback.UploadDateEpochMillis > feedbackDate {
				feedbackDate = drill.Feedback.UploadDateEpochMillis
			}
		}
		if feedbackDate > 0 {
			entries = append(entries, &InboxEntry{
				CreationDateEpochMillis: feedbackDate,
				Actor:                   coachActor,
				Type:                    FeedbackProvided,
				Metadata: map[string]string{
					"sessionNumber": fmt.Sprintf("%v", sess.SessionNumber),
				},
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreationDateEpochMillis > entries[j].CreationDateEpochMillis
	})

	return entries, nil
}
