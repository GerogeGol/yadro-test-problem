package service_test

import (
	"fmt"
	"math"
	"testing"

	memqueue "github.com/GerogeGol/yadro-test-problem/domain/queue/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/service"
	"github.com/GerogeGol/yadro-test-problem/domain/service/event"
	"github.com/GerogeGol/yadro-test-problem/domain/store"
	memstore "github.com/GerogeGol/yadro-test-problem/domain/store/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/test"
)

func TestArriveEvent(t *testing.T) {
	t.Run("client arrives", func(t *testing.T) {
		s := service.NewService(dummyClub())

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		test.AssertTrue(t, event.IsEmpty(gotEvent))
	})
	t.Run("client arrives at closed time ", func(t *testing.T) {
		s := service.NewService(dummyClub())

		arrivalEvent := event.NewArrivalEvent(store.NewDayTime(21, 0), dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertErrorEvent(t, gotEvent, service.NotOpenYet)
	})
	t.Run("client arrives 2 times", func(t *testing.T) {
		s := service.NewService(dummyClub())

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		test.AssertTrue(t, event.IsEmpty(gotEvent))

		gotEvent = s.ServeEvent(arrivalEvent)
		assertErrorEvent(t, gotEvent, service.YouShallNotPass)
	})
}

func TestSitDownEvent(t *testing.T) {
	t.Run("client arrives and sit down", func(t *testing.T) {
		s := service.NewService(dummyClub())

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitDownEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitDownEvent)
		assertNoErrorEvent(t, gotEvent)
		assertEmptyEvent(t, gotEvent)
	})
	t.Run("client arrives, sits down and changes sit", func(t *testing.T) {
		s := service.NewService(service.NewComputerClub(2, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitDownEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, 1)
		gotEvent = s.ServeEvent(sitDownEvent)
		assertNoErrorEvent(t, gotEvent)

		sitDownEvent = event.NewSitDownEvent(dummyDayTime, dummyClient, 2)
		gotEvent = s.ServeEvent(sitDownEvent)
		assertNoErrorEvent(t, gotEvent)
		assertEmptyEvent(t, gotEvent)
	})

	t.Run("client arrives, sits down to busy table", func(t *testing.T) {
		s := service.NewService(dummyClub())

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitDownEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitDownEvent)
		assertNoErrorEvent(t, gotEvent)

		newClient := "NewClient"
		arrivalEvent = event.NewArrivalEvent(dummyDayTime, newClient)
		gotEvent = s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitDownEvent = event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitDownEvent)
		assertErrorEvent(t, gotEvent, service.PlaceIsBusy)
	})

	t.Run("client that not in club sit downt. should get ClientUnknown", func(t *testing.T) {
		s := service.NewService(dummyClub())

		sitDownEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent := s.ServeEvent(sitDownEvent)
		assertErrorEvent(t, gotEvent, service.ClientUnknown)
	})
}

func TestWaitEvent(t *testing.T) {
	t.Run("3 clients. every one waits after another. first sits", func(t *testing.T) {
		s := service.NewService(service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		waitEvent := event.NewWaitEvent(dummyDayTime, dummyClient)
		gotEvent = s.ServeEvent(waitEvent)
		assertErrorEvent(t, gotEvent, service.ICanWaitNoLonger)

		sitDownEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitDownEvent)
		assertNoErrorEvent(t, gotEvent)

		newClient := "NewClient"
		arrivalEvent = event.NewArrivalEvent(dummyDayTime, newClient)
		gotEvent = s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		waitEvent = event.NewWaitEvent(dummyDayTime, newClient)
		gotEvent = s.ServeEvent(waitEvent)
		assertNoErrorEvent(t, gotEvent)
		test.AssertTrue(t, event.IsEmpty(gotEvent))

		newClient2 := "NewClient2"
		arrivalEvent = event.NewArrivalEvent(dummyDayTime, newClient2)
		gotEvent = s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		waitEvent = event.NewWaitEvent(dummyDayTime, newClient2)
		gotEvent = s.ServeEvent(waitEvent)

		outLeaveEvent, ok := gotEvent.(*event.OutLeaveEvent)
		test.AssertTrue(t, ok)
		test.AssertEqual(t, outLeaveEvent.Id(), event.OutLeaveEventId)
		test.AssertEqual(t, outLeaveEvent.Client(), newClient2)
	})

}

func TestLeaveEvent(t *testing.T) {
	t.Run("client just leave", func(t *testing.T) {
		s := service.NewService(service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		leaveEvent := event.NewLeaveEvent(dummyDayTime, dummyClient)
		gotEvent = s.ServeEvent(leaveEvent)
		assertEmptyEvent(t, gotEvent)
	})

	t.Run("client arrive, sit down and leave. his place taken by other client", func(t *testing.T) {
		s := service.NewService(service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitEvent)
		assertNoErrorEvent(t, gotEvent)

		leaveEvent := event.NewLeaveEvent(dummyDayTime, dummyClient)
		gotEvent = s.ServeEvent(leaveEvent)
		assertEmptyEvent(t, gotEvent)

		newClient := "NewClient"
		arrivalEvent = event.NewArrivalEvent(dummyDayTime, newClient)
		gotEvent = s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitEvent = event.NewSitDownEvent(dummyDayTime, newClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitEvent)
		assertNoErrorEvent(t, gotEvent)
	})

	t.Run("client arrive, sit down and leave. his place taken by client in queue", func(t *testing.T) {
		s := service.NewService(service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitEvent)
		assertNoErrorEvent(t, gotEvent)

		newClient := "NewClient"
		arrivalEvent = event.NewArrivalEvent(dummyDayTime, newClient)
		gotEvent = s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		waitClient := event.NewWaitEvent(dummyDayTime, newClient)
		gotEvent = s.ServeEvent(waitClient)
		assertNoErrorEvent(t, gotEvent)

		leaveEvent := event.NewLeaveEvent(dummyDayTime, dummyClient)
		gotEvent = s.ServeEvent(leaveEvent)

		outSitDownEvent, ok := gotEvent.(*event.OutSitDownEvent)
		test.AssertTrue(t, ok)
		test.AssertEqual(t, outSitDownEvent.Client(), newClient)
		test.AssertEqual(t, outSitDownEvent.Table(), dummyTableNumber)
	})
}

func TestServiceClose(t *testing.T) {
	t.Run("no clients stayed unitl closing", func(t *testing.T) {
		s := service.NewService(dummyClub())

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitEvent)
		assertNoErrorEvent(t, gotEvent)

		leaveEvent := event.NewLeaveEvent(dummyDayTime, dummyClient)
		gotEvent = s.ServeEvent(leaveEvent)
		assertNoErrorEvent(t, gotEvent)

		clients, err := s.Close()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(clients), 0)
	})

	t.Run("one client stayed unitl closing", func(t *testing.T) {
		s := service.NewService(dummyClub())

		arrivalEvent := event.NewArrivalEvent(dummyDayTime, dummyClient)
		gotEvent := s.ServeEvent(arrivalEvent)
		assertNoErrorEvent(t, gotEvent)

		sitEvent := event.NewSitDownEvent(dummyDayTime, dummyClient, dummyTableNumber)
		gotEvent = s.ServeEvent(sitEvent)
		assertNoErrorEvent(t, gotEvent)

		clients, err := s.Close()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(clients), 1)
	})

	t.Run("clients stayed unitl closing", func(t *testing.T) {
		s := service.NewService(service.NewComputerClub(9, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))
		clientsCount := 9
		var clientNames []string
		for i := 1; i <= clientsCount; i++ {
			client := fmt.Sprintf("Client%d", i)

			arrivalEvent := event.NewArrivalEvent(dummyDayTime, client)
			gotEvent := s.ServeEvent(arrivalEvent)
			assertNoErrorEvent(t, gotEvent)

			sitEvent := event.NewSitDownEvent(dummyDayTime, client, i)
			gotEvent = s.ServeEvent(sitEvent)
			assertNoErrorEvent(t, gotEvent)
			clientNames = append(clientNames, client)
		}

		clients, err := s.Close()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(clients), clientsCount)
		for i, client := range clients {
			test.AssertEqual(t, client.Client(), clientNames[i])
		}
	})
}

func TestServiceProfit(t *testing.T) {
	t.Run("no clients seated", func(t *testing.T) {
		tablesCount := 9
		s := service.NewService(service.NewComputerClub(tablesCount, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))
		clientsCount := tablesCount
		for i := 1; i <= clientsCount; i++ {
			client := fmt.Sprintf("Client%d", i)

			arrivalEvent := event.NewArrivalEvent(dummyDayTime, client)
			gotEvent := s.ServeEvent(arrivalEvent)
			assertNoErrorEvent(t, gotEvent)
		}

		tables, err := s.Profit()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(tables), tablesCount)
		for _, table := range tables {
			test.AssertEqual(t, table.Profit, 0)
			test.AssertEqual(t, table.WorkingTime, 0)
		}
	})

	t.Run("all clients seated and leaved at close", func(t *testing.T) {
		tablesCount := 9
		s := service.NewService(service.NewComputerClub(tablesCount, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))
		clientsCount := tablesCount
		for i := 1; i <= clientsCount; i++ {
			client := fmt.Sprintf("Client%d", i)

			arrivalEvent := event.NewArrivalEvent(dummyDayTime, client)
			gotEvent := s.ServeEvent(arrivalEvent)
			assertNoErrorEvent(t, gotEvent)

			sitEvent := event.NewSitDownEvent(dummyDayTime, client, i)
			gotEvent = s.ServeEvent(sitEvent)
			assertNoErrorEvent(t, gotEvent)

			leaveEvent := event.NewLeaveEvent(dummyCloseTime, client)
			gotEvent = s.ServeEvent(leaveEvent)
			assertNoErrorEvent(t, gotEvent)
		}

		tables, err := s.Profit()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(tables), tablesCount)
		for _, table := range tables {
			workingTime := dummyCloseTime.Sub(dummyOpenTime.Time)
			test.AssertEqual(t, table.Profit, math.Ceil(workingTime.Hours())*dummyMoneyPerHour)
			test.AssertEqual(t, table.WorkingTime, workingTime)
		}
	})

	t.Run("all clients seating unitl closing", func(t *testing.T) {
		tablesCount := 9
		s := service.NewService(service.NewComputerClub(tablesCount, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue()))
		clientsCount := tablesCount
		for i := 1; i <= clientsCount; i++ {
			client := fmt.Sprintf("Client%d", i)

			arrivalEvent := event.NewArrivalEvent(dummyDayTime, client)
			gotEvent := s.ServeEvent(arrivalEvent)
			assertNoErrorEvent(t, gotEvent)

			sitEvent := event.NewSitDownEvent(dummyDayTime, client, i)
			gotEvent = s.ServeEvent(sitEvent)
			assertNoErrorEvent(t, gotEvent)
		}
		_, err := s.Close()
		test.AssertNoError(t, err)

		tables, err := s.Profit()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(tables), tablesCount)
		for _, table := range tables {
			workingTime := dummyCloseTime.Sub(dummyOpenTime.Time)
			test.AssertEqual(t, table.Profit, math.Ceil(workingTime.Hours())*dummyMoneyPerHour)
			test.AssertEqual(t, table.WorkingTime, workingTime)
		}
	})
}

func assertEmptyEvent(t testing.TB, e event.Event) {
	t.Helper()
	if !event.IsEmpty(e) {
		t.Fatalf("got error event: %#v", e)
	}
}

func assertNoErrorEvent(t testing.TB, e event.Event) {
	t.Helper()
	err, ok := e.(*event.ErrorEvent)
	if ok {
		t.Fatalf("got error event: %q", err.Err())
	}
}

func assertErrorEvent(t testing.TB, e event.Event, err error) {
	t.Helper()
	errEvent, ok := e.(*event.ErrorEvent)
	if !ok {
		t.Fatalf("expected error event got: %#v", e)
	}
	test.AssertError(t, errEvent.Err(), err)
}
