package session

import (
	"../../dao/drills"
	"../../dao/session"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type GetSessionsRequest struct {
	UserId string `json:"userId"`
}

type GetSessionsResponse struct {
	Sessions []*Session `json:"sessions"`
}

type Session struct {
	UserId string `json:"userId"`
	SessionNumber int `json:"sessionNumber"`
	Drills []*Drill `json:"drills"`
}

type Drill struct {
	Drill drills.Drill `json:"drill"`
	Submission *session.Submission `json:"submission"`
	Feedback *session.Feedback `json:"feedback"`
	Tips []string `json:"tips"`
	Repetitions int `json:"repetitions"`
	DurationMinutes int `json:"durationMinutes"`
}

func GetSessions(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetSessionsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	sessions, err := session.GetSessions(request.UserId)
	if err != nil {
		log.Printf("Error while getting sessions: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	sessionsWithDrills, err := getSessionsWithDrills(sessions)
	if err != nil {
		log.Printf("Error while getting drills for sessions: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	rspBody, err := json.Marshal(GetSessionsResponse{
		Sessions: sessionsWithDrills,
	})
	if err != nil {
		log.Printf("Error while marshalling response: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &events.APIGatewayProxyResponse{
		Body:       string(rspBody),
		StatusCode: http.StatusAccepted,
	}
}

func getSessionsWithDrills(sessions []*session.Session) ([]*Session, error) {
	sessionsWithDrills := make([]*Session, 0)
	for _, sess := range sessions {
		sessionWithDrill, err := getSessionWithDrills(sess)
		if err != nil {
			return nil, err
		}
		sessionsWithDrills = append(sessionsWithDrills, sessionWithDrill)
	}
	return sessionsWithDrills, nil
}

func getSessionWithDrills(sess *session.Session) (*Session, error) {
	drillIds := getDrillIds(sess.Drills)
	drillDetails, err := drills.BatchGetDrills(drillIds)
	if err != nil {
		return nil, err
	}

	fullDrills := make([]*Drill, 0)
	for _, drill := range sess.Drills {
		fullDrills = append(fullDrills, &Drill{
			Drill:       *drillDetails[drill.DrillId],
			Submission:  drill.Submission,
			Feedback:    drill.Feedback,
			Tips:        drill.Tips,
			Repetitions: drill.Repetitions,
			DurationMinutes: drill.DurationMinutes,
		})
	}

	return &Session{
		UserId:        sess.UserId,
		SessionNumber: sess.SessionNumber,
		Drills:        fullDrills,
	}, nil
}

func getDrillIds(drills []*session.Drill) []string {
	drillIds := make([]string, 0)
	for _, drill := range drills {
		drillIds = append(drillIds, drill.DrillId)
	}
	return drillIds
}
