package typeutil

import "golang.org/x/exp/constraints"

// 返回两者中的最小值
func Min[T constraints.Integer](a, b T) T {
	if a <= b {
		return a
	}
	return b
}

// 返回三者中的最小值
func Min3[T constraints.Integer](a, b, c T) T {
	if a <= b {
		if a <= c {
			return a
		}
		return c

	} else {
		if b <= c {
			return b
		}
		return c
	}
}
