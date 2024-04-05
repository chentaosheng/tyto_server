package memutil

// 环形缓冲区，支持自动扩展容量
type RingBuffer[T any] struct {
	buff  []T // 实际的内存缓冲区
	size  int // 当前元素个数
	first int // 数据开始位置、读位置
	last  int // 数据结束位置、写位置
}

func NewRingBuffer[T any](initCapacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buff:  make([]T, initCapacity),
		size:  0,
		first: 0,
		last:  0,
	}
}

// 添加元素到缓冲区尾部
func (rb *RingBuffer[T]) Push(v T) {
	if (rb.size + 1) > len(rb.buff) {
		rb.grow(len(rb.buff) + 1)
	}

	rb.buff[rb.last] = v
	rb.last = (rb.last + 1) % len(rb.buff)
	rb.size++
}

// 从缓冲区头部弹出元素
func (rb *RingBuffer[T]) Pop() (T, bool) {
	if rb.size == 0 {
		var t T
		return t, false
	}

	v := rb.buff[rb.first]
	rb.first = (rb.first + 1) % len(rb.buff)
	rb.size--
	return v, true
}

// 当前元素个数
func (rb *RingBuffer[T]) Len() int {
	return rb.size
}

// 当前缓冲区容量
func (rb *RingBuffer[T]) Cap() int {
	return len(rb.buff)
}

// 判断缓冲区是否为空
func (rb *RingBuffer[T]) Empty() bool {
	return rb.size == 0
}

// 清空缓冲区
func (rb *RingBuffer[T]) Clear() {
	rb.size = 0
	rb.first = 0
	rb.last = 0
}

// 计算新的缓冲区容量
func (rb *RingBuffer[T]) calculateNewCapacity(old, need int) int {
	const (
		MAX_DOUBLE_SIZE   = 2 * 1024
		MAX_ADD_HALF_SIZE = 15 * 1024
	)

	cap := old
	for cap <= need {
		if cap < MAX_DOUBLE_SIZE {
			cap = cap * 2

		} else if cap < MAX_ADD_HALF_SIZE {
			cap = cap + cap/2

		} else {
			cap = cap + cap/4
		}
	}

	return cap
}

// 扩展缓冲区容量
func (rb *RingBuffer[T]) grow(need int) {
	cap := rb.calculateNewCapacity(len(rb.buff), need)
	newBuff := make([]T, cap)

	if rb.first < rb.last {
		copy(newBuff, rb.buff[rb.first:rb.last])
	} else {
		copy(newBuff, rb.buff[rb.first:])
		copy(newBuff[len(rb.buff)-rb.first:], rb.buff[:rb.last])
	}

	rb.buff = newBuff
	rb.first = 0
	rb.last = rb.size
}
