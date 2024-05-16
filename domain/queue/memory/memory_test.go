package memqueue_test

import (
	"testing"

	"github.com/GerogeGol/yadro-test-problem/domain/queue"
	memqueue "github.com/GerogeGol/yadro-test-problem/domain/queue/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/test"
)

func TestQueue(t *testing.T) {
	q := memqueue.NewQueue()

	t.Run("push", func(t *testing.T) {
		q.Push("val")
		val, ok := q.Top()

		test.AssertTrue(t, ok)
		test.AssertEqual(t, val, "val")
	})
	t.Run("pop", func(t *testing.T) {
		err := q.Pop()
		test.AssertNoError(t, err)

		_, ok := q.Top()
		test.AssertFalse(t, ok)

		err = q.Pop()
		test.AssertError(t, err, queue.QueueIsEmpty)
	})
}
