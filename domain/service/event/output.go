package event

import (
	"fmt"

	"github.com/GerogeGol/yadro-test-problem/domain/store"
)

var OutLeaveEventId = 11
var OutSitDownEventId = 12
var ErrorEventId = 13

type ErrorEvent struct {
	Event
	err error
}

func NewErrorEvent(t store.DayTime, err error) *ErrorEvent {
	return &ErrorEvent{&BaseEvent{t, ErrorEventId}, err}
}

func (e *ErrorEvent) Err() error {
	return e.err
}

func (e *ErrorEvent) String() string {
	return fmt.Sprintf("%s %d %s", e.Time(), e.Id(), e.Err())
}

func IsError(e Event) bool {
	_, ok := e.(*ErrorEvent)
	return ok
}

type OutSitDownEvent struct {
	Event
	client string
	table  int
}

func NewOutSitDownEvent(t store.DayTime, client string, table int) *OutSitDownEvent {
	return &OutSitDownEvent{Event: &BaseEvent{t, OutSitDownEventId}, client: client, table: table}
}

func (e *OutSitDownEvent) Table() int {
	return e.table
}

func (e *OutSitDownEvent) Client() string {
	return e.client
}

func (e *OutSitDownEvent) String() string {
	return fmt.Sprintf("%s %d %s %d", e.Time(), e.Id(), e.Client(), e.Table())
}

type OutLeaveEvent struct {
	Event
	client string
}

func NewOutLeaveEvent(t store.DayTime, client string) *OutLeaveEvent {
	return &OutLeaveEvent{Event: &BaseEvent{t, OutLeaveEventId}, client: client}
}

func (e *OutLeaveEvent) Client() string {
	return e.client
}

func (e *OutLeaveEvent) String() string {
	return fmt.Sprintf("%s %d %s", e.Time(), e.Id(), e.Client())
}
