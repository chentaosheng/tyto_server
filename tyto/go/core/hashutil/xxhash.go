package hashutil

import (
	"github.com/cespare/xxhash/v2"
	"tyto/core/tyto"
)

// 计算字符串的xxhash值
func XxHashSum(ctx tyto.Context, format string, a ...interface{}) string {
	h := xxhash.New()
	return hashSum(ctx, h, format, a...)
}

func XxHashSumToBase64(ctx tyto.Context, format string, a ...interface{}) string {
	h := xxhash.New()
	return hashSumToBase64(ctx, h, format, a...)
}
