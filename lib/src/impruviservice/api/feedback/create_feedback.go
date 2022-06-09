package feedback

import (
	"../../dao/coaches"
	"../../dao/players"
	"../../dao/session"
	"../../notification"
	"../converter"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

type CreateFeedbackRequest struct {
	CoachId       string `json:"coachId"`
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
}

func CreateFeedback(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request CreateFeedbackRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = session.CreateFeedback(request.SessionNumber, request.CoachId, request.DrillId)
	if err != nil {
		return converter.InternalServiceError("Error while creating feedback: %v\n", err)
	}

	sendNotifications(request.CoachId, request.PlayerId)

	return converter.Success(nil)
}

func sendNotifications(coachId, playerId string) {
	coach, err := coaches.GetCoachById(coachId)
	player, err := players.GetPlayerById(playerId)
	if err == nil {
		notification.Notify(fmt.Sprintf("%v %v submitted feedback for %v %v", coach.FirstName, coach.LastName, player.FirstName, player.LastName))
	} else {
		// don't fail the request just because we failed to send the notifications
		log.Printf("Error while sending text notification: %v\n", err)
	}
}
