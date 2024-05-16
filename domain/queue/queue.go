package queue

import "fmt"

var QueueIsEmpty = fmt.Errorf("queue is empty")

type Queue interface {
	Pop() error
	Push(value string)
	Top() (string, bool)
	Len() int
}
