package memstore

import (
	"fmt"
	"time"

	"github.com/GerogeGol/yadro-test-problem/domain/store"
)

type MemoryStore struct {
	clients map[string]store.Client
	tables  map[int]store.Table
}

func NewStore() *MemoryStore {
	return &MemoryStore{
		clients: map[string]store.Client{},
		tables:  map[int]store.Table{},
	}
}

func (m *MemoryStore) AddClient(clientName string) error {
	m.clients[clientName] = store.Client{Name: clientName}
	return nil
}

func (m *MemoryStore) RemoveClient(clientName string) error {
	delete(m.clients, clientName)
	return nil
}

func (m *MemoryStore) IsClientExists(clientName string) (bool, error) {
	_, ok := m.clients[clientName]
	return ok, nil
}

func (m *MemoryStore) UpdateClientTable(clientName string, tableNumber int) error {
	if exists, _ := m.IsClientExists(clientName); !exists {
		return fmt.Errorf("MemoryStore.UpdateClientTable: %w", store.ClientUnknown)
	}

	client := m.clients[clientName]
	client.Table = tableNumber
	m.clients[clientName] = client
	return nil
}

func (m *MemoryStore) UpdateClientPlayingSince(clientName string, t store.DayTime) error {
	if exists, _ := m.IsClientExists(clientName); !exists {
		return fmt.Errorf("MemoryStore.UpdateClientTable: %w", store.ClientUnknown)
	}

	client := m.clients[clientName]
	client.PlayingSince = t
	m.clients[clientName] = client
	return nil
}

func (m *MemoryStore) Client(clientName string) (client store.Client, err error) {
	client, ok := m.clients[clientName]
	if !ok {
		return client, fmt.Errorf("MemoryStore.Client: %w", store.ClientUnknown)
	}
	return client, nil
}

func (m *MemoryStore) Table(tableNumber int) (table store.Table, err error) {
	table, _ = m.tables[tableNumber]
	return table, nil
}

func (m *MemoryStore) IsTableBusy(tableNumber int) (bool, error) {
	table, _ := m.tables[tableNumber]
	return table.IsBusy, nil
}

func (m *MemoryStore) UpdateTableBusy(tableNumber int, isBusy bool) error {
	table := m.tables[tableNumber]
	table.IsBusy = isBusy
	m.tables[tableNumber] = table
	return nil
}

func (m *MemoryStore) UpdateTableWorkingTime(tableNumber int, newTime time.Duration) error {
	table := m.tables[tableNumber]
	table.WorkingTime = newTime
	m.tables[tableNumber] = table
	return nil
}

func (m *MemoryStore) UpdateTableProfit(tableNumber int, newProfit float64) error {
	table := m.tables[tableNumber]
	table.Profit = newProfit
	m.tables[tableNumber] = table
	return nil
}
