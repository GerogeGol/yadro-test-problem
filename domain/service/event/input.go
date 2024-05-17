package event

import (
	"github.com/GerogeGol/yadro-test-problem/domain/store"
)

var ArrivalEventId = 1
var SitDownEventId = 2
var WaitEventId = 3
var LeaveEventId = 4

type InputEvent interface {
	Event
	Client() string
}

type inputEvent struct {
	Event
	client string
}

func newInputEvent(t store.DayTime, id int, client string) *inputEvent {
	return &inputEvent{
		Event:  &BaseEvent{t, id},
		client: client,
	}
}

func (e *inputEvent) Client() string {
	return e.client
}

type ArriveEvent struct {
	InputEvent
}

func NewArrivalEvent(t store.DayTime, client string) *ArriveEvent {
	return &ArriveEvent{
		InputEvent: newInputEvent(t, ArrivalEventId, client),
	}
}

type SitDownEvent struct {
	InputEvent
	table int
}

func NewSitDownEvent(t store.DayTime, client string, table int) *SitDownEvent {
	return &SitDownEvent{
		InputEvent: newInputEvent(t, SitDownEventId, client),
		table:      table,
	}
}

func (e *SitDownEvent) Table() int {
	return e.table
}

type WaitEvent struct {
	InputEvent
}

func NewWaitEvent(t store.DayTime, client string) *WaitEvent {
	return &WaitEvent{InputEvent: newInputEvent(t, WaitEventId, client)}
}

type LeaveEvent struct {
	InputEvent
}

func NewLeaveEvent(t store.DayTime, client string) *LeaveEvent {
	return &LeaveEvent{InputEvent: newInputEvent(t, LeaveEventId, client)}
}
