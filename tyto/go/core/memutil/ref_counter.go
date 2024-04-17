package memutil

import "sync/atomic"

// 引用计数器
type RefCounter struct {
	ref atomic.Int32
}

// 增加引用计数
func (rc *RefCounter) IncRef() int32 {
	return rc.ref.Add(1)
}

// 减少引用计数
func (rc *RefCounter) DecRef() int32 {
	return rc.ref.Add(-1)
}

// 获取引用计数
func (rc *RefCounter) RefCount() int32 {
	return rc.ref.Load()
}
