package invitationcode

import (
	"impruviService/dao/coach"
	"impruviService/exceptions"
	coachFacade "impruviService/facade/coach"
	invitationCodeFacade "impruviService/facade/invitationcode"
	"log"
)

type ValidateInvitationCodeRequest struct {
	InvitationCode string `json:"invitationCode"`
	ExpoPushToken  string `json:"expoPushToken"`
}

type ValidateInvitationCodeResponse struct {
	Coach *coaches.CoachDB `json:"coach"`
}

func ValidateInvitationCode(request *ValidateInvitationCodeRequest) (*ValidateInvitationCodeResponse, error) {
	log.Printf("ValidateInvitationCodeRequest: %+v\n", request)
	err := validateValidateCodeRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid ValidateInvitationCodeRequest: %v\n", err)
		return nil, err
	}

	coach, err := invitationCodeFacade.ValidateCode(request.InvitationCode)
	if err != nil {
		return nil, err
	}

	if coach.NotificationId != request.ExpoPushToken {
		coach, err = coachFacade.UpdateNotificationId(coach, request.ExpoPushToken)
		if err != nil {
			return nil, err
		}
	}
	return &ValidateInvitationCodeResponse{
		Coach: coach,
	}, nil
}

func validateValidateCodeRequest(request *ValidateInvitationCodeRequest) error {
	if request.InvitationCode == "" {
		return exceptions.InvalidRequestError{Message: "InvitationCode cannot be null/empty"}
	}
	return nil
}
