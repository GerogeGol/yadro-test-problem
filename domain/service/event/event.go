package event

import "github.com/GerogeGol/yadro-test-problem/domain/store"

type Event interface {
	Time() store.DayTime
	Id() int
}

type BaseEvent struct {
	time store.DayTime
	id   int
}

func (e *BaseEvent) Time() store.DayTime {
	return e.time
}

func (e *BaseEvent) Id() int {
	return e.id
}

var EmptyEvent = &BaseEvent{}

func IsEmpty(e Event) bool {
	return e.Id() == 0
}
