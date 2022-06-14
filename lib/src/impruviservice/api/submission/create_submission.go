package submission

import (
	"../../dao/players"
	"../../dao/session"
	"../../notification"
	"../converter"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
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

	sendNotifications(request.PlayerId)

	return converter.Success(nil)
}

func sendNotifications(playerId string) {
	player, err := players.GetPlayerById(playerId)
	if err == nil {
		notification.Notify(fmt.Sprintf("%v %v submitted a video!", player.FirstName, player.LastName))
		notification.Publish()
	} else {
		// don't fail the request just because we failed to send the notifications
		log.Printf("Error while getting user by id: %v\n", err)
	}
}
