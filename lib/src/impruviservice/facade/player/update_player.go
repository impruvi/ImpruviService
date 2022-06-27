package profile

import (
	"../../dao/players"
	"../../dao/session"
	"../../model"
	"../../util"
	"errors"
	"log"
	"sort"
)

func UpdatePlayer(player *players.Player) error {
	currentPlayer, err := players.GetPlayerById(player.PlayerId)
	if err != nil {
		return err
	}

	err = players.PutPlayer(player)
	if err != nil {
		return err
	}

	log.Printf("Old availability: %v\n", currentPlayer.Availability)
	log.Printf("New availability: %v\n", player.Availability)

	if didAvailabilityChange(currentPlayer.Availability, player.Availability) {
		log.Printf("Availability changed!")
		if len(player.Availability) == 0 {
			log.Printf("Availability cannot be empty!")
			return errors.New("availability cannot be empty")
		}

		return UpdateSessionDatesForPlayer(player)
	}

	return nil
}

func UpdateSessionDatesForPlayerByPlayerId(playerId string) error {
	player, err := players.GetPlayerById(playerId)
	if err != nil {
		return err
	}
	return UpdateSessionDatesForPlayer(player)
}

func UpdateSessionDatesForPlayer(player *players.Player) error {
	sessions, err := session.GetSessions(player.PlayerId)
	if err != nil {
		return err
	}

	log.Printf("Sessions: %v\n", sessions)
	sessionsToChange := make([]*session.Session, 0)
	for _, sess := range sessions {
		if shouldChangeSessionDate(sess) {
			log.Printf("Should change session: %v\n", sess)
			sessionsToChange = append(sessionsToChange, sess)
		}
	}

	log.Printf("Sessions to change: %v\n", sessionsToChange)
	return updateSessionDates(sessionsToChange, hasSessionStartedToday(sessions), player.Availability)
}

func updateSessionDates(sessionsToChange []*session.Session, hasStartedSessionToday bool, availability []string) error {
	log.Printf("Updating session dates: %v, hasStartedSessionToday: %v, availability: %v\n", sessionsToChange, hasStartedSessionToday, availability)

	availabilitySet := make(map[int]bool, 0)
	for _, availableDay := range availability {
		availabilitySet[util.ConvertDayOfWeekToInt(availableDay)] = true
	}

	log.Printf("Availability set: %v\n", availabilitySet)

	var nextAvailableDate *model.Date
	if hasStartedSessionToday {
		nextAvailableDate = getNextAvailableDate(util.GetCurrentDate().GetNextDate(), availabilitySet)
	} else {
		nextAvailableDate = getNextAvailableDate(util.GetCurrentDate(), availabilitySet)
	}

	log.Printf("Next available date: %v\n", nextAvailableDate)

	sort.Slice(sessionsToChange, func(i, j int) bool {
		return sessionsToChange[i].SessionNumber < sessionsToChange[j].SessionNumber
	})

	log.Printf("sessionsToChange sorted: %v\n", sessionsToChange)

	for _, sess := range sessionsToChange {
		sess.Date = nextAvailableDate
		log.Printf("Updating session: %v\n", sess)
		err := session.PutSession(sess)
		if err != nil {
			return err
		}
		nextAvailableDate = getNextAvailableDate(nextAvailableDate.GetNextDate(), availabilitySet)
		log.Printf("Next available date: %v\n", nextAvailableDate)
	}

	return nil
}

func getNextAvailableDate(date *model.Date, availability map[int]bool) *model.Date {
	_, ok := availability[date.GetDayOfWeek()]
	for !ok {
		date = date.GetNextDate()
		_, ok = availability[date.GetDayOfWeek()]
	}
	return date
}

func shouldChangeSessionDate(sess *session.Session) bool {
	if sess.Date == nil {
		return true
	}

	currentDate := util.GetCurrentDate()
	log.Printf("Current date: %v\n", currentDate)

	// session date is before current date. Don't change date of sessions that already passed
	if util.CompareDates(sess.Date, currentDate) == -1 {
		log.Printf("Date is before current date: %v\n", sess.Date)
		return false
	}

	// if player has already started session, don't change date of session
	if hasStartedSession(sess) {
		log.Printf("Has started session: %v\n", sess.Date)
		return false
	}

	return true
}

func hasSessionStartedToday(sessions []*session.Session) bool {
	currentDate := util.GetCurrentDate()
	for _, sess := range sessions {
		if sess.Date != nil && util.CompareDates(sess.Date, currentDate) == 0 && hasStartedSession(sess) {
			return true
		}
	}
	return false
}

func didAvailabilityChange(oldAvailability, newAvailability []string) bool {
	oldAvailabilitySet := make(map[string]bool, 0)
	newAvailabilitySet := make(map[string]bool, 0)

	for _, dayOfWeek := range oldAvailability {
		oldAvailabilitySet[dayOfWeek] = true
	}
	for _, dayOfWeek := range newAvailability {
		newAvailabilitySet[dayOfWeek] = true
	}

	if len(oldAvailabilitySet) != len(newAvailabilitySet) {
		return true
	}
	for dayOfWeek := range newAvailabilitySet {
		if _, ok := oldAvailabilitySet[dayOfWeek]; !ok {
			return true
		}
	}
	return false
}

func hasStartedSession(sess *session.Session) bool {
	for _, drill := range sess.Drills {
		if drill.Submission != nil && drill.Submission.VideoUploadDateEpochMillis > 0 {
			return true
		}
	}
	return false
}
