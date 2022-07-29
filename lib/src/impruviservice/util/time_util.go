package util

import (
	"time"
)

const OneHourInMilliseconds int64 = 60 * 60 * 1000
const TwentyFourHoursInMilliseconds int64 = 24 * 60 * 60 * 1000
const TwelveHoursInSeconds int64 = 12 * 60 * 60

func GetCurrentTimeEpochMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
