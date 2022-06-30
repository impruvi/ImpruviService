package feedback

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

	err = session.CreateFeedback(request.SessionNumber, request.PlayerId, request.DrillId)
	if err != nil {
		return converter.InternalServiceError("Error while creating feedback: %v\n", err)
	}

	thisSession, err := session.GetSession(request.PlayerId, request.SessionNumber)
	if err != nil {
		// don't fail the request just because we failed to send the notifications
		log.Printf("Error while getting session by id to notify on submission: %v\n", err)
	} else {
		if sessionUtil.IsSessionFeedbackComplete(thisSession) {
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
	coach, err := coaches.GetCoachById(player.CoachId)
	if err != nil {
		// don't fail the request just because we failed to send the notifications
		log.Printf("Error while getting coach by id to notify on submission: %v\n", err)
	}
	// Add this if we collect player numbers and they authorize text messages
	sendPlayerTextNotification(coach.FirstName)
	sendPlayerPushNotification(coach.FirstName, player.NotificationId)
}

func sendPlayerTextNotification(coachName string) {
	notification.Notify(fmt.Sprintf("Coach %v submitted feedback on your session!", coachName))
}

func sendPlayerPushNotification(coachName string, notificationId string) {
	if notificationId != "" {
		notification.Publish(
			fmt.Sprintf("Coach %v submitted feedback!", coachName),
			fmt.Sprintf("Review your feedback before the next session"),
			notificationId)
	}
}
