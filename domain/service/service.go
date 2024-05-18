package service

import (
	"fmt"
	"sort"

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
		arrivalEvent, ok := e.(*event.ArriveEvent)
		if !ok {
			return event.NewErrorEvent(e.Time(), fmt.Errorf("Service.ServeEvent: cant interpret event to ArrivalEvent"))
		}
		err := s.cc.Arrive(arrivalEvent.Time(), arrivalEvent.Client())
		if err != nil {
			return event.NewErrorEvent(arrivalEvent.Time(), err)
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

func (s *Service) Close() ([]event.OutLeaveEvent, error) {
	clients, err := s.cc.Close()
	if err != nil {
		return nil, err
	}

	var events []event.OutLeaveEvent
	for _, c := range clients {
		events = append(events, *event.NewOutLeaveEvent(s.cc.CloseTime, c.Name))
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Client() <= events[j].Client()

	})
	return events, nil
}

func (s *Service) Profit() ([]TableInfo, error) {
	return s.cc.TablesInfo()
}
