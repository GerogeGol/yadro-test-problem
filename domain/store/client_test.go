package store_test

import (
	"testing"

	"github.com/GerogeGol/yadro-test-problem/domain/store"
	"github.com/GerogeGol/yadro-test-problem/domain/test"
)

func TestClient(t *testing.T) {
	t.Run("get payment value", func(t *testing.T) {
		client := store.Client{PlayingSince: store.NewDayTime(0, 0)}
		end := store.NewDayTime(10, 0)

		payment := client.Payment(end, 1)
		test.AssertEqual(t, payment, 10)
	})
}
