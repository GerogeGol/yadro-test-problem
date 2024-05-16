package store_test

import (
	"testing"
	"time"

	"github.com/GerogeGol/yadro-test-problem/domain/store"
	"github.com/GerogeGol/yadro-test-problem/domain/test"
)

func TestTable(t *testing.T) {
	t.Run("add working time", func(t *testing.T) {
		table := store.Table{}

		table.AddWorkingTime(time.Duration(10 * time.Hour))
		table.AddWorkingTime(time.Duration(20 * time.Second))

		want := time.Duration(10*time.Hour + 20*time.Second)
		test.AssertEqual(t, table.WorkingTime, want)
	})

	t.Run("add profit", func(t *testing.T) {
		table := store.Table{}

		table.AddProfit(10)
		table.AddProfit(100)

		want := 110.0
		test.AssertEqual(t, table.Profit, want)
	})
}
