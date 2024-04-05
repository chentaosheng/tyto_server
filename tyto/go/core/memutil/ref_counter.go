package memutil

import "sync/atomic"

// 引用计数器
type RefCounter[T any] struct {
	obj       T
	ref       atomic.Int32
	destroyer func(*T)
}

// 创建引用计数器对象
func NewRefCounter[T any](obj *T, destroyer func(*T)) *RefCounter[T] {
	return &RefCounter[T]{
		obj:       *obj,
		ref:       atomic.Int32{},
		destroyer: destroyer,
	}
}

// 增加引用计数
func (c *RefCounter[T]) IncRef() int32 {
	return c.ref.Add(1)
}

// 减少引用计数，当为0时，自动释放资源
func (c *RefCounter[T]) DecRef() int32 {
	v := c.ref.Add(-1)
	if v == 0 && c.destroyer != nil {
		c.destroyer(&c.obj)
	}
	return v
}

// 获取引用计数
func (c *RefCounter[T]) RefCount() int32 {
	return c.ref.Load()
}

// 获取对象
func (c *RefCounter[T]) Object() *T {
	return &c.obj
}
