package store

import "time"

type DayTime struct{ time.Time }

func NewDayTime(hour int, minute int) DayTime {
	return DayTime{time.Date(1, 1, 1, hour, minute, 0, 0, time.UTC)}
}

func (d DayTime) String() string {
	return d.Time.Format("15:04")
}
