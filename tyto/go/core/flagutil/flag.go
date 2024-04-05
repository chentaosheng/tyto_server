package flagutil

import "golang.org/x/exp/constraints"

// 是否拥有标记
func HasFlag[T constraints.Unsigned](bitset T, flag T) bool {
	return bitset&flag == flag
}

// 添加标记
func AddFlag[T constraints.Unsigned](bitset T, flag T) T {
	return bitset | flag
}

// 移除标记
func RemoveFlag[T constraints.Unsigned](bitset T, flag T) T {
	return bitset & ^flag
}
