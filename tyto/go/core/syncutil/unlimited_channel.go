package syncutil

import "tyto/core/memutil"

// 无限容量的channel
type UnlimitedChannel[T any] struct {
	in     chan T
	out    chan T
	buff   memutil.RingBuffer[T]
	length chan int32
}

// initCapacity: 初始容量
func NewUnlimitedChannel[T any](initCapacity int32) *UnlimitedChannel[T] {
	// 只是为了避免channel阻塞
	chanSize := 64

	ch := &UnlimitedChannel[T]{
		in:     make(chan T, chanSize),
		out:    make(chan T, chanSize),
		buff:   *memutil.NewRingBuffer[T](initCapacity),
		length: make(chan int32),
	}

	go ch.run()
	return ch
}

func (ch *UnlimitedChannel[T]) In() chan<- T {
	return ch.in
}

func (ch *UnlimitedChannel[T]) Out() <-chan T {
	return ch.out
}

func (ch *UnlimitedChannel[T]) Close() {
	close(ch.in)
}

func (ch *UnlimitedChannel[T]) Len() int32 {
	return <-ch.length
}

func (ch *UnlimitedChannel[T]) Empty() bool {
	return ch.Len() == 0
}

func (ch *UnlimitedChannel[T]) run() {
	var (
		in   chan T
		out  chan T
		next T
	)

	// init
	in = ch.in

	for in != nil || out != nil {
		select {
		case v, ok := <-in:
			if ok {
				ch.buff.Push(v)
			} else {
				in = nil
			}

		case out <- next:
			ch.buff.Pop()

		case ch.length <- ch.buff.Len():
		}

		if !ch.buff.Empty() {
			next, _ = ch.buff.Front()
			out = ch.out

		} else {
			out = nil
		}
	}

	close(ch.out)
	close(ch.length)
}
