package randutil

import "github.com/bits-and-blooms/bitset"

// 选择记录，用于selector
type selectRecord struct {
	bitset      bitset.BitSet
	selectCount int32
	maxCount    int32
}

func newSelectRecord(maxCount int32) *selectRecord {
	return &selectRecord{
		bitset:      *bitset.New(uint(maxCount)),
		selectCount: 0,
		maxCount:    maxCount,
	}
}

func (s *selectRecord) SelectCount() int32 {
	return s.selectCount
}

func (s *selectRecord) IsSelected(index int32) bool {
	if index < 0 || index >= s.maxCount {
		return false
	}

	return s.bitset.Test(uint(index))
}

func (s *selectRecord) CanSelectAny() bool {
	return s.selectCount < s.maxCount
}

func (s *selectRecord) Select(index int32) {
	if index < 0 || index >= s.maxCount {
		return
	}

	if s.bitset.Test(uint(index)) {
		return
	}

	s.bitset.Set(uint(index))
	s.selectCount++
}

func (s *selectRecord) Reset() {
	s.bitset.ClearAll()
	s.selectCount = 0
}
