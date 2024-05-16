package store

import "time"

type Table struct {
	IsBusy      bool
	WorkingTime time.Duration
	Profit      float64
}

func (t *Table) AddWorkingTime(d time.Duration) {
	t.WorkingTime += d
}

func (t *Table) AddProfit(p float64) {
	t.Profit += p
}
