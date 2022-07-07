package inbox

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	"impruviService/dao/players"
	"impruviService/dao/session"
	coachFacade "impruviService/facade/coach"
	"sort"
)

type GetInboxForPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetInboxForPlayerResponse struct {
	Entries []*Entry `json:"entries"`
}

type EntryType string

const (
	NewSessionAdded  EntryType = "NEW_SESSION_ADDED"
	FeedbackProvided EntryType = "FEEDBACK_PROVIDED"
	Message          EntryType = "MESSAGE"
)

type Entry struct {
	CreationDateEpochMillis int64             `json:"creationDateEpochMillis"`
	Actor                   *Actor            `json:"actor"`
	Type                    EntryType         `json:"type"`
	Metadata                map[string]string `json:"metadata"`
}

type Actor struct {
	Name              string `json:"name"`
	ImageFileLocation string `json:"image"`
}

func GetInboxForPlayer(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetInboxForPlayerRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	sessions, err := session.GetSessions(request.PlayerId)
	if err != nil {
		return converter.InternalServiceError("Error getting inbox for player: %v\n: %v\n", request.PlayerId, err)
	}

	player, err := players.GetPlayerById(request.PlayerId)
	if err != nil {
		return converter.InternalServiceError("Error getting player by playerId: %v\n: %v\n", request.PlayerId, err)
	}
	coach, err := coachFacade.GetCoachById(player.CoachId)
	coachActor := &Actor{
		Name:              fmt.Sprintf(coach.FirstName),
		ImageFileLocation: coach.Headshot.FileLocation,
	}

	entries := make([]*Entry, 0)
	entries = append(entries, &Entry{
		CreationDateEpochMillis: player.CreationDateEpochMillis,
		Actor:                   coachActor,
		Type:                    Message,
		Metadata: map[string]string{
			"messageContent": fmt.Sprintf("Hey %v, welcome to Impruvi! I'm excited to get started training.", player.FirstName),
		},
	})

	for _, sess := range sessions {
		entries = append(entries, &Entry{
			CreationDateEpochMillis: sess.CreationDateEpochMillis,
			Actor:                   coachActor,
			Type:                    NewSessionAdded,
			Metadata: map[string]string{
				"sessionNumber": fmt.Sprintf("%v", sess.SessionNumber),
			},
		})
		var feedbackDate int64 = 0
		for _, drill := range sess.Drills {
			if drill.Feedback != nil && drill.Feedback.VideoUploadDateEpochMillis > feedbackDate {
				feedbackDate = drill.Feedback.VideoUploadDateEpochMillis
			}
		}
		if feedbackDate > 0 {
			entries = append(entries, &Entry{
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

	return converter.Success(GetInboxForPlayerResponse{
		Entries: entries,
	})
}
