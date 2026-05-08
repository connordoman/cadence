package interpreter

import (
	"time"
)

func firstWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, loc *time.Location) time.Time {
	d := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	offset := (int(weekday) - int(d.Weekday()) + 7) % 7
	return d.AddDate(0, 0, offset)
}

func lastWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, loc *time.Location) time.Time {
	d := time.Date(year, month+1, 0, 0, 0, 0, 0, loc)
	offset := (int(d.Weekday()) - int(weekday) + 7) % 7
	return d.AddDate(0, 0, -offset)
}
