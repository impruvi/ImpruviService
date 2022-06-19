package model

import "time"

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

func (d *Date) GetNextDate() *Date {
	date := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.Local)
	nextDate := date.AddDate(0, 0, 1)
	return &Date{
		Year:  nextDate.Year(),
		Month: int(nextDate.Month()),
		Day:   nextDate.Day(),
	}
}

func (d *Date) GetDayOfWeek() int {
	date := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.Local)
	return int(date.Weekday())
}
