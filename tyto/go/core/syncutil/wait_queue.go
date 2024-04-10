package syncutil

import (
	"sync"
	"tyto/core/memutil"
)

// 等待队列，会阻塞等待Pop操作
type WaitQueue[T any] struct {
	cond  *sync.Cond
	queue memutil.RingBuffer[T]
}

func NewWaitQueue[T any](initCapacity int32) *WaitQueue[T] {
	return &WaitQueue[T]{
		cond:  sync.NewCond(&SpinLock{}),
		queue: *memutil.NewRingBuffer[T](initCapacity),
	}
}

// 弹出队列头部元素
func (q *WaitQueue[T]) Pop() T {
	q.cond.L.Lock()

	for q.queue.Empty() {
		q.cond.Wait()
	}

	v, _ := q.queue.Pop()

	q.cond.L.Unlock()

	return v
}

// 弹出最多len(vs)个元素
func (q *WaitQueue[T]) PopSome(vs []T) int32 {
	q.cond.L.Lock()

	for q.queue.Empty() {
		q.cond.Wait()
	}

	count := q.queue.PopSome(vs)

	q.cond.L.Unlock()

	return count
}

// 添加元素到队列尾部
func (q *WaitQueue[T]) Push(v T) {
	q.cond.L.Lock()
	q.queue.Push(v)
	q.cond.Signal()
	q.cond.L.Unlock()
}

func (q *WaitQueue[T]) PushAll(vs []T) {
	q.cond.L.Lock()
	q.queue.PushAll(vs)
	q.cond.Signal()
	q.cond.L.Unlock()
}

// 返回队列长度
func (q *WaitQueue[T]) Len() int32 {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.queue.Len()
}

// 判断队列是否为空
func (q *WaitQueue[T]) Empty() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.queue.Empty()
}

// 清空队列
func (q *WaitQueue[T]) Clear() {
	q.cond.L.Lock()
	q.queue.Clear()
	q.cond.Signal()
	q.cond.L.Unlock()
}
