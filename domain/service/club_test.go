package service_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	memqueue "github.com/GerogeGol/yadro-test-problem/domain/queue/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/service"
	"github.com/GerogeGol/yadro-test-problem/domain/store"
	memstore "github.com/GerogeGol/yadro-test-problem/domain/store/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/test"
)

var dummyDayTime = test.DummyDayTime
var dummyOpenTime = test.DummyDayTime
var dummyCloseTime = store.NewDayTime(19, 0)
var dummyClient = test.DummyClient
var dummyTableNumber = test.DummyTableNumber
var dummyComputersCount = 2
var dummyMoneyPerHour = 1.0

func TestArrive(t *testing.T) {
	t.Run("client arrived", func(t *testing.T) {
		club := dummyClub()

		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)
	})

	t.Run("client is already in computer club. should return YouShallNotPass", func(t *testing.T) {
		club := dummyClub()
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)
		err = club.Arrive(dummyDayTime, dummyClient)
		test.AssertError(t, err, service.YouShallNotPass)

	})

	t.Run("client arrived but club is closed. should return NotOpenYet", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, store.NewDayTime(9, 0), store.NewDayTime(19, 0), memstore.NewStore(), memqueue.NewQueue())
		err := club.Arrive(store.NewDayTime(8, 9), dummyClient)
		test.AssertError(t, err, service.NotOpenYet)

		err = club.Arrive(store.NewDayTime(19, 0), dummyClient)
		test.AssertError(t, err, service.NotOpenYet)
	})
}

func TestSitDown(t *testing.T) {
	t.Run("client sat down at a free table", func(t *testing.T) {
		club := dummyClub()
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertNoError(t, err)
	})

	t.Run("client change table to free table", func(t *testing.T) {
		club := dummyClub()
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, 1)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, 2)
		test.AssertNoError(t, err)
	})

	t.Run("client change table to free table. New client takes his previous seat", func(t *testing.T) {
		club := dummyClub()
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, 1)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, 2)
		test.AssertNoError(t, err)

		newClient := "NewClient"
		err = club.Arrive(dummyDayTime, newClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, 1)
		test.AssertNoError(t, err)
	})

	t.Run("client sat down at an occupied table", func(t *testing.T) {
		club := dummyClub()
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertError(t, err, service.PlaceIsBusy)
	})

	t.Run("client change table to an occupied table", func(t *testing.T) {
		club := dummyClub()
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertNoError(t, err)

		newClient := "NewClient"

		err = club.Arrive(dummyDayTime, newClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, newClient, 2)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, newClient, dummyTableNumber)
		test.AssertError(t, err, service.PlaceIsBusy)
	})

	t.Run("client not in a club", func(t *testing.T) {
		club := dummyClub()

		err := club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertError(t, err, service.ClientUnknown)
	})
}

func TestWait(t *testing.T) {
	t.Run("client wait", func(t *testing.T) {
		club := dummyClub()
		_ = club.Arrive(dummyDayTime, dummyClient)

		isWaiting, err := club.Wait(dummyDayTime, dummyClient)
		test.AssertError(t, err, service.ICanWaitNoLonger)
		test.AssertFalse(t, isWaiting)
	})

	t.Run("client wating for an empty table", func(t *testing.T) {
		club := service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue())

		_ = club.Arrive(dummyDayTime, dummyClient)
		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		newClient := "NewClient"
		_ = club.Arrive(dummyDayTime, newClient)
		isWaiting, err := club.Wait(dummyDayTime, newClient)
		test.AssertNoError(t, err)
		test.AssertTrue(t, isWaiting)
	})

	t.Run("client no wating for an empty table", func(t *testing.T) {
		club := service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue())

		_ = club.Arrive(dummyDayTime, dummyClient)
		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		newClient := "NewClient"
		_ = club.Arrive(dummyDayTime, newClient)
		isWaiting, err := club.Wait(dummyDayTime, newClient)
		test.AssertTrue(t, isWaiting)

		newClient2 := "newClient2"

		_ = club.Arrive(dummyDayTime, newClient2)
		isWaiting, err = club.Wait(dummyDayTime, newClient2)
		test.AssertNoError(t, err)
		test.AssertFalse(t, isWaiting)
	})
}

func TestLeave(t *testing.T) {
	t.Run("client arrive and leave", func(t *testing.T) {
		club := dummyClub()
		_ = club.Arrive(dummyDayTime, dummyClient)

		_, _, err := club.Leave(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)
	})

	t.Run("client leave", func(t *testing.T) {
		club := dummyClub()

		_, _, err := club.Leave(dummyDayTime, dummyClient)
		test.AssertError(t, err, service.ClientUnknown)
	})
	t.Run("client arrive and sit and leave. hus place taken by other client", func(t *testing.T) {
		club := dummyClub()
		_ = club.Arrive(dummyDayTime, dummyClient)

		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		_, occupied, err := club.Leave(dummyCloseTime, dummyClient)
		test.AssertNoError(t, err)
		test.AssertFalse(t, occupied)

		table, err := club.Info(dummyTableNumber)
		test.AssertNoError(t, err)
		playingTime := dummyCloseTime.Sub(dummyDayTime.Time)
		test.AssertEqual(t, table.WorkingTime, playingTime)
		test.AssertEqual(t, table.Profit, math.Ceil(playingTime.Hours())*dummyMoneyPerHour)

		newClient := "newClient"
		_ = club.Arrive(dummyDayTime, newClient)

		err = club.SitDown(dummyDayTime, newClient, dummyTableNumber)
		test.AssertNoError(t, err)
	})
	t.Run("client leave while one is waiting", func(t *testing.T) {
		club := service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue())
		_ = club.Arrive(dummyDayTime, dummyClient)
		err := club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		newClient := "NewClient"

		_ = club.Arrive(dummyDayTime, newClient)
		wait, err := club.Wait(dummyDayTime, newClient)
		test.AssertNoError(t, err)
		test.AssertTrue(t, wait)

		client, occupied, err := club.Leave(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)
		test.AssertTrue(t, occupied)
		test.AssertEqual(t, client.Name, newClient)
		test.AssertEqual(t, client.Table, dummyTableNumber)
	})
}

func TestClose(t *testing.T) {
	t.Run("all clients leaved before closing", func(t *testing.T) {
		club := dummyClub()

		_ = club.Arrive(dummyDayTime, dummyClient)
		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		_, _, _ = club.Leave(dummyDayTime, dummyClient)

		leftClients, err := club.Close()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(leftClients), 0)
	})
	t.Run("all clients stayed to closing", func(t *testing.T) {
		club := dummyClub()

		_ = club.Arrive(dummyDayTime, dummyClient)
		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		leftClients, err := club.Close()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(leftClients), 1)
	})
	t.Run("all playing and arrived clients stayed to closing", func(t *testing.T) {
		club := service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue())

		_ = club.Arrive(dummyDayTime, dummyClient)
		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		newClient := "NewClient"
		_ = club.Arrive(dummyDayTime, newClient)

		leftClients, err := club.Close()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(leftClients), 2)
	})
	t.Run("all playing, arrived and waiting clients stayed to closing", func(t *testing.T) {
		club := service.NewComputerClub(1, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue())

		_ = club.Arrive(dummyDayTime, dummyClient)
		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		newClient := "NewClient"
		_ = club.Arrive(dummyDayTime, newClient)

		newClient2 := "NewClient2"
		_ = club.Arrive(dummyDayTime, newClient2)
		wait, err := club.Wait(dummyDayTime, newClient2)
		test.AssertNoError(t, err)
		test.AssertTrue(t, wait)

		leftClients, err := club.Close()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(leftClients), 3)
	})
}

func TestProfit(t *testing.T) {
	t.Run("zero profit", func(t *testing.T) {
		computersCount := 10
		club := service.NewComputerClub(computersCount, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue())
		tables, err := club.TablesInfo()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(tables), computersCount)

		for _, table := range tables {
			test.AssertEqual(t, table.Profit, 0)
			test.AssertEqual(t, table.WorkingTime, time.Duration(0))
		}
	})

	t.Run("max profit", func(t *testing.T) {
		computersCount := 10
		club := service.NewComputerClub(computersCount, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue())
		for i := 1; i <= computersCount; i++ {
			client := fmt.Sprintf("Client %d", i)
			_ = club.Arrive(dummyOpenTime, client)
			_ = club.SitDown(dummyOpenTime, client, i)
		}

		for i := 1; i <= computersCount; i++ {
			client := fmt.Sprintf("Client %d", i)
			_, _, _ = club.Leave(dummyCloseTime, client)
		}

		tables, err := club.TablesInfo()
		test.AssertNoError(t, err)
		test.AssertEqual(t, len(tables), computersCount)

		workingTime := (dummyCloseTime.Time.Sub(dummyOpenTime.Time))
		for _, table := range tables {
			test.AssertEqual(t, table.Profit, math.Ceil(workingTime.Hours())*dummyMoneyPerHour)
			test.AssertEqual(t, table.WorkingTime, workingTime)
		}
	})
}

func dummyClub() *service.ComputerClub {
	return service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyOpenTime, dummyCloseTime, memstore.NewStore(), memqueue.NewQueue())
}
