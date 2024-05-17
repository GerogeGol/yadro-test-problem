package service

import (
	"errors"
	"fmt"

	"github.com/GerogeGol/yadro-test-problem/domain/queue"
	"github.com/GerogeGol/yadro-test-problem/domain/store"
)

var YouShallNotPass = errors.New("YouShallNotPass")
var NotOpenYet = errors.New("NotOpenYet")
var PlaceIsBusy = errors.New("PlaceIsBusy")
var ClientUnknown = errors.New("ClientUnknown")
var ICanWaitNoLonger = errors.New("ICanWaitNoLonger!")

type ComputerClub struct {
	ComputerCount int
	busyComputers int
	MoneyPerHour  float64
	OpenTime      store.DayTime
	CloseTime     store.DayTime
	store         store.Store
	queue         queue.Queue
}

func NewComputerClub(computerCount int, moneyPerHour float64, openTime store.DayTime, closeTime store.DayTime, store store.Store, queue queue.Queue) *ComputerClub {
	return &ComputerClub{
		ComputerCount: computerCount,
		OpenTime:      openTime,
		CloseTime:     closeTime,
		MoneyPerHour:  moneyPerHour,
		store:         store,
		queue:         queue,
	}
}

func (cc *ComputerClub) Arrive(t store.DayTime, client string) error {
	if t.Compare(cc.OpenTime.Time) == -1 || t.Compare(cc.CloseTime.Time) >= 0 {
		return NotOpenYet
	}

	exists, err := cc.store.IsClientExists(client)
	if err != nil {
		return fmt.Errorf("ComputerClub.Arrive: %w", err)
	}

	if exists {
		return YouShallNotPass
	}

	cc.store.AddClient(client)
	return nil
}

func (cc *ComputerClub) SitDown(t store.DayTime, clientName string, tableNumber int) error {
	if tableNumber <= 0 || tableNumber > cc.ComputerCount {
		return fmt.Errorf("ComputerClub.SitDown: incorrect tableNumber")
	}

	exists, err := cc.store.IsClientExists(clientName)
	if err != nil {
		return fmt.Errorf("ComputerClub.SitDown: %w", err)
	}
	if !exists {
		return ClientUnknown
	}

	isBusy, err := cc.store.IsTableBusy(tableNumber)
	if err != nil {
		return fmt.Errorf("ComputerClub.SitDown: %w", err)
	}

	if isBusy {
		return PlaceIsBusy
	}

	client, err := cc.store.Client(clientName)
	if err != nil {
		return fmt.Errorf("ComputerClub.SitDown: %w", err)
	}

	if client.Table == 0 {
		cc.setClientTable(t, clientName, tableNumber)

	} else {
		cc.changeClientTable(t, clientName, tableNumber)
	}
	return nil
}

func (cc *ComputerClub) Wait(t store.DayTime, clientName string) (bool, error) {
	if cc.queue.Len() >= cc.ComputerCount {
		return false, nil
	}

	if cc.busyComputers < cc.ComputerCount {
		return false, ICanWaitNoLonger
	}

	cc.queue.Push(clientName)
	return true, nil
}

func (cc *ComputerClub) Leave(t store.DayTime, clientName string) (seatedClient store.Client, occupied bool, err error) {
	exists, err := cc.store.IsClientExists(clientName)

	if err != nil {
		err = fmt.Errorf("ComputerClub.SitDown: %w", err)
		return
	}
	if !exists {
		err = ClientUnknown
		return
	}

	client, err := cc.store.Client(clientName)
	if err != nil {
		return
	}

	if err = cc.clientLeave(t, client); err != nil {
		return
	}

	if cc.queue.Len() == 0 {
		cc.busyComputers--
		return
	}

	waitClient, _ := cc.queue.Top()
	if err = cc.queue.Pop(); err != nil {
		return
	}

	if err = cc.setClientTable(t, waitClient, client.Table); err != nil {
		return
	}

	seatedClient.Name = waitClient
	seatedClient.Table = client.Table
	return seatedClient, true, nil
}

func (cc *ComputerClub) Close() ([]store.Client, error) {
	var leavedClients []store.Client
	for cc.queue.Len() != 0 {
		name, _ := cc.queue.Top()

		if err := cc.queue.Pop(); err != nil {
			return nil, err
		}

		if _, _, err := cc.Leave(cc.CloseTime, name); err != nil {
			return nil, err
		}

		leavedClients = append(leavedClients, store.Client{Name: name})
	}

	clients, err := cc.store.Clients()
	if err != nil {
		return nil, err
	}

	for _, client := range clients {
		cc.Leave(cc.CloseTime, client.Name)
		leavedClients = append(leavedClients, client)
	}

	return leavedClients, nil
}

func (cc *ComputerClub) Info(tableNumber int) (store.Table, error) {
	return cc.store.Table(tableNumber)
}

func (cc *ComputerClub) TablesInfo() ([]store.Table, error) {
	var tables []store.Table
	for i := 1; i <= cc.ComputerCount; i++ {
		table, err := cc.Info(i)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (cc *ComputerClub) setClientTable(t store.DayTime, clientName string, tableNumber int) error {
	if err := cc.store.UpdateClientTable(clientName, tableNumber); err != nil {
		return err
	}

	if err := cc.store.UpdateClientPlayingSince(clientName, t); err != nil {
		return err
	}

	if err := cc.store.UpdateTableBusy(tableNumber, true); err != nil {
		return err
	}

	cc.busyComputers++
	return nil
}

func (cc *ComputerClub) changeClientTable(t store.DayTime, clientName string, tableNumber int) error {
	client, err := cc.store.Client(clientName)
	if err != nil {
		return err
	}

	if err = cc.store.UpdateTableBusy(client.Table, false); err != nil {
		return err
	}

	if err = cc.store.UpdateClientTable(clientName, tableNumber); err != nil {
		return err
	}

	if err = cc.store.UpdateTableBusy(tableNumber, true); err != nil {
		return err
	}

	return nil
}

func (cc *ComputerClub) clientLeave(t store.DayTime, client store.Client) error {
	if err := cc.store.RemoveClient(client.Name); err != nil {
		return err
	}

	table, err := cc.store.Table(client.Table)
	if err != nil {
		return err
	}

	playingTime := client.PlayingTime(t)
	payment := client.Payment(t, cc.MoneyPerHour)

	if err = cc.store.UpdateTableBusy(client.Table, false); err != nil {
		return err
	}
	if err = cc.store.UpdateTableWorkingTime(client.Table, table.WorkingTime+playingTime); err != nil {
		return err
	}
	if err = cc.store.UpdateTableProfit(client.Table, table.Profit+payment); err != nil {
		return err
	}
	return err
}
