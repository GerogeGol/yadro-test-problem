package service

import (
	"fmt"

	"github.com/GerogeGol/yadro-test-problem/domain/service/event"
)

type Service struct {
	cc *ComputerClub
}

func NewService(cc *ComputerClub) *Service {
	return &Service{cc: cc}
}

func (s *Service) ServeEvent(e event.InputEvent) event.Event {
	switch e.Id() {
	case event.ArrivalEventId:
		err := s.cc.Arrive(e.Time(), e.Client())
		if err != nil {
			return event.NewErrorEvent(e.Time(), err)
		}
	case event.SitDownEventId:
		sitDownEvent, ok := e.(*event.SitDownEvent)
		if !ok {
			return event.NewErrorEvent(e.Time(), fmt.Errorf("Service.ServeEvent: cant interpret event to SitDownEvent"))
		}

		err := s.cc.SitDown(sitDownEvent.Time(), sitDownEvent.Client(), sitDownEvent.Table())
		if err != nil {
			return event.NewErrorEvent(e.Time(), err)
		}

		return event.NewOutSitDownEvent(sitDownEvent.Time(), sitDownEvent.Client(), sitDownEvent.Table())
	case event.WaitEventId:
		waitEvent, ok := e.(*event.WaitEvent)
		if !ok {
			return event.NewErrorEvent(e.Time(), fmt.Errorf("Service.ServeEvent: cant interpret event to WaitEvent"))
		}

		isWaiting, err := s.cc.Wait(waitEvent.Time(), waitEvent.Client())
		if err != nil {
			return event.NewErrorEvent(e.Time(), err)
		}

		if !isWaiting {
			return event.NewOutLeaveEvent(waitEvent.Time(), waitEvent.Client())
		}
	case event.LeaveEventId:
		leaveEvent, ok := e.(*event.LeaveEvent)
		if !ok {
			return event.NewErrorEvent(e.Time(), fmt.Errorf("Service.ServeEvent: cant interpret event to LeaveEvent"))
		}

		client, occupied, err := s.cc.Leave(leaveEvent.Time(), leaveEvent.Client())
		if err != nil {
			return event.NewErrorEvent(e.Time(), err)
		}

		if occupied {
			return event.NewOutSitDownEvent(e.Time(), client.Name, client.Table)
		}
	}
	return event.EmptyEvent
}
