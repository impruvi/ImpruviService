package inbox

import (
	"fmt"
	sessionDao "impruviService/dao/session"
	coachFacade "impruviService/facade/coach"
	playerFacade "impruviService/facade/player"
	"impruviService/model"
	"sort"
)

func GetInboxForPlayer(playerId string) ([]*model.InboxEntry, error) {
	sessions, err := sessionDao.GetSessions(playerId)
	if err != nil {
		return nil, err
	}

	player, err := playerFacade.GetPlayerById(playerId)
	if err != nil {
		return nil, err
	}
	coach, err := coachFacade.GetCoachById(player.CoachId)
	coachActor := &model.InboxEntryActor{
		FirstName:         fmt.Sprintf(coach.FirstName),
		LastName:          fmt.Sprintf(coach.LastName),
		ImageFileLocation: coach.Headshot.FileLocation,
	}

	entries := make([]*model.InboxEntry, 0)
	entries = append(entries, &model.InboxEntry{
		CreationDateEpochMillis: player.CreationDateEpochMillis,
		Actor:                   coachActor,
		Type:                    model.Message,
		Metadata: map[string]string{
			"messageContent": fmt.Sprintf("Hey %v, welcome to Impruvi! I'm excited to get started training.", player.FirstName),
		},
	})

	for _, sess := range sessions {
		entries = append(entries, &model.InboxEntry{
			CreationDateEpochMillis: sess.CreationDateEpochMillis,
			Actor:                   coachActor,
			Type:                    model.NewSessionAdded,
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
			entries = append(entries, &model.InboxEntry{
				CreationDateEpochMillis: feedbackDate,
				Actor:                   coachActor,
				Type:                    model.FeedbackProvided,
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
