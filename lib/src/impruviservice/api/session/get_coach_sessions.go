package session

import (
	"../../dao/users"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type GetCoachSessionsRequest struct {
	UserId string `json:"userId"`
}

type GetCoachSessionsResponse struct {
	PlayerSessions []*PlayerSessions `json:"playerSessions"`
}

type PlayerSessions struct {
	User     *users.User `json:"user"`
	Sessions []*Session  `json:"sessions"`
}

func GetCoachSessions(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	var request GetPlayerSessionsRequest
	var err = json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		log.Printf("Error unmarshalling request: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	playersForCoach, err := users.GetPlayersForCoach(request.UserId)
	log.Printf("Players for coach: %v\n", playersForCoach)

	playerSessions := make([]*PlayerSessions, 0)
	for _, player := range playersForCoach {
		sessionsWithDrills, err := getSessionsWithDrillsForUser(player.UserId)
		if err != nil {
			log.Printf("Error while getting sessions with drills for user: %v. %v\n", player, err)
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}
		}

		playerSessions = append(playerSessions, &PlayerSessions{
			User:     player,
			Sessions: sessionsWithDrills,
		})
	}

	rspBody, err := json.Marshal(GetCoachSessionsResponse{
		PlayerSessions: playerSessions,
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
