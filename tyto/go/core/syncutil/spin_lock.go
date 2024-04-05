package syncutil

import (
	"math/rand/v2"
	"runtime"
	"sync/atomic"
)

// 对退避计算时，最大的偏移值
const MAX_BACKOFF_SHIFT = 8

// 自旋锁
type SpinLock struct {
	flag atomic.Bool
}

// 加锁
func (lock *SpinLock) Lock() {
	// 冲突次数
	conflictCount := 0

	for {
		for lock.flag.Load() {
			runtime.Gosched()
		}

		if !lock.flag.CompareAndSwap(false, true) {
			shift := conflictCount
			if conflictCount > MAX_BACKOFF_SHIFT {
				shift = MAX_BACKOFF_SHIFT
			}

			conflictCount++

			// 退避
			relaxCount := rand.IntN(1 << shift)
			for i := 0; i < relaxCount; i++ {
				runtime.Gosched()
			}
			continue

		} else {
			// 加锁成功
			break
		}
	}
}

// 尝试加锁
func (lock *SpinLock) TryLock() {
	lock.flag.CompareAndSwap(false, true)
}

// 解锁
func (lock *SpinLock) Unlock() {
	lock.flag.Store(false)
}
