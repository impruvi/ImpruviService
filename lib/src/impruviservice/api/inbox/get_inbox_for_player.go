package inbox

import (
	inboxFacade "impruviService/facade/inbox"
)

type GetInboxForPlayerRequest struct {
	PlayerId string `json:"playerId"`
}

type GetInboxForPlayerResponse struct {
	Entries []*inboxFacade.InboxEntry `json:"entries"`
}

func GetInboxForPlayer(request *GetInboxForPlayerRequest) (*GetInboxForPlayerResponse, error) {
	entries, err := inboxFacade.GetInboxForPlayer(request.PlayerId)
	if err != nil {
		return nil, err
	}
	return &GetInboxForPlayerResponse{
		Entries: entries,
	}, nil
}
