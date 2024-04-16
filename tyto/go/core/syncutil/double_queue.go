package syncutil

import (
	"sync"
	"tyto/core/memutil"
)

// 双缓冲队列，适合多生产者单消费者场景
type DoubleQueue[T any] struct {
	cond       *sync.Cond
	writeQueue *memutil.RingBuffer[T]
	readQueue  *memutil.RingBuffer[T]
}

func NewDoubleQueue[T any](initCapacity int32) *DoubleQueue[T] {
	return &DoubleQueue[T]{
		cond:       sync.NewCond(&SpinLock{}),
		writeQueue: memutil.NewRingBuffer[T](initCapacity),
		readQueue:  memutil.NewRingBuffer[T](initCapacity),
	}
}

// 弹出队列头部元素
func (q *DoubleQueue[T]) Pop() T {
	if q.readQueue.Empty() {
		q.cond.L.Lock()
		for q.writeQueue.Empty() {
			q.cond.Wait()
		}

		q.readQueue, q.writeQueue = q.writeQueue, q.readQueue

		q.cond.L.Unlock()
	}

	v, _ := q.readQueue.Pop()

	return v
}

// 添加元素到队列尾部
func (q *DoubleQueue[T]) Push(v T) {
	q.cond.L.Lock()
	q.writeQueue.Push(v)
	q.cond.Signal()
	q.cond.L.Unlock()
}

func (q *DoubleQueue[T]) PushAll(vs []T) {
	q.cond.L.Lock()
	q.writeQueue.PushAll(vs)
	q.cond.Signal()
	q.cond.L.Unlock()
}
