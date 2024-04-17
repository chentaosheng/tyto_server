package memutil

// 计数为0时，用于回收资源的函数
type RefDestroyer[T any] func(*RefObject[T])

// 引用计数器
type RefObject[T any] struct {
	object    T
	counter   RefCounter
	destroyer RefDestroyer[T]
}

// 创建引用计数器对象
// 注意：初始ref为0，返回后，需要手动调用IncRef()增加引用计数
func NewRefObject[T any](obj T, destroyer RefDestroyer[T]) *RefObject[T] {
	return &RefObject[T]{
		object:    obj,
		counter:   RefCounter{},
		destroyer: destroyer,
	}
}

// 增加引用计数
func (ro *RefObject[T]) IncRef() int32 {
	return ro.counter.IncRef()
}

// 减少引用计数，当为0时，自动释放资源
func (ro *RefObject[T]) DecRef() int32 {
	v := ro.counter.DecRef()
	if v == 0 && ro.destroyer != nil {
		ro.destroyer(ro)
	}
	return v
}

// 获取引用计数
func (ro *RefObject[T]) RefCount() int32 {
	return ro.counter.RefCount()
}

// 获取对象
func (ro *RefObject[T]) Object() *T {
	return &ro.object
}
