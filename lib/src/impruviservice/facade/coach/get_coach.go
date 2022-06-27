package coach

import (
	"../../dao/coaches"
	"../../files"
	"../../model"
)

func GetCoachById(coachId string) (*coaches.Coach, error) {
	coach, err := coaches.GetCoachById(coachId)
	if err != nil {
		return nil, err
	}
	if coach.Headshot != nil && coach.Headshot.UploadDateEpochMillis > 0 {
		coach.Headshot.FileLocation = files.GetHeadshotFileLocation(model.Coach, coachId).URL
	}
	return coach, nil
}
