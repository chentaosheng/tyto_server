package memutil

// 环形缓冲区，支持自动扩展容量
type RingBuffer[T any] struct {
	buff []T   // 实际的内存缓冲区
	size int32 // 当前元素个数
	head int32 // 数据开始位置、读位置
	tail int32 // 数据结束位置、写位置
}

func NewRingBuffer[T any](initCapacity int32) *RingBuffer[T] {
	return &RingBuffer[T]{
		buff: make([]T, initCapacity),
		size: 0,
		head: 0,
		tail: 0,
	}
}

// 添加元素到缓冲区尾部
func (rb *RingBuffer[T]) Push(v T) {
	if (rb.size + 1) > rb.Cap() {
		rb.grow(rb.size + 1)
	}

	rb.buff[rb.tail] = v
	rb.tail = (rb.tail + 1) % rb.Cap()
	rb.size++
}

// 批量添加元素到缓冲区尾部
func (rb *RingBuffer[T]) PushAll(vs []T) {
	if len(vs) <= 0 {
		return
	}

	if (rb.size + int32(len(vs))) > rb.Cap() {
		rb.grow(rb.size + int32(len(vs)))
	}

	if rb.tail < rb.head {
		copy(rb.buff[rb.tail:], vs)
		rb.tail = (rb.tail + int32(len(vs))) % rb.Cap()

	} else {
		remain := rb.Cap() - rb.tail
		if int32(len(vs)) <= remain {
			copy(rb.buff[rb.tail:], vs)
			rb.tail = (rb.tail + int32(len(vs))) % rb.Cap()

		} else {
			copy(rb.buff[rb.tail:], vs[:remain])
			copy(rb.buff, vs[remain:])
			rb.tail = int32(len(vs)) - remain
		}
	}

	rb.size += int32(len(vs))
}

// 从缓冲区头部读取元素
func (rb *RingBuffer[T]) Pop() (T, bool) {
	if rb.size == 0 {
		var zero T
		return zero, false
	}

	v := rb.buff[rb.head]
	rb.head = (rb.head + 1) % rb.Cap()
	rb.size--
	return v, true
}

// 从缓充区读取最多len(vs)个元素
func (rb *RingBuffer[T]) PopSome(vs []T) int32 {
	count := int32(len(vs))
	if count <= 0 {
		return 0
	}

	if count >= rb.size {
		count = rb.size
	}

	if rb.head < rb.tail {
		copy(vs, rb.buff[rb.head:rb.head+count])
		rb.head = (rb.head + count) % rb.Cap()

	} else {
		remain := rb.Cap() - rb.head
		if count <= remain {
			copy(vs, rb.buff[rb.head:rb.head+count])
			rb.head = (rb.head + count) % rb.Cap()

		} else {
			copy(vs, rb.buff[rb.head:])
			copy(vs[remain:], rb.buff[:count-remain])
			rb.head = count - remain
		}
	}

	rb.size -= count
	return count
}

// 当前元素个数
func (rb *RingBuffer[T]) Len() int32 {
	return rb.size
}

// 当前缓冲区容量
func (rb *RingBuffer[T]) Cap() int32 {
	return int32(len(rb.buff))
}

// 判断缓冲区是否为空
func (rb *RingBuffer[T]) Empty() bool {
	return rb.size == 0
}

// 清空缓冲区
func (rb *RingBuffer[T]) Clear() {
	rb.size = 0
	rb.head = 0
	rb.tail = 0
}

// 计算新的缓冲区容量
func (rb *RingBuffer[T]) calculateNewCapacity(old, need int32) int32 {
	const (
		MAX_DOUBLE_SIZE   = 2 * 1024
		MAX_ADD_HALF_SIZE = 15 * 1024
	)

	capacity := old
	for capacity <= need {
		if capacity < MAX_DOUBLE_SIZE {
			capacity = capacity * 2

		} else if capacity < MAX_ADD_HALF_SIZE {
			capacity = capacity + capacity/2

		} else {
			capacity = capacity + capacity/4
		}
	}

	return capacity
}

// 扩展缓冲区容量
func (rb *RingBuffer[T]) grow(need int32) {
	capacity := rb.calculateNewCapacity(rb.Cap(), need)
	newBuff := make([]T, capacity)

	if rb.head < rb.tail {
		copy(newBuff, rb.buff[rb.head:rb.tail])
	} else {
		copy(newBuff, rb.buff[rb.head:])
		copy(newBuff[rb.Cap()-rb.head:], rb.buff[:rb.tail])
	}

	rb.buff = newBuff
	rb.head = 0
	rb.tail = rb.size
}
