package memutil

import "sync/atomic"

// 引用计数器
type RefCounter[T any] struct {
	obj       T
	ref       atomic.Int32
	destroyer func(*RefCounter[T])
}

// 创建引用计数器对象
// 注意：初始ref为0，返回后，需要手动调用IncRef()增加引用计数
func NewRefCounter[T any, F ~func(*RefCounter[T])](obj T, destroyer F) *RefCounter[T] {
	return &RefCounter[T]{
		obj:       obj,
		ref:       atomic.Int32{},
		destroyer: destroyer,
	}
}

// 增加引用计数
func (rc *RefCounter[T]) IncRef() int32 {
	return rc.ref.Add(1)
}

// 减少引用计数，当为0时，自动释放资源
func (rc *RefCounter[T]) DecRef() int32 {
	v := rc.ref.Add(-1)
	if v == 0 && rc.destroyer != nil {
		rc.destroyer(rc)
	}
	return v
}

// 获取引用计数
func (rc *RefCounter[T]) RefCount() int32 {
	return rc.ref.Load()
}

// 获取对象
func (rc *RefCounter[T]) Object() *T {
	return &rc.obj
}
