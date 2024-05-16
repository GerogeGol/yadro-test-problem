package store

import (
	"fmt"
	"time"
)

var ClientUnknown = fmt.Errorf("ClientUnknown")

type Store interface {
	AddClient(clientName string) error
	IsClientExists(clientName string) (bool, error)
	UpdateClientTable(clientName string, tableNumber int) error
	UpdateClientPlayingSince(clientName string, t DayTime) error
	Client(clientName string) (Client, error)
	RemoveClient(clientName string) error

	Table(tableNumber int) (Table, error)
	IsTableBusy(tableNumber int) (bool, error)
	UpdateTableBusy(tableNumber int, isBusy bool) error
	UpdateTableWorkingTime(tableNumber int, workingTime time.Duration) error
	UpdateTableProfit(tableNumber int, profit float64) error
}
