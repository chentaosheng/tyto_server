package syncutil

import (
	"sync"
	"tyto/core/memutil"
)

// 等待队列，会阻塞等待Pop操作
type WaitQueue[T any] struct {
	cond *sync.Cond
	buff memutil.RingBuffer[T]
}

func NewWaitQueue[T any](initCapacity int32) *WaitQueue[T] {
	return &WaitQueue[T]{
		cond: sync.NewCond(&SpinLock{}),
		buff: *memutil.NewRingBuffer[T](initCapacity),
	}
}

// 弹出队列头部元素
func (q *WaitQueue[T]) Pop() T {
	q.cond.L.Lock()

	for q.buff.Empty() {
		q.cond.Wait()
	}

	v, _ := q.buff.Pop()

	q.cond.L.Unlock()

	return v
}

// 添加元素到队列尾部
func (q *WaitQueue[T]) Push(v T) {
	q.cond.L.Lock()
	q.buff.Push(v)
	q.cond.L.Unlock()
	q.cond.Signal()
}

// 返回队列长度
func (q *WaitQueue[T]) Len() int32 {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.buff.Len()
}

// 判断队列是否为空
func (q *WaitQueue[T]) Empty() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.buff.Empty()
}

// 清空队列
func (q *WaitQueue[T]) Clear() {
	q.cond.L.Lock()
	q.buff.Clear()
	q.cond.L.Unlock()
	q.cond.Signal()
}
