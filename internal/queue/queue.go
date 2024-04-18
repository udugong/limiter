package queue

import "time"

//go:generate mockgen -source=./queue.go -package=queuemocks -destination=mocks/queue.mock.go BoundedQueue
type BoundedQueue interface {
	Enqueue(val time.Time) error // 入队
	Dequeue() (time.Time, error) // 出队
	Peek() (time.Time, error)    // 查看对头元素
	IsFull() bool                // 是否队满
}
