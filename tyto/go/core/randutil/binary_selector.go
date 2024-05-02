package randutil

import (
	"sort"
	"tyto/core/tyto"
)

type binaryElement[T any] struct {
	index  int32 // 在集合中的索引
	weight int32 // 当前元素的权重
	sum    int64 // 从0项到当前项的权重和
	value  T     // 实际的数据
}

// 二分方式进行随机选择，时间复杂度O(log2N)，
// 元素数量较少时，效率较高
// 一般用于掉落物品、抽奖等场景
type BinarySelector[T any] struct {
	elements []binaryElement[T]
	total    int64
	record   *selectRecord
}

func NewBinarySelector[T any]() Selector[T] {
	return &BinarySelector[T]{
		elements: make([]binaryElement[T], 0),
		total:    0,
		record:   nil,
	}
}

func (s *BinarySelector[T]) Add(ctx tyto.Context, value T, weight int32) {
	if weight <= 0 {
		ctx.Logger().Error("invalid weight:", weight)
		return
	}

	s.elements = append(s.elements, binaryElement[T]{
		index:  0,
		weight: weight,
		sum:    0,
		value:  value,
	})
}

// 构建选择列表
func (s *BinarySelector[T]) BuildList() {
	s.total = 0
	Shuffle(s.elements)

	for i, elem := range s.elements {
		s.total += int64(elem.weight)
		elem.sum = s.total
		elem.index = int32(i)
	}

	s.record = newSelectRecord(int32(len(s.elements)))
}

// 允许重复选择
func (s *BinarySelector[T]) Select(ctx tyto.Context) (T, bool) {
	if s.total <= 0 || len(s.elements) <= 0 {
		ctx.Logger().Error("no element for select")
		var zero T
		return zero, false
	}

	rate := Int64N(s.total)

	index := sort.Search(len(s.elements), func(i int) bool {
		return s.elements[i].sum > rate
	})

	return s.elements[index].value, true
}

// 不会重复选择，首次调用前需要调用Reset()方法重置选择记录
func (s *BinarySelector[T]) UniqueSelect(ctx tyto.Context) (T, bool) {
	if s.total <= 0 || len(s.elements) <= 0 || !s.record.CanSelectAny() {
		ctx.Logger().Warn("no element or all selected:", s.elements)
		var zero T
		return zero, false
	}

	const MAX_RETRY_COUNT = 24
	for i := 0; i < MAX_RETRY_COUNT; i++ {
		rate := Int64N(s.total)

		index := int32(sort.Search(len(s.elements), func(i int) bool {
			return s.elements[i].sum > rate
		}))

		if s.record.IsSelected(index) {
			continue
		}

		s.record.Select(index)
		return s.elements[index].value, true
	}

	ctx.Logger().Warn("unique select failed, too few elements to select from:", s.elements)

	return s.forceSelect(ctx)
}

// 重置选择记录
func (s *BinarySelector[T]) Reset() {
	s.record.Reset()
}

// 忽略概率，强制选择一个元素，不会重复选择
// 所以，随机列表不能太小，尽量避免调用这个方法
func (s *BinarySelector[T]) forceSelect(ctx tyto.Context) (T, bool) {
	index := Int32N(int32(len(s.elements)))

	// 随机到的位置
	if !s.record.IsSelected(index) {
		s.record.Select(index)
		return s.elements[index].value, true
	}

	for i := int32(1); i < int32(len(s.elements)); i++ {
		left := index - i
		right := index + i

		// 没有可选的元素
		if left < 0 && right >= int32(len(s.elements)) {
			break
		}

		// 左边的元素
		if left >= 0 && !s.record.IsSelected(left) {
			s.record.Select(left)
			return s.elements[left].value, true
		}

		// 右边的元素
		if right < int32(len(s.elements)) && !s.record.IsSelected(right) {
			s.record.Select(right)
			return s.elements[right].value, true
		}
	}

	ctx.Logger().Error("force select failed, no available element:", s.elements)

	// 没有可选的元素
	var zero T
	return zero, false
}
