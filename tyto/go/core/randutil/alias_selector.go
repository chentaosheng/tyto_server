package randutil

import "tyto/core/tyto"

// 无效别名索引
const INVALID_ALIAS_INDEX int32 = -1

type aliasElement[T any] struct {
	index  int32 // 在集合中的索引
	weight int32 // 当前元素的权重
	prob   int64 // 当前元素的概率
	value  T     // 实际的数据
}

// 别名法进行随机选择，时间复杂度O(1)
// 元素数量与性能无关，效率稳定
// 一般用于掉落物品、抽奖等场景
type AliasSelector[T any] struct {
	elements []aliasElement[T]
	alias    []int32
	average  int64
	record   *selectRecord
}

func NewAliasSelector[T any]() Selector[T] {
	return &AliasSelector[T]{
		elements: make([]aliasElement[T], 0),
		alias:    make([]int32, 0),
		average:  0,
		record:   nil,
	}
}

func (s *AliasSelector[T]) Add(ctx tyto.Context, value T, weight int32) {
	if weight <= 0 {
		ctx.Logger().Error("invalid weight:", weight)
		return
	}

	s.elements = append(s.elements, aliasElement[T]{
		index:  0,
		weight: weight,
		value:  value,
		prob:   0,
	})
}

// 构建选择列表
func (s *AliasSelector[T]) BuildList() {
	Shuffle(s.elements)

	s.record = newSelectRecord(int32(len(s.elements)))

	var (
		minWeight int32 = -1
		maxWeight int32 = -1
	)

	s.average = 0

	for i, elem := range s.elements {
		elem.index = int32(i)
		// 放大权重，避免计算时有小数
		elem.prob = int64(elem.weight) * int64(len(s.elements))
		// 放大后，平均值即为weight总和
		s.average += int64(elem.weight)
		// 初始化alias
		s.alias = append(s.alias, INVALID_ALIAS_INDEX)

		// 计算最小值
		if minWeight < 0 || elem.weight < minWeight {
			minWeight = elem.weight
		}

		// 计算最大值
		if maxWeight < 0 || elem.weight > maxWeight {
			maxWeight = elem.weight
		}
	}

	if minWeight == maxWeight {
		// 所有元素权重相同，无需计算
		return
	}

	var (
		small []int32
		large []int32
	)

	// 计算alias
	// 将所有元素分为两组，一组权重小于平均值，一组权重大于等于平均值
	for i, elem := range s.elements {
		if elem.prob < s.average {
			small = append(small, int32(i))
		} else {
			large = append(large, int32(i))
		}
	}

	// 依次填充alias
	for len(small) > 0 && len(large) > 0 {
		less := small[0]
		small = small[1:]

		more := large[0]
		large = large[1:]

		s.alias[less] = more
		s.elements[more].prob -= s.average - s.elements[less].prob

		if s.elements[more].prob < s.average {
			small = append(small, more)
		} else {
			large = append(large, more)
		}
	}

	// 对剩余的元素进行修正
	for len(small) > 0 {
		less := small[0]
		small = small[1:]
		s.elements[less].prob = s.average
	}
	for len(large) > 0 {
		more := large[0]
		large = large[1:]
		s.elements[more].prob = s.average
	}
}

// 允许重复选择
func (s *AliasSelector[T]) Select(ctx tyto.Context) (T, bool) {
	if s.average <= 0 || len(s.elements) <= 0 {
		ctx.Logger().Error("no element for select")
		var zero T
		return zero, false
	}

	index := Int32N(int32(len(s.elements)))
	aliasIndex := s.alias[index]

	// 该组元素只有一个
	if aliasIndex == INVALID_ALIAS_INDEX {
		return s.elements[index].value, true
	}

	// 该组元素有两个
	rand := Int64N(s.average)
	// 概率在index范围内，则选择index，否则选择aliasIndex
	newIndex := aliasIndex
	if rand < s.elements[index].prob {
		newIndex = index
	}

	return s.elements[newIndex].value, true
}

// 不会重复选择，首次调用前需要调用Reset()方法重置选择记录
func (s *AliasSelector[T]) UniqueSelect(ctx tyto.Context) (T, bool) {
	if s.average <= 0 || len(s.elements) <= 0 || !s.record.CanSelectAny() {
		ctx.Logger().Error("no element or all selected:", s.elements)
		var zero T
		return zero, false
	}

	const MAX_RETRY_COUNT = 24
	for i := 0; i < MAX_RETRY_COUNT; i++ {
		index := Int32N(int32(len(s.elements)))
		aliasIndex := s.alias[index]

		// 该组元素只有一个
		if aliasIndex == INVALID_ALIAS_INDEX {
			if s.record.IsSelected(index) {
				continue
			}

			s.record.Select(index)
			return s.elements[index].value, true
		}

		// 该组元素有两个
		rand := Int64N(s.average)
		// 概率在index范围内，则选择index，否则选择aliasIndex
		newIndex := aliasIndex
		if rand < s.elements[index].prob {
			newIndex = index
		}

		if s.record.IsSelected(newIndex) {
			continue
		}

		s.record.Select(newIndex)
		return s.elements[newIndex].value, true
	}

	ctx.Logger().Warn("unique select failed, too few elements to select from:", s.elements)

	return s.forceSelect(ctx)
}

// 重置选择记录
func (s *AliasSelector[T]) Reset() {
	s.record.Reset()
}

// 忽略概率，强制选择一个元素，不会重复选择
// 所以，随机列表不能太小，尽量避免调用这个方法
func (s *AliasSelector[T]) forceSelect(ctx tyto.Context) (T, bool) {
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
