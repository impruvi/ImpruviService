package dynamic

import (
	"impruviService/util"
	"math"
)

// We send reminders at the 12-hour mark and the 1-hour mark. We then want another invocation
// of this function in order to notify us of any coaches who miss the feedback deadline
func getNewWaitSeconds(hoursRemaining int, lastSubmissionUploadDateEpochMillis int64) int64 {
	overdueTimeEpochMillis := lastSubmissionUploadDateEpochMillis + util.TwentyFourHoursInMilliseconds
	currentTimeEpochMillis := util.GetCurrentTimeEpochMillis()
	if hoursRemaining > 1 {
		return (overdueTimeEpochMillis - util.OneHourInMilliseconds - currentTimeEpochMillis) / 1000
	} else {
		return (overdueTimeEpochMillis - currentTimeEpochMillis) / 1000
	}
}

func getHoursRemaining(periodStartTimeEpochMillis int64) int {
	currentTimeEpochMillis := util.GetCurrentTimeEpochMillis()
	millisecondsRemaining := periodStartTimeEpochMillis + util.TwentyFourHoursInMilliseconds - currentTimeEpochMillis
	return int(math.Round(float64(millisecondsRemaining) / float64(util.OneHourInMilliseconds)))
}
