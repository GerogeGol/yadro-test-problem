package service_test

import (
	"testing"

	memqueue "github.com/GerogeGol/yadro-test-problem/domain/queue/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/service"
	"github.com/GerogeGol/yadro-test-problem/domain/store"
	memstore "github.com/GerogeGol/yadro-test-problem/domain/store/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/test"
)

var dummyDayTime = test.DummyDayTime
var dummyClient = test.DummyClient
var dummyTableNumber = test.DummyTableNumber
var dummyComputersCount = 2
var dummyMoneyPerHour = 1.0

func TestArrive(t *testing.T) {
	t.Run("client arrived", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())

		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)
	})

	t.Run("client is already in computer club. should return YouShallNotPass", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.Arrive(dummyDayTime, dummyClient)
		test.AssertError(t, err, service.YouShallNotPass)
	})

	t.Run("client arrived but club is closed. should return NotOpenYet", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, store.NewDayTime(9, 0), memstore.NewStore(), memqueue.NewQueue())
		err := club.Arrive(store.NewDayTime(8, 9), dummyClient)
		test.AssertError(t, err, service.NotOpenYet)
	})
}

func TestSitDown(t *testing.T) {
	t.Run("client sat down at a free table", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertNoError(t, err)
	})

	t.Run("client change table to free table", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, 1)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, 2)
		test.AssertNoError(t, err)
	})

	t.Run("client change table to free table. New client takes his previous seat", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
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
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
		err := club.Arrive(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertNoError(t, err)

		err = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertError(t, err, service.PlaceIsBusy)
	})

	t.Run("client change table to an occupied table", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
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
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())

		err := club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)
		test.AssertError(t, err, service.ClientUnknown)
	})
}

func TestWait(t *testing.T) {
	t.Run("client wait", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
		_ = club.Arrive(dummyDayTime, dummyClient)

		isWaiting, err := club.Wait(dummyDayTime, dummyClient)
		test.AssertError(t, err, service.ICanWaitNoLonger)
		test.AssertFalse(t, isWaiting)
	})

	t.Run("client wating for an empty table", func(t *testing.T) {
		club := service.NewComputerClub(1, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())

		_ = club.Arrive(dummyDayTime, dummyClient)
		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		newClient := "NewClient"
		_ = club.Arrive(dummyDayTime, newClient)
		isWaiting, err := club.Wait(dummyDayTime, newClient)
		test.AssertNoError(t, err)
		test.AssertTrue(t, isWaiting)
	})

	t.Run("client no wating for an empty table", func(t *testing.T) {
		club := service.NewComputerClub(1, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())

		_ = club.Arrive(dummyDayTime, dummyClient)
		_ = club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		newClient := "NewClient"
		_ = club.Arrive(dummyDayTime, newClient)
		_, err := club.Wait(dummyDayTime, newClient)
		test.AssertNoError(t, err)

		newClient2 := "newClient2"

		_ = club.Arrive(dummyDayTime, newClient2)
		isWaiting, err := club.Wait(dummyDayTime, newClient2)
		test.AssertNoError(t, err)
		test.AssertFalse(t, isWaiting)
	})
}

func TestLeave(t *testing.T) {
	t.Run("client arrive and leave", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
		_ = club.Arrive(dummyDayTime, dummyClient)

		tableNumber, err := club.Leave(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)
		test.AssertEqual(t, tableNumber, 0)
	})
	t.Run("client arrive and leave", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())

		_, err := club.Leave(dummyDayTime, dummyClient)
		test.AssertError(t, err, service.ClientUnknown)
	})
	t.Run("client arrive and sit and leave", func(t *testing.T) {
		club := service.NewComputerClub(dummyComputersCount, dummyMoneyPerHour, dummyDayTime, memstore.NewStore(), memqueue.NewQueue())
		_ = club.Arrive(dummyDayTime, dummyClient)

		err := club.SitDown(dummyDayTime, dummyClient, dummyTableNumber)

		tableNumber, err := club.Leave(dummyDayTime, dummyClient)
		test.AssertNoError(t, err)
		test.AssertEqual(t, tableNumber, dummyTableNumber)
	})
}
