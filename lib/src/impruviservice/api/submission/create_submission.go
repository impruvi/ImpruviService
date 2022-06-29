package submission

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"impruviService/api/converter"
	sessionUtil "impruviService/api/session"
	"impruviService/dao/coaches"
	"impruviService/dao/players"
	"impruviService/dao/session"
	"impruviService/notification"
	"log"
)

type CreateSubmissionRequest struct {
	PlayerId      string `json:"playerId"`
	SessionNumber int    `json:"sessionNumber"`
	DrillId       string `json:"drillId"`
}

func CreateSubmission(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request CreateSubmissionRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		return converter.BadRequest("Error unmarshalling request: %v\n", err)
	}

	err = session.CreateSubmission(request.SessionNumber, request.PlayerId, request.DrillId)
	if err != nil {
		return converter.InternalServiceError("Error while creating submission: %v\n", err)
	}
	thisSession, err := session.GetSession(request.PlayerId, request.SessionNumber)
	if err != nil {
		// don't fail the request just because we failed to send the notifications
		log.Printf("Error while getting session by id to notify on submission: %v\n", err)
	} else {
		if sessionUtil.IsSessionSubmissionComplete(thisSession) {
			sendNotifications(request.PlayerId)
		}
	}

	return converter.Success(nil)
}

func sendNotifications(playerId string) {
	player, err := players.GetPlayerById(playerId)
	if err != nil {
		// don't fail the request just because we failed to send the notifications
		log.Printf("Error while getting player by id to notify on submission: %v\n", err)
	}
	sendCoachTextNotification(player.FirstName)
	sendCoachPushNotification(player.CoachId, player.FirstName)
}

func sendCoachTextNotification(firstName string) {
	notification.Notify(fmt.Sprintf("%v completed a session!", firstName))
}

func sendCoachPushNotification(coachId string, playerName string) {
	coach, err := coaches.GetCoachById(coachId)
	if err != nil {
		// don't fail the request just because we failed to send the notifications
		log.Printf("Error while getting coach by id: %v\n", err)
	}
	if coach.NotificationId != "" {
		notification.Publish(
			fmt.Sprintf("%v completed a session!", playerName),
			fmt.Sprintf("You have 24 hours to submit feedback"),
			coach.NotificationId)
	}
}
