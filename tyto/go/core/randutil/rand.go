package randutil

import (
	"fmt"
	"math/rand/v2"
)

// 范围[min, max)
func Int32R(min, max int32) int32 {
	if min >= max {
		panic(fmt.Errorf("random range error, min: %d, max: %d", min, max))
	}

	return rand.Int32N(max-min) + min
}

// 范围[0, max)
func Int32N(max int32) int32 {
	return rand.Int32N(max)
}

func Int32() int32 {
	return rand.Int32()
}

// 范围[min, max)
func Uint32R(min, max uint32) uint32 {
	if min >= max {
		panic(fmt.Errorf("random range error, min: %d, max: %d", min, max))
	}

	return rand.Uint32N(max-min) + min
}

// 范围[0, max)
func Uint32N(max uint32) uint32 {
	return rand.Uint32N(max)
}

func Uint32() uint32 {
	return rand.Uint32()
}

// 范围[min, max)
func Int64R(min, max int64) int64 {
	if min >= max {
		panic(fmt.Errorf("random range error, min: %d, max: %d", min, max))
	}

	return rand.Int64N(max-min) + min
}

// 范围[0, max)
func Int64N(max int64) int64 {
	return rand.Int64N(max)
}

func Int64() int64 {
	return rand.Int64()
}

// 范围[min, max)
func Uint64R(min, max uint64) uint64 {
	if min >= max {
		panic(fmt.Errorf("random range error, min: %d, max: %d", min, max))
	}

	return rand.Uint64N(max-min) + min
}

// 范围[0, max)
func Uint64N(max uint64) uint64 {
	return rand.Uint64N(max)
}

func Uint64() uint64 {
	return rand.Uint64()
}

func Float32() float32 {
	return rand.Float32()
}

func Float64() float64 {
	return rand.Float64()
}
