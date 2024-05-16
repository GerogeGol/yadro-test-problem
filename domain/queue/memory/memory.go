package memqueue

import (
	"fmt"

	"github.com/GerogeGol/yadro-test-problem/domain/queue"
)

type Queue struct {
	q []string
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Push(s string) {
	q.q = append(q.q, s)
}

func (q *Queue) Top() (string, bool) {
	if len(q.q) == 0 {
		return "", false
	}
	return q.q[0], true
}

func (q *Queue) Pop() error {
	if len(q.q) == 0 {
		return fmt.Errorf("Queue.Pop: %w", queue.QueueIsEmpty)
	}
	q.q = q.q[1:]
	return nil
}

func (q *Queue) Len() int {
	return len(q.q)
}
