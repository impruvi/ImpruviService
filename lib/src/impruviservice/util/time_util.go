package util

import "time"

func GetCurrentTimeEpochMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
