package idutil

import "tyto/core/memutil"

const (
	// 无效id值
	INVALID_ID int32 = 0
)

// ID池
type Pool struct {
	nextId       int32                     // 下一个可分配的ID
	maxId        int32                     // 最大ID值
	minFreeCount int32                     // 最小空闲ID数量
	freeList     memutil.RingBuffer[int32] // 空闲ID列表
}

// 创建一个ID池，ID范围为[1, maxId)
// minFreeCount: 最小空闲ID数量，空闲列表中的ID超过此数量则直接重用旧ID
func NewPool(maxId int32, minFreeCount int32) *Pool {
	if maxId <= 0 {
		return nil
	}

	if minFreeCount > maxId {
		minFreeCount = 0
	}

	return &Pool{
		nextId:       1,
		maxId:        maxId,
		minFreeCount: minFreeCount,
		freeList:     *memutil.NewRingBuffer[int32](minFreeCount * 2),
	}
}

func (p *Pool) exhausted() bool {
	return p.nextId >= p.maxId && p.freeList.Len() == 0
}

func (p *Pool) reuse() int32 {
	id, ok := p.freeList.Pop()
	if !ok {
		return INVALID_ID
	}
	return id
}

// 获取一个ID
func (p *Pool) Get() int32 {
	if p.exhausted() {
		// 没有可分配的id
		return INVALID_ID
	}

	if p.nextId >= p.maxId {
		// 没有可分配的id，但有可重用的id
		return p.reuse()
	}

	if p.minFreeCount > 0 && p.freeList.Len() > p.minFreeCount {
		// 超过最小空闲id数量，直接重用旧id
		return p.reuse()
	}

	// 分配新id
	id := p.nextId
	p.nextId++

	return id
}

// 释放一个ID
func (p *Pool) Put(id int32) {
	if id <= 0 {
		return
	}
	p.freeList.Push(id)
}
