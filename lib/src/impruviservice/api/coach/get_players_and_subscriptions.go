package coach

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"log"
	"sync"
)

type GetPlayersAndSubscriptionsRequest struct {
	CoachId string `json:"coachId"`
}

type GetPlayersAndSubscriptionsResponse struct {
	PlayerAndSubscriptions []*PlayerAndSubscription `json:"playersAndSubscriptions"`
}

type PlayerAndSubscription struct {
	Player       *playerFacade.Player       `json:"player"`
	Subscription *stripeFacade.Subscription `json:"subscription"`
}

func GetPlayersAndSubscriptions(request *GetPlayersAndSubscriptionsRequest) (*GetPlayersAndSubscriptionsResponse, error) {
	log.Printf("GetPlayersAndSubscriptionsRequest: %+v\n", request)
	err := validateGetPlayersRequiringTrainingsRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetPlayersAndSubscriptionsRequest: %v\n", err)
		return nil, err
	}

	players, err := playerFacade.GetPlayersForCoach(request.CoachId)
	if err != nil {
		return nil, err
	}

	playerAndSubscriptions := make([]*PlayerAndSubscription, 0)
	playerAndSubscriptionsChan := make(chan *PlayerAndSubscription)

	var wg1 sync.WaitGroup
	wg1.Add(len(players))
	for _, player := range players {
		go func(player *playerFacade.Player) {
			defer wg1.Done()
			subscription, err := stripeFacade.GetSubscription(player.StripeCustomerId)
			if err != nil {
				log.Printf("Error while getting subscription for player: %+v. error %v\n", player, err)
				// TODO: handle this more gracefully
				panic(err)
			}

			playerAndSubscriptionsChan <- &PlayerAndSubscription{
				Player:       player,
				Subscription: subscription,
			}
		}(player)
	}

	go func() {
		wg1.Wait()
		close(playerAndSubscriptionsChan)
	}()

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for playerAndSubscription := range playerAndSubscriptionsChan {
			playerAndSubscriptions = append(playerAndSubscriptions, playerAndSubscription)
		}
	}()
	wg2.Wait()

	return &GetPlayersAndSubscriptionsResponse{
		PlayerAndSubscriptions: playerAndSubscriptions,
	}, nil
}

func validateGetPlayersRequiringTrainingsRequest(request *GetPlayersAndSubscriptionsRequest) error {
	if request.CoachId == "" {
		return exceptions.InvalidRequestError{Message: "CoachId cannot be null/empty"}
	}
	return nil
}
