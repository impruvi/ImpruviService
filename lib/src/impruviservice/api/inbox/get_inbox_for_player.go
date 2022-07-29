package inbox

import (
	"impruviService/exceptions"
	inboxFacade "impruviService/facade/inbox"
	"log"
)

type GetInboxForPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetInboxForPlayerResponse struct {
	Entries []*inboxFacade.InboxEntry `json:"entries"`
}

func GetInboxForPlayer(request *GetInboxForPlayerRequest) (*GetInboxForPlayerResponse, error) {
	log.Printf("GetInboxForPlayerRequest: %+v\n", request)
	err := validateGetInboxForPlayerRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetInboxForPlayerRequest: %v\n", err)
		return nil, err
	}

	entries, err := inboxFacade.GetInboxForPlayer(request.PlayerId)
	if err != nil {
		return nil, err
	}
	return &GetInboxForPlayerResponse{
		Entries: entries,
	}, nil
}

func validateGetInboxForPlayerRequest(request *GetInboxForPlayerRequest) error {
	if request.PlayerId == "" {
		return exceptions.InvalidRequestError{Message: "PlayerId cannot be null/empty"}
	}
	return nil
}
