package subscription

import (
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"log"
)


type GetSubscriptionHistoryRequest struct {
	PlayerId string `json:"playerId"`
}

type GetSubscriptionHistoryResponse struct {
	Subscriptions []*stripeFacade.Subscription `json:"subscriptions"`
}

func GetSubscriptionHistory(request *GetSubscriptionHistoryRequest) (*GetSubscriptionHistoryResponse, error) {
	player, err := playerFacade.GetPlayerById(request.PlayerId)
	if err != nil {
		log.Printf("Error getting player by id: %v. Error: %v\n", request.PlayerId, err)
	}
	subscriptions, err := stripeFacade.ListSubscriptions(player.StripeCustomerId)
	if err != nil {
		log.Printf("Error listing subscriptions for player: %+v. Error: %v\n", player, err)
	}

	return &GetSubscriptionHistoryResponse{Subscriptions: subscriptions}, nil
}