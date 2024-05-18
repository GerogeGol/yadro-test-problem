package store

import "time"

type DayTime struct{ time.Time }

func NewDayTime(hour int, minute int) DayTime {
	return DayTime{time.Date(0, 0, 0, hour, minute, 0, 0, time.Local)}
}

func (d DayTime) String() string {
	return d.Time.Format("15:04")
}
