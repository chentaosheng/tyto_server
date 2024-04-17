package flagutil

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"unsafe"
)

func checkIndex[TIndex constraints.Integer](index TIndex, max TIndex) bool {
	if index < 0 || index >= max {
		panic(fmt.Errorf("index is out of range, value: %d, max: %d", index, max))
	}
	return true
}

// 第index位是否为1，index从0开始计算
func HasIndex[T constraints.Unsigned, TIndex constraints.Integer](bitset T, index TIndex) bool {
	if !checkIndex(index, TIndex(unsafe.Sizeof(bitset)*8)) {
		return false
	}

	flag := T(1) << index
	return bitset&flag == flag
}

// 将第index位设置为1
func AddIndex[T constraints.Unsigned, TIndex constraints.Integer](bitset T, index TIndex) T {
	if !checkIndex(index, TIndex(unsafe.Sizeof(bitset)*8)) {
		return bitset
	}

	flag := T(1) << index
	return bitset | flag
}

// 将第index位设置为0
func RemoveIndex[T constraints.Unsigned, TIndex constraints.Integer](bitset T, index TIndex) T {
	if !checkIndex(index, TIndex(unsafe.Sizeof(bitset)*8)) {
		return bitset
	}

	flag := T(1) << index
	return bitset & ^flag
}
