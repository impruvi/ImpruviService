package invitationcode

import (
	"fmt"
	coachDao "impruviService/dao/coach"
	invitationCodeDao "impruviService/dao/invitationcode"
	"impruviService/exceptions"
	coachFacade "impruviService/facade/coach"
	"impruviService/model"
	"log"
	"strings"
)

func ValidateCode(invitationCode string) (*coachDao.CoachDB, error) {
	invitationCode = strings.TrimSpace(invitationCode)
	invitationCodeEntry, err := invitationCodeDao.GetInvitationCodeEntry(invitationCode)

	if err != nil {
		log.Printf("Error: %v\n", err)
		if _, ok := err.(exceptions.ResourceNotFoundError); ok {
			return nil, exceptions.NotAuthorizedError{Message: fmt.Sprintf("Invalid invitation code: %v\n", invitationCode)}
		} else {
			return nil, err
		}
	}
	log.Printf("Invitation code entry: %v\n", invitationCodeEntry)

	if invitationCodeEntry.UserType != model.Coach {
		log.Printf("Invitation codes are no longer supported for users of type: %v. user: %+v\n", invitationCodeEntry.UserType, invitationCodeEntry)
		return nil, exceptions.NotAuthorizedError{Message: fmt.Sprintf("Invalid invitation code: %v\n", invitationCode)}
	}

	return coachFacade.GetCoachById(invitationCodeEntry.UserId)
}
