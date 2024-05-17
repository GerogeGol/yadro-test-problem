package memstore_test

import (
	"testing"
	"time"

	"github.com/GerogeGol/yadro-test-problem/domain/store"
	memstore "github.com/GerogeGol/yadro-test-problem/domain/store/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/test"
)

var dummyClient = test.DummyClient
var dummyDayTime = test.DummyDayTime
var dummyTableNumber = test.DummyTableNumber

func TestAddClient(t *testing.T) {

	t.Run("add client to store", func(t *testing.T) {
		m := memstore.NewStore()

		_ = m.AddClient(dummyClient)
		exists, _ := m.IsClientExists(dummyClient)

		test.AssertTrue(t, exists)
	})
}

func TestRemoveClient(t *testing.T) {
	t.Run("remove client to ", func(t *testing.T) {
		m := memstore.NewStore()

		_ = m.AddClient(dummyClient)
		exists, _ := m.IsClientExists(dummyClient)
		test.AssertTrue(t, exists)

		_ = m.RemoveClient(dummyClient)
		exists, _ = m.IsClientExists(dummyClient)
		test.AssertFalse(t, exists)
	})
}

func TestUpdateClientTable(t *testing.T) {
	t.Run("update client table", func(t *testing.T) {
		m := memstore.NewStore()

		_ = m.AddClient(dummyClient)

		err := m.UpdateClientTable(dummyClient, dummyTableNumber)
		test.AssertNoError(t, err)

		client, err := m.Client(dummyClient)
		test.AssertNoError(t, err)
		test.AssertEqual(t, client.Table, dummyTableNumber)
	})
	t.Run("update client table, but client does not exists. Should return ClientDoesNotExists error", func(t *testing.T) {
		m := memstore.NewStore()

		err := m.UpdateClientTable(dummyClient, dummyTableNumber)
		test.AssertError(t, err, store.ClientDoesNotExist)
	})
}

func TestUpdateClientPlayingSince(t *testing.T) {
	t.Run("update client table", func(t *testing.T) {
		m := memstore.NewStore()

		_ = m.AddClient(dummyClient)

		playingSince := store.NewDayTime(10, 0)
		err := m.UpdateClientPlayingSince(dummyClient, playingSince)
		test.AssertNoError(t, err)

		client, err := m.Client(dummyClient)
		test.AssertNoError(t, err)
		test.AssertEqual(t, client.PlayingSince, playingSince)
	})
	t.Run("update client table, but client does not exists. Should return ClientDoesNotExists error", func(t *testing.T) {
		m := memstore.NewStore()

		err := m.UpdateClientPlayingSince(dummyClient, dummyDayTime)
		test.AssertError(t, err, store.ClientDoesNotExist)
	})
}

func TestUpdateTableBusy(t *testing.T) {
	t.Run("update table busy", func(t *testing.T) {
		m := memstore.NewStore()

		busy, _ := m.IsTableBusy(dummyTableNumber)
		test.AssertFalse(t, busy)

		_ = m.UpdateTableBusy(dummyTableNumber, true)
		busy, _ = m.IsTableBusy(dummyTableNumber)
		test.AssertTrue(t, busy)

		_ = m.UpdateTableBusy(dummyTableNumber, false)
		busy, _ = m.IsTableBusy(dummyTableNumber)
		test.AssertFalse(t, busy)
	})
}

func TestUpdateTableWorkingTime(t *testing.T) {
	t.Run("update table time", func(t *testing.T) {
		m := memstore.NewStore()

		addDur := time.Duration(10 * time.Hour)
		_ = m.UpdateTableWorkingTime(dummyTableNumber, addDur)

		table, _ := m.Table(dummyTableNumber)
		test.AssertEqual(t, table.WorkingTime, addDur)
	})
}

func TestUpdateTableProfit(t *testing.T) {
	t.Run("update table time", func(t *testing.T) {
		m := memstore.NewStore()

		addProfit := 100.0
		_ = m.UpdateTableProfit(dummyTableNumber, addProfit)

		table, _ := m.Table(dummyTableNumber)
		test.AssertEqual(t, table.Profit, addProfit)
	})
}
