package util

import (
	"../model"
	"time"
)

func GetCurrentTimeEpochMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetCurrentDate() *model.Date {
	year, month, day := time.Now().Date()
	return &model.Date{
		Year:  year,
		Month: int(month),
		Day:   day,
	}
}

// CompareDates returns -1 if date1 < date2, 0 if date1 == date2 and 1 if date1 > date2
func CompareDates(date1, date2 *model.Date) int {
	if date1.Year < date2.Year {
		return -1
	} else if date2.Year < date1.Year {
		return 1
	} else if date1.Month < date2.Month {
		return -1
	} else if date2.Month < date1.Month {
		return 1
	} else if date1.Day < date2.Day {
		return -1
	} else if date2.Day < date1.Day {
		return 1
	} else {
		return 0
	}
}

func ConvertDayOfWeekToInt(dayOfWeek string) int {
	switch dayOfWeek {
	case "SUNDAY":
		return 0
	case "MONDAY":
		return 1
	case "TUESDAY":
		return 2
	case "WEDNESDAY":
		return 3
	case "THURSDAY":
		return 4
	case "FRIDAY":
		return 5
	default:
		return 6
	}
}
